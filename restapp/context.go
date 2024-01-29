package restapp

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/goplus"
	"strings"
)

type irisServer struct {
	ctx iris.Context
}

type ContextOption struct {
	CheckAuth *bool
	TenantId  *string
}
type ContextOptions func(option *ContextOption)

const (
	Authorization = "Authorization"
)

var (
	DefaultAuthToken = ""
)

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
	return NewContext(ictx, func(opt *ContextOption) {
		opt.CheckAuth = goplus.PBool(false)
		opt.TenantId = nil
	})
}

func newContextOption(opts ...ContextOptions) *ContextOption {
	opt := &ContextOption{}
	for _, item := range opts {
		if item != nil {
			item(opt)
		}
	}
	return opt
}

func NewContext(ictx iris.Context, opts ...ContextOptions) (newCtx context.Context, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	var pCtx context.Context = ictx
	if ictx == nil {
		pCtx = context.Background()
	}

	opt := newContextOption(opts...)
	for _, fun := range opts {
		fun(opt)
	}

	// 添加 日志 上下文
	newCtx = logs.NewContext(pCtx)

	if ictx != nil {
		// 添加 ServerHeader 上下文
		newCtx = appctx.NewServerContext(newCtx, &irisServer{ictx})
		newCtx = appctx.NewIrisContext(newCtx, ictx)
	}

	//添加 租户 上下文
	if opt.TenantId != nil {
		newCtx = appctx.NewTenantContext(newCtx, goplus.String(opt.TenantId, ""))
	}

	// 添加 Header 上下文
	header := newHeader(ictx)
	newCtx = appctx.NewHeaderContext(newCtx, header)

	//添加 用户认证 上下文
	newCtx, _, err = NewAuthTokenContext(newCtx, header[Authorization], goplus.Bool(opt.CheckAuth))
	if err != nil {
		return nil, err
	}

	return newCtx, err
}

func newHeader(ictx iris.Context) appctx.Header {
	var header appctx.Header
	if ictx != nil {
		header = appctx.Header(ictx.Request().Header)
	}
	if header == nil {
		header = appctx.Header{}
	}

	isHave := false
	if val, ok := header[Authorization]; ok {
		for _, s := range val {
			if strings.Trim(s, " ") != "" {
				isHave = true
				break
			}
		}
	}
	if !isHave && DefaultAuthToken != "" {
		header[Authorization] = []string{DefaultAuthToken}
	}

	return header
}

func NewAuthTokenContext(parent context.Context, headerValues []string, checkAuth bool) (newCtx context.Context, authToken string, err error) {
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
			return nil, token, errors.New("Header Authorization is null")
		} else {
			return parent, "", nil
		}
	}
	newCtx, err = appctx.NewAuthContext(parent, token)
	return newCtx, token, err
}

func (i *irisServer) SetResponseHeader(key string, value string) {
	i.ctx.ResponseWriter().Header().Set(key, value)
}

func (i *irisServer) URLParamDefault(name, def string) string {
	return i.ctx.URLParamDefault(name, def)
}
