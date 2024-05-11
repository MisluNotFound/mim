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
