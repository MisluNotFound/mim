package redis

import (
	"context"
	"mim/db"
	"mim/internal/logic/dao"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// 用户突然下线，将消息记录到离线消息数量中
func StoreOfflineMessage(msg dao.Message) error {

	return nil
}

func GetUserSessions(uid int64) ([]string, error) {
	userSessions := getUserSession(uid)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	return db.RDB.ZRevRange(ctx, userSessions, 0, -1).Result()
}

// 记录离线消息数量
func AddUnReadCount(uid int64, senderID int64) error {
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
		pipe.Incr(ctx, offlineKey)
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
		pipe.Set(ctx, offlineCount, 0, 0)
		_, err := pipe.Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	return nil
}
