package utils

import (
	"context"
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

func ConnectToRedis() error {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GetConfig().RedisHost, config.GetConfig().RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	if pong != "PONG" {
		return fmt.Errorf("expected PONG, got %s", pong)
	}

	rcache = redis.NewRedisCache(rdb)

	return nil
}
