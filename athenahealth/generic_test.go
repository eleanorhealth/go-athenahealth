package athenahealth

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalIntStr(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name     string
		input    string
		expected StrInt
		isError  bool
	}{
		{
			name:  "Success case - string message",
			input: `{"status": "My name is Athena and I don't understand JSON types"}`,
			expected: StrInt{
				IntVal: 0,
				StrVal: "My name is Athena and I don't understand JSON types",
			},
			isError: false,
		},
		{
			name:  "Success case - number as a string",
			input: `{"status": "42"}`,
			expected: StrInt{
				IntVal: 42,
				StrVal: "42",
			},
			isError: false,
		},
		{
			name:  "Success case - number as a number",
			input: `{"status": 42}`,
			expected: StrInt{
				IntVal: 42,
				StrVal: "42",
			},
			isError: false,
		},
		{
			name:  "Error case - bool value",
			input: `{"status": true}`,
			expected: StrInt{
				IntVal: 0,
				StrVal: "",
			},
			isError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var response StatusResponse
			err := json.Unmarshal([]byte(tc.input), &response)

			if tc.isError {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.expected.IntVal, response.Status.IntVal)
				assert.Equal(tc.expected.StrVal, response.Status.StrVal)
			}
		})
	}
}
