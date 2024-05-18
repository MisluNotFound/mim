// ConnectRpc服务调用的client
package wsrpc

import (
	"github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
)

var connectRpc client.XClient

func InitWsRpc() {
	d, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost:8084", "")
	if err != nil {
		zap.L().Error("init WsRpc client failed: ", zap.Error(err))
		return
	}
	zap.L().Info("init WsRpc client success")
	connectRpc = client.NewXClient("ConnectRpc", client.Failtry, client.RandomSelect, d, client.DefaultOption)
}
