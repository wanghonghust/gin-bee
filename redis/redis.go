package redis

import (
	goredis "github.com/go-redis/redis/v8"
)

var RedisCli goredis.Client

func getRedisClient() goredis.Client {
	redisCl := goredis.NewClient(&goredis.Options{
		Addr:     "121.4.61.20:6379",
		Password: "Emergency520",
		DB:       2,
	})
	return *redisCl
}

func init() {
	RedisCli = getRedisClient()
}
