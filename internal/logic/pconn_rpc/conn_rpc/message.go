package wsrpc

import (
	"context"
	"mim/pkg/proto"
)

func PushMessage(req *proto.PushMessageReq) {
	resp := proto.MessageResp{}
	connectRpc.Call(context.Background(), "PushMessage", req, &resp)
}
