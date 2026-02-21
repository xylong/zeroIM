package websocket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	sync.RWMutex

	routes map[string]HandleFunc // 路由
	addr   string

	opt            *serverOption
	authentication Authentication
	pattern        string

	connToUser map[*Conn]string
	userToConn map[string]*Conn

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

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

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
	conn := NewConn(s, w, r)
	if conn == nil {
		http.Error(w, "websocket upgrade failed", http.StatusInternalServerError)
		return
	}

	s.addConn(conn, uid)
}

func (s *Server) addConn(conn *Conn, userId string) {
	s.Lock()
	oldConn, hadOld := s.userToConn[userId]
	s.connToUser[conn] = userId
	s.userToConn[userId] = conn
	s.Unlock()

	// 同用户重复登录时关闭旧连接并从映射移除
	if hadOld && oldConn != nil && oldConn != conn {
		s.cleanupConn(oldConn)
	}
}

// GetConn 获取单个用户的连接
func (s *Server) GetConn(uid string) *Conn {
	s.RLock()
	defer s.RUnlock()
	return s.userToConn[uid]
}

// GetConnections 获取多个用户的连接（可能包含 nil，表示该用户未连接）
func (s *Server) GetConnections(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}
	s.RLock()
	defer s.RUnlock()
	connections := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		connections = append(connections, s.userToConn[uid])
	}
	return connections
}

// GetUid 获取单个连接的 uid
func (s *Server) GetUid(conn *Conn) string {
	s.RLock()
	defer s.RUnlock()
	return s.connToUser[conn]
}

// GetUserIds 获取多个连接对应的 uid
func (s *Server) GetUserIds(cons ...*Conn) []string {
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

	conns := s.GetConnections(userIds...)
	for _, conn := range conns {
		if err := conn.Send(message); err != nil {
			return err
		}
	}

	return nil
}

// cleanupConn 删除连接
func (s *Server) cleanupConn(conn *Conn) {
	// 安全检查
	if conn == nil {
		return
	}

	// 删除映射(不阻塞锁尽早释放，可以防止新的消息被路由到已关闭的连接)
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
