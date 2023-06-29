package cache

import (
	"fmt"
	"nsproxy/cache"
	"nsproxy/cache/redis"
	"nsproxy/config"

	goredis "github.com/redis/go-redis/v9"
)

var rcache cache.Cache

func GetCache() cache.Cache {
	return rcache
}

func ConnectToRedis() {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GetConfig().RedisHost, config.GetConfig().RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rcache = redis.NewRedisCache(rdb)
}
