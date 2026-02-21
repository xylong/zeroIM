package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/jsonx"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultReadLimit  = 512 * 1024 // 单条消息最大 512KB，防止恶意大包
	defaultPingPeriod = 30 * time.Second
	defaultPongWait   = 60 * time.Second
)

// Conn ws连接
type Conn struct {
	conn   *websocket.Conn
	closed atomic.Bool // 幂等关闭

	server            *Server
	mu                sync.Mutex
	idle              atomic.Int64 // 最后活跃时间,使用atomic.Int64，无任何 mutex
	maxConnectionIdle time.Duration

	sender chan *Message
	done   chan struct{}
}

func NewConn(server *Server, w http.ResponseWriter, r *http.Request) *Conn {
	wsConn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.Errorf("upgrade http conn error: %v", err)
		return nil
	}

	c := &Conn{
		conn:              wsConn,
		server:            server,
		sender:            make(chan *Message, 100),
		done:              make(chan struct{}),
		maxConnectionIdle: defaultMaxConnectionIdle,
	}

	c.touch()
	wsConn.SetPingHandler(func(appData string) error {
		c.touch()
		return nil // gorilla 自动回复 pong
	})

	go c.keepalive()
	go c.readPump()
	go c.writePump()
	return c
}

// keepalive 心跳检测
func (c *Conn) keepalive() {
	ticker := time.NewTicker(time.Second * 30) // todo 动态配置
	defer ticker.Stop()

	idleTimeout := c.maxConnectionIdle.Nanoseconds()
	for {
		select {
		case <-ticker.C:
			if time.Now().UnixNano()-c.idle.Load() > idleTimeout {
				c.Close()
				return
			}
		case <-c.done:
			return
		}
	}
}

// readPump 读取消息
func (c *Conn) readPump() {
	defer c.Close()
	c.conn.SetReadLimit(defaultReadLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(defaultPongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.touch()
		_ = c.conn.SetReadDeadline(time.Now().Add(defaultPongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNoStatusReceived,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway) {
				c.server.Errorf("unexpected ws close: %v", err)
			}
			return
		}

		c.touch()
		var message Message
		if err := jsonx.Unmarshal(data, &message); err != nil {
			c.server.Errorf("ws unmarshal message err: %v, msg: %s", err, string(data))
			continue
		}

		switch message.FrameType {
		case FramePing:
			_ = c.Send(&Message{
				FrameType: FramePing,
			})
		case FrameData:
			handler, ok := c.server.routes[message.Method]
			if ok {
				handler(c.server, c, &message)
			} else {
				if err := c.Send(&Message{
					FrameType: FrameData,
					Method:    message.Method,
					Data:      "method not found",
				}); err != nil {
					c.server.Errorf("send message err: %v", err)
				}
			}
		}
	}
}

// writePump 发送消息(串行写)
func (c *Conn) writePump() {
	defer c.Close()

	for {
		select {
		case msg := <-c.sender:
			data, _ := jsonx.Marshal(msg)
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
			c.touch()
		case <-c.done:
			return
		}
	}
}

// Send 发送消息
func (c *Conn) Send(message *Message) error {
	if c.closed.Load() {
		return errors.New("connection closed")
	}

	select {
	case c.sender <- message:
		return nil
	case <-c.done:
		return errors.New("connection closed")
	}
}

// Close 关闭ws连接
func (c *Conn) Close() {
	if !c.closed.CompareAndSwap(false, true) {
		return
	}
	close(c.done)
	_ = c.conn.Close()
	c.server.cleanupConn(c)
}

// touch 更新最后活跃时间
func (c *Conn) touch() {
	c.idle.Store(time.Now().UnixNano())
}

// getIdle 获取最后活跃时间
func (c *Conn) getIdle() time.Time {
	return time.Unix(0, c.idle.Load())
}
