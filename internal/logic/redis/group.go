package redis

import (
	"context"
	"mim/db"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func JoinGroup(userID, groupID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := getGroupUser(groupID)
	return db.RDB.SAdd(ctx, key, userID).Err()
}

func LeaveGroup(userID, gourpID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := getGroupUser(gourpID)
	uidstr := strconv.FormatInt(userID, 10)
	return db.RDB.SRem(ctx, key, uidstr).Err()
}

// 获取群成员状态
func GetUsersInfo(groupID int64) (map[int64]UserInfo, error) {
	groupKey := getGroupUser(groupID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	// 所有群成员
	users := db.RDB.SMembers(ctx, groupKey).Val()
	var userKeys []string

	// 键和结果集
	userInfos := make(map[int64]UserInfo)
	for _, userID := range users {
		userKeys = append(userKeys, getOnlineUser(userID))
		uid, _ := strconv.ParseInt(userID, 10, 64)
		userInfos[uid] = UserInfo{}
	}

	var cmds []redis.Cmder
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()
		for _, key := range userKeys {
			pipe.HGetAll(ctx, key)
		}

		var err error
		cmds, err = pipe.Exec(ctx)
		return err
	}, userKeys...)

	if err != nil {
		zap.L().Error("GetUsersInfo failed: ", zap.Error(err))
		return nil, err
	}

	for _, cmd := range cmds {
		results, err := cmd.(*redis.MapStringStringCmd).Result()
		if err == redis.Nil {
			return nil, err
		}
		if len(results) == 0 {
			continue
		}

		id, _ := strconv.ParseInt(results["user_id"], 10, 64)
		sid, _ := strconv.Atoi(results["server_id"])
		bid, _ := strconv.Atoi(results["bucket_id"])
		userInfos[id] = UserInfo{
			UserID:   id,
			ServerID: sid,
			BucketID: bid,
		}
	}

	return userInfos, nil
}
