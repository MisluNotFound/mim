package handlers

import (
	"mim/internal/api/rpc"
	"mim/pkg/code"
	"mim/pkg/proto"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

func NewGroup(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := ParamNewGroup{}

	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := proto.NewGroupReq{
		OwnerID:     uid,
		GroupName:   p.GroupName,
		Description: p.Description,
	}

	code, data, err := rpc.NewGroup(&req)
	if err != nil {
		zap.L().Error("NewGroup() failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func JoinGroup(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := ParamJoinGroup{}

	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := proto.JoinGroupReq{
		UserID:  uid,
		GroupID: p.GroupID,
	}

	code, data, err := rpc.JoinGroup(&req)
	if err != nil {
		zap.L().Error("JoinGroup() failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func FindGroup(c *gin.Context) {
	p := ParamFindGroup{}

	if err := c.ShouldBindBodyWith(&p, binding.JSON); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := proto.FindGroupReq{
		GroupID: p.GroupID,
	}

	code, data, err := rpc.FindGroup(&req)
	if err != nil {
		zap.L().Error("FindGroup() failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func LeaveGroup(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := ParamLeaveGroup{}
	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := proto.LeaveGroupReq{
		UserID:  uid,
		GroupID: p.GroupID,
	}

	code, err := rpc.LeaveGroup(&req)
	if err != nil {
		zap.L().Error("LeaveGroup() failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, nil)
}

func GetGroups(c *gin.Context) {
	uid := c.GetInt64("userId")
	req := &proto.GetGroupsReq{
		UserID: uid,
	}

	code, data, err := rpc.GetGroups(req)
	if err != nil {
		zap.L().Error("LeaveGroup() failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}
