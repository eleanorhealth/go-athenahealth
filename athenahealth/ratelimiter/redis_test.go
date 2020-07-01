package ratelimiter

import (
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedis_Allowed(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	rateLimiter := NewRedis(redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	}), 1, 1)

	retryAfter, err := rateLimiter.Allowed(true)
	assert.Zero(retryAfter)
	assert.NoError(err)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			rateLimiter.Allowed(true)
			wg.Done()
		}()
	}

	wg.Wait()

	retryAfterPreview, err := rateLimiter.Allowed(true)
	assert.NotZero(retryAfterPreview)
	assert.IsType(ErrRateExceeded, err)

	retryAfterProd, err := rateLimiter.Allowed(false)
	assert.Zero(retryAfterProd)
	assert.NoError(err)

	time.Sleep(retryAfterPreview)

	retryAfter, err = rateLimiter.Allowed(true)
	assert.Zero(retryAfter)
	assert.NoError(err)
}
