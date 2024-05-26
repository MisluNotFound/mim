package logicrpc

import (
	"context"
	"mim/pkg/proto"
)

func StoreOffline(req *proto.MessageReq) error {
	resp := &proto.MessageResp{}
	err := logicPRpc.Call(context.Background(), "Auth", req, resp)
	if err != nil {
		return err
	}

	return nil
}
