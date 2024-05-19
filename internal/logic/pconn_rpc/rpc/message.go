package prpc

import (
	"context"
	"mim/internal/logic/dao"
	"mim/internal/logic/pconn_rpc/messaging"
	"mim/internal/logic/redis"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

// 接收connect层发送的消息 并交给message层处理
func (r *PRpc) SendMessage(ctx context.Context, req *proto.MessageReq, resp *proto.MessageResp) error {
	zap.L().Info("logic server receive message: ", zap.Any("msg", req))
	messaging.Receiver.Queue <- req
	return nil
}

func (r *PRpc) PullMessage(ctx context.Context, req *proto.PullMessageReq, resp *proto.PullMessageResp) error {
	resp.Code = code.CodeSuccess

	// redis获取消息seq数组
	seqs, err := redis.GetMessages(req.UserID, req.TargetID, req.Size, req.LastSeq)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	// mysql获取消息内容
	messages, err := dao.GetMessages(seqs)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	// 包装
	data := map[int64][]dao.Message{} 
	for _, msg := range messages {
        targetID := msg.TargetID
        data[targetID] = append(data[targetID], msg)
    }

	resp.Messages = data
	return nil
}
