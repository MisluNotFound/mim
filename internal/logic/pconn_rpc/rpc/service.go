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
	// 在线用户列表删除
	
	return nil
}
