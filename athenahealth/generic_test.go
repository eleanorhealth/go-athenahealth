package athenahealth

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestBodyToErrorString(t *testing.T) {
	assert := assert.New(t)

	t.Run("raw string response - returns an error", func(t *testing.T) {
		input := io.NopCloser(strings.NewReader(`"My name is Athena and I don't understand JSON"`))
		err := BodyToErrorString(input)
		assert.Error(err)
		assert.Equal("My name is Athena and I don't understand JSON", err.Error())
	})
	t.Run("invalid json - returns a different error", func(t *testing.T) {
		input := io.NopCloser(strings.NewReader(""))
		err := BodyToErrorString(input)
		assert.Error(err)
		assert.Equal("unexpected end of JSON input", err.Error())
	})
}
