package ctx_pkg

import (
	"context"
	"time"
)

type CtxPkg struct {
}

func New() *CtxPkg {
	return &CtxPkg{}
}

func (c *CtxPkg) Background() context.Context {
	return context.Background()
}

func (c *CtxPkg) WithValue(parent context.Context, key, val any) context.Context {
	return context.WithValue(parent, key, val)
}

func (c *CtxPkg) WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

func (c *CtxPkg) WithCancel(parent context.Context) (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(parent)
}
