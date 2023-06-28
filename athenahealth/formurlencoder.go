package athenahealth

import (
	"io"
	"net/url"
	"sort"
)

const (
	defaultFormURLEncoderBufferSize = 512
)

type formURLEncoder struct {
	entries map[string][]io.Reader
}

func NewFormURLEncoder() *formURLEncoder {
	return &formURLEncoder{
		entries: make(map[string][]io.Reader),
	}
}

func (f *formURLEncoder) Add(key string, value io.Reader) {
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
		for _, reader := range f.entries[key] {

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

				pr, pw := io.Pipe()

				go func() {
					for {
						buf := make([]byte, defaultFormURLEncoderBufferSize)
						n, err := reader.Read(buf)
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

				return nil
			}()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
