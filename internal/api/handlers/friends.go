package handlers

import (
	"mim/internal/api/rpc"
	"mim/pkg/code"
	"mim/pkg/proto"

	"github.com/gin-gonic/gin"
)

func AddFriend(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := ParamAddFriend{}
	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.AddFriendReq{
		UserID:   uid,
		FriendID: p.FriendID,
	}

	code, user, err := rpc.AddFriend(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, user)
}

func GetFriends(c *gin.Context) {
	uid := c.GetInt64("userId")

	req := &proto.GetFriendsReq{
		UserID: uid,
	}

	code, users, err := rpc.GetFriends(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, users)
}

func RemoveFriend(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamRemoveFriend{}
	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.RemoveFriendReq{
		UserID:   uid,
		FriendID: p.FriendID,
	}

	code, err := rpc.RemoveFriend(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, nil)
}

func UpdateFriendRemark(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamUpdateFriendRemark{}
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.UpdateFriendRemarkReq{
		UserID:   uid,
		FriendID: p.FriendID,
		Name:     p.Name,
	}

	code, err := rpc.UpdateFriendRemark(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, nil)
}

func FindFriend(c *gin.Context) {
	p := &ParamFindFriend{}

	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.FindFriendReq{
		UserID: p.UserID,
	}

	code, data, err := rpc.FindFriend(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}
