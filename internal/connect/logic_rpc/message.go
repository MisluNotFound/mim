package logicrpc

import (
	"context"
	"mim/pkg/proto"
)

func StoreOffline(req *proto.OfflineMessageReq) error {
	resp := &proto.MessageResp{}
	err := logicSRpc.Call(context.Background(), "StoreOffline", req, resp)
	if err != nil {
		return err
	}

	return nil
}
