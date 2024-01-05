package restapp

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/auth"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type serverContext struct {
	ctx iris.Context
}

type tenantCtxKey struct {
}

type ContextOptions struct {
	checkAuth *bool
	tenantId  *string
}

const (
	HttpHeadAuthorization = "Authorization"
)

var (
	DefaultAuthToken = ""
)

func NewContextOptions(opts ...*ContextOptions) *ContextOptions {
	o := &ContextOptions{}
	for _, item := range opts {
		if item.checkAuth != nil {
			o.checkAuth = item.checkAuth
		}
		if item.tenantId != nil {
			o.tenantId = item.tenantId
		}
	}
	return o
}

// NewContextNoAuth
//
//	@Description:
//	@param ictx
//	@return newCtx
//	@return err
func NewContextNoAuth(ictx iris.Context) (newCtx context.Context, err error) {
	return NewContext(ictx, NewContextOptions().SetCheckAuth(false))
}

func NewContextTenantId(ictx iris.Context, tenantId string) (newCtx context.Context, err error) {
	return NewContext(ictx, NewContextOptions().SetTenantId(tenantId))
}

func NewContext(ictx iris.Context, opts ...*ContextOptions) (newCtx context.Context, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	var pCtx context.Context = ictx
	if ictx == nil {
		pCtx = context.Background()
	}
	opt := NewContextOptions(opts...)
	newCtx = logs.NewContext(pCtx)

	metadata := make(map[string]string)
	if ictx != nil {
		header := ictx.Request().Header
		for k, v := range header {
			metadata[k] = v[0]
		}
	}

	if opt.tenantId != nil {
		newCtx = newTenantContext(newCtx, opt.TenantId())
	}

	token := metadata[HttpHeadAuthorization]
	newCtx, authToken, err := newAuthTokenContext(newCtx, token, opt.CheckAuth())
	if err != nil {
		return nil, err
	}
	if token != authToken {
		metadata[HttpHeadAuthorization] = authToken
	}

	if ictx != nil {
		serverCtx := newIrisContext(ictx)
		newCtx = ddd_context.NewContext(newCtx, metadata, serverCtx)
	}

	return newCtx, err
}

func newAuthTokenContext(parent context.Context, token string, checkAuth bool) (newCtx context.Context, authToken string, err error) {
	newCtx = parent
	if token == "" && DefaultAuthToken != "" {
		token = DefaultAuthToken
	}
	if token == "" {
		if checkAuth {
			return nil, token, errors.New("Authorization is null")
		} else {
			return parent, "", nil
		}
	}
	newCtx, err = auth.NewContext(parent, token)
	return newCtx, token, err
}

func newIrisContext(ctx iris.Context) ddd_context.ServerContext {
	return &serverContext{
		ctx: ctx,
	}
}

func newTenantContext(parent context.Context, tenantId string) context.Context {
	return context.WithValue(parent, tenantCtxKey{}, tenantId)
}

// GetTenantId
//
//	@Description: 根据上下文取得租户ID
//	@param ctx
//	@return string
//	@return bool
func GetTenantId(ctx context.Context) (string, bool) {
	val := ctx.Value(tenantCtxKey{})
	if val == nil {
		return "", false
	}
	tenantId, ok := val.(string)
	return tenantId, ok
}

func NewLoggerContext(ctx context.Context) context.Context {
	return logs.NewContext(ctx)
}

// GetIrisContext
//
//	@Description: 获取Iris上下文
//	@param ctx
//	@return iris.Context
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

func (o *ContextOptions) TenantId() string {
	if o.tenantId != nil {
		return *o.tenantId
	}
	return ""
}

func (o *ContextOptions) SetTenantId(val string) *ContextOptions {
	o.tenantId = &val
	return o
}
