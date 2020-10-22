package tokencacher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault_IncrRequests(t *testing.T) {
	assert := assert.New(t)

	stats := NewDefault()
	err := stats.IncrRequests()
	assert.NoError(err)
}
