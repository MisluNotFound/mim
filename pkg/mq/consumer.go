package mq

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Consumer struct {
	rabbitMQ   *RabbitMQ
	exchange   string
	queueName  string
	routingKey string
	Tag        string
}

func NewConsumer(rabbitMQ *RabbitMQ, exchange, queueName, routingKey string) *Consumer {
	c := &Consumer{
		rabbitMQ:   rabbitMQ,
		exchange:   exchange,
		queueName:  queueName,
		routingKey: routingKey,
	}

	return c
}

func (c *Consumer) Work(isWork bool, handler func(<-chan amqp.Delivery)) {
	c.rabbitMQ.DeclareExchange(c.exchange, "direct", true, false, false, false, nil)
	queue, _ := c.rabbitMQ.DeclareQueue(c.queueName, true, false, false, false, nil)
	c.rabbitMQ.BindQueue(queue.Name, c.routingKey, c.exchange, false, nil)

	fmt.Printf("consumer is running queue: %s\n", c.queueName)
	if isWork {
		// 设置 QoS，预取计数为1，确保每个消费者一次只处理一个消息
		c.rabbitMQ.lock.Lock()
		err := c.rabbitMQ.ConsumerChannel.Qos(1, 0, false)
		if err != nil {
			zap.L().Error("set qos failed: ", zap.Error(err))
		}
		c.rabbitMQ.lock.Unlock()
	}

	msgs, err := c.rabbitMQ.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		zap.L().Error("consume message failed: ", zap.Error(err))
		return
	}

	// 调用 handler 处理消息
	go handler(msgs)
}
