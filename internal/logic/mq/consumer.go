package mq

import (
	"encoding/json"
	"fmt"
	"mim/internal/logic/dao"
	"mim/pkg/mq"
	"mim/pkg/proto"
	"strconv"

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

		fmt.Println("logic consumer receive message", p)
		target, _ := strconv.ParseInt(p.TargetID, 10, 64)
		msg := dao.Message{
			Seq:      p.Seq,
			SenderID: p.SenderID,
			TargetID: target,
			Content:  p.Body,
			Type:     p.Media,
			URL:      p.URL,
		}
		switch p.Type {
		case TypeSingle:
			go singleHandler(&msg)
		case TypeGroup:
			go groupHandler(&msg)
		case TypeAck:
			go ackHandler(&msg)
		case TypeErr:
			go errHandler(&msg)
		}
	}
}
