package connect

import (
	"mim/internal/connect/rpc"
	"mim/internal/connect/websocket"
)

func InitConnect() {
	buckets := make([]*websocket.Bucket, 50)
	for i, _ := range buckets {
		buckets[i] = websocket.NewBucket()
	}

	websocket.Default = websocket.NewServer(buckets)
	go rpc.InitConnectRpc()
	go websocket.InitWebsocket()
}
