package websocket

import (
	"container/heap"
	"encoding/json"
	logicrpc "mim/internal/connect/logic_rpc"
	"mim/pkg/mq"
	"mim/pkg/proto"
	"strconv"
	"sync"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Bucket struct {
	ID       int
	cHeap    clientHeap
	capacity int
	clients  map[int64]*Client
	lock     sync.RWMutex
	Consumer *mq.Consumer
}

func (b *Bucket) RemoveClient(id int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	delete(b.clients, id)
}

func NewBucket(rabbitMQ *mq.RabbitMQ, serverID, bucketID int) *Bucket {
	exchangeName := strconv.Itoa(serverID)
	routingKey := strconv.Itoa(bucketID)
	b := &Bucket{
		ID:       bucketID,
		cHeap:    make(clientHeap, 0),
		clients:  make(map[int64]*Client, 50),
		capacity: 50,
		lock:     sync.RWMutex{},
		Consumer: mq.NewConsumer(rabbitMQ, exchangeName, exchangeName+routingKey, routingKey),
	}
	heap.Init(&b.cHeap)
	return b

}

func consumeMessage(messages <-chan amqp.Delivery) {
	for d := range messages {
		zap.L().Info("connect consumer receive", zap.Any("msg", d.Body))
		req := &proto.PushMessageReq{}
		json.Unmarshal(d.Body, req)
		c, ok := Default.GetUser(req.TargetID)
		if !ok {
			//
			req := &proto.OfflineMessageReq{
				SenderID: req.SenderID,
				TargetID: req.TargetID,
				Seq:      req.Seq,
			}
			logicrpc.StoreOffline(req)
		}

		c.Channel <- d.Body
		d.Ack(false)
	}
}

func (b *Bucket) LRU() {
	oldestClient := heap.Pop(&b.cHeap).(*Client)
	oldestClient.offline()
	delete(b.clients, oldestClient.ID)
}
