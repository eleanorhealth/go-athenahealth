package athenahealth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type errorReader struct {
	err error
}

func (er *errorReader) Read(p []byte) (int, error) {
	return 0, er.err
}

type slowReader struct {
	r               io.Reader
	maxBytesPerRead int64
	sleepPerRead    time.Duration
}

func (sr *slowReader) Read(p []byte) (n int, err error) {
	if int64(len(p)) > sr.maxBytesPerRead {
		p = p[0:sr.maxBytesPerRead]
	}
	n, err = sr.r.Read(p)
	time.Sleep(sr.sleepPerRead)
	return
}

func Test_formURLEncoder_Encode_table(t *testing.T) {
	fileContents := make([]byte, 10000)
	rand.Read(fileContents)

	fileContentsBase64 := url.QueryEscape(base64.StdEncoding.EncodeToString(fileContents))

	errBadRead := errors.New("bad read")

	tests := []struct {
		name    string
		fue     *formURLEncoder
		wantW   string
		wantErr error
		ctxFn   func() (context.Context, func())
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
				fue.AddReader("document", bytes.NewReader(fileContents))
				return fue
			}(),
			wantW:   fmt.Sprintf("document=%s", fileContentsBase64),
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
		{
			name: "context cancellation",
			fue: func() *formURLEncoder {
				fue := NewFormURLEncoder()
				fue.AddReader("file", &slowReader{
					r:               bytes.NewReader(fileContents),
					maxBytesPerRead: 1,
					sleepPerRead:    time.Second,
				})
				return fue
			}(),
			wantW:   "file=",
			wantErr: context.Canceled,
			ctxFn: func() (context.Context, func()) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
				cancel()
				return ctx, cancel
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.ctxFn != nil {
				var cancel func()
				ctx, cancel = tt.ctxFn()
				defer cancel()
			}
			w := &bytes.Buffer{}
			if err := tt.fue.Encode(ctx, w); !errors.Is(err, tt.wantErr) {
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
	err = fue.Encode(context.Background(), b)
	assert.NoError(err)
	assert.Equal(fmt.Sprintf("count=10&doc=%s&%s=%s", url.QueryEscape(base64.StdEncoding.EncodeToString(docBytes)), url.QueryEscape("str!"), url.QueryEscape("hello world!")), b.String())
}
