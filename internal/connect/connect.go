package connect

import (
	"mim/internal/connect/rpc"
	logicrpc "mim/internal/connect/rpc/logic_rpc"
	"mim/internal/connect/websocket"
)

func InitConnect() {
	buckets := make([]*websocket.Bucket, 50)
	for i := range buckets {
		buckets[i] = websocket.NewBucket()
	}

	websocket.Default = websocket.NewServer(buckets, 1)
	go rpc.InitConnectRpc()
	go websocket.InitWebsocket()
	go logicrpc.InitLogicRpc()
}
