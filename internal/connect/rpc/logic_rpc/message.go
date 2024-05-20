package logicrpc

import (
	"context"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func SendMessage(req *proto.MessageReq) {
	resp := &proto.MessageResp{}
	// zap.L().Info("logic client receive message: ", zap.Any("msg", req))
	if err := logicPRpc.Call(context.Background(), "SendMessage", req, resp); err != nil {
		zap.L().Error("SendMessage() failed: ", zap.Error(err))
		return 
	}
}
