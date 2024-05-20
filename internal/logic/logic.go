package logic

import (
	wsrpc "mim/internal/logic/pconn_rpc/conn_rpc"
	"mim/internal/logic/pconn_rpc/messaging"
	prpc "mim/internal/logic/pconn_rpc/rpc"
	"mim/internal/logic/redis"
	"mim/internal/logic/rpc"
)

func InitLogic() {
	go rpc.InitLogicRpc()
	go prpc.InitLogicRpc()
	go wsrpc.InitWsRpc()
	messaging.Receiver = messaging.NewReceiver(10)
	redis.Close()
	go messaging.Receiver.Start()
}
