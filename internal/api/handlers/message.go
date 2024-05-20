package handlers

import (
	"mim/internal/api/rpc"
	"mim/pkg/code"
	"mim/pkg/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PullMessage(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := ParamPullMessage{}

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error(err.Error())
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.PullMessageReq{
		UserID: uid,
		TargetID: p.TargetID,
		LastSeq: p.LastSeq,
		Size: p.Size,
	}
	
	code, data, err := rpc.PullMessage(req)
	if err != nil {
		zap.L().Error("PullMessage() Failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}
	
	ResponseSuccess(c, data)
}
