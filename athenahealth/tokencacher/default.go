package tokencacher

import (
	"context"
	"sync"
	"time"
)

type Default struct {
	token     string
	expiresAt time.Time

	lock sync.Mutex
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) Get(ctx context.Context) (string, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if len(d.token) == 0 {
		return "", ErrTokenNotExist
	}

	if time.Now().After(d.expiresAt) {
		return "", ErrTokenExpired
	}

	return d.token, nil
}

func (d *Default) Set(ctx context.Context, token string, expiresAt time.Time) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.token = token
	d.expiresAt = expiresAt

	return nil
}
