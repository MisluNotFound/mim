package redis

import (
	"context"
	"mim/db"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// 记录离线消息数量
func AddUnReadCount(uid int64, senderID int64, seq int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	offlineKey := getSessionOfflineCount(uid, senderID)
	sessionKey := getUserSession(uid)
	senderStr := strconv.FormatInt(senderID, 10)
	if err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.Pipeline()

		pipe.ZAdd(ctx, sessionKey, redis.Z{
			Member: senderStr,
			Score:  1,
		})

		pipe.SetNX(ctx, offlineKey, seq, 0)
		_, err := pipe.Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func MarkAsRead(uid, senderID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	userSession := getUserSession(uid)
	offlineCount := getSessionOfflineCount(uid, senderID)
	if err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.Pipeline()

		pipe.ZAdd(ctx, userSession, redis.Z{
			Member: senderID,
			Score:  0,
		})
		pipe.Del(ctx, offlineCount)
		_, err := pipe.Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	return nil
}

// 获取某个会话的lastRead
func GetLastRead(uid, senderID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	offlineCount := getSessionOfflineCount(uid, senderID)
	seqStr, err := db.RDB.Get(ctx, offlineCount).Result()
	if err != nil {
		return 0, err
	}

	seq, _ := strconv.ParseInt(seqStr, 10, 64)
	return seq, nil
}
