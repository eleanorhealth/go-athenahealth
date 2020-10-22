package stats

import "context"

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) IncrRequests(ctx context.Context) error {
	return nil
}
