package db

import (
	"mim/setting"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRDB(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr: setting.Conf.RedisConfig.Addr,
	})
}