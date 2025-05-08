package server

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	swagger3 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/swagger/v3"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/k6server/js/modules"
	"sync"
)

type Server struct {
	app      *iris.Application
	vu       modules.VU
	services map[string]*Service
	swagger  *swagger3.Swagger
	cfg      common.IEnvConfig
	mux      sync.Mutex
}

type AddServiceOption struct {
	Name    string
	Desc    string
	Service *goja.Object
}

func NewServer(app *iris.Application, vu modules.VU, cfg common.IEnvConfig) *Server {
	return &Server{app: app, vu: vu, services: make(map[string]*Service), swagger: swagger3.NewSwagger(), cfg: cfg}
}

func (e *Server) Run() {
	if err := e.initSwagger(); err != nil {
		panic(err)
	}
}

func (e *Server) initSwagger() error {
	sb := NewSwaggerBuilder()
	swagger, err := sb.Build(e)
	if err != nil {
		return err
	}
	e.swagger = swagger
	return nil
}

func (e *Server) Swagger() *swagger3.Swagger {
	return e.swagger
}

func (e *Server) ReadJson(ictx iris.Context, data ...any) *common.Result[any] {
	var v any
	if len(data) > 0 {
		v = data[0]
	} else {
		v = make(map[string]any)
	}

	err := ictx.ReadJSON(&v)
	return common.NewResult[any](v, err)
}

func (e *Server) AddService(opts *AddServiceOption) {
	if opts == nil {
		panic(errors.New("opts is nil"))
	}
	if opts.Service == nil {
		panic(errors.New("opts.Service is nil"))
	}
	service := &Service{
		Service: opts.Service,
		Name:    opts.Name,
		Desc:    opts.Desc,
		server:  e,
	}
	e.services[opts.Name] = service

}

func (e *Server) newObject(m map[string]any) (*goja.Object, error) {
	return NewObject(e.vu.Runtime(), m)
}

func (e *Server) Get(opt *HandleOptions) {
	opt.Method = GET
	e.Handle(opt)
}

func (e *Server) Post(opt *HandleOptions) {
	opt.Method = POST
	e.Handle(opt)
}

func (e *Server) Put(opt *HandleOptions) {
	opt.Method = PUT
	e.Handle(opt)
}

func (e *Server) Delete(opt *HandleOptions) {
	opt.Method = DELETE
	e.Handle(opt)
}

func (e *Server) Handle(opt *HandleOptions) {
	if opt.Handle == nil {
		e.app.Logger().Error("未在%s上设置处理函数", opt.Path)
		return
	}
	if opt.Method == "" {
		opt.Method = GET
	}

	method := opt.Method.String()
	e.app.Handle(method, opt.Path, func(ictx iris.Context) {
		e.handle(ictx, opt)
	})
}

func (e *Server) handle(ictx iris.Context, opt *HandleOptions) {
	e.mux.Lock()
	defer e.mux.Unlock()

	var ctx context.Context
	var err error

	defer func() {
		_ = catchError(ictx, err, recover())
	}()

	ctx, err = restapp.NewContext(ictx)
	if err != nil {
		return
	}
	wctx := NewWebContext(ctx, ictx, e.vu)
	params := e.GetParams(wctx, opt.Params, opt.ParamsUrl)
	res := opt.Handle(wctx, params)
	if err, ok := res.(error); ok {
		setError(ictx, err)
		return
	}
	if res != nil {
		wctx.WriteJson(res)
	}
}

func (e *Server) GetUrlParams(ctx iris.Context) map[string]any {
	params := make(map[string]any)
	// 获取所有路径参数
	pathParams := ctx.Params()
	for _, param := range pathParams.Store {
		params[param.Key] = param.Value
	}

	// 获取所有查询参数
	queryParams := ctx.URLParams()
	for key, value := range queryParams {
		params[key] = value
	}

	return params
}

// GetParams
//
//	@Description: 获取请求的参数，aParamsURL的优先级最高，当为空时aParams参数生效。
//	@param wctx 请求的web上下文
//	@param aParams  通过对象定义的参数类型
//	@param aParamsURL  通过url定义的参数类型,
//	@return map[string]any 参数
func (e *Server) GetParams(wctx *WebContext, aParams map[string]RequestParam, aParamsURL string) map[string]any {
	var err error
	var params = aParams
	if aParamsURL != "" {
		if params == nil {
			params = map[string]RequestParam{}
		}
		schemaFileUrl := aParamsURL
		urlPars := e.GetUrlParams(wctx.ictx)
		if (urlPars != nil) && (len(urlPars) > 0) {
			schemaFileUrl = stringutils.ReplacePlaceholders(aParamsURL, urlPars)
		}

		bytes := _fs.ReadFile(schemaFileUrl)
		if err = jsonutils.Unmarshal(bytes, &params); err != nil {
			panic(err)
		}
	}

	data := map[string]any{}
	ictx := wctx.ictx

	var bodyData any = nil
	if len(params) > 0 {
		for key, v := range params {
			if v.Schema != nil {
				v.Schema.Name = key
			}
			var val any
			switch v.In {
			case InParamTypeURL.String():
				val = ictx.URLParam(key)
			case InParamTypePath.String():
				val = ictx.Params().Get(key)
			case InParamTypeBody.String():
				if v.Schema != nil && bodyData == nil {
					obj := wctx.ReadObject(v.Schema)
					bodyData = obj
				}
				val = bodyData
			case InParamTypeFormValue.String():
				val = wctx.FormValue(key, v.Required)
			case InParamTypeFormObject.String():
				val = wctx.FormObject(key, v.Required, v.Schema)
			case InParamTypeFormFile.String():
				val = wctx.FormFile(key)
			default:
				panic(fmt.Sprintf("The requested parameter [%s] type [%s] is incorrect, please use url,path,body,formValue", key, v.In))
			}

			data[key] = val

			if key == InParamTypeURL.String() || key == InParamTypePath.String() {
				if v.Required && val == "" {
					err = errors.New("The requested parameter %s cannot be empty", key)
					panic(err)
				} else {
					val, err = types.Convert(v.Type, v)
					if err != nil {
						panic(err)
					}
					data[key] = v
				}
			}

		}
	}
	return data
}

func (e *Server) AddHandles(handlers ...*HandleOptions) error {
	for _, opt := range handlers {
		if opt != nil {
			e.Handle(opt)
		}
	}
	return nil
}
