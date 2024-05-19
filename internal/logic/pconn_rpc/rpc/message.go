package prpc

import (
	"context"
	"mim/internal/logic/pconn_rpc/messaging"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

// 接收connect层发送的消息 并交给message层处理
func (r *PRpc) SendMessage(ctx context.Context, req *proto.MessageReq, resp *proto.MessageResp) error {
	zap.L().Info("logic server receive message: ", zap.Any("msg", req))
	messaging.Receiver.Queue <- req
	return nil
}
