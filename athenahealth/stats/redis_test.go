package stats

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedis_IncrRequests(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	stats := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), "")

	err = stats.IncrRequests(context.Background())
	assert.NoError(err)

	err = stats.IncrRequests(context.Background())
	assert.NoError(err)

	assert.Equal("2", s.HGet(RedisDefaultKey, "requests"))
	assert.NoError(err)
}
