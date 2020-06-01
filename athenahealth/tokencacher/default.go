package tokencacher

import (
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

func (d *Default) Get() (string, error) {
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

func (d *Default) Set(token string, expiresAt time.Time) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.token = token
	d.expiresAt = expiresAt

	return nil
}
