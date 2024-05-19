package rpc

import (
	"context"
	"encoding/json"
	"mim/internal/connect/websocket"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func (cr *ConnectRpc) PushMessage(ctx context.Context, req *proto.PushMessageReq, resp *proto.PushMessageResp) error {
	zap.L().Info("ws server pushMessage() receive message", zap.Any("msg", req))
	c, ok := websocket.Default.GetUser(req.TargetID)
	if !ok {
		// 用户突然下线
		// 应该告诉logic层转存离线消息
		resp.IsOffline = true
		return nil
	}

	msg, _ := json.Marshal(req)
	c.Channel <- msg

	return nil
}
