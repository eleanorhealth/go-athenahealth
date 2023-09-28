package tokencacher

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedis_Get(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	expectedToken := "foo"

	s.Set(RedisDefaultKey, expectedToken)
	s.SetTTL(RedisDefaultKey, time.Minute*1)

	cacher := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), "")

	token, err := cacher.Get(context.Background())

	assert.Equal(expectedToken, token)
	assert.NoError(err)
}

func TestRedis_Get_ErrTokenNotExist(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	cacher := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), "")

	token, err := cacher.Get(context.Background())

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenNotExist))
}

func TestRedis_Set(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	cacher := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), "")

	expectedToken := "foo"
	err = cacher.Set(context.Background(), expectedToken, time.Now().Add(time.Minute*1))

	assert.NoError(err)

	token, _ := s.Get(RedisDefaultKey)
	ttl := s.TTL(RedisDefaultKey)

	assert.Equal(expectedToken, token)
	assert.True(time.Now().Add(time.Second * ttl).After(time.Now()))
}
