package tokencacher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefault_Get(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = "foo"
	cacher.expiresAt = time.Now().Add(time.Minute * 1)

	token, err := cacher.Get()

	assert.Equal(cacher.token, token)
	assert.Nil(err)
}

func TestDefault_Get_ErrTokenNotExist(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = ""
	cacher.expiresAt = time.Now().Add(time.Minute * 1)

	token, err := cacher.Get()

	assert.Empty(token)
	assert.Error(err)
	assert.IsType(ErrTokenNotExist, err)
}

func TestDefault_Get_ErrTokenExpired(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()
	cacher.token = ""
	cacher.expiresAt = time.Now().Add(-time.Minute * 1)

	token, err := cacher.Get()

	assert.Empty(token)
	assert.Error(err)
	assert.IsType(ErrTokenExpired, err)
}

func TestDefault_Set(t *testing.T) {
	assert := assert.New(t)

	cacher := NewDefault()

	token := "foo"
	expiresAt := time.Now().Add(time.Minute * 1)

	err := cacher.Set(token, expiresAt)

	assert.Equal(token, cacher.token)
	assert.True(expiresAt.Equal(cacher.expiresAt))
	assert.Nil(err)
}
