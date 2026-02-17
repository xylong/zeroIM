package websocket

import "github.com/gorilla/websocket"

// Route ws路由
type Route struct {
	Method  string
	Handler HandleFunc
}

type HandleFunc func(server *Server, conn *websocket.Conn, message *Message)
