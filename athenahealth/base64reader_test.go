package athenahealth

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"math/rand"
	"testing"
)

func TestBase64Reader_Read(t *testing.T) {

	chars := "abcdefghijklmnopqrstuvwxyz 1234567890"
	maxBufChars := ""
	for i := 0; i < 511; i++ {
		maxBufChars += string(chars[rand.Intn(len(chars))])
	}
	encodedMaxBufChars := base64.StdEncoding.EncodeToString([]byte(maxBufChars))

	testCases := []struct {
		name        string
		input       []byte
		expected    []byte
		expectedErr error
	}{
		{
			name:        "Empty input",
			input:       []byte{},
			expected:    []byte{},
			expectedErr: nil,
		},
		{
			name:        "Input encoded",
			input:       []byte("Hello World! This is a long string."),
			expected:    []byte("SGVsbG8gV29ybGQhIFRoaXMgaXMgYSBsb25nIHN0cmluZy4="),
			expectedErr: nil,
		},
		{
			name:        "max",
			input:       []byte(maxBufChars),
			expected:    []byte(encodedMaxBufChars),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewReader(tc.input)
			base64Reader := newBase64Reader(reader)

			encodedBytes, err := io.ReadAll(base64Reader)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error '%v', but got '%v'", tc.expectedErr, err)
			}

			encodedData := string(encodedBytes)
			if encodedData != string(tc.expected) {
				t.Errorf("Expected encoded data '%s', but got '%s'", tc.expected, encodedData)
			}

			decodedBytes, err := base64.StdEncoding.DecodeString(string(encodedBytes))
			if err != nil {
				t.Errorf("decoding encoded bytes")
			}

			if !bytes.Equal(decodedBytes, tc.input) {
				t.Errorf("bytes not equal")
			}
		})
	}
}
