package rpc

import (
	"context"
	"mim/internal/logic/dao"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func PullMessage(req *proto.PullMessageReq) (code.ResCode, map[int64][]dao.Message, error) {
	resp := &proto.PullMessageResp{}

	if err := logicRpc.Call(context.Background(), "PullMessage", req, resp); err != nil {
		zap.L().Error("PullMessage() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, nil, err
	}

	return resp.Code, resp.Messages, nil
}
