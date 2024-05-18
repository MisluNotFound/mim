package messaging

import (
	"fmt"
	"mim/internal/logic/redis"
	"mim/pkg/proto"
	"time"
)

type MessageReceiver struct {
	Queue chan *proto.MessageReq
}

const (
	TypePong = 1 + iota
	TypeSingle
	TypeGroup
	TypeAck
)

var Receiver MessageReceiver

func NewReceiver(size int) *MessageReceiver {
	return &MessageReceiver{
		Queue: make(chan *proto.MessageReq, size),
	}
}

func (mr *MessageReceiver) Start() {
	go mr.handleMessage()
	go mr.listenUnAckMessage()
}

func (mr *MessageReceiver) handleMessage() {
	for msg := range mr.Queue {
		m := &redis.Message{
			SenderID: msg.SenderID,
			TargetID: msg.TargetID,
			Body:     msg.Body,
		}

		switch msg.Type {
		case TypeSingle:
			singleHandler(m)
		case TypeGroup:
			groupHandler(m)
		case TypePong:
			pongHandler(m)
		case TypeAck:
			ackHandler(m)
		}
	}
}

func (mr *MessageReceiver) listenUnAckMessage() {
	for {
		msgs := redis.GetUnAckMessage()
		for _, m := range *msgs {
			// 直接重发
			fmt.Println(m)
		}
		time.Sleep(time.Second * 10)
	}
}
