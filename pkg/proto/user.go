// 定义各层之间rpc通信的消息类型
package proto

import "mim/pkg/code"

type SignUpReq struct {
	Username   string
	Password   string
	RePassword string
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

type AuthReq struct {
	Token string
}

type AuthResp struct {
	Code     code.ResCode
	UserID   int64
	Username string
}