package restapp

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type irisServer struct {
	ctx iris.Context
}

type ContextOptions struct {
	checkAuth *bool
	tenantId  *string
}

const (
	Authorization = "Authorization"
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

func NewLoggerContext(ctx context.Context) context.Context {
	return logs.NewContext(ctx)
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

func NewContext(ictx iris.Context, opts ...*ContextOptions) (newCtx context.Context, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	var pCtx context.Context = ictx
	if ictx == nil {
		pCtx = context.Background()
	}
	opt := NewContextOptions(opts...)

	// 添加 日志 上下文
	newCtx = logs.NewContext(pCtx)

	header := make(map[string][]string)
	if ictx != nil {
		for k, v := range ictx.Request().Header {
			header[k] = v
		}
		// 添加 Header 上下文
		newCtx = appctx.NewHeaderContext(newCtx, header)
		// 添加 ServerHeader 上下文
		newCtx = appctx.NewServerContext(newCtx, &irisServer{ictx})
		// 添加 iris.context 上下文
		newCtx = appctx.NewIrisContext(newCtx, ictx)
	}

	//添加 租户 上下文
	if opt.tenantId != nil {
		newCtx = appctx.NewTenantContext(newCtx, opt.TenantId())
	}

	//添加 用户认证 上下文
	newCtx, _, err = newAuthTokenContext(newCtx, header[Authorization], opt.CheckAuth())
	if err != nil {
		return nil, err
	}

	return newCtx, err
}

func newAuthTokenContext(parent context.Context, headerValues []string, checkAuth bool) (newCtx context.Context, authToken string, err error) {
	token := ""
	for _, v := range headerValues {
		if len(v) > 0 {
			token = v
		}
	}

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
	newCtx, err = appctx.NewAuthContext(parent, token)
	return newCtx, token, err
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

func (i *irisServer) SetResponseHeader(key string, value string) {
	i.ctx.ResponseWriter().Header().Set(key, value)
}

func (i *irisServer) URLParamDefault(name, def string) string {
	return i.ctx.URLParamDefault(name, def)
}
