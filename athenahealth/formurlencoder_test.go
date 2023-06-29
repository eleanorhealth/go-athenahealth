package athenahealth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errorReader struct {
	err error
}

func (er *errorReader) Read(p []byte) (int, error) {
	return 0, er.err
}

func Test_formURLEncoder_Encode_table(t *testing.T) {
	fileContents := make([]byte, 10000)
	rand.Read(fileContents)

	fileContentsBase64 := base64.StdEncoding.EncodeToString(fileContents)
	fileContentsBase64AndQueryEscaped := url.QueryEscape(fileContentsBase64)

	errBadRead := errors.New("bad read")

	tests := []struct {
		name    string
		fue     *formURLEncoder
		wantW   string
		wantErr error
	}{
		{
			name:    "empty",
			fue:     func() *formURLEncoder { return NewFormURLEncoder() }(),
			wantW:   "",
			wantErr: nil,
		},
		{
			name: "multiple readers",
			fue: func() *formURLEncoder {
				fue := NewFormURLEncoder()
				fue.AddString("str2", "#$%#%^")
				fue.AddString("str1", "test @#$")
				fue.AddString("str1", "!@#%  ()*")
				return fue
			}(),
			wantW:   "str1=test+%40%23%24&str1=%21%40%23%25++%28%29%2A&str2=%23%24%25%23%25%5E",
			wantErr: nil,
		},
		{
			name: "document",
			fue: func() *formURLEncoder {
				fue := NewFormURLEncoder()
				fue.AddReader("document", newBase64Reader(bytes.NewReader(fileContents)))
				return fue
			}(),
			wantW:   fmt.Sprintf("document=%s", fileContentsBase64AndQueryEscaped),
			wantErr: nil,
		},
		{
			name: "error",
			fue: func() *formURLEncoder {
				fue := NewFormURLEncoder()
				fue.AddReader("error", &errorReader{errBadRead})
				fue.AddString("str1", "test @#$")
				fue.AddString("str1", "!@#%  ()*")
				fue.AddString("str2", "#$%#%^")
				return fue
			}(),
			wantW:   "error=",
			wantErr: errBadRead,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.fue.Encode(w); !errors.Is(err, tt.wantErr) {
				t.Errorf("formURLEncoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("formURLEncoder.Encode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_formURLEncoder_Encode(t *testing.T) {
	assert := assert.New(t)

	docBytes := make([]byte, 50)
	_, err := rand.Read(docBytes)
	assert.NoError(err)

	fue := NewFormURLEncoder()
	fue.AddReader("doc", bytes.NewReader(docBytes))
	fue.AddInt("count", 10)
	fue.AddString("str!", "hello world!")

	b := bytes.NewBuffer(nil)
	err = fue.Encode(b)
	assert.NoError(err)
	assert.Equal(fmt.Sprintf("count=10&doc=%s&%s=%s", url.QueryEscape(string(docBytes)), url.QueryEscape("str!"), url.QueryEscape("hello world!")), b.String())
}
