package proto

import (
	"mim/internal/logic/dao"
	"mim/pkg/code"
)

type AddFriendReq struct {
	UserID   int64
	FriendID int64
}

type AddFriendResp struct {
	Code   code.ResCode
	Friend dao.User
}

type RemoveFriendReq struct {
	UserID   int64
	FriendID int64
}

type RemoveFriendResp struct {
	Code code.ResCode
}

type GetFriendsReq struct {
	UserID int64
}

type GetFriendsReps struct {
	Code    code.ResCode
	Friends []dao.FriendInfo
}

type UpdateFriendRemarkReq struct {
	UserID   int64
	FriendID int64
	Name     string
}

type UpdateFriendRemarkResp struct {
	Code code.ResCode
}

type FindFriendReq struct {
	UserID int64
}

type FindFriendResp struct {
	Code code.ResCode
	User dao.User
}
