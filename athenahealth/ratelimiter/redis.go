package ratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

const redisKeyPreview = "athena_rate_limit:preview"
const redisKeyProd = "athena_rate_limit:prod"

const defaultRatePerSecPreview = 5
const defaultRatePerSecProd = 100

type Redis struct {
	client  *redis.Client
	limiter *redis_rate.Limiter

	ratePreivew int
	rateProd    int
}

func NewRedis(client *redis.Client, ratePreview, rateProd int) *Redis {
	if client == nil {
		panic("client is nil")
	}

	if ratePreview <= 0 {
		ratePreview = defaultRatePerSecPreview
	}

	if rateProd <= 0 {
		rateProd = defaultRatePerSecProd
	}

	r := &Redis{
		client:  client,
		limiter: redis_rate.NewLimiter(client),

		ratePreivew: ratePreview,
		rateProd:    rateProd,
	}

	return r
}

func (r *Redis) Allowed(ctx context.Context, preview bool) (time.Duration, error) {
	var key string
	var limit redis_rate.Limit

	if preview {
		key = redisKeyPreview
		limit = redis_rate.PerSecond(r.ratePreivew)
	} else {
		key = redisKeyProd
		limit = redis_rate.PerSecond(r.rateProd)
	}

	res, err := r.limiter.Allow(ctx, key, limit)
	if err != nil {
		return 0, err
	}

	if res.RetryAfter > 0 {
		return res.RetryAfter, ErrRateExceeded
	}

	return 0, nil
}
