package websocket

import (
	"encoding/json"
	"mim/pkg/mq"
	"mim/pkg/proto"
	"strconv"
	"sync"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Bucket struct {
	ID       int
	clients  map[int64]*Client
	lock     sync.RWMutex
	Consumer *mq.Consumer
}

func (b *Bucket) RemoveClient(id int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	delete(b.clients, id)
}

func NewBucket(rabbitMQ *mq.RabbitMQ, serverID, bucketID int) (b *Bucket) {
	b = new(Bucket)
	b.ID = bucketID
	exchange := strconv.Itoa(serverID)
	routingKey := strconv.Itoa(b.ID)
	queueName := exchange + routingKey
	b.Consumer = mq.NewConsumer(rabbitMQ, exchange, queueName, routingKey)
	b.clients = make(map[int64]*Client, 50)

	return b
}

func (b *Bucket) LRU() {

}

func consumeMessage(messages <-chan amqp.Delivery) {
	for d := range messages {
		zap.L().Info("connect consumer receive", zap.Any("msg", d.Body))
		req := &proto.PushMessageReq{}
		json.Unmarshal(d.Body, req)
		c, ok := Default.GetUser(req.TargetID)
		if !ok {
			
		}

		c.Channel <- d.Body
		d.Ack(false)
	}
}
