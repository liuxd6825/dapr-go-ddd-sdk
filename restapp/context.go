package restapp

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type serverContext struct {
	ctx iris.Context
}

func NewContext(ictx iris.Context) context.Context {
	loggerCtx := logs.NewContext(ictx, _logger)
	metadata := make(map[string]string, 0)
	serverCtx := newIrisContext(ictx)
	if ictx != nil {
		header := ictx.Request().Header
		for k, v := range header {
			metadata[k] = v[0]
		}
	}
	return ddd_context.NewContext(loggerCtx, metadata, serverCtx)
}

func NewLoggerContext(ctx context.Context) context.Context {
	return logs.NewContext(ctx, _logger)
}

func NewContext2(ctx context.Context) context.Context {
	return logs.NewContext(ctx, _logger)
}

func newIrisContext(ctx iris.Context) ddd_context.ServerContext {
	return &serverContext{
		ctx: ctx,
	}
}

func GetIrisContext(ctx context.Context) iris.Context {
	v := ddd_context.GetServerContext(ctx)
	if s, ok := v.(*serverContext); ok {
		return s.ctx
	}
	return nil
}

func (s *serverContext) SetResponseHeader(name string, value string) {
	s.ctx.Header(name, value)
}

func (s *serverContext) URLParamDefault(name, def string) string {
	return s.ctx.URLParamDefault(name, def)
}
