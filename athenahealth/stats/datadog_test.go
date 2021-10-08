package stats

import (
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	statsd.ClientInterface
	incrFn func(name string, tags []string, rate float64) error
}

func (m *mockClient) Incr(name string, tags []string, rate float64) error {
	return m.incrFn(name, tags, rate)
}

func TestDatadog_Request(t *testing.T) {
	assert := assert.New(t)

	client := &mockClient{}

	client.incrFn = func(name string, tags []string, rate float64) error {
		assert.Equal("path:GET /patients/{id}", tags[0])
		return nil
	}

	datadog := NewDatadog(client)

	err := datadog.Request("get", "/patients/123")
	assert.NoError(err)
}

func TestRemoveIDsFromPath(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("/patients/{id}", removeIDsFromPath("/patients/123"))
	assert.Equal("/patients/{id}", removeIDsFromPath("/patients/123"))
	assert.Equal("/patients/{id}/foo/{id}", removeIDsFromPath("/patients/123/foo/1"))
	assert.Equal("/patients/{id}/foo/{id}/", removeIDsFromPath("/patients/123/foo/1/"))
}
