package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault_Request(t *testing.T) {
	assert := assert.New(t)

	stats := NewDefault()
	err := stats.Request()
	assert.NoError(err)
}

func TestDefault_ResponseSuccess(t *testing.T) {
	assert := assert.New(t)

	stats := NewDefault()
	err := stats.ResponseSuccess()
	assert.NoError(err)
}

func TestDefault_ResponseError(t *testing.T) {
	assert := assert.New(t)

	stats := NewDefault()
	err := stats.ResponseError()
	assert.NoError(err)
}
