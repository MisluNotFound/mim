package rpc

import (
	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
)

var logicRpc client.XClient

func InitAPIRpc() {
	d, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost:8081", "")
	if err != nil {
		zap.L().Error("init api rpc failed: ", zap.Error(err))
	}
	zap.L().Info("init api rpc success")
	logicRpc = client.NewXClient("LogicRpc", client.Failtry, client.RandomSelect, d, client.DefaultOption)
}
