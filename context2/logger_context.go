package context2

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

func NewLoggerContext(ctx context.Context) context.Context {
	return logs.NewContext(ctx)
}
