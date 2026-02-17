package websocket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

type Server struct {
	sync.RWMutex

	routes map[string]HandleFunc // 路由
	addr   string

	authentication Authentication
	connToUser     map[*websocket.Conn]string
	userToConn     map[string]*websocket.Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		routes: make(map[string]HandleFunc),
		addr:   addr,

		authentication: new(authentication),
		connToUser:     make(map[*websocket.Conn]string),
		userToConn:     make(map[string]*websocket.Conn),

		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Logger: logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			s.Errorf("server handle websocket recover err: %v", err)
		}
	}()

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade: %v", err)
		return
	}

	go s.handleConn(conn)
}

func (s *Server) addConn(conn *websocket.Conn, r *http.Request) {
	uid := s.authentication.UserId(r)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// GetConn 获取单个用户的连接
func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	return s.userToConn[uid]
}

// GetConnections 获取多个用户的连接
func (s *Server) GetConnections(uids ...string) []*websocket.Conn {
	if len(uids) == 0 {
		return nil
	}
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	connections := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		connections = append(connections, s.userToConn[uid])
	}
	return connections
}

// GetUid 获取单个用户的uid
func (s *Server) GetUid(conn *websocket.Conn) string {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	return s.connToUser[conn]
}

// GetUids 获取多个用户的uid
func (s *Server) GetUids(conns ...*websocket.Conn) []string {
	if len(conns) == 0 {
		return nil
	}
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	var uids []string
	for _, conn := range conns {
		uids = append(uids, s.connToUser[conn])
	}
	return uids
}

// Close 关闭单个连接
// todo：使用原子状态标记（推荐），不能简单地“先删映射再 Close”，如果 Close 阻塞或失败，连接已从映射移除，无法重试清理，给每个连接附加一个 关闭状态标志，配合 sync.Once 确保只关闭一次
func (s *Server) Close(conn *websocket.Conn) {
	// 安全检查
	if conn == nil {
		return
	}

	// 先删除映射，不阻塞锁
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

	// 再关闭连接(锁外执行，避免持有锁时阻塞)
	if err := conn.Close(); err != nil {
		s.Errorf("ws close conn err: %v", err)
	}
}

func (s *Server) handleConn(conn *websocket.Conn) {
	for {
		// 读取消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("ws read message err: %v", err)
			return
		}
		// 解析消息
		var message Message
		if err = jsonx.Unmarshal(msg, &message); err != nil {
			s.Errorf("ws unmarshal message err: %v, msg", err, string(msg))
			return
		}
		// 按路由处理
		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, &message)
		} else {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("ws no route %v", message.Method)))
		}
	}
}

// AddRoutes 添加路由
func (s *Server) AddRoutes(routes []Route) {
	for _, route := range routes {
		s.routes[route.Method] = route.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.ServerWs)
	_ = http.ListenAndServe(s.addr, nil)
}

func (s *Server) Stop() {
	fmt.Println("stop")
}
