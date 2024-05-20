package redis

import (
	"context"
	"errors"
	"mim/db"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var prefixOnlineUser = "online-user:"

func AddOnlineUser(uid int64, sid int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := prefixOnlineUser + strconv.FormatInt(uid, 10)
	info := make(map[string]interface{})
	info["server_id"] = sid
	info["user_id"] = uid
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
	u.ServerID = sid
	u.UserID = id
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