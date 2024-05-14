package websocket

import (
	"encoding/json"
	"mim/setting"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	TypePong = 1 + iota
	TypeSingle
	TypeGroup
	TypeAck
)

// client维持用户的ws连接
type Client struct {
	Conn      *websocket.Conn
	channel   chan []byte
	done      chan struct{}
	ID        int64
	Username  string
	server    *Server
	next      *Client
	pre       *Client
	HeartBeat time.Time
	lock      sync.RWMutex
	IsUse     bool
}

type Message struct {
	Body   []byte
	Seq    int
	Type   int // 1单聊 2群聊
	From   int64
	Target int64
}

func NewClient(id int64, username string, size int) *Client {
	return &Client{
		channel: make(chan []byte, size),
		ID: id,
		Username: username,
		HeartBeat: time.Now(),
	}
}

func (c *Client) writeProc() {
	// 发送消息协程
	// 1. 心跳检测
	ticker := time.NewTicker(setting.Conf.WsConfig.TickerPeriod)
	defer func() {
		c.done <- struct{}{}
		ticker.Stop()
		c.server.getBucket(c.ID).RemoveClient(c.ID)
		c.Conn.Close()
	}()
	// 2. 接收消息

	for {
		select {
		case msg, ok := <-c.channel:
			// 接收到消息之后，设置响应时间
			c.Conn.SetWriteDeadline(time.Now().Add(setting.Conf.WsConfig.WriteDeadline))
			if !ok {
				zap.L().Error("write message to client failed, ", zap.Any("client", c.ID))
				c.Conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			var err error
			// 失败重试
			for i := 0; i < setting.Conf.WsConfig.MaxRetries; i++ {
				c.lock.Lock()
				err = c.Conn.WriteMessage(websocket.BinaryMessage, msg)
				if err == nil {
					return
				}
				c.lock.Unlock()
			}

			if err != nil {
				zap.L().Error("write message failed: ", zap.Error(err), zap.Any("client", c.ID))
				return
			}

			c.HeartBeat = time.Now()
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(setting.Conf.WsConfig.WriteDeadline))
			c.lock.Lock()
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				zap.L().Error("write ping message failed: ", zap.Error(err), zap.Any("client", c.ID))
				return
			}
			c.lock.Unlock()
		case <-c.done:
			zap.L().Error("write routine was closed by read routine")
			return
		}
	}
}

func (c *Client) readProc() {
	defer func() {
		c.server.getBucket(c.ID).RemoveClient(c.ID)
		c.Conn.Close()
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			_, msg, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.done <- struct{}{}
					zap.L().Error("read message from client failed: ", zap.Error(err), zap.Any("client", c.ID))
					return
				}
			}

			c.HandleMessage(msg)
		}
	}
}
func (c *Client) HandleMessage(msg []byte) {
	if len(msg) == 0 {
		return
	}

	m := Message{}
	if err := json.Unmarshal(msg, &m); err != nil {
		zap.L().Error("unmarshal message failed: ", zap.Error(err), zap.Any("client", c.ID), zap.Any("msg content", msg))
		c.done <- struct{}{}
		return
	}

	zap.L().Debug("message content: ", zap.Any("msg", msg), zap.Any("client", c.ID))
	uid := m.Target

	target, ok := c.server.GetUser(uid)
	if !ok {
		// relay
	}
	target.channel <- m.Body
}
