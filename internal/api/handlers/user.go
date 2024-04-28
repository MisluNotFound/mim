// 调用rpc中的方法由其转发给logic
package handlers

import (
	"mim/internal/api/rpc"
	"mim/pkg/code"
	"mim/pkg/proto"

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
	}

	code, token, err := rpc.SignUp(&req)
	if err != nil {
		zap.L().Error("")
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
