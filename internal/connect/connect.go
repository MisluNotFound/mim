package connect

import (
	logicrpc "mim/internal/connect/logic_rpc"
	"mim/internal/connect/websocket"
	"mim/setting"
)

func InitConnect() {
	websocket.Default = websocket.NewServer(setting.Conf.WsConfig.WSServers[0].BucketSize,
		setting.Conf.WsConfig.WSServers[0].ID,
		setting.Conf.MQConfig.URL)
	go websocket.InitWebsocket()
	go logicrpc.InitLogicRpc()
}
