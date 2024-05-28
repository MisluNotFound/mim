package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mim/db"
	"mim/internal/logic/dao"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetSessionID(senderID, targetID int64) string {
	senderStr := strconv.FormatInt(senderID, 10)
	targetStr := strconv.FormatInt(targetID, 10)
	var sessionID string
	if senderID > targetID {
		sessionID += targetStr + senderStr
	} else {
		sessionID += senderStr + targetStr
	}

	return sessionID
}

func getUserSession(uid int64) string {
	uidStr := strconv.FormatInt(uid, 10)

	return prefixSession + uidStr
}

// 获取最新消息id
func getStart(uid, targetID int64) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	sessionID := GetSessionID(uid, targetID)
	key := getUserSession(uid)
	start := db.RDB.ZScore(ctx, key, sessionID).Val()

	return int64(start)
}

func StoreRedisMessage(msg dao.Message, status string) {
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

	if status == "offline" {
		offlineMessageID := prefixMessageOffline + targetStr

		val, _ := json.Marshal(msg)
		db.RDB.SAdd(ctx, offlineMessageID, msg.SenderID).Result()

		offlineList := prefixMessageOfflineList + senderStr
		db.RDB.LPush(ctx, offlineList, val).Result()

		return
	}

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
		val, _ := json.Marshal(msg)
		if _, err := pipe.ZAdd(ctx, key, redis.Z{Member: val, Score: float64(msg.Seq)}).Result(); err != nil {
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

		// 执行事务
		_, err = pipe.Exec(ctx)
		return err
	}, senderSession, targetSession)

	if err != nil {
		zap.L().Error("store message to redis failed: ", zap.Error(err))
	}

}

// 用户突然下线，将消息存到离线列表中
func StoreOfflineMessage(msg dao.Message) error {
	senderStr := strconv.FormatInt(msg.SenderID, 10)
	targetStr := strconv.FormatInt(msg.TargetID, 10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	offlineMessageID := prefixMessageOffline + targetStr
	offlineList := prefixMessageOfflineList + senderStr
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()

		val, _ := json.Marshal(msg)
		if _, err := pipe.SAdd(ctx, offlineMessageID, msg.SenderID).Result(); err != nil {
			return err
		}

		if _, err := pipe.LPush(ctx, offlineList, val).Result(); err != nil {
			return err
		}

		return nil
	}, offlineMessageID, offlineList)

	return err
}

// 暂时用来获取未ack的消息
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

func GetMessages(uid, targetID, start int64, size int) ([]dao.Message, error) {
	uidStr := strconv.FormatInt(uid, 10)
	targetStr := strconv.FormatInt(targetID, 10)
	session := prefixMessage
	if uidStr > targetStr {
		session = targetStr + uidStr
	} else {
		session = uidStr + targetStr
	}

	// 用户刚上线，拉取所有会话中的部分消息
	if targetID == 0 {
		return getAllMessage(uidStr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	if start == 0 {
		start = getStart(uid, targetID)
	}

	// 拉取特定会话中的部分消息
	result, err := db.RDB.ZRevRangeByScoreWithScores(ctx, session, &redis.ZRangeBy{
		Max:    fmt.Sprintf("%d", start),
		Min:    "-inf",
		Offset: 0,
		Count:  int64(size),
	}).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var messages []dao.Message
	for _, r := range result {
		m := dao.Message{}
		json.Unmarshal([]byte(r.Member.(string)), &m)
		messages = append(messages, m)
	}

	return messages, nil
}

func GetUserSessions(uid int64) ([]string, error) {
	userSessions := getUserSession(uid)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	return db.RDB.ZRevRange(ctx, userSessions, 0, -1).Result()
}

func getAllMessage(uid string) ([]dao.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	var sessions []string

	// 获取sessions
	userSession := prefixSession + uid
	sessions = db.RDB.ZRange(ctx, userSession, 0, -1).Val()
	fmt.Println("sessions", sessions)
	// 根据session获取messages
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

	// 可以优化 慢
	var messages []dao.Message
	for _, cmd := range cmds {
		result, err := cmd.(*redis.StringSliceCmd).Result()
		if err != nil {
			zap.L().Error("get result error: ", zap.Error(err))
			return nil, err
		}

		for _, r := range result {
			m := dao.Message{}
			json.Unmarshal([]byte(r), &m)
			messages = append(messages, m)
		}
	}

	return messages, nil
}

func GetOfflineMessages(uid int64) ([]dao.Message, error) {
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

	var cmds []redis.Cmder
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) (err error) {
		pipe := tx.TxPipeline()

		for _, k := range keys {
			if err = pipe.LPop(ctx, k).Err(); err != nil {
				return err
			}
		}

		cmds, err = pipe.Exec(ctx)
		return
	}, keys...)

	if err != nil {
		return nil, err
	}

	var messages []dao.Message
	for _, cmd := range cmds {
		result, err := cmd.(*redis.StringCmd).Result()
		if err != nil {
			zap.L().Error("get result error: ", zap.Error(err))
			return nil, err
		}

		m := dao.Message{}
		json.Unmarshal([]byte(result), &m)
		messages = append(messages, m)
	}

	go storeAsRead(messages)

	return messages, nil
}

func storeAsRead(messages []dao.Message) {
	if len(messages) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	session := GetSessionID(messages[0].SenderID, messages[0].TargetID)
	key := prefixMessage + session
	db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()

		for _, m := range messages {
			val, _ := json.Marshal(m)
			if err := pipe.ZAdd(ctx, key, redis.Z{
				Member: val,
				Score:  float64(m.Seq),
			}).Err(); err != nil {
				return err
			}
		}

		pipe.Exec(ctx)
		return nil
	}, key)
}

// 未命中 回写
func WriteBack(uid int64, messages []dao.Message) error {
	if len(messages) == 0 {
		return nil
	}

	// 获取会话
	session := GetSessionID(messages[0].SenderID, messages[0].TargetID)
	// 放
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := prefixMessage + session
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()

		for _, m := range messages {
			val, _ := json.Marshal(m)
			if err := pipe.ZAdd(ctx, key, redis.Z{
				Member: val,
				Score:  float64(m.Seq),
			}).Err(); err != nil {
				return err
			}
		}

		return nil
	}, key)

	if err != nil {
		zap.L().Error("write message back failed: ", zap.Error(err))
		return err
	}

	return nil
}