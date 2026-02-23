package iface

import (
	"github.com/gorilla/websocket"
	"time"
)

type WsConn interface {
	Start()
	Stop()
	GetConn() *websocket.Conn
	GetConnId() uint32
	Send(data []byte) error
	GetIdle() time.Time
}

// HandleFunc 连接业务处理函数
type HandleFunc func(conn *websocket.Conn, data []byte, length int) error
