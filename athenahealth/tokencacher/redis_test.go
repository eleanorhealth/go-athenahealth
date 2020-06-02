package tokencacher

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
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

	token, err := cacher.Get()

	assert.Equal(expectedToken, token)
	assert.Nil(err)
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

	token, err := cacher.Get()

	assert.Empty(token)
	assert.Error(err)
	assert.IsType(ErrTokenNotExist, err)
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
	err = cacher.Set(expectedToken, time.Now().Add(time.Minute*1))

	assert.Nil(err)

	token, _ := s.Get(RedisDefaultKey)
	ttl := s.TTL(RedisDefaultKey)

	assert.Equal(expectedToken, token)
	assert.True(time.Now().Add(time.Second * ttl).After(time.Now()))
}
