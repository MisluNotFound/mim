package rpc

import (
	"context"
	"mim/internal/logic/dao"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func AddFriend(req *proto.AddFriendReq) (code.ResCode, dao.User, error) {
	resp := &proto.AddFriendResp{}

	if err := logicRpc.Call(context.Background(), "AddFriend", req, resp); err != nil {
		zap.L().Error("AddFriend() call logic failed: ", zap.Error(err))
		return resp.Code, dao.User{}, err
	}

	return resp.Code, resp.Friend, nil
}

func GetFriends(req *proto.GetFriendsReq) (code.ResCode, interface{}, error) {
	resp := &proto.GetFriendsReps{}

	if err := logicRpc.Call(context.Background(), "GetFriends", req, resp); err != nil {
		zap.L().Error("GetFriends() call logic failed: ", zap.Error(err))
		return resp.Code, []dao.User{}, err
	}

	return resp.Code, resp.Friends, nil
}

func RemoveFriend(req *proto.RemoveFriendReq) (code.ResCode, error) {
	resp := &proto.RemoveFriendResp{}

	if err := logicRpc.Call(context.Background(), "RemoveFriend", req, resp); err != nil {
		zap.L().Error("RemoveFriend() call logic failed: ", zap.Error(err))
		return resp.Code, err
	}

	return resp.Code, nil
}

func UpdateFriendRemark(req *proto.UpdateFriendRemarkReq) (code.ResCode, error) {
	resp := &proto.UpdateFriendRemarkResp{}

	if err := logicRpc.Call(context.Background(), "UpdateFriendRemark", req, resp); err != nil {
		zap.L().Error("UpdateFriendRemark() call logic failed: ", zap.Error(err))
		return resp.Code, err
	}

	return resp.Code, nil
}

func FindFriend(req *proto.FindFriendReq) (code.ResCode, interface{}, error) {
	resp := &proto.FindFriendResp{}

	if err := logicRpc.Call(context.Background(), "FindFriend", req, resp); err != nil {
		zap.L().Error("FindFriend() call logic failed: ", zap.Error(err))
		return resp.Code, nil, err
	}

	return resp.Code, resp.User, nil
}
