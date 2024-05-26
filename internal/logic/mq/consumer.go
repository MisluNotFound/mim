package mq

import (
	"encoding/json"
	"fmt"
	"mim/internal/logic/dao"
	"mim/internal/logic/redis"
	"mim/pkg/proto"
	"mim/pkg/snowflake"
	"strconv"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
	TypePong = 1 + iota
	TypeSingle
	TypeGroup
	TypeAck

	StatusOnline  = "online"
	StatusOffline = "offline"
)

func consumeMessage(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 解析消息
		p := proto.MessageReq{}
		json.Unmarshal(d.Body, &p)
		msg := dao.Message{
			SenderID: p.SenderID,
			TargetID: p.TargetID,
			Content:  p.Body,
		}
		fmt.Println("logic consumer receive message", d.Body)
		switch p.Type {
		case TypeSingle:
			go singleHandler(&msg)
			d.Ack(false)
		case TypeGroup:

		case TypeAck:

		case TypePong:

		}
	}
}

func singleHandler(msg *dao.Message) {
	info, err := redis.GetUserInfo(msg.TargetID)
	if err != nil {
		zap.L().Error("singleHandler() failed: ", zap.Error(err))
	}
	msg.Seq = snowflake.GenID()

	fmt.Println(info)
	var status string
	if info.UserID == 0 {
		status = StatusOffline
	} else {
		status = StatusOnline
		go pushInMQ(info.ServerID, info.BucketID, *msg)
	}

	zap.L().Info("singleHandler receive message", zap.Any("msg", msg))
	go asyncSaveMessage(msg, status)
}

func groupHandler(msg *redis.Message) {

}

func pongHandler(msg *redis.Message) {

}

func ackHandler(msg *redis.Message) {

}

func pushInMQ(serverID, bucketID int, msg dao.Message) {
	req := proto.PushMessageReq{
		SenderID: msg.SenderID,
		TargetID: msg.TargetID,
		Body:     msg.Content,
		Seq:      msg.Seq,
	}

	body, _ := json.Marshal(req)
	exchange := strconv.Itoa(serverID)
	routingKey := strconv.Itoa(bucketID)
	queueName := exchange + routingKey
	GetPublisher().PublishMessage(body, exchange, routingKey, queueName)
}

func asyncSaveMessage(msg *dao.Message, status string) {
	redis.StoreRedisMessage(*msg, status)
	dao.StoreMysqlMessage(&dao.Message{
		SenderID: msg.SenderID,
		TargetID: msg.TargetID,
		Content:  msg.Content,
		Seq:      msg.Seq,
	})
}
