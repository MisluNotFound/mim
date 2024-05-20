// 为ws层提供接口
package prpc

import (
	"context"
	"mim/internal/logic/redis"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

func (r *PRpc) Online(ctx context.Context, req *proto.OnlineReq, resp *proto.OnlineResp) error {
	resp.Code = code.CodeSuccess
	if err := redis.AddOnlineUser(req.UserID, req.ServerID); err != nil {
		zap.L().Info("logic Online() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (r *PRpc) Offline(ctx context.Context, req *proto.OfflineReq, resp *proto.OfflineResp) error {
	resp.Code = code.CodeSuccess

	if err := redis.RemoveOnlineUser(req.UserID); err != nil {
		zap.L().Error("logic Offline() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}
