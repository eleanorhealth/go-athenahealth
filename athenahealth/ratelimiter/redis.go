package ratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

const redisPreviewKey = "athena_rate_limit:preview"
const redisProdKey = "athena_rate_limit:prod"

type Redis struct {
	client  *redis.Client
	limiter *redis_rate.Limiter
}

func NewRedis(client *redis.Client) *Redis {
	if client == nil {
		panic("client is nil")
	}

	r := &Redis{
		client:  client,
		limiter: redis_rate.NewLimiter(client),
	}

	return r
}

func (r *Redis) Allowed(preview bool) (time.Duration, error) {
	var key string
	var limit *redis_rate.Limit

	if preview {
		key = redisPreviewKey
		limit = redis_rate.PerSecond(5)
	} else {
		key = redisProdKey
		limit = redis_rate.PerSecond(100)
	}

	res, err := r.limiter.Allow(context.Background(), key, limit)
	if err != nil {
		return 0, err
	}

	if res.RetryAfter > 0 {
		return res.RetryAfter, ErrRateExceeded
	}

	return 0, nil
}
