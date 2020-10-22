package tokencacher

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const RedisDefaultKey = "athena_stats"

type Redis struct {
	client *redis.Client
	key    string
}

func NewRedis(client *redis.Client, key string) *Redis {
	if client == nil {
		panic("client is nil")
	}

	r := &Redis{
		client: client,
		key:    key,
	}

	if len(r.key) == 0 {
		r.key = RedisDefaultKey
	}

	return r
}

func (r *Redis) IncrRequests() error {
	_, err := r.client.HIncrBy(context.Background(), r.key, "requests", 1).Result()
	return err
}
