package restapp

import (
	"context"
	"strings"

	"github.com/kataras/iris/v12"
	appctx2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
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
	TestToken        = `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJEdXhtLUp3dC1Ub2tlbiIsImV4cCI6MTcwMzc1Mjk3NCwidXNlciI6eyJ0ZW5hbnROYW1lIjoidGVzdCIsIm5hbWUiOiJ0ZXN0IiwidGVuYW50SWQiOiJ0ZXN0IiwidGVuYW50QWNjb3VudCI6InRlc3QiLCJpZCI6IjE3MjY0Nzk0NDEyNTUwNTEyNjQiLCJ1c2VyVHlwZSI6IlRFTkFOVF9BRE1JTiIsImFjY291bnQiOiJ0ZXN0Iiwic3RhdHVzIjoiVVNFSU5HIn0sImNsaWVudF9pZCI6IjA5OGY2YmNkNDYyMWQzNzNjYWRlNGU4MzI2MjdiNGY2In0.s_kHa3pKt6XehbsL7E9PJqywM_pxbbq6V2zHyZCJmDk`
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
		tenantId := "test"
		opt.CheckAuth = gp.PBool(false)
		opt.TenantId = &tenantId
	})
}

// NewTestContext
//
//	@Description:
//	@param ctx
//	@return newCtx
//	@return err
func NewTestContext(ctx context.Context, opts ...ContextOptions) (newCtx context.Context, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	var pCtx context.Context = ctx

	// 添加 日志 上下文
	newCtx = logs.NewContext(pCtx)

	opt := newContextOption(opts...)
	for _, fun := range opts {
		fun(opt)
	}

	//添加 租户 上下文
	if opt.TenantId != nil {
		newCtx = appctx2.NewTenantContext(newCtx, gp.String(opt.TenantId, ""))
	}

	newCtx, err = appctx2.NewAuthContext(newCtx, TestToken)
	if err != nil {
		return nil, err
	}

	return newCtx, err
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
		newCtx = appctx2.NewServerContext(newCtx, &irisServer{ictx})
		newCtx = appctx2.NewIrisContext(newCtx, ictx)
	}

	//添加 租户 上下文
	if opt.TenantId != nil {
		newCtx = appctx2.NewTenantContext(newCtx, gp.String(opt.TenantId, ""))
	}

	// 添加 Header 上下文
	header := newHeader(ictx)
	newCtx = appctx2.NewHeaderContext(newCtx, header)

	//添加 用户认证 上下文
	newCtx, _, err = NewAuthTokenContext(newCtx, header[Authorization], gp.Bool(opt.CheckAuth))
	if err != nil {
		return nil, err
	}

	return newCtx, err
}

func newHeader(ictx iris.Context) appctx2.Header {
	var header appctx2.Header
	if ictx != nil {
		header = appctx2.Header(ictx.Request().Header)
	}
	if header == nil {
		header = appctx2.Header{}
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
	newCtx, err = appctx2.NewAuthContext(parent, token)
	return newCtx, token, err
}

func (i *irisServer) SetResponseHeader(key string, value string) {
	i.ctx.ResponseWriter().Header().Set(key, value)
}

func (i *irisServer) URLParamDefault(name, def string) string {
	return i.ctx.URLParamDefault(name, def)
}
