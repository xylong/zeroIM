package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"
	"zeroIM/apps/im/ws/internal/svc"
	iface2 "zeroIM/apps/im/ws/internal/websocket/iface"
)

type Server struct {
	Name string
	IP   string
	Port int

	upper *websocket.Upgrader
	auth  iface2.Authentication
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewServer(opts ...ServerOption) iface2.WsServer {
	sp := newServerOption(opts...)
	server := &Server{
		Name: "ws server",
		IP:   "0.0.0.0",
		Port: 10090,

		upper: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		auth: sp.auth,
	}
	return server
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	// 先鉴权，防止握手消耗
	if !s.auth.Auth(w, r) {
		http.Error(w, "ws auth failed，access denied", http.StatusUnauthorized)
		return
	}
	uid := s.auth.UserId(r)
	if uid == "" {
		http.Error(w, "user id missing", http.StatusForbidden)
		return
	}

	// 升级为 websocket 连接
	wsConn, err := s.upper.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade http conn error: %v", err)
		return
	}
	if wsConn == nil {
		http.Error(w, "websocket upgrade failed", http.StatusInternalServerError)
		return
	}

	conn := NewConn(wsConn, uint32(time.Now().UnixMilli()), nil)
	go conn.Start()
}

// Start 启动服务
func (s *Server) Start() {
	http.HandleFunc("/ws", s.Serve)
	if err := http.ListenAndServe(s.getAddr(), nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Errorf("ws server listen err: %v", err)
		return
	}
}

func (s *Server) getAddr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

func (s *Server) Stop() {

}

func (s *Server) Run() {
	s.Start()
}
