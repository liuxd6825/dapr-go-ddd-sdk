package restapp

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/auth"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

var JwtKey = "#@!{[duXm-serVice-t0ken]},.(10086)$!"

type serverContext struct {
	ctx iris.Context
}

type ContextOptions struct {
	checkAuth *bool
}

func NewContextOptions(opts ...ContextOptions) *ContextOptions {
	o := &ContextOptions{}
	for _, item := range opts {
		if item.checkAuth != nil {
			o.checkAuth = item.checkAuth
		}
	}
	return o
}

func (o *ContextOptions) CheckAuth() bool {
	if o.checkAuth != nil {
		return *o.checkAuth
	}
	return true
}

func (o *ContextOptions) SetCheckAuth(val bool) *ContextOptions {
	o.checkAuth = &val
	return o
}

func NewContext(ictx iris.Context, opts ...ContextOptions) (newCtx context.Context, err error) {
	opt := NewContextOptions(opts...)
	newCtx = logs.NewContext(ictx, _logger)
	metadata := make(map[string]string)
	serverCtx := newIrisContext(ictx)
	if ictx != nil {
		header := ictx.Request().Header
		for k, v := range header {
			metadata[k] = v[0]
		}
	}
	newCtx = ddd_context.NewContext(newCtx, metadata, serverCtx)

	jwt, jwtOk := metadata["Authorization"]
	if !jwtOk && opt.CheckAuth() {
		return nil, errors.New("Authorization is null")
	}

	if jwtOk {
		newCtx, err = auth.NewContext(newCtx, jwt, JwtKey)
	}

	return newCtx, err
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
