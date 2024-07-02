package mq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Publisher struct {
	rabbitMQ *RabbitMQ
}

func NewPublishers(rabbitMQ *RabbitMQ, publisherNum int32) []*Publisher {
	var publishers []*Publisher
	for i := 0; i < int(publisherNum); i++ {
		publisher := NewPublisher(rabbitMQ)
		publishers = append(publishers, publisher)
	}

	return publishers
}

// NewPublisher 创建并返回一个新的 Publisher 实例
func NewPublisher(rabbitMQ *RabbitMQ) *Publisher {
	p := &Publisher{
		rabbitMQ: rabbitMQ,
	}

	return p
}

// PublishMessage 发布消息到指定的交换机
func (p *Publisher) PublishMessage(body []byte, exchange, routingKey, queueName string) error {
	fmt.Printf("publisher send message to %s %s \n", exchange, queueName)
	err := p.rabbitMQ.PublisherChannel.Publish(
		exchange,   
		routingKey, 
		false,      
		false,      
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		return err
	}

	return nil
}
