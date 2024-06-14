package rpc

import (
	"context"
	"mim/internal/logic/dao"
	"mim/pkg/code"
	"mim/pkg/proto"
)

func (l *LogicRpc) UpdateMyName(ctx context.Context, req *proto.UpdateMyNameReq, resp *proto.UpdateMyNameResp) error {
	resp.Code = code.CodeSuccess

	err := dao.UpdateMyName(req.UserID, req.GroupID, req.Name)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}
