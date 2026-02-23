package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
	"zeroIM/apps/im/ws/internal/websocket/iface"
)

const (
	defaultReadLimit  = 512 * 1024 // 单条消息最大 512KB，防止恶意大包
	defaultPingPeriod = 30 * time.Second
	defaultPongWait   = 60 * time.Second
)

type Conn struct {
	// websocket连接
	WC *websocket.Conn

	// 当前连接的id
	ConnId uint32

	// 当前连接状态
	closed atomic.Bool

	// 当前连接绑定的处理业务方法
	handleAPI iface.HandleFunc

	// 当前连接最后一次活跃时间
	idle atomic.Int64

	// 退出信号
	ExitChan chan struct{}
}

func NewConn(conn *websocket.Conn, connId uint32, callback iface.HandleFunc) *Conn {
	c := &Conn{
		WC:        conn,
		ConnId:    connId,
		handleAPI: callback,
		ExitChan:  make(chan struct{}, 1),
	}

	return c
}

func (c *Conn) readPump() {
	defer c.Stop()

	for {
		// 读取消息
		_, data, err := c.WC.ReadMessage()
		if err != nil {
			fmt.Println("read data err", err)
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNoStatusReceived,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway) {
				continue
			}
		}

		// 调用handle方法
		if err := c.handleAPI(c.WC, data, len(data)); err != nil {
			fmt.Printf("connid=%v,handle err=%v \n", c.ConnId, err)
			break
		}
	}
}

func (c *Conn) Start() {
	fmt.Printf("conn start,connId=%v \n", c.ConnId)
	go c.readPump()
}

func (c *Conn) Stop() {
	if !c.closed.CompareAndSwap(false, true) {
		return
	}
	_ = c.WC.Close()
	close(c.ExitChan)
}

func (c *Conn) GetConn() *websocket.Conn {
	return nil
}

func (c *Conn) GetConnId() uint32 {
	return 0
}

func (c *Conn) Send(data []byte) error {
	return nil
}

func (c *Conn) touch() {
	c.idle.Store(time.Now().UnixNano())
}

func (c *Conn) GetIdle() time.Time {
	return time.Unix(0, c.idle.Load())
}
