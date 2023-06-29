package redis

import (
	"context"
	"nsproxy/cache"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func (r *redisCache) Get(key string) (interface{}, error) {
	return r.client.Get(context.Background(), key).Result()
}
func (r *redisCache) Set(key string, value interface{}, TTL time.Duration) error {
	return r.client.Set(context.Background(), key, value, TTL).Err()

}

func NewRedisCache(client *redis.Client) cache.Cache {
	return &redisCache{client: client}

}
