package ratelimiter

import (
	"time"
)

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) Allowed(preview bool) (time.Duration, error) {
	return 0, nil
}
