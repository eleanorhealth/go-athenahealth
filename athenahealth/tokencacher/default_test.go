package tokencacher

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefault_Get(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = "foo"
	cacher.expiresAt = time.Now().Add(time.Minute * 1)

	token, err := cacher.Get(context.Background())

	assert.Equal(cacher.token, token)
	assert.NoError(err)
}

func TestDefault_Get_ErrTokenNotExist(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = ""
	cacher.expiresAt = time.Now().Add(time.Minute * 1)

	token, err := cacher.Get(context.Background())

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenNotExist))
}

func TestDefault_Get_ErrTokenExpired(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = "foo"
	cacher.expiresAt = time.Now().Add(-time.Minute * 1)

	token, err := cacher.Get(context.Background())

	assert.Empty(token)
	assert.Error(err)
	assert.True(errors.Is(err, ErrTokenExpired))
}

func TestDefault_Set(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()

	token := "foo"
	expiresAt := time.Now().Add(time.Minute * 1)

	err := cacher.Set(context.Background(), token, expiresAt)

	assert.Equal(token, cacher.token)
	assert.True(expiresAt.Equal(cacher.expiresAt))
	assert.NoError(err)
}
