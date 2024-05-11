package rpc

import (
	"context"
	"mim/pkg/code"
	"mim/pkg/proto"
)

func (l *LogicRpc) NewGroup(ctx context.Context, req *proto.NewGroupReq, resp *proto.NewGroupResp) {
	resp.Code = code.CodeSuccess

}

func (l *LogicRpc) JoinGroup(ctx context.Context, req *proto.NewGroupReq, resp *proto.NewGroupResp) {

}

func (l *LogicRpc) FindGroup(ctx context.Context, req *proto.NewGroupReq, resp *proto.NewGroupResp) {

}

func (l *LogicRpc) LeaveGroup(ctx context.Context, req *proto.NewGroupReq, resp *proto.NewGroupResp) {

}
