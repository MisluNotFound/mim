package mq

import (
	"encoding/json"
	"fmt"
	"mim/internal/logic/dao"
	"mim/internal/logic/redis"
	"mim/pkg/mq"
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

		case TypePong:

		}
	}
}

func singleHandler(msg *dao.Message) {
	info, err := redis.GetUserInfo(msg.TargetID)
	if err != nil {
		zap.L().Error("singleHandler() failed: ", zap.Error(err))
		return
	}
	msg.Seq = snowflake.GenID()

	req := proto.PushMessageReq{
		SenderID: msg.SenderID,
		TargetID: msg.TargetID,
		Body:     msg.Content,
		Seq:      msg.Seq,
	}

	var status string
	if info.UserID == 0 {
		status = StatusOffline
	} else {
		status = StatusOnline
		go pushInMQ(info.ServerID, info.BucketID, req)
	}

	zap.L().Info("singleHandler receive message", zap.Any("msg", msg))
	go asyncSaveSingleMessage(*msg, status)
}

func groupHandler(msg *dao.Message) {
	fmt.Println("groupHandler receive message")
	msg.Seq = snowflake.GenID()
	// 先查询所有成员的状态
	userInfos, err := redis.GetUsersInfo(msg.TargetID)
	if err != nil {
		zap.L().Error("singleHandler() failed: ", zap.Error(err))
		return
	}
	// 根据成员状态异步处理
	for _, u := range userInfos {
		if u.UserID == 0 {
			// 告诉用户这个群有离线消息
			if u.IsNotice == 0 {
				redis.NoticeOfflineMessage(u.UserID, msg.TargetID)
			}
		} else if u.UserID != msg.SenderID {
			realSender := make(map[string]int64)
			realSender["realSender"] = msg.SenderID
			req := proto.PushMessageReq{
				TargetID: u.UserID,
				SenderID: msg.TargetID,
				Seq:      msg.Seq,
				Body:     msg.Content,
				Extra:    msg.SenderID,
			}
			pushInMQ(u.ServerID, u.BucketID, req)
		}
	}

	go asyncSaveGroupMessage(*msg)
}

func pongHandler(msg *dao.Message) {

}

func ackHandler(msg *dao.Message) {

}

func pushInMQ(serverID, bucketID int, req proto.PushMessageReq) {
	body, _ := json.Marshal(req)
	exchange := strconv.Itoa(serverID)
	routingKey := strconv.Itoa(bucketID)
	queueName := exchange + routingKey
	GetPublisher().PublishMessage(body, exchange, routingKey, queueName)
}

func asyncSaveSingleMessage(msg dao.Message, status string) {
	redis.StoreRedisMessage(msg, status)
	dao.StoreMysqlMessage(&msg)
}

func asyncSaveGroupMessage(msg dao.Message) {
	redis.StoreGroupMessage(msg)
	dao.StoreMysqlMessage(&msg)
}
