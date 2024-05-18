package logicrpc

import (
	"context"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func SendMessage(req *proto.MessageReq) {
	resp := &proto.MessageResp{}

	if err := logicPRpc.Call(context.Background(), "SendMessage", req, resp); err != nil {
		zap.L().Error("SendMessage() failed: ", zap.Error(err))
		return 
	}
}
