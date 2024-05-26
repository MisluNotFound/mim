package logicrpc

import (
	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
)

var logicSRpc client.XClient
var logicPRpc client.XClient

func InitLogicRpc() {
	sd, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost:8081", "")
	if err != nil {
		zap.L().Error("init connect rpc failed: ", zap.Error(err))
		return
	}
	pd, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost:8083", "")

	if err != nil {
		zap.L().Error("init connect rpc failed: ", zap.Error(err))
		return
	}

	logicPRpc = client.NewXClient("PRpc", client.Failtry, client.RandomSelect, pd, client.DefaultOption)
	logicSRpc = client.NewXClient("LogicRpc", client.Failtry, client.RandomSelect, sd, client.DefaultOption)
	zap.L().Info("init connect rpc success")
}
