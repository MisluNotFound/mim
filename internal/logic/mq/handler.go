package mq

import (
	"encoding/json"
	"fmt"
	"mim/internal/logic/dao"
	"mim/internal/logic/redis"
	"mim/pkg/proto"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func singleHandler(msg *dao.Message) {
	info, err := redis.GetUserInfo(msg.TargetID)
	if err != nil {
		zap.L().Error("singleHandler() failed: ", zap.Error(err))
		return
	}
	// msg.Seq = snowflake.GenID()

	req := proto.PushMessageReq{
		SenderID: msg.SenderID,
		TargetID: msg.TargetID,
		Body:     msg.Content,
		Type:     msg.Type,
		URL:      msg.URL,
		Seq:      msg.Seq,
	}

	fmt.Println("userinfo", info)
	if info.UserID == 0 {
		// 记录离线消息
		go redis.AddUnReadCount(msg.TargetID, msg.SenderID, msg.Seq)
	} else {
		go pushInMQ(info.ServerID, info.BucketID, req)
	}

	zap.L().Info("singleHandler receive message", zap.Any("msg", msg))
	go asyncSaveMessage(*msg)
}

func groupHandler(msg *dao.Message) {
	fmt.Println("groupHandler receive message")
	// msg.Seq = snowflake.GenID()
	fmt.Println(msg)
	msg.IsGroup = true
	// 先查询所有成员的状态
	userInfos, err := redis.GetUsersInfo(msg.TargetID)
	if err != nil {
		zap.L().Error("singleHandler() failed: ", zap.Error(err))
		return
	}
	fmt.Println("userInfos", userInfos)
	// 根据成员状态异步处理
	for id, u := range userInfos {
		fmt.Println("userinfo		", u)
		// 用户离线则info为空
		if u.UserID == 0 {
			// 记录离线消息
			redis.AddUnReadCount(id, msg.TargetID, msg.Seq)
		} else if u.UserID != msg.SenderID {
			realSender := make(map[string]int64)
			realSender["realSender"] = msg.SenderID
			req := proto.PushMessageReq{
				// 群成员
				TargetID: u.UserID,
				// 表示是群的消息
				SenderID: msg.TargetID,
				Seq:      msg.Seq,
				Body:     msg.Content,
				Timer:    time.Now(),
				Type:     msg.Type,
				URL:      msg.URL,
				// 谁发的
				Extra: realSender,
			}
			go pushInMQ(u.ServerID, u.BucketID, req)
		}
	}

	go asyncSaveMessage(*msg)
}

func ackHandler(msg *dao.Message) {
	err := redis.AckMessage(msg.TargetID, msg.Seq)
	fmt.Println("ack handler receive message", msg)
	if err != nil {
		zap.L().Error("ackHandler failed: ", zap.Error(err), zap.Any("seq", msg.Seq))
	}
}

func errHandler(msg *dao.Message) {
	err := redis.ErrMessage(msg.TargetID, msg.Seq)
	if err != nil {
		zap.L().Error("ackHandler failed: ", zap.Error(err), zap.Any("seq", msg.Seq))
	}
}

func pushInMQ(serverID, bucketID int, req proto.PushMessageReq) {
	body, _ := json.Marshal(req)
	exchange := strconv.Itoa(serverID)
	routingKey := strconv.Itoa(bucketID)
	queueName := exchange + routingKey
	GetPublisher().PublishMessage(body, exchange, routingKey, queueName)
}

func asyncSaveMessage(msg dao.Message) {
	dao.StoreMysqlMessage(&msg)
}
