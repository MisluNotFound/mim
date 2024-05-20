package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mim/db"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	prefixSession            = "session:"              // 记录所有会话 prefix+userid: sender+target
	prefixMessage            = "message-list:"         // 记录会话中的所有消息 prefix+sender+target: seq
	prefixMessageOffline     = "message-offline:"      // 索引 记录哪些用户给目标用户发送了离线消息 prefix+uid: sender
	prefixMessageOfflineList = "message-offline-list:" // 索引 记录离线消息列表 prefix+senderid: seq
	prefixMessageList        = "message-list"          // 记录所有未接收的消息 unack
)

func StoreRedisMessage(msg Message) {
	senderStr := strconv.FormatInt(msg.SenderID, 10)
	targetStr := strconv.FormatInt(msg.TargetID, 10)
	senderSession := prefixSession + senderStr
	targetSession := prefixSession + targetStr
	var sessionID string
	if msg.SenderID > msg.TargetID {
		sessionID += targetStr + senderStr
	} else {
		sessionID += senderStr + targetStr
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
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
		key := prefixMessage + sessionID
		if _, err := pipe.ZAdd(ctx, key, redis.Z{Member: msg.Seq, Score: float64(msg.Seq)}).Result(); err != nil {
			return err
		}

		// 添加到未接收消息列表中
		key = prefixMessageList
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
			fmt.Println(offlineMessageID)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	unackedMessages, err := db.RDB.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", threshold),
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

func GetMessages(uid, targetID, start int64, size int) ([]int64, error) {
	uidStr := strconv.FormatInt(uid, 10)
	targetStr := strconv.FormatInt(targetID, 10)
	session := prefixMessage
	if uidStr > targetStr {
		session = targetStr + uidStr
	} else {
		session = uidStr + targetStr
	}

	if targetID == 0 {
		return getAllMessage(uidStr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	messages, err := db.RDB.ZRangeByScoreWithScores(ctx, session, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", start),
		Max:    "+inf",
		Offset: 0,
		Count:  int64(size),
	}).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	messageIDs := make([]int64, len(messages))
	for i, msg := range messages {
		messageID, err := strconv.ParseInt(msg.Member.(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse message ID: %v", err)
		}
		messageIDs[i] = messageID
	}

	return messageIDs, nil
}

func getAllMessage(uid string) ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	var sessions []string

	// 获取sessions
	userSession := prefixSession + uid
	sessions = db.RDB.ZRange(ctx, userSession, 0, -1).Val()
	fmt.Println("sessions", sessions)
	// 根据session获取seqs
	var keys []string
	for _, session := range sessions {
		keys = append(keys, prefixMessage+session)
	}

	var cmds []redis.Cmder
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) (err error) {
		pipe := tx.TxPipeline()

		for _, k := range keys {
			if err := pipe.ZRevRange(ctx, k, 0, 9).Err(); err != nil {
				return err
			}
		}

		cmds, err = pipe.Exec(ctx)
		return
	}, keys...)

	if err != nil {
		return nil, err
	}

	var seqs []int64
	for _, cmd := range cmds {
		result, err := cmd.(*redis.StringSliceCmd).Result()
		if err != nil {
			zap.L().Error("get result error: ", zap.Error(err))
			return nil, err
		}

		for _, r := range result {
			seq, _ := strconv.ParseInt(r, 10, 64)
			seqs = append(seqs, seq)
		}
	}

	return seqs, nil
}

func GetOfflineMessages(uid int64) ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	// 获取哪些用户给我发送过离线消息
	var sessions []string
	offlineSession := prefixMessageOffline + strconv.FormatInt(uid, 10)
	sessions = db.RDB.SMembers(ctx, offlineSession).Val()
	// 根据会话查询所有未读消息
	var keys []string
	for _, session := range sessions {
		keys = append(keys, prefixMessageOfflineList+session)
	}
	fmt.Println("keys", keys)
	var cmds []redis.Cmder
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) (err error) {
		pipe := tx.TxPipeline()

		for _, k := range keys {
			if err = pipe.RPop(ctx, k).Err(); err != nil {
				return err
			}
		}

		cmds, err = pipe.Exec(ctx)
		return
	}, keys...)

	if err != nil {
		return nil, err
	}

	var seqs []int64
	for _, cmd := range cmds {
		result, err := cmd.(*redis.StringCmd).Result()
		if err != nil {
			zap.L().Error("get result error: ", zap.Error(err))
			return nil, err
		}

		seq, _ := strconv.ParseInt(result, 10, 64)
		seqs = append(seqs, seq)
	}

	return seqs, nil
}
