package server

import (
	"context"
	"time"
)

type ContextPkg struct {
}

func NewContextPkg() *ContextPkg {
	return &ContextPkg{}
}

func (c *ContextPkg) Background() context.Context {
	return context.Background()
}

func (c *ContextPkg) WithValue(parent context.Context, key, val any) context.Context {
	return context.WithValue(parent, key, val)
}

func (c *ContextPkg) WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

func (c *ContextPkg) WithCancel(parent context.Context) (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(parent)
}
