package restapi

import (
	context2 "context"
	"fmt"

	"reflect"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
)

type ApiController struct {
	app      *iris.Application
	rootPath string
	routes   []*router.Route
	apiName  string
	apiCtl   any
}

type HandleType int

const (
	HandleType_API   HandleType = iota // 服务API
	HandleType_Event                   // Event事件
	HandleType_CDC                     // 数据库数据变化
)

type ApiFunc func(ctx context2.Context, ictx *context.Context, params any) (any, error)

type APIController interface {
	NewAPIController(app *iris.Application) *ApiController
}

type RenderViewFunc func(ctx context2.Context, ictx iris.Context, fileName string, viewData ...map[string]any) error

type CallOptions struct {
	IsAuthentication *bool                        // 是否进行身份认证
	ParamsInBody     *bool                        // 强制要求参数入HttpBody中获取
	InitMethod       func(callMethod *CallMethod) // 初始化函数
	GetFieldValue    func(ictx *context.Context, parentObject any, fieldType reflect.StructField, fieldValue reflect.Value) (outFieldValue any, ok bool, err error)
	Before           func(ictx *context.Context, params any) (any, context2.Context, error)
	After            func(ctx context2.Context, ictx *context.Context, resData any, err error) (any, error)
}

var apis = types.NewCMap[any]()

var View RenderViewFunc = func(ctx context2.Context, ictx iris.Context, fileName string, viewData ...map[string]any) error {
	if fileName == "" {
		fileName = ictx.Request().URL.Path
	}
	vData := make([]any, len(viewData))
	for i, data := range viewData {
		vData[i] = data
	}
	return ictx.View(fileName, vData...)
}

func NewCallOptions(opts ...CallOptions) CallOptions {
	o := CallOptions{}
	for _, i := range opts {
		if i.InitMethod != nil {
			o.InitMethod = i.InitMethod
		}
		if i.GetFieldValue != nil {
			o.GetFieldValue = i.GetFieldValue
		}
		if i.Before != nil {
			o.Before = i.Before
		}
		if i.After != nil {
			o.After = i.After
		}
		if i.ParamsInBody != nil {
			o.ParamsInBody = i.ParamsInBody
		}
		if i.IsAuthentication != nil {
			o.IsAuthentication = i.IsAuthentication
		}
	}
	return o
}

func (o *CallOptions) GetIsAuthentication() bool {
	if o.IsAuthentication == nil {
		return false
	}
	return *o.IsAuthentication
}

func WithParamsInBody(paramsInBody bool) CallOptions {
	return CallOptions{
		ParamsInBody: &paramsInBody,
	}
}

func RegisterController(app *iris.Application, api APIController) {
	apiCtl := api.NewAPIController(app)
	if apiCtl == nil {
		panic("api controller is nil")
	}
	apiName := apiCtl.apiName
	if apis.Has(apiName) {
		panic(fmt.Sprintf("ApiName %s has already been initialized", apiName))
	}
	apis.Add(apiName, api)
}

func GetAPI(apiName string) (any, bool) {
	return apis.Get(apiName)
}

func NewController(app *iris.Application, rootPath string, apiName string, apiController any) *ApiController {
	return &ApiController{
		app:      app,
		rootPath: rootPath,
		apiName:  apiName,
		apiCtl:   apiController,
		routes:   make([]*router.Route, 0),
	}
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
	if c.apiCtl == nil {
		return nil, errors.New("controller is nil")
	}
	method, err := NewCallMethod(c.apiCtl, handlerName)
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
		GetFieldValue: func(ictx *context.Context, object any, fieldType reflect.StructField, fieldValue reflect.Value) (outParam any, ok bool, err error) {
			switch fieldType.Name {
			case "ValueCols":
				valueCols := ictx.URLParam("value-cols")
				var val []*store.ValueCol
				if valueCols != "" {
					val = store.NewValueColsWidth(valueCols)
				}
				return val, true, nil
			case "GroupCols":
				groupCols := ictx.URLParam("group-cols")
				var val []*store.GroupCol
				if groupCols != "" {
					val = store.NewGroupColsWidthString(groupCols)
				}
				return val, true, nil
			case "GroupKeys":
				groupKeys := ictx.URLParam("group-keys")
				var val []any
				if groupKeys != "" {
					val = store.NewGroupKeysWidthString(groupKeys)
				}
				return val, true, nil
			}
			return nil, false, nil
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

func (c *ApiController) CDCHandle(path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.callMethod2(iris.MethodPost, path, handlerName, HandleType_CDC, false, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) Handle(method string, path string, handlerName string, opts ...CallOptions) *router.Route {
	r := c.call(method, path, handlerName, opts...)
	c.addRouter(r)
	return r
}

func (c *ApiController) View(path string, handlerName string, opts ...CallOptions) *router.Route {
	return c.callMethod2(iris.MethodGet, path, handlerName, HandleType_API, true, opts...)
}

func (c *ApiController) call(method string, path string, handlerName string, opts ...CallOptions) *router.Route {
	return c.callMethod(method, path, handlerName, HandleType_API, opts...)
}

func (c *ApiController) callEventHandle(path string, handlerName string, opts ...CallOptions) *router.Route {
	return c.callMethod(iris.MethodPost, path, handlerName, HandleType_Event, opts...)
}

func (c *ApiController) cdcHandle(path string, handlerName string, opts ...CallOptions) *router.Route {
	return c.callMethod(iris.MethodPost, path, handlerName, HandleType_CDC, opts...)
}

func (c *ApiController) callMethod(method string, path string, handlerName string, handleType HandleType, opts ...CallOptions) *router.Route {
	return c.callMethod2(method, path, handlerName, handleType, false, opts...)
}

func (c *ApiController) callView(method string, path string, handlerName string, handleType HandleType, opts ...CallOptions) *router.Route {
	return c.callMethod2(method, path, handlerName, handleType, true, opts...)
}

func (c *ApiController) callMethod2(method string, path string, handlerName string, handleType HandleType, isViewHandle bool, opts ...CallOptions) *router.Route {
	backCtx := context2.Background()
	callMethod, err := c.newCallMethod(handlerName)
	if err != nil {
		ctlType := reflect.TypeOf(c.apiCtl)
		if ctlType.Kind() == reflect.Ptr {
			ctlType = ctlType.Elem()
		}
		pkgPath := ctlType.PkgPath()
		typeName := ctlType.Name()
		err = errors.New("%s %s.%s() func error: %s ", pkgPath, typeName, handlerName, err.Error())
		panic(err)
	}

	if !isViewHandle {
		path = c.getPath(path)
	}

	r := c.app.Handle(method, path, func(ictx *context.Context) {
		gp.Try(func() error {
			params, ctx, err := c.getParams(ictx, callMethod, handleType, opts...)
			if err != nil {
				if handleType == HandleType_Event {
					logs.Error(context2.Background(), logs.Fields{"type": "event", "urlPath": path, "handlerName": handlerName, "error": err.Error()})
					return nil
				}
				return err
			}
			logs.Info(backCtx, logs.Fields{"method": method, "path": path, "params": params})
			data, err := callMethod.Call(ctx, ictx, params)
			for _, opt := range opts {
				if opt.After != nil {
					data, err = opt.After(ctx, ictx, data, err)
				}
			}
			if err == nil && callMethod.OutData >= 0 {
				//logs.Info(backCtx, logs.Fields{"method": method, "path": path, "handlerName": handlerName, "data": data})
				err = SetOKJsonData(ictx, data)
			}
			return err
		}).Catch(func(err error) {
			SetError(ictx, err)
			logs.Error(backCtx, logs.Fields{"method": method, "path": path, "error": err})
		})
	})
	c.addRouter(r)
	return r
}

func (c *ApiController) getParams(ictx *context.Context, callMethod *CallMethod, handleType HandleType, opts ...CallOptions) (params any, rctx context2.Context, err error) {
	// CallMethod初始化
	for _, opt := range opts {
		if opt.InitMethod != nil {
			opt.InitMethod(callMethod)
		}
	}

	if callMethod.InParams >= 0 && !callMethod.CloseInParams {
		inParamsType := callMethod.Method.Type().In(callMethod.InParams)
		params, rctx, err = c.GetParams(ictx, inParamsType, handleType, opts...)
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
func (c *ApiController) GetParams(ictx *context.Context, paramType reflect.Type, handleType HandleType, opts ...CallOptions) (params any, rctx context2.Context, err error) {
	paramsValue, err := reflectutils.New(paramType)
	if err != nil {
		return nil, nil, errors.New("create params error: %s ", err.Error())
	}
	params = paramsValue.Interface()

	if handleType == HandleType_Event {
		params, rctx, err = GetEventParams(ictx, params)
		return params, rctx, err
	} else if handleType == HandleType_CDC {
		params, rctx, err = GetCDCParams(ictx, params)
		return params, rctx, err
	}
	if paramType != nil {
		rctx, err = c.GetCtx(ictx)
		if err != nil {
			return nil, nil, err
		}

		params, err = GetWebParams(ictx, params, []string{"base", "command"}, opts...)
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
	ctx, err := c.newContext(context2.Background(), ictx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (c *ApiController) newContext(parent context2.Context, ictx *context.Context) (ctx context2.Context, err error) {
	ctx = parent
	app := env.GetEnv().App
	appId := app.AppId
	appName := app.AppName
	if app.ProdMode {
		ctx = appctx.NewWebContext(ctx, ictx)
	} else {
		ctx, err = restapp.NewTestContext(context2.Background())
		if err != nil {
			return nil, err
		}
	}

	ctx = appctx.NewAppContext(ctx, appId, appName)
	return ctx, nil
}
