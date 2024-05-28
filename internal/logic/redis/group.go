package redis

import (
	"context"
	"encoding/json"
	"errors"
	"mim/db"
	"mim/internal/logic/dao"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func JoinGroup(userID, groupID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := getGroupUser(groupID)
	return db.RDB.HSet(ctx, key, userID, 0).Err()
}

func LeaveGroup(userID, gourpID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := getGroupUser(gourpID)
	uidstr := strconv.FormatInt(userID, 10)
	return db.RDB.HDel(ctx, key, uidstr).Err()
}

// 获取群成员状态
func GetUsersInfo(groupID int64) (map[int64]UserInfo, error) {
	groupKey := getGroupUser(groupID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	users := db.RDB.HGetAll(ctx, groupKey).Val()
	var userKeys []string

	// 键和结果集
	userInfos := make(map[int64]UserInfo)
	for userID, isNOtice := range users {
		userKeys = append(userKeys, getOnlineUser(userID))
		is, _ := strconv.Atoi(isNOtice)
		uid, _ := strconv.ParseInt(userID, 10, 64)
		userInfos[uid] = UserInfo{
			IsNotice: is,
		}
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

// 有离线群消息
func NoticeOfflineMessage(userID, groupID int64) error {
	keyOffline := getGroupOffline(userID)
	keyGroup := getGroupUser(groupID)
	uidStr := strconv.FormatInt(userID, 10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()

		// 告诉用户这个群有未读消息
		pipe.SAdd(ctx, keyOffline, groupID)
		// 已经通知
		info := make(map[string]interface{})
		info[uidStr] = 1
		pipe.HSet(ctx, keyGroup, info)
		_, err := pipe.Exec(ctx)

		return err
	}, keyOffline, keyGroup)

	return err
}

// 用户上线了，清除通知状态
func MarkedRead(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*500)
	defer cancel()

	key := getGroupOffline(userID)
	uidStr := strconv.FormatInt(userID, 10)
	// 获取哪些群有离线消息
	groups := db.RDB.SMembers(ctx, key).Val()
	var groupKeys []string
	for i := 0; i < len(groups); i++ {
		groupKeys = append(groupKeys, prefixGroupUser+groups[i])
	}
	// 更新在这些群里的状态
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.Pipeline()
		for _, k := range groupKeys {
			info := make(map[string]interface{})
			info[uidStr] = 0
			pipe.HSet(ctx, k, info)
		}

		_, err := pipe.Exec(ctx)
		return err
	}, groupKeys...)

	return err
}

func StoreGroupMessage(msg dao.Message) {
	// 存入群消息缓存
	key := getGroupMessage(msg.TargetID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	msgStr, _ := json.Marshal(msg)
	db.RDB.ZAdd(ctx, key, redis.Z{Member: msgStr, Score: float64(msg.Seq)})
}

// 用于用户查询历史记录
func GetGroupMessage(userID, groupID, start int64, size int) ([]dao.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	key := getGroupMessage(groupID)
	// 获取用户入群时间
	earlyKey := getEarlyMessage(userID)
	earlySeq := db.RDB.Get(ctx, earlyKey).Val()
	seq, _ := strconv.ParseInt(earlySeq, 10, 64)
	var ss string
	if start == 0 {
		ss = "+inf"
	} else {
		ss = strconv.FormatInt(start, 10)
	}

	// 不合法
	if start < seq {
		return nil, nil
	}

	// 尝试获取
	result, err := db.RDB.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Max:    ss,
		Min:    earlySeq,
		Offset: 0,
		Count:  int64(size),
	}).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	// 未命中
	if len(result) < size {
		dao.GetMessages(0, groupID, start, size, 2)
	}
	var messages []dao.Message
	for _, r := range result {
		m := dao.Message{}
		json.Unmarshal([]byte(r.Member.(string)), &m)
		messages = append(messages, m)
	}

	return messages, nil
}

// 用于用户上线时拉取离线消息
func GetGroupOfflineMessage(userID int64) ([]dao.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	// 获取哪些群聊有消息
	key := getGroupOffline(userID)
	var groups []string
	err := db.RDB.Watch(ctx, func(tx *redis.Tx) (err error) {
		pipe := tx.Pipeline()

		cmd := pipe.SMembers(ctx, key)
		pipe.Del(ctx, key).Err()
		_, err = pipe.Exec(ctx)

		groups, _ = cmd.Result()
		return
	}, key)

	if err != nil {
		return nil, err
	}

	// 获取最后一条消息
	lastMessageKey := getLastMessage(userID)
	lastMessages := db.RDB.HGetAll(ctx, lastMessageKey).Val()

	// 尝试获取
	var cmds []redis.Cmder
	var groupKeys []string
	for _, group := range groups {
		groupKeys = append(groupKeys, prefixGroupMessage+group)
	}
	err = db.RDB.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.Pipeline()

		for i := 0; i < len(groupKeys); i++ {
			// 获取从last开始到最新的所有消息
			pipe.ZRevRangeByScoreWithScores(ctx, groupKeys[i], &redis.ZRangeBy{
				Min: lastMessages[groups[i]],
				Max: "+inf",
			})
		}

		cmds, err = pipe.Exec(ctx)
		return err
	}, groupKeys...)

	if err != nil {
		return nil, err
	}

	// 解析所有消息
	var messages []dao.Message
	for i, cmd := range cmds {
		result, err := cmd.(*redis.ZSliceCmd).Result()
		if err != nil {
			return nil, err
		}

		last, _ := strconv.ParseFloat(lastMessages[groups[i]], 64)
		// 未命中
		if result[len(result)-1].Score != last {
			gid, _ := strconv.ParseInt(groups[i], 10, 64)
			msgs, err := dao.GetMessages(0, gid, int64(last), -1, 2)
			if err != nil {
				return nil, err
			}

			messages = append(messages, msgs...)
			continue
		} else {
			for _, r := range result {
				msg := dao.Message{}
				json.Unmarshal([]byte(r.Member.(string)), &msg)
				messages = append(messages, msg)
			}
		}
	}

	return messages, nil
}
