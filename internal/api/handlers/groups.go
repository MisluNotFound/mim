package handlers

import (
	"fmt"
	"mim/internal/api/rpc"
	"mim/pkg/code"
	"mim/pkg/proto"
	"strconv"

	"github.com/gin-gonic/gin"
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

	if err := c.ShouldBindQuery(&p); err != nil {
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

	gid, _ := strconv.ParseInt(p.GroupID, 10, 64)
	req := proto.LeaveGroupReq{
		UserID:  uid,
		GroupID: gid,
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
	fmt.Println(data)
	ResponseSuccess(c, data)
}

func GetMembers(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamGetMembers{}

	if err := c.ShouldBindQuery(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	gid, _ := strconv.ParseInt(p.GroupID, 10, 64)
	req := &proto.GetMembersReq{
		UserID:  uid,
		GroupID: gid,
	}
	code, data, err := rpc.GetMembers(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func GetRole(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamGetRole{}
	if err := c.ShouldBindQuery(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	gid, _ := strconv.ParseInt(p.GroupID, 10, 64)
	req := &proto.GetRoleReq{
		UserID:  uid,
		GroupID: gid,
	}

	code, data, err := rpc.GetRole(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func UpdateGroupPhoto(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamUpdateGroupPhoto{}
	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	gid, _ := strconv.ParseInt(p.GroupID, 10, 64)
	req := &proto.UpdateGroupPhotoReq{
		GroupID: gid,
		UserID:  uid,
		Avatar:  p.Avatar,
	}

	code, err := rpc.UpdateGroupPhoto(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, nil)
}

func RemoveMember(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamRemoveMember{}
	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	gid, _ := strconv.ParseInt(p.GroupID, 10, 64)
	memberID, _ := strconv.ParseInt(p.MemberID, 10, 64)
	req := &proto.RemoveMemberReq{
		GroupID: gid,
		UserID:  uid,
		MemberID: memberID,
	}

	code, err := rpc.RemoveMember(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, nil)
}