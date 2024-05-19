package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"mim/db"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	prefixSession        = "session:"					// 记录所有会话 prefix+senderid+targetid
	prefixMessage  = "message-list:"					// 记录会话中的所有消息
	prefixMessageOffline = "message-offline:"			// 索引 记录哪些用户给目标用户发送了离线消息
	prefixMessageOfflineList = "message-offline-list:"	// 索引 记录离线消息列表 prefix+senderid: seq
	prefixMessageList = "message-list"					// 记录所有未接收的消息 unack
)

func StoreRedisMessage(msg Message) {
	senderStr := strconv.FormatInt(msg.SenderID, 10)
	targetStr := strconv.FormatInt(msg.TargetID, 10)
	senderSession := prefixSession + senderStr
	targetSession := prefixSession + targetStr
	sessionID := prefixMessage + senderStr + targetStr

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 1000)
	defer cancel()
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()
	
		// 更新会话状态
		if _, err := pipe.ZAdd(ctx, senderSession, redis.Z{Member: sessionID, Score: float64(msg.Seq)}).Result(); err != nil {
			return err
		}
		if _, err := pipe.ZAdd(ctx, targetSession, redis.Z{Member: sessionID, Score: float64(msg.Seq)}).Result(); err != nil {
			return err
		}
	
		// 插入消息列表中
		if _, err := pipe.ZAdd(ctx, sessionID, redis.Z{Member: msg.Seq, Score: float64(msg.Seq)}).Result(); err != nil {
			return err
		}
	
		// 添加到未接收消息列表中
		key := prefixMessageList
		val, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		if _, err := pipe.ZAdd(ctx, key, redis.Z{Member: val, Score: float64(time.Now().Unix())}).Result(); err != nil {
			return err
		}
	
		// 处理离线消息
		if msg.Status == "offline" {
			offlineMessageID := prefixMessageOffline + targetStr
			if _, err := pipe.SAdd(ctx, offlineMessageID, msg.SenderID).Result(); err != nil {
				return err
			}
			offlineList := prefixMessageOfflineList + senderStr
			if _, err := pipe.LPush(ctx, offlineList, msg.Seq).Result(); err != nil {
				return err
			}
		}
	
		// 执行事务
		_, err = pipe.Exec(ctx)
		return err
	}, senderSession, targetSession)
	
	if err != nil {
		zap.L().Error("store message to redis failed: ", zap.Error(err))
	}

}

func StoreOfflineMessage(msg Message) error {
	senderStr := strconv.FormatInt(msg.SenderID, 10)
	targetStr := strconv.FormatInt(msg.TargetID, 10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 1000)
	defer cancel()
	
	offlineMessageID := prefixMessageOffline + targetStr
	offlineList := prefixMessageOfflineList + senderStr
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()
		
		if _, err := pipe.SAdd(ctx, offlineMessageID, msg.SenderID).Result(); err != nil {
			return err
		}

		if _, err := pipe.LPush(ctx, offlineList, msg.Seq).Result(); err != nil {
			return err
		}

		return nil
	}, offlineMessageID, offlineList)
	
	return err
}


func GetUnAckMessage() *[]Message {
	key := prefixMessageList
	now := float64(time.Now().Unix())
	threshold := now - 10 // 10秒之前的时间
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 500)
	defer cancel()

	unackedMessages, err := db.RDB.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    fmt.Sprintf("%f", threshold),
	}).Result()

	if err != nil {
		zap.L().Error("GetUnAckMessge() failed: ", zap.Error(err))
		return nil
	}

	msgs := []Message{}
	for _, m := range unackedMessages {
		msg := Message{}
		json.Unmarshal([]byte(m), &msg)
		msgs = append(msgs, msg)
	}

	return &msgs
}