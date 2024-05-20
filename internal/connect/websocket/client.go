package websocket

import (
	"encoding/json"
	logicrpc "mim/internal/connect/rpc/logic_rpc"
	"mim/pkg/proto"
	"mim/setting"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// client维持用户的ws连接
type Client struct {
	Conn      *websocket.Conn
	Channel   chan []byte
	done      chan struct{}
	ID        int64
	Username  string
	server    *Server
	HeartBeat time.Time
	lock      sync.RWMutex
	IsUse     bool
}

func NewClient(id int64, username string, size int) *Client {
	return &Client{
		Channel:   make(chan []byte, size),
		ID:        id,
		Username:  username,
		HeartBeat: time.Now(),
		done:      make(chan struct{}),
	}
}

func (c *Client) writeProc() {
	// 发送消息协程
	// 1. 心跳检测
	ticker := time.NewTicker(setting.Conf.WsConfig.TickerPeriod)
	defer func() {
		c.done <- struct{}{}
		ticker.Stop()
		c.offline()
	}()
	// 2. 接收消息

	for {
		select {
		case msg, ok := <-c.Channel:
			// 接收到消息之后，设置响应时间
			c.Conn.SetWriteDeadline(time.Now().Add(setting.Conf.WsConfig.WriteDeadline))
			if !ok {
				zap.L().Error("write message to client failed, ", zap.Any("client", c.ID))
				c.Conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			zap.L().Info("read msg from channel", zap.Any("msg: ", msg))

			var err error
			// TODO 失败重试
			c.lock.Lock()
			err = c.Conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				zap.L().Error("write message failed: ", zap.Error(err))
				return 
			}
			c.lock.Unlock()

			c.HeartBeat = time.Now()
		case <-ticker.C:
			// c.Conn.SetWriteDeadline(time.Now().Add(setting.Conf.WsConfig.WriteDeadline))
			// c.lock.Lock()
			// if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			// 	zap.L().Error("write ping message failed: ", zap.Error(err), zap.Any("client", c.ID))
			// 	return
			// }
			// c.lock.Unlock()
			zap.L().Info("ticker ping")
		case <-c.done:
			zap.L().Error("write routine was closed by read routine")
			return
		}
	}
}

func (c *Client) readProc() {
	defer func() {
		c.offline()
	}()

	for {
		select {
		case <-c.done:
			zap.L().Error("read routine was closed by write routine")
			return
		default:
			_, msg, err := c.Conn.ReadMessage()
			if err != nil {
				zap.L().Error("read message from client failed: ", zap.Error(err), zap.Any("client", c.ID))
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.done <- struct{}{}
					return
				}
			}

			c.handleMessage(msg)
		}
	}
}

func (c *Client) handleMessage(msg []byte) {
	if len(msg) == 0 {
		return
	}

	req := &proto.MessageReq{}
	if err := json.Unmarshal(msg, &req); err != nil {
		zap.L().Error("unmarshal message failed: ", zap.Error(err), zap.Any("client", c.ID), zap.Any("msg content", msg))
		c.done <- struct{}{}
		return
	}

	logicrpc.SendMessage(req)
}

func (c *Client) offline() {
	req := &proto.OfflineReq{
		UserID: c.ID,
	}

	zap.L().Info("send offline request", zap.Any("client:", req.UserID))
	if err := logicrpc.Offline(req); err != nil {
		zap.L().Error("failed to notify server of offline status: ", zap.Error(err), zap.Any("client", c.ID))
	}

	// 释放资源
	c.server.getBucket(c.ID).RemoveClient(c.ID)
	c.Conn.WriteJSON("close conn")
	c.Conn.Close()
}
