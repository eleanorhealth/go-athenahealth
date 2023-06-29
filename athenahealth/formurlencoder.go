package athenahealth

import (
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sort"
	"strconv"
)

const (
	defaultFormURLEncoderBufferSize = 512
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

func (f *formURLEncoder) Encode(w io.Writer) error {
	keys := make([]string, 0, len(f.entries))
	for k := range f.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	isFirstEntry := true
	for _, key := range keys {
		for _, val := range f.entries[key] {

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
					pr, pw := io.Pipe()

					go func() {
						for {
							buf := make([]byte, defaultFormURLEncoderBufferSize)
							n, err := v.Read(buf)
							if err != nil {
								//nolint
								pw.CloseWithError(err)
								return
							}

							_, err = pw.Write([]byte(url.QueryEscape(string(buf[:n]))))
							if err != nil {
								//nolint
								pw.CloseWithError(err)
								return
							}
						}
					}()

					_, err = io.Copy(w, pr)
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
					return fmt.Errorf("invalid form url encoder value type '%s' for key %s", reflect.TypeOf(v).String(), key)
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
