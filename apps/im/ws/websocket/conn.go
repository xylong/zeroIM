package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const defaultReadLimit = 512 * 1024 // 单条消息最大 512KB，防止恶意大包

// Conn ws连接
type Conn struct {
	conn   *websocket.Conn
	closed atomic.Bool // 是否关闭

	server            *Server
	mu                sync.Mutex
	idle              time.Time
	maxConnectionIdle time.Duration
	done              chan struct{}
}

func NewConn(server *Server, w http.ResponseWriter, r *http.Request) *Conn {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.Errorf("upgrade http conn error: %v", err)
		return nil
	}

	conn.SetReadLimit(defaultReadLimit)

	return &Conn{
		conn:              conn,
		server:            server,
		idle:              time.Now(),
		done:              make(chan struct{}),
		maxConnectionIdle: defaultMaxConnectionIdle,
	}
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer idleTimer.Stop()

	for {
		select {
		case <-idleTimer.C:
			c.mu.Lock()
			idle := c.idle
			if idle.IsZero() {
				c.mu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}

			val := c.maxConnectionIdle - time.Since(idle)
			c.mu.Unlock()
			if val <= 0 {
				c.server.Close(c)
				return
			}
		case <-c.done:
			return
		}
	}
}

// ReadMessage 读取消息
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.conn.ReadMessage()
	c.idle = time.Now()
	return
}

// WriteMessage 发送消息
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	if c.closed.Load() {
		return errors.New("connection closed")
	}
	err := c.conn.WriteMessage(messageType, data)
	c.idle = time.Now()
	return err
}

// Close 关闭ws连接
func (c *Conn) Close() {
	close(c.done)
	c.server.Close(c)
}
