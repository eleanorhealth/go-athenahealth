package athenahealth

import (
	"errors"
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

func newFormURLEncoder() *formURLEncoder {
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

	count := 0
	for _, key := range keys {
		for _, reader := range f.entries[key] {

			err := func() error {
				keyEscaped := url.QueryEscape(key)
				if count > 0 {
					_, err := w.Write([]byte("&"))
					if err != nil {
						return err
					}
				}
				_, err := w.Write([]byte(keyEscaped))
				if err != nil {
					return err
				}

				_, err = w.Write([]byte("="))
				if err != nil {
					return err
				}

				pr, pw := io.Pipe()
				errCh := make(chan error, 1)
				defer close(errCh)

				go func() {
					for {
						buf := make([]byte, defaultFormURLEncoderBufferSize)
						n, err := reader.Read(buf)
						if err != nil {
							if errors.Is(err, io.EOF) {
								errCh <- nil
								break
							}
							errCh <- err
							break
						}

						_, err = pw.Write([]byte(url.QueryEscape(string(buf[:n]))))
						if err != nil {
							errCh <- err
							break
						}
					}

					pw.Close()
				}()

				_, err = io.Copy(w, pr)
				if err != nil {
					return err
				}

				if err := <-errCh; err != nil {
					return err
				}

				count++

				return nil
			}()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
