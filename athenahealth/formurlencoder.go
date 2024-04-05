package athenahealth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sort"
	"strconv"
)

type formURLEncoder struct {
	entries map[string][]any
}

func NewFormURLEncoder() *formURLEncoder {
	return &formURLEncoder{
		entries: make(map[string][]any),
	}
}

func (f *formURLEncoder) AddString(key string, value string) {
	f.entries[key] = append(f.entries[key], value)
}

func (f *formURLEncoder) AddInt(key string, value int) {
	f.entries[key] = append(f.entries[key], value)
}

func (f *formURLEncoder) AddReader(key string, value io.Reader) {
	f.entries[key] = append(f.entries[key], value)
}

// Encode encodes the values into “URL encoded” form
// ("bar=baz&foo=quux") sorted by key.
func (f *formURLEncoder) Encode(ctx context.Context, w io.Writer) error {
	keys := make([]string, 0, len(f.entries))
	for k := range f.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	isFirstEntry := true
	for _, key := range keys {
		for valIdx, val := range f.entries[key] {

			err := func() error {
				if isFirstEntry {
					isFirstEntry = false
				} else {
					_, err := w.Write([]byte("&"))
					if err != nil {
						return err
					}
				}

				keyEscaped := url.QueryEscape(key)
				_, err := w.Write([]byte(keyEscaped))
				if err != nil {
					return err
				}

				_, err = w.Write([]byte("="))
				if err != nil {
					return err
				}

				switch v := val.(type) {
				case io.Reader:
					encoder := base64.NewEncoder(base64.StdEncoding, &urlQueryEscapeWriter{w})

					err = copyCtx(ctx, encoder, v)
					if err != nil {
						return err
					}

					err = encoder.Close()
					if err != nil {
						return err
					}

				case string:
					_, err = w.Write([]byte(url.QueryEscape(v)))
					if err != nil {
						return err
					}

				case int:
					_, err = w.Write([]byte(url.QueryEscape(strconv.Itoa(v))))
					if err != nil {
						return err
					}

				default:
					return fmt.Errorf("invalid form url encoder value type '%s' for key %s[%d]", reflect.TypeOf(v).String(), key, valIdx)
				}

				return nil
			}()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

func copyCtx(ctx context.Context, dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, readerFunc(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return src.Read(p)
		}
	}))
	return err
}

type urlQueryEscapeWriter struct {
	io.Writer
}

func (w *urlQueryEscapeWriter) Write(p []byte) (int, error) {
	escaped := url.QueryEscape(string(p))

	return w.Writer.Write([]byte(escaped))
}
