package rpc

import (
	"context"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func PullMessage(req *proto.PullMessageReq) (code.ResCode, interface{}, error) {
	resp := &proto.PullMessageResp{}

	if err := logicRpc.Call(context.Background(), "PullMessage", req, resp); err != nil {
		zap.L().Error("PullMessage() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, nil, err
	}

	return resp.Code, resp.Messages, nil
}

func PullOfflineMessage(req *proto.PullOfflineMessageReq) (code.ResCode, interface{}, error) {
	resp := &proto.PullMessageResp{}

	if err := logicRpc.Call(context.Background(), "PullOfflineMessage", req, resp); err != nil {
		zap.L().Error("PullOfflineMessage() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, nil, err
	}

	return resp.Code, resp.Messages, nil
}

func GetUnReadCount(req *proto.GetUnReadCountReq) (code.ResCode, interface{}, error) {
	resp := &proto.GetUnReadResp{}

	if err := logicRpc.Call(context.Background(), "GetUnReadCount", req, resp); err != nil {
		zap.L().Error("PullOfflineMessage() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, nil, err
	}

	return resp.Code, resp.SessionInfo, nil
}

func PullErrMessage(req *proto.PullErrMessageReq) (code.ResCode, interface{}, error) {
	resp := &proto.PullErrMessageResp{}

	if err := logicRpc.Call(context.Background(), "PullErrMessage", req, resp); err != nil {
		zap.L().Error("PullErrMessage() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, nil, err
	}

	return resp.Code, resp.Messages, nil
}
