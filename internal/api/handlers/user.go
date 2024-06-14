// 调用rpc中的方法由其转发给logic
package handlers

import (
	"mim/internal/api/rpc"
	"mim/pkg/code"
	"mim/pkg/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignUp(c *gin.Context) {
	p := ParamSignUp{}
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("signUp() failed: invalid params")
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := proto.SignUpReq{
		Username: p.Username,
		Password: p.Password,
		Avatar:   p.Avatar,
	}

	code, token, err := rpc.SignUp(&req)
	if err != nil {
		zap.L().Error("sign up failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, token)
}

func SignIn(c *gin.Context) {
	p := ParamSignIn{}
	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := proto.SignInReq{
		Username: p.Username,
		Password: p.Password,
	}

	code, token, err := rpc.SignIn(&req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, token)
}

func GetInfo(c *gin.Context) {
	uid := c.GetInt64("userId")

	req := &proto.GetInfoReq{
		UserID: uid,
	}

	code, user, err := rpc.GetInfo(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, user)
}

func NearbyOpen(c *gin.Context) {
	p := ParamNearby{}
	if err := c.ShouldBindJSON(&p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.NearbyReq{
		UserID:    c.GetInt64("userId"),
		Longitude: p.Longitude,
		Latitude:  p.Latitude,
	}

	code, data, err := rpc.NearBy(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func UpdatePhoto(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamUpdatePhoto{}
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.UpdatePhotoReq{
		UserID: uid,
		Avatar: p.Avatar,
	}

	code, err := rpc.UpdatePhoto(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, nil)
}

func UpdatePassword(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamUpdatePassword{}
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.UpdatePasswordReq{
		UserID:      uid,
		OldPassword: p.OldPassword,
		NewPassword: p.NewPassword,
	}

	code, err := rpc.UpdatePassword(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, nil)
}

func UpdateName(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := &ParamUpdateName{}
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	fmt.Println(p)
	req := &proto.UpdateNameReq{
		UserID: uid,
		Name:   p.Name,
	}

	code, err := rpc.UpdateName(req)
	if err != nil {
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, nil)
}
