package rpc

import (
	"context"
	"fmt"
	"mim/pkg/code"
	"mim/pkg/proto"

	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
)

var logicRpc client.XClient

func InitAPIRpc() {
	d, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost:8081", "")
	if err != nil {
		zap.L().Error("init api rpc failed: ", zap.Error(err))
	}
	logicRpc = client.NewXClient("LogicRpc", client.Failtry, client.RandomSelect, d, client.DefaultOption)
}

func SignUp(req *proto.SignUpReq) (code.ResCode, string, error) {
	resp := &proto.SignUpResp{}

	err := logicRpc.Call(context.Background(), "SignUp", req, resp)
	if err != nil {
		zap.L().Error("call logic")
		return code.CodeServerBusy, "", err
	}
	zap.L().Info(resp.Token)
	return resp.Code, resp.Token, nil
}

func SignIn(req *proto.SignInReq) (code.ResCode, string, error) {
	resp := &proto.SignInResp{}

	err := logicRpc.Call(context.Background(), "SignIn", req, resp)
	if err != nil {
		return code.CodeServerBusy, "", err
	}
	fmt.Println("api: ", resp.Token)
	return resp.Code, resp.Token, err
}
