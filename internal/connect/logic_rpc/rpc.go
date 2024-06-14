package logicrpc

import (
	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
)

var logicSRpc client.XClient

func InitLogicRpc() {
	sd, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost:8081", "")
	if err != nil {
		zap.L().Error("init connect rpc failed: ", zap.Error(err))
		return
	}

	logicSRpc = client.NewXClient("LogicRpc", client.Failtry, client.RandomSelect, sd, client.DefaultOption)
	zap.L().Info("init connect rpc success")
}
