package websocket

// Route ws路由
type Route struct {
	Method  string
	Handler HandleFunc
}

// HandleFunc 处理函数
type HandleFunc func(server *Server, conn *Conn, message *Message)
