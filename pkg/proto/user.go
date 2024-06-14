// 定义各层之间rpc通信的消息类型
package proto

import (
	"mim/internal/logic/dao"
	"mim/pkg/code"
)

type SignUpReq struct {
	Username   string
	Password   string
	RePassword string
	Avatar     string
}

type SignUpResp struct {
	Code  code.ResCode
	Token string
}

type SignInReq struct {
	Username string
	Password string
}

type SignInResp struct {
	Code  code.ResCode
	Token string
}

type UpdatePhotoReq struct {
	UserID int64
	Avatar string
}

type UpdatePasswordReq struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

type UpdatePasswordResp struct {
	Code code.ResCode
}

type UpdatePhotoResp struct {
	Code code.ResCode
}

type UpdateNameReq struct {
	UserID int64
	Name   string
}

type UpdateNameResp struct {
	Code code.ResCode
}

type AuthReq struct {
	Token string
}

type AuthResp struct {
	Code     code.ResCode
	UserID   int64
	Username string
}

type NearbyReq struct {
	UserID    int64
	Longitude float64
	Latitude  float64
}

type NearByResp struct {
	Code  code.ResCode
	Users []dao.User
}

type GetInfoReq struct {
	UserID int64
}

type GetInfoResp struct {
	Code code.ResCode
	User *dao.User
}

type RecentSessionReq struct {
	UserID int64
}

type RecentSessionResp struct {
}
