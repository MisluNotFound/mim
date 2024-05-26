package mq

import (
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn             *amqp.Connection
	ConsumerChannel  *amqp.Channel
	PublisherChannel *amqp.Channel
	lock             sync.RWMutex
}

// 连接到 RabbitMQ
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	cch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	pch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Conn:             conn,
		ConsumerChannel:  cch,
		PublisherChannel: pch,
	}, nil
}

// 关闭连接和通道
func (r *RabbitMQ) Close() {
	if r.ConsumerChannel != nil {
		r.ConsumerChannel.Close()
	}

	if r.PublisherChannel != nil {
		r.PublisherChannel.Close()
	}

	if r.Conn != nil {
		r.Conn.Close()
	}
}

// 声明交换机
func (r *RabbitMQ) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	err := r.PublisherChannel.ExchangeDeclare(
		name,       // 交换机名称
		kind,       // 交换机类型
		durable,    // 持久化
		autoDelete, // 自动删除
		internal,   // 内部使用
		noWait,     // 阻塞
		args,       // 额外参数
	)

	return err
}

// 声明队列
func (r *RabbitMQ) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := r.PublisherChannel.QueueDeclare(
		name,       // 队列名称
		durable,    // 持久化
		autoDelete, // 自动删除
		exclusive,  // 独占
		noWait,     // 阻塞
		args,       // 额外参数
	)

	return q, err
}

// 绑定队列到交换机
func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string, noWait bool, args amqp.Table) error {
	err := r.PublisherChannel.QueueBind(
		queueName,    // 队列名称
		routingKey,   // 路由键
		exchangeName, // 交换机名称
		noWait,       // 阻塞
		args,         // 额外参数
	)

	return err
}

// 发布消息
func (r *RabbitMQ) Publish(exchange, routingKey string, mandatory, immediate bool, msg amqp.Publishing) error {
	err := r.PublisherChannel.Publish(
		exchange,   // 交换机名称
		routingKey, // 路由键
		mandatory,  // 必须
		immediate,  // 立即
		msg,
	)

	return err
}

// 消费消息
func (r *RabbitMQ) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	msgs, err := r.ConsumerChannel.Consume(
		queue,     // 队列名称
		consumer,  // 消费者标签
		autoAck,   // 自动应答
		exclusive, // 独占
		noLocal,   // 非本地
		noWait,    // 阻塞
		args,      // 额外参数
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
