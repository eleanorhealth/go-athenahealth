package ratelimiter

import (
	"context"
	"time"
)

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) Allowed(ctx context.Context, preview bool) (time.Duration, error) {
	return 0, nil
}
