package mq

import (
	"encoding/json"
	"fmt"
	"mim/internal/logic/dao"
	"mim/pkg/mq"
	"mim/pkg/proto"

	"github.com/streadway/amqp"
)

const (
	TypePong = 1 + iota
	TypeSingle
	TypeGroup
	TypeAck
	TypeErr

	StatusOnline  = "online"
	StatusOffline = "offline"
)

var consumers []*mq.Consumer

func consumeMessage(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 解析消息
		p := proto.MessageReq{}
		json.Unmarshal(d.Body, &p)
		msg := dao.Message{
			SenderID: p.SenderID,
			TargetID: p.TargetID,
			Content:  p.Body,
			Type:     p.Media,
			URL:      p.URL,
		}
		fmt.Println("logic consumer receive message", d.Body)
		switch p.Type {
		case TypeSingle:
			go singleHandler(&msg)
			d.Ack(false)
		case TypeGroup:
			go groupHandler(&msg)
			d.Ack(false)
		case TypeAck:
			go ackHandler(&msg)
		case TypePong:

		case TypeErr:
			go errHandler(&msg)
		}
	}
}
