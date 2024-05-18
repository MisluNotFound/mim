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
	prefixSession        = "session:"
	prefixMessage  = "message-list:"
	prefixMessageOffline = "message-offline:"			// 索引 记录哪些用户给目标用户发送了离线消息
	prefixMessageOfflineList = "message-offline-list:"	// 索引 记录离线消息列表
	prefixMessageList = "message-list"					// 记录所有未接收的消息
)

func StoreRedisMessage(msg Message) {
	senderStr := strconv.FormatInt(msg.SenderID, 10)
	targetStr := strconv.FormatInt(msg.TargetID, 10)
	senderSession := prefixSession + senderStr
	targetSession := prefixSession + targetStr
	sessionID := prefixMessage + senderStr + targetStr

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 1000)
	defer cancel()
	db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()
		
		// 更新会话状态	
		pipe.ZAdd(ctx, senderSession, redis.Z{
			Member: sessionID,
			Score: float64(msg.Seq),
		})

		pipe.ZAdd(ctx, targetSession, redis.Z{
			Member: sessionID,
			Score: float64(msg.Seq),
		})

		// 插入消息列表中
		pipe.ZAdd(ctx, sessionID, redis.Z{
			Member: msg.Seq,
			Score: float64(msg.Seq),
		})

		// 添加到未接收消息列表中
		key := prefixMessageList
		val, _ := json.Marshal(msg)
		pipe.ZAdd(ctx, key, redis.Z{
			Member: val,
			Score: float64(time.Now().Unix()),
		})

		if msg.Status == "offline" {
			offlineMessageID := prefixMessageOffline + targetStr
			pipe.SAdd(ctx, offlineMessageID, msg.SenderID)
			offlineList := prefixMessageOfflineList + senderStr
			pipe.LPush(ctx, offlineList, msg.Seq)
		}

		return nil
	}, senderSession, targetSession)

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