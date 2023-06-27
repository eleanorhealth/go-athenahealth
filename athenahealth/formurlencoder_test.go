package athenahealth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"testing"
)

type errorReader struct {
	err error
}

func (er *errorReader) Read(p []byte) (int, error) {
	return 0, er.err
}

func Test_formURLEncoder_Encode(t *testing.T) {
	fileContents := make([]byte, 10000)
	rand.Read(fileContents)

	fileContentsBase64 := base64.StdEncoding.EncodeToString(fileContents)
	fileContentsBase64AndQueryEscaped := url.QueryEscape(fileContentsBase64)

	errBadRead := errors.New("bad read")

	tests := []struct {
		name    string
		items   map[string][]io.Reader
		wantW   string
		wantErr error
	}{
		{
			name:    "empty",
			items:   map[string][]io.Reader{},
			wantW:   "",
			wantErr: nil,
		},
		{
			name: "multiple readers",
			items: map[string][]io.Reader{
				// output should be sorted by key
				"str2": {
					strings.NewReader("#$%#%^"),
				},
				"str1": {
					strings.NewReader("test @#$"),
					strings.NewReader("!@#%  ()*"),
				},
			},
			wantW:   "str1=test+%40%23%24&str1=%21%40%23%25++%28%29%2A&str2=%23%24%25%23%25%5E",
			wantErr: nil,
		},
		{
			name: "document",
			items: map[string][]io.Reader{
				"document": {
					newBase64Reader(bytes.NewReader(fileContents)),
				},
			},
			wantW:   fmt.Sprintf("document=%s", fileContentsBase64AndQueryEscaped),
			wantErr: nil,
		},
		{
			name: "error",
			items: map[string][]io.Reader{
				"error": {
					&errorReader{errBadRead},
				},
				"str1": {
					strings.NewReader("test @#$"),
					strings.NewReader("!@#%  ()*"),
				},
				"str2": {
					strings.NewReader("#$%#%^"),
				},
			},
			wantW:   "error=",
			wantErr: errBadRead,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &formURLEncoder{
				entries: tt.items,
			}
			w := &bytes.Buffer{}
			if err := f.Encode(w); !errors.Is(err, tt.wantErr) {
				t.Errorf("formURLEncoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("formURLEncoder.Encode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
