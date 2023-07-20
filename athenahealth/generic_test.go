package athenahealth

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalIntStr(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name         string
		errorValue   error
		expected     StatusResponse
		input        string
		returnsError bool
	}{
		{
			name:       "Success case - string message",
			input:      `{"status": "my name is Athena and I don't understand JSON types"}`,
			errorValue: errors.New("my name is Athena and I don't understand JSON types"),
			expected: StatusResponse{
				StrVal:  "my name is Athena and I don't understand JSON types",
				IntVal:  0,
				IsValid: true,
				IsError: true,
			},
			returnsError: false,
		},
		{
			name:       "Success case - raw integer",
			input:      `42`,
			errorValue: nil,
			expected: StatusResponse{
				StrVal:  "42",
				IntVal:  42,
				IsValid: true,
				IsError: false,
			},
			returnsError: false,
		},
		{
			name:       "Success case - number as a string",
			errorValue: nil,
			input:      `{"status": "42"}`,
			expected: StatusResponse{
				IntVal:  42,
				IsError: false,
				IsValid: true,
				StrVal:  "42",
			},
			returnsError: false,
		},
		{
			name:       "Error case - bool value",
			input:      `{"status": true}`,
			errorValue: errors.New("StatusResponse is neither an int [0] or a string []"),
			expected: StatusResponse{
				StrVal:  "",
				IntVal:  0,
				IsValid: false,
				IsError: true,
			},
			returnsError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := &StatusResponse{}
			err := json.Unmarshal([]byte(tc.input), &response)

			if tc.returnsError {
				assert.Error(err)
				assert.False(response.IsValid)
				assert.Equal(tc.errorValue.Error(), response.GetError().Error())
			} else {
				assert.NoError(err)
				assert.Equal(tc.expected.IntVal, response.IntVal)
				assert.Equal(tc.expected.IsError, response.IsError)
				assert.Equal(tc.expected.IsValid, response.IsValid)
				assert.Equal(tc.expected.StrVal, response.StrVal)
			}
		})
	}
}
