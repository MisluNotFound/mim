package mq

import (
	"mim/pkg/mq"

	"go.uber.org/zap"
)

var rabbitMQ *mq.RabbitMQ

func InitMQ(rabbitMQURL string, exchangeName, queueName, routingKey string, consumerNum int, publisherNum int32) {
	var err error
	rabbitMQ, err = mq.NewRabbitMQ(rabbitMQURL)
	if err != nil {
		zap.L().Error("consumer connect mq failed: ", zap.Error(err))
		return
	}

	// consumerNum个消费者去connect层的exchangeName的queueName消费消息
	StartConsumers(exchangeName, queueName, routingKey, consumerNum)
	StartPublishers(publisherNum)
}

func StartConsumers(exchangeName, queueName, routingKey string, consumerNum int) {
	for i := 0; i < consumerNum; i++ {
		consumer := mq.NewConsumer(rabbitMQ, exchangeName, queueName, routingKey)
		go consumer.Work(true, consumeMessage)
	}
}

func StartPublishers(publisherNum int32) {
	publishers = mq.NewPublishers(rabbitMQ, publisherNum)
}
