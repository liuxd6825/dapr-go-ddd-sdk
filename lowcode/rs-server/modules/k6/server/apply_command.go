package server

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

func ApplyCommand(ctx context.Context, agg any, cmd ddd.Command, opts ...*ddd.ApplyCommandOptions) (err error) {
	return ddd.ApplyCommand(ctx, agg, cmd, opts...)
}
