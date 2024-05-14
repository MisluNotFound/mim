package rpc

import (
	"context"
	"errors"
	"mim/pkg/code"
	"mim/pkg/proto"

	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
)

var connectRpc client.XClient

func InitConnectRpc() {
	d, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost:8081", "")
	if err != nil {
		zap.L().Error("init connect rpc failed: ", zap.Error(err))
		return
	}
	connectRpc = client.NewXClient("LogicRpc", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	zap.L().Info("init connect rpc success")
}

func Auth(req *proto.AuthReq) (int64, string, error) {
	resp := &proto.AuthResp{}
	err := connectRpc.Call(context.Background(), "Auth", req, resp)
	if err != nil {
		return -1, "", err
	}

	if resp.Code != code.CodeSuccess {
		return -1, "", errors.New("invalid token")
	}

	return resp.UserID, resp.Username, nil
}

func GetGroup(req *proto.FindGroupsReq) (*[]int64, error) {
	resp := &proto.FindGroupsResp{}

	err := connectRpc.Call(context.Background(), "FindGroups", req, resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != code.CodeSuccess {
		return nil, errors.New("server busy")
	}

	return resp.Groups, nil
}