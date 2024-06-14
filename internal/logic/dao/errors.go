package dao

import "errors"

var (
	ErrorUserExist          = errors.New("用户已存在")
	ErrorUserNotExist       = errors.New("用户不存在")
	ErrorInvalidPassword    = errors.New("用户名或密码错误")
	ErrorInvalidID          = errors.New("无效的ID")
	ErrorGroupNotExist      = errors.New("群不存在")
	ErrorGroupAlreadyJoined = errors.New("已经加入该群")
	ErrorNotJoinGroup       = errors.New("未加入该群")
	ErrorFriendAlreadyAdd   = errors.New("已添加好友")
	ErrorFriendNotExist     = errors.New("未添加好友")
	ErrorPermissionDenied   = errors.New("没有权限")
)
