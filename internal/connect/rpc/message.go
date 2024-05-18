package rpc

import (
	"context"
	"encoding/json"
	"mim/internal/connect/websocket"
	"mim/pkg/proto"
)

func (cr *ConnectRpc) PushMessage(ctx context.Context, req *proto.PushMessageReq, resp *proto.PushMessageResp) error {
	c, ok := websocket.Default.GetUser(req.TargetID)
	if !ok {
		// 用户突然下线
		// 应该告诉logic层转存离线消息
		resp.IsOffline = true
	}

	msg, _ := json.Marshal(req)
	c.Channel <- msg
	
	return nil
}
