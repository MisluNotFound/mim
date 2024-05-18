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
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 500)
	defer cancel()

	key := prefixOnlineUser + strconv.FormatInt(uid, 10)
	info := make(map[string]interface{})
	info["server_id"] = sid
	err = db.RDB.HMSet(ctx, key, info).Err()
	
	return 
}

func GetUserInfo(uid int64) (u UserInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 500)
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
	u.ServerID = sid
	u.UserID = uid
	return
}