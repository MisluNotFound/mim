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
		resp.IsOffline = true
		return nil
	}

	msg, err := json.Marshal(req)
	if err != nil {
		zap.L().Error("marshal message failed: ", zap.Error(err))
		return err
	}
	c.Channel <- msg

	return nil
}
