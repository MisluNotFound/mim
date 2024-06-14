package rpc

import (
	"context"
	"mim/internal/logic/dao"
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

func NearBy(req *proto.NearbyReq) (code.ResCode, []dao.User, error) {
	resp := &proto.NearByResp{}

	err := logicRpc.Call(context.Background(), "NearBy", req, resp)
	if err != nil {
		zap.L().Error("NearBy() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, nil, err
	}
	return resp.Code, resp.Users, err
}

func GetInfo(req *proto.GetInfoReq) (code.ResCode, *dao.User, error) {
	resp := &proto.GetInfoResp{}

	err := logicRpc.Call(context.Background(), "GetInfo", req, resp)
	if err != nil {
		zap.L().Error("GetInfo() call logic failed: ", zap.Error(err))
		return resp.Code, nil, err
	}
	return resp.Code, resp.User, err

}

func UpdatePhoto(req *proto.UpdatePhotoReq) (code.ResCode, error) {
	resp := &proto.UpdatePhotoResp{}

	err := logicRpc.Call(context.Background(), "UpdatePhoto", req, resp)
	if err != nil {
		zap.L().Error("UpdatePhoto() call logic failed: ", zap.Error(err))
		return resp.Code, err
	}

	return resp.Code, nil
}

func UpdatePassword(req *proto.UpdatePasswordReq) (code.ResCode, error) {
	resp := &proto.UpdatePasswordResp{}

	err := logicRpc.Call(context.Background(), "UpdatePassword", req, resp)
	if err != nil {
		zap.L().Error("UpdatePassword() call logic failed: ", zap.Error(err))
		return resp.Code, err
	}

	return resp.Code, nil
}

func UpdateName(req *proto.UpdateNameReq) (code.ResCode, error) {
	resp := &proto.UpdateNameResp{}

	err := logicRpc.Call(context.Background(), "UpdateName", req, resp)
	if err != nil {
		zap.L().Error("UpdateName() call logic failed: ", zap.Error(err))
		return resp.Code, err
	}

	return resp.Code, nil
}
