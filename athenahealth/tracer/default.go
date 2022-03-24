package tracer

import "context"

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) Before(ctx context.Context, method, path string) context.Context {
	return ctx
}

func (d *Default) After(ctx context.Context) {
}
