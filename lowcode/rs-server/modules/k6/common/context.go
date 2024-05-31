package common

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
func (c *ContextPkg) TODO() context.Context {
	return context.TODO()
}

type WithCancelResult struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func (c *ContextPkg) WithCancel(parent context.Context) *WithCancelResult {
	ctx, cancel := context.WithCancel(parent)
	return &WithCancelResult{Ctx: ctx, Cancel: cancel}
}

type WithTimeoutResult struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func (c *ContextPkg) WithTimeout(parent context.Context, timeout time.Duration) *WithTimeoutResult {
	ctx, cancel := context.WithTimeout(parent, timeout)
	return &WithTimeoutResult{Ctx: ctx, Cancel: cancel}
}
