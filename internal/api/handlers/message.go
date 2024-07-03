package handlers

import (
	"math"
	"mim/internal/api/rpc"
	"mim/pkg/code"
	"mim/pkg/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PullMessage(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := ParamPullMessage{}

	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error(err.Error())
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	if p.LastSeq == 0 {
		p.LastSeq = math.MaxInt64
	}

	req := &proto.PullMessageReq{
		UserID:    uid,
		SessionID: p.SessionID,
		LastSeq:   p.LastSeq,
		Size:      p.Size,
		IsGroup:   p.IsGroup,
	}

	code, data, err := rpc.PullMessage(req)
	if err != nil {
		zap.L().Error("PullMessage() Failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func PullOfflineMessage(c *gin.Context) {
	uid := c.GetInt64("userId")

	p := ParamPullOfflineMessage{}

	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("PullOfflineMessage() Failed: ", zap.Error(err))
		ResponseError(c, code.CodeInvalidParam)
		return
	}

	req := &proto.PullOfflineMessageReq{
		UserID:    uid,
		IsGroup:   p.IsGroup,
		SessionID: p.SessionID,
	}

	code, data, err := rpc.PullOfflineMessage(req)
	if err != nil {
		zap.L().Error("PullOfflineMessage() Failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func GetUnReadCount(c *gin.Context) {
	uid := c.GetInt64("userId")

	req := &proto.GetUnReadCountReq{
		UserID: uid,
	}

	code, data, err := rpc.GetUnReadCount(req)
	if err != nil {
		zap.L().Error("GetUnReadMessage() failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}

func PullErrMessage(c *gin.Context) {
	uid := c.GetInt64("userId")

	req := &proto.PullErrMessageReq{
		UserID: uid,
	}

	code, data, err := rpc.PullErrMessage(req)
	if err != nil {
		zap.L().Error("GetUnReadMessage() failed: ", zap.Error(err))
		ResponseError(c, code)
		return
	}

	ResponseSuccess(c, data)
}
