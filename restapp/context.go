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
	metadata := make(map[string]string, 0)
	header := ictx.Request().Header
	for k, v := range header {
		metadata[k] = v[0]
	}
	serverCtx := newServerContext(ictx)
	loggerCtx := logs.NewContext(ictx, _logger)
	return ddd_context.NewContext(loggerCtx, metadata, serverCtx)
}

func newServerContext(ctx iris.Context) ddd_context.ServerContext {
	return &serverContext{
		ctx: ctx,
	}
}

func (s *serverContext) SetResponseHeader(name string, value string) {
	s.ctx.Header(name, value)
}

func (s *serverContext) URLParamDefault(name, def string) string {
	return s.ctx.URLParamDefault(name, def)
}
