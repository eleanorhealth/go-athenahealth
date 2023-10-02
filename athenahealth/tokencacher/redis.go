package tokencacher

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const RedisDefaultKey = "athena_token"

type Redis struct {
	client redis.UniversalClient
	key    string
}

func NewRedis(client redis.UniversalClient, key string) *Redis {
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

func (r *Redis) Get(ctx context.Context) (string, error) {
	val, err := r.client.Get(context.Background(), r.key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenNotExist
		}

		return "", err
	}

	return val, nil
}

func (r *Redis) Set(ctx context.Context, token string, expiresAt time.Time) error {
	_, err := r.client.Set(context.Background(), r.key, token, time.Second*time.Duration(expiresAt.Unix()-time.Now().Unix())).Result()

	return err
}
