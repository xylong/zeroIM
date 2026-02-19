package websocket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
)

const defaultReadLimit = 512 * 1024 // 单条消息最大 512KB，防止恶意大包

type Server struct {
	sync.RWMutex

	routes map[string]HandleFunc // 路由
	addr   string

	authentication Authentication
	pattern        string

	connToUser map[*websocket.Conn]string
	userToConn map[string]*websocket.Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, options ...ServerOption) *Server {
	sp := newServerOption(options...)

	return &Server{
		routes: make(map[string]HandleFunc),
		addr:   addr,

		authentication: sp.auth,
		pattern:        sp.pattern,

		connToUser: make(map[*websocket.Conn]string),
		userToConn: make(map[string]*websocket.Conn),

		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		Logger: logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	// 先鉴权，防止握手消耗
	if !s.authentication.Auth(w, r) {
		http.Error(w, "ws auth failed，access denied", http.StatusUnauthorized)
		return
	}
	uid := s.authentication.UserId(r)
	if uid == "" {
		http.Error(w, "user id missing", http.StatusForbidden)
		return
	}

	// 升级为 websocket 连接
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("websocket upgrade error: %v", err)
		return
	}

	conn.SetReadLimit(defaultReadLimit)
	s.addConn(conn, uid)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.Errorf("handleConn panic: %v", err)
			}
			s.Close(conn) // 统一清理
		}()

		s.handleConn(conn)
	}()
}

func (s *Server) addConn(conn *websocket.Conn, userId string) {

	s.Lock()
	oldConn, hadOld := s.userToConn[userId]
	s.connToUser[conn] = userId
	s.userToConn[userId] = conn
	s.Unlock()

	// 同用户重复登录时关闭旧连接并从映射移除
	if hadOld && oldConn != nil && oldConn != conn {
		s.Close(oldConn)
	}
}

// GetConn 获取单个用户的连接
func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RLock()
	defer s.RUnlock()
	return s.userToConn[uid]
}

// GetConnections 获取多个用户的连接（可能包含 nil，表示该用户未连接）
func (s *Server) GetConnections(uids ...string) []*websocket.Conn {
	if len(uids) == 0 {
		return nil
	}
	s.RLock()
	defer s.RUnlock()
	connections := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		connections = append(connections, s.userToConn[uid])
	}
	return connections
}

// GetUid 获取单个连接的 uid
func (s *Server) GetUid(conn *websocket.Conn) string {
	s.RLock()
	defer s.RUnlock()
	return s.connToUser[conn]
}

// GetUserIds 获取多个连接对应的 uid
func (s *Server) GetUserIds(cons ...*websocket.Conn) []string {
	if len(cons) == 0 {
		return nil
	}
	s.RLock()
	defer s.RUnlock()

	userIds := make([]string, 0, len(cons))
	for _, conn := range cons {
		userIds = append(userIds, s.connToUser[conn])
	}
	return userIds
}

func (s *Server) GetAllUserIds() []string {
	s.RLock()
	defer s.RUnlock()
	userIds := make([]string, 0, len(s.connToUser))
	for _, uid := range s.connToUser {
		userIds = append(userIds, uid)
	}
	return userIds
}

func (s *Server) SendByUserId(message *Message, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}

	return s.Send(message, s.GetConnections(userIds...)...)
}

// Send 发送消息
func (s *Server) Send(message *Message, conns ...*websocket.Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := jsonx.Marshal(message)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

// Close 关闭单个连接
// todo：使用原子状态标记（推荐），不能简单地“先删映射再 Close”，如果 Close 阻塞或失败，连接已从映射移除，无法重试清理，给每个连接附加一个 关闭状态标志，配合 sync.Once 确保只关闭一次
func (s *Server) Close(conn *websocket.Conn) {
	// 安全检查
	if conn == nil {
		return
	}

	// 删除映射(不阻塞锁尽早释放，可以防止新的消息被路由到已关闭的连接)
	{
		s.Lock()
		uid, exists := s.connToUser[conn]
		// 连接未注册或已被清理，直接返回
		if !exists {
			s.Unlock()
			return
		}
		delete(s.connToUser, conn)

		// 防止空 uid 导致误删
		if uid != "" {
			if existingConn, ok := s.userToConn[uid]; ok && existingConn == conn {
				delete(s.userToConn, uid)
			}
		}
		s.Unlock()
	}

	// 关闭连接，可能耗时(锁外执行，避免持有锁时阻塞)
	if err := conn.Close(); err != nil {
		s.Errorf("ws close conn err: %v", err)
	}
}

func (s *Server) handleConn(conn *websocket.Conn) {
	defer s.Close(conn) // 无论正常退出还是异常退出，都从映射中移除并关闭连接

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("ws read message err: %v", err)
			return
		}
		var message Message
		if err = jsonx.Unmarshal(msg, &message); err != nil {
			s.Errorf("ws unmarshal message err: %v, msg: %s", err, string(msg))
			return
		}
		s.RLock()
		handler, ok := s.routes[message.Method]
		s.RUnlock()
		if ok {
			handler(s, conn, &message)
		} else {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("ws no route %v", message.Method)))
		}
	}
}

// AddRoutes 添加路由（应在 Start 之前调用，或与 handleConn 并发安全）
func (s *Server) AddRoutes(routes []Route) {
	s.Lock()
	defer s.Unlock()
	for _, route := range routes {
		s.routes[route.Method] = route.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.pattern, s.ServerWs)
	if err := http.ListenAndServe(s.addr, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Errorf("ws server listen err: %v", err)
	}
}

func (s *Server) Stop() {
	fmt.Println("stop")
}
