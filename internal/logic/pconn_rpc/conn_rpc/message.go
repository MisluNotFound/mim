package wsrpc

import (
	"context"
	"mim/internal/logic/dao"
	"mim/internal/logic/redis"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func PushMessage(req *proto.PushMessageReq) {
	resp := proto.PushMessageResp{}
	zap.L().Info("ws client receive message", zap.Any("msg", req))
	if err := connectRpc.Call(context.Background(), "PushMessage", req, &resp); err != nil {
		zap.L().Error("PushMessage() failed: ", zap.Error(err))
	}

	if resp.IsOffline {
		redis.StoreOfflineMessage(dao.Message{
			SenderID: req.SenderID,
			TargetID: req.TargetID,
			Seq:      req.Seq,
			Content:  req.Body,
		})
	}
}
