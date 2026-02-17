package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync/atomic"
)

type wsConn struct {
	conn   *websocket.Conn
	closed atomic.Bool
}

func (c *wsConn) WriteMessage(mt int, data []byte) error {
	if c.closed.Load() {
		return errors.New("connection closed")
	}
	return c.conn.WriteMessage(mt, data)
}

// Close 关闭ws连接
func (c *wsConn) Close() error {
	if c.closed.CompareAndSwap(false, true) {
		return c.conn.Close()

	}
	// 已经关闭，避免重复关闭
	return nil
}
