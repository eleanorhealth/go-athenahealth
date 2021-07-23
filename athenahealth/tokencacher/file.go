package tokencacher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type File struct {
	path string

	lock sync.Mutex
}

type fileCache struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewFile(path string) *File {
	if len(path) == 0 {
		panic("path required")
	}

	return &File{
		path: path,
	}
}

func (f *File) Get(ctx context.Context) (string, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	_, err := os.Stat(f.path)
	if os.IsNotExist(err) {
		err = ioutil.WriteFile(f.path, nil, 0600)
		if err != nil {
			return "", err
		}
	}

	contents, err := ioutil.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	if len(contents) == 0 {
		return "", ErrTokenNotExist
	}

	c := &fileCache{}
	err = json.Unmarshal(contents, c)
	if err != nil {
		return "", fmt.Errorf("Error unmarshaling token: %s", err)
	}

	if time.Now().After(c.ExpiresAt) {
		return "", ErrTokenExpired
	}

	return c.Token, nil
}

func (f *File) Set(ctx context.Context, token string, expiresAt time.Time) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	c := &fileCache{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.path, b, 0600)
	if err != nil {
		return err
	}

	return nil
}
