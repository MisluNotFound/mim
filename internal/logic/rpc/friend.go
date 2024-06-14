package rpc

import (
	"context"
	"mim/internal/logic/dao"
	"mim/pkg/code"
	"mim/pkg/proto"
)

func (r *LogicRpc) AddFriend(ctx context.Context, req *proto.AddFriendReq, resp *proto.AddFriendResp) error {
	resp.Code = code.CodeSuccess
	u, ok, err := dao.FindUserByID(req.FriendID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}
	if !ok {
		resp.Code = code.CodeUserNotExist
		return dao.ErrorUserNotExist
	}

	f, err := dao.GetFriend(req.UserID, req.FriendID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if f.UserA != 0 || f.UserB != 0 {
		resp.Code = code.CodeAlreadyAdd
		return dao.ErrorFriendAlreadyAdd
	}

	if err := dao.AddFriend(req.UserID, req.FriendID); err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	u.Password = ""
	resp.Friend = *u

	return nil
}

func (r *LogicRpc) RemoveFriend(ctx context.Context, req *proto.RemoveFriendReq, resp *proto.RemoveFriendResp) error {
	resp.Code = code.CodeSuccess

	f, err := dao.GetFriend(req.UserID, req.FriendID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}
	if f.UserA == 0 || f.UserB == 0 {
		resp.Code = code.CodeFriendNotExist
		return dao.ErrorFriendNotExist
	}

	if err = dao.DeleteFriend(req.UserID, req.FriendID); err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (r *LogicRpc) GetFriends(ctx context.Context, req *proto.GetFriendsReq, resp *proto.GetFriendsReps) error {
	resp.Code = code.CodeSuccess

	friends, err := dao.GetFriends(req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	var ids []int64
	var remarks map[int64]string = make(map[int64]string)
	for _, f := range friends {
		var id int64
		if f.UserA != req.UserID {
			id = f.UserA
			remarks[id] = f.BtoA
		} else {
			id = f.UserB
			remarks[id] = f.AtoB
		}
		ids = append(ids, id)
	}

	users, err := dao.GetFriendsInfo(ids)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	var friendInfos []dao.FriendInfo
	for _, u := range users {
		friendInfos = append(friendInfos, dao.FriendInfo{
			Info:   u,
			Remark: remarks[u.ID],
		})
	}
	resp.Friends = friendInfos
	return nil
}

func (r *LogicRpc) UpdateFriendRemark(ctx context.Context, req *proto.UpdateFriendRemarkReq, resp *proto.UpdateFriendRemarkResp) error {
	resp.Code = code.CodeSuccess

	friend, err := dao.GetFriend(req.UserID, req.FriendID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if friend.UserA == 0 || friend.UserB == 0 {
		resp.Code = code.CodeFriendNotExist
		return dao.ErrorFriendNotExist
	}

	if friend.UserA == req.UserID {
		friend.AtoB = req.Name
	} else {
		friend.BtoA = req.Name
	}

	err = dao.UpdateFriendRemark(friend)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (r *LogicRpc) FindFriend(ctx context.Context, req *proto.FindFriendReq, resp *proto.FindFriendResp) error {
	resp.Code = code.CodeSuccess

	u, ok, err := dao.FindUserByID(req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeUserNotExist
		return dao.ErrorUserNotExist
	}

	resp.User = *u
	return nil
}
