package rpc

import (
	"context"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func SignUp(req *proto.SignUpReq) (code.ResCode, string, error) {
	resp := &proto.SignUpResp{}

	err := logicRpc.Call(context.Background(), "SignUp", req, resp)
	if err != nil {
		zap.L().Error("SignUp() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, "", err
	}
	return resp.Code, resp.Token, nil
}

func SignIn(req *proto.SignInReq) (code.ResCode, string, error) {
	resp := &proto.SignInResp{}

	err := logicRpc.Call(context.Background(), "SignIn", req, resp)
	if err != nil {
		zap.L().Error("SignIn() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, "", err
	}
	return resp.Code, resp.Token, err
}
