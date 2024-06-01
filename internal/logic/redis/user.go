package redis

import (
	"context"
	"errors"
	"mim/db"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func AddOnlineUser(uid int64, sid, bucketId int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := prefixOnlineUser + strconv.FormatInt(uid, 10)
	info := make(map[string]interface{})
	info["server_id"] = sid
	info["user_id"] = uid
	info["bucket_id"] = bucketId
	err = db.RDB.HSet(ctx, key, info).Err()

	return
}

func RemoveOnlineUser(uid int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := prefixOnlineUser + strconv.FormatInt(uid, 10)
	err = db.RDB.Del(ctx, key).Err()
	return
}

func GetUserInfo(uid int64) (u UserInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := prefixOnlineUser + strconv.FormatInt(uid, 10)
	infos, err := db.RDB.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return
	}

	sid, _ := strconv.Atoi(infos["server_id"])
	id, _ := strconv.ParseInt(infos["user_id"], 10, 64)
	bid, _ := strconv.Atoi(infos["bucket_id"])

	u.ServerID = sid
	u.UserID = id
	u.BucketID = bid
	return
}

func Close() {
	pattern := prefixOnlineUser + "*"
	ctx := context.Background()

	iter := db.RDB.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		db.RDB.Del(ctx, key)
	}
}

// 记录lastAck
func AckMessage(uid, seq int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	// 先找出user的lastErr和lastAck
	key := getAckMessage(uid)
	seqs, err := db.RDB.HGetAll(ctx, key).Result()
	if !errors.Is(err, redis.Nil) {
		return err
	}
	lastAck, _ := strconv.ParseInt(seqs["last_ack"], 10, 64)
	lastErr, _ := strconv.ParseInt(seqs["last_err"], 10, 64)
	// 做判断
	if lastErr == 0 && seq > lastAck {
		lastAck = seq
	} else {
		if seq < lastErr {
			lastAck = seq
		}
	}
	// 更新
	seqs["last_ack"] = strconv.FormatInt(lastAck, 10)
	if err != db.RDB.HMSet(ctx, key, seqs).Err() {
		return err
	}

	return nil
}

// 记录lastErr
func ErrMessage(uid, seq int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	// 先找出user的lastErr和lastAck
	key := getAckMessage(uid)
	seqs, err := db.RDB.HGetAll(ctx, key).Result()
	if !errors.Is(err, redis.Nil) {
		return err
	}
	lastAck, _ := strconv.ParseInt(seqs["last_ack"], 10, 64)
	lastErr, _ := strconv.ParseInt(seqs["last_err"], 10, 64)
	if lastErr == 0 {
		lastErr = seq
	} else if seq <= lastErr {
		lastErr = seq
	} else if seq > lastAck {
		lastAck = seq - 1
	}

	return db.RDB.HSet(ctx, key, "last_err", lastErr, "last_ack", lastAck).Err()
}

// 获取会话中的离线消息数
func GetUnReadCount(uid int64) (map[int64]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	// 其实是对方id
	key := getUserSession(uid)
	sessions, err := db.RDB.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "1",
		Max: "1",
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	// 构造会话key
	var offlineKeys []string
	for _, s := range sessions {
		offlineKeys = append(offlineKeys, getSessionOfflineCount(uid, s))
	}
	var cmder []redis.Cmder
	if err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.Pipeline()

		for _, k := range offlineKeys {
			pipe.Get(ctx, k)
		}

		var err error
		cmder, err = pipe.Exec(ctx)
		return err
	}, offlineKeys...); err != nil {
		return nil, err
	}

	counts := make(map[int64]int, 0)
	for i, cmd := range cmder {
		if val, ok := cmd.(*redis.StringCmd); ok {
			resultStr := val.Val()
			sessionStr := sessions[i]
			r, _ := strconv.Atoi(resultStr)
			s, _ := strconv.ParseInt(sessionStr, 10, 64)
			counts[s] = r
		}
	}
	return counts, nil
}

func GetLastAck(uid int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := getAckMessage(uid)
	ackStr, err := db.RDB.HGet(ctx, key, "last_ack").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return -1, nil
		}
		return -1, err
	}

	lastAck, _ := strconv.ParseInt(ackStr, 10, 64)
	return lastAck, nil
}

func ClearErrAck(uid int64, lastAck int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := getAckMessage(uid)
	return db.RDB.HSet(ctx, key, "last_ack", lastAck, "last_err", "+inf").Err()
}
