package rpc

import (
	"context"
	"mim/internal/logic/dao"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func JoinGroup(req *proto.JoinGroupReq) (code.ResCode, *dao.Group, error) {
	resp := &proto.JoinGroupResp{}

	err := logicRpc.Call(context.Background(), "JoinGroup", req, resp)
	if err != nil {
		zap.L().Error("JoinGroup() call logic failed: ", zap.Error(err))
		return code.CodeServerBusy, nil, err
	}

	return resp.Code, resp.Group, nil
}

func NewGroup(req *proto.NewGroupReq) (code.ResCode, *dao.Group, error) {
	resp := &proto.NewGroupResp{}

	err := logicRpc.Call(context.Background(), "NewGroup", req, resp)
	if err != nil {
		zap.L().Error("NewGroup() call logic failed: ", zap.Error(err))
		return resp.Code, nil, err
	}

	return resp.Code, resp.Group, nil
}

func FindGroup(req *proto.FindGroupReq) (code.ResCode, *dao.Group, error) {
	resp := &proto.FindGroupResp{}

	err := logicRpc.Call(context.Background(), "FindGroup", req, resp)
	if err != nil {
		zap.L().Error("FindGroup() call logic failed: ", zap.Error(err))
		return resp.Code, nil, err
	}

	return resp.Code, resp.Group, nil
}

func LeaveGroup(req *proto.LeaveGroupReq) (code.ResCode, error) {
	resp := &proto.LeaveGroupResp{}

	err := logicRpc.Call(context.Background(), "LeaveGroup", req, resp)
	if err != nil {
		zap.L().Error("LeaveGroup() call logic failed: ", zap.Error(err))
		return resp.Code, err
	}

	return resp.Code, nil
}

func GetGroups(req *proto.GetGroupsReq) (code.ResCode, []dao.Group, error) {
	resp := &proto.GetGroupsResp{}

	err := logicRpc.Call(context.Background(), "GetGroups", req, resp)
	if err != nil {
		zap.L().Error("LeaveGroup() call logic failed: ", zap.Error(err))
		return resp.Code, []dao.Group{}, err
	}

	return resp.Code, resp.Groups, nil
}
