package context2

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/auth"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

var JwtKey = "#@!{[duXm-serVice-t0ken]},.(10086)$!"

type serverContext struct {
	ctx iris.Context
}

func NewContext(ictx iris.Context, tenantId string, opts ...Options) (newCtx context.Context, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	opt := NewOptions(opts...)
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
	// 日志上下文
	newCtx = NewLoggerContext(newCtx)
	// 租房上下文
	newCtx = NewTenantContext(newCtx, tenantId)

	return newCtx, err
}

func newIrisContext(ctx iris.Context) ddd_context.ServerContext {
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
