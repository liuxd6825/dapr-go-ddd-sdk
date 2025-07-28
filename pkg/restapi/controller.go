package restapi

import (
	context2 "context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"reflect"
	"strings"
)

type ApiController struct {
	app      *iris.Application
	rootPath string
	routes   []*router.Route
	ctl      any
}

type ApiFunc func(ctx context2.Context, ictx *context.Context, params any) (any, error)

type Controller interface {
	InitController(app *iris.Application) error
}

type CallOptions struct {
	InitMethod func(callMethod *CallMethod)
	Before     func(ictx *context.Context, params any) (any, context2.Context, error)
	After      func(ctx context2.Context, ictx *context.Context, resData any, err error) (any, error)
}

func InitController(app *iris.Application, controller Controller) {
	if err := controller.InitController(app); err != nil {
		panic(err)
	}
}

func NewController(app *iris.Application, rootPath string, ctl any) *ApiController {
	return &ApiController{app: app, rootPath: rootPath, ctl: ctl, routes: make([]*router.Route, 0)}
}

func (c *ApiController) getPath(path string) string {
	root := c.rootPath
	if ok := strings.HasPrefix(path, "/"); ok {
		path = path[1:]
	}
	if ok := strings.HasSuffix(c.rootPath, "/"); ok {
		root = root[:len(root)-1]
	}
	return root + "/" + path
}

func (c *ApiController) addRouter(router *router.Route) {
	c.routes = append(c.routes, router)
}

func (c *ApiController) newCallMethod(handlerName string) (callMethod *CallMethod, err error) {
	if c.ctl == nil {
		return nil, errors.New("controller is nil")
	}
	method, err := NewCallMethod(c.ctl, handlerName)
	if err != nil {
		return nil, errors.New("get api func error: %s ", err.Error())
	}
	return method, nil
}

func (c *ApiController) GetOne(path string, handlerName string, opts ...CallOptions) *router.Route {
	opts = append(opts, CallOptions{
		After: func(ctx context2.Context, ictx *context.Context, data any, err error) (any, error) {
			if err != nil {
				return nil, err
			}
			if reflectutils.IsNil(data) {
				ictx.StatusCode(404)
				return nil, NotFoundError()
			}
			return data, nil
		},
	})
	r := c.call(iris.MethodGet, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) GetData(path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.call(iris.MethodGet, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) GetPaging(path string, handlerName string, opts ...CallOptions) *router.Route {
	opts = append(opts, CallOptions{
		InitMethod: func(method *CallMethod) {
			method.CloseInParams = true
		},
		Before: func(ictx *context.Context, param any) (any, context2.Context, error) {
			request, err := GetFindPagingRequest(ictx)
			if err != nil {
				return nil, nil, err
			}
			ctx := appctx.NewWebContext(context2.Background(), ictx)
			return request, ctx, nil
		},
	})
	r := c.call(iris.MethodGet, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) GetList(path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.call(iris.MethodGet, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) Delete(path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.call(iris.MethodDelete, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) Put(path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.call(iris.MethodPut, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) Post(path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.call(iris.MethodPost, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) EventHandle(path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.callEventHandle(path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) Handle(method string, path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.call(method, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) call(method string, path string, handlerName string, opts ...CallOptions) *router.Route {
	return c.callMethod(method, path, handlerName, false, opts...)
}

func (c *ApiController) callEventHandle(path string, handlerName string, opts ...CallOptions) *router.Route {
	return c.callMethod(iris.MethodPost, path, handlerName, true, opts...)
}

func (c *ApiController) callMethod(method string, path string, handlerName string, isEventHandle bool, opts ...CallOptions) *router.Route {
	callMethod, err := c.newCallMethod(handlerName)
	if err != nil {
		ctlType := reflect.TypeOf(c.ctl)
		if ctlType.Kind() == reflect.Ptr {
			ctlType = ctlType.Elem()
		}
		pkgPath := ctlType.PkgPath()
		typeName := ctlType.Name()
		err = errors.New("%s %s.%s() func error: %s ", pkgPath, typeName, handlerName, err.Error())
		panic(err)
	}
	path = c.getPath(path)
	r := c.app.Handle(method, path, func(ictx *context.Context) {
		gp.Try(func() error {
			params, ctx, err := c.getParams(ictx, callMethod, isEventHandle, opts...)
			if err != nil {
				return err
			}
			data, err := callMethod.Call(ctx, ictx, params)
			for _, opt := range opts {
				if opt.After != nil {
					data, err = opt.After(ctx, ictx, data, err)
				}
			}
			if err == nil && callMethod.OutData >= 0 {
				err = SetOKData(ictx, data)
			}
			return err
		}).Catch(func(err error) {
			SetError(ictx, err)
		})
	})
	c.addRouter(r)
	return r
}

func (c *ApiController) getParams(ictx *context.Context, callMethod *CallMethod, isEventHandle bool, opts ...CallOptions) (params any, rctx context2.Context, err error) {
	// CallMethod初始化
	for _, opt := range opts {
		if opt.InitMethod != nil {
			opt.InitMethod(callMethod)
		}
	}

	if callMethod.InParams >= 0 && !callMethod.CloseInParams {
		inParamsType := callMethod.Method.Type().In(callMethod.InParams)
		params, rctx, err = c.GetParams(ictx, inParamsType, isEventHandle)
		if err != nil {
			return nil, nil, err
		}
	}
	//logs.Info(ctx, logs.Fields{"method:": ictx.Method(), "url:": ictx.Request().URL, "param": params})
	// 后期处理
	for _, opt := range opts {
		if opt.Before != nil {
			params, rctx, err = opt.Before(ictx, params)
			break
		}
	}
	if err == nil && rctx == nil {
		rctx = appctx.NewWebContext(context2.Background(), ictx)
	}
	return params, rctx, err
}

// GetParams 获取参数
// ictx: iris请求上下文
// paramType: 参数类型
// isEventHandle: 是否是事件处理
func (c *ApiController) GetParams(ictx *context.Context, paramType reflect.Type, isEventHandle bool) (params any, rctx context2.Context, err error) {
	if isEventHandle {
		paramsValue, err := reflectutils.New(paramType)
		if err != nil {
			return nil, nil, errors.New("create params error: %s ", err.Error())
		}
		params = paramsValue.Interface()
		params, rctx, err = GetEventParams(ictx, params)
		return params, rctx, err
	}
	if paramType != nil {
		rctx, err = c.GetCtx(ictx)
		if err != nil {
			return nil, nil, err
		}

		paramsValue, err := reflectutils.New(paramType)
		if err != nil {
			return nil, nil, errors.New("create params error: %s ", err.Error())
		}
		params = paramsValue.Interface()
		params, err = GetWebParams(ictx, params, "base", "command")
		if err != nil {
			return nil, nil, err
		}
		if val, ok := params.(*map[string]any); ok {
			params = *val
		}
	}
	return params, rctx, nil
}

func (c *ApiController) GetCtx(ictx *context.Context) (context2.Context, error) {
	ctx, err := c.newContext(ictx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (c *ApiController) newContext(ictx *context.Context) (context2.Context, error) {
	ctx, err := restapp.NewTestContext(context2.Background())
	return ctx, err
}
