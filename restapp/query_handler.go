package restapp

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

type BaseQueryHandler struct {
}

func NewBaseQueryHandler() *BaseQueryHandler {
	return &BaseQueryHandler{}
}

func (h *BaseQueryHandler) DoSession(ctx context.Context, structNameFunc func() string, event ddd.Event, fun func(ctx context.Context) error) error {
	return applog.DoEventLog(ctx, structNameFunc, event, applog.RunFuncName(1), func() error {
		return ddd_repository.StartSession(ctx, NewSession(), func(ctx context.Context) error {
			return fun(ctx)
		})
	})
}

func (h *BaseQueryHandler) Do(ctx context.Context, structNameFunc func() string, event ddd.Event, fun func(ctx context.Context) error) error {
	return applog.DoEventLog(ctx, structNameFunc, event, applog.RunFuncName(1), func() error {
		return fun(ctx)
	})
}
