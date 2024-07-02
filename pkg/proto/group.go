package proto

import (
	"mim/internal/logic/dao"
	"mim/pkg/code"
)

type NewGroupReq struct {
	OwnerID     int64
	GroupName   string
	Description string
}

type NewGroupResp struct {
	Code  code.ResCode
	Group *dao.Group
}

type JoinGroupReq struct {
	UserID  int64
	GroupID int64
}

type JoinGroupResp struct {
	Code  code.ResCode
	Group *dao.Group
}

type FindGroupReq struct {
	GroupID int64
}

type FindGroupResp struct {
	Code  code.ResCode
	Group *dao.Group
}

type LeaveGroupReq struct {
	UserID  int64
	GroupID int64
}

type LeaveGroupResp struct {
	Code code.ResCode
}

type FindGroupsReq struct {
	UserID int64
}

type FindGroupsResp struct {
	Code   code.ResCode
	Groups *[]int64
}

type GetGroupsReq struct {
	UserID int64
}

type GetGroupsResp struct {
	Code   code.ResCode
	Groups []dao.Group
}

type GetMembersReq struct {
	UserID  int64
	GroupID int64
}

type UserInfo struct {
	Member dao.User
	Role   string
}

type GetMembersResp struct {
	Code    code.ResCode
	Members []UserInfo
}

type UpdateGroupNameReq struct {
	UserID  int64
	GroupID int64
	Name    string
}

type UpdateGroupNameResp struct {
	Code code.ResCode
}

type UpdateGroupPhotoReq struct {
	UserID  int64
	GroupID int64
	Avatar  string
}

type UpdateGroupPhotoResp struct {
	Code code.ResCode
}

type UpdateMyNameReq struct {
	UserID  int64
	GroupID int64
	Name    string
}

type UpdateMyNameResp struct {
	Code code.ResCode
}

type GetRoleReq struct {
	UserID  int64
	GroupID int64
}

type GetRoleResp struct {
	Code code.ResCode
	Role string
}

type RemoveMemberReq struct {
	UserID   int64
	MemberID int64
	GroupID int64
}

type RemoveMemberResp struct {
	Code code.ResCode
}