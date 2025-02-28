package restapi

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/jsonschema/v6"
	"strings"
)

type ApiHandle struct {
	server          element.Server
	service         element.Service
	tag             *element.FuncTag
	fun             element.Func
	paramsTypeCache *types.CMap[common.ParamsType]
	schemaCache     *types.CMap[*jsonschema.Schema]
	config          *RestConfig
}

func Register(service element.Service, fun element.Func, tag *element.FuncTag) {
	h := NewApiHandle(service, tag, fun)
	if err := h.Initialize(); err != nil {
		panic(err)
	}
}

func NewApiHandle(service element.Service, tag *element.FuncTag, fun element.Func) *ApiHandle {
	cfg := NewRestConfig(tag)
	return &ApiHandle{
		server:          service.Server(),
		service:         service,
		tag:             tag,
		fun:             fun,
		config:          cfg,
		paramsTypeCache: types.NewCMap[common.ParamsType](),
		schemaCache:     types.NewCMap[*jsonschema.Schema](),
	}
}

func (s *ApiHandle) Initialize() error {
	url := s.config.AbsUrl()
	if url == "" {
		url = s.service.Config().Url() + s.config.Url()
	}
	method := strings.ToUpper(s.config.Method())
	s.service.Server().App().Handle(method, url, s.handle)
	return nil
}

func (s *ApiHandle) handle(ictx iris.Context) {
	var ctx context.Context
	var err error

	defer func() {
		err = utils.RecoverError(err, recover())
		if err != nil {
			utils.SetError(ictx, err)
			logs.Errorfmt(ctx, "", "%s error:%s", ictx.Request().RequestURI, err.Error())
		}
	}()
	logs.Infof(ctx, "", nil, "request %s %s", ictx.Request().Method, ictx.Request().RequestURI)

	ctx, err = restapp.NewContext(ictx)
	if err != nil {
		return
	}
	wctx := s.server.Factory().NewWebContext(ctx, ictx)
	params := s.GetParamsValue(wctx)
	s.run(wctx.Ctx(), wctx, params)
}

func (s *ApiHandle) run(ctx context.Context, wctx element.WebContext, params any) {
	values := &element.ApiRunValues{
		Server:     s.server,
		Self:       s.service,
		WebContext: wctx,
		WorkPath:   s.service.WorkPath(),
	}
	tenantId := wctx.GetTenantId()
	funcName := s.fun.Config().FuncName
	val, err := s.service.Funcs().Run(ctx, funcName, values, true, func(vm *goja.Runtime) error {
		_ = vm.Set("params", params)
		_ = vm.Set("wctx", wctx)
		_ = vm.Set("ctx", wctx.Ctx())
		_ = vm.Set("tenantId", tenantId)
		return nil
	})
	if err != nil {
		wctx.SetError(err)
		return
	}
	if val != nil {
		if err, ok := val.(error); ok {
			wctx.SetError(err)
		} else if v, ok := val.(goja.Value); ok {
			wctx.WriteJson(v.Export())
		} else {
			wctx.WriteJson(val)
		}
	}
}

// GetParamsValue
//
//	@Description: 获取请求的参数
//	@param wctx 请求的web上下文
//	@param aParams  通过对象定义的参数类型
//	@param cfgUrl  通过url定义的参数类型,
//	@return map[string]any 参数
func (s *ApiHandle) GetParamsValue(wctx element.WebContext) map[string]any {
	var err error

	// 获取参数文件与参数类型
	paramsTypeFileName, paramsType := s.GetParamsType(wctx.ICtx())
	if paramsType == nil {
		return nil
	}

	if paramsTypeFileName == "/definition/params/findPaging.json" {
		return wctx.GetFindPaging().AsMap()
	}

	// 要返回的值
	data := map[string]any{}
	ictx := wctx.ICtx()

	for key, v := range paramsType {
		key = strings.ToLower(key)
		var val any
		switch v.In {
		case InParamTypeURL.String():
			val = ictx.URLParam(key)
		case InParamTypePath.String():
			val = ictx.Params().Get(key)
		case InParamTypeBody.String():
			if v.Schema != nil {
				schema := v.Schema.Init(paramsTypeFileName, s.server.SchemaLoader())
				val = wctx.ReadObject(schema)
			} else {
				panic(errors.New("paramType %s is no schema defined", key))
			}
		case InParamTypeFormValue.String():
			val = wctx.FormValue(key, v.Required)
		case InParamTypeFormObject.String():
			if v.Schema != nil {
				schema := v.Schema.Init(paramsTypeFileName, s.server.SchemaLoader())
				val = wctx.FormObject(key, v.Required, schema)
			}
		case InParamTypeFormFile.String():
			val = wctx.FormFile(key)
		default:
			panic(fmt.Sprintf("The requested parameter [%s] type [%s] is incorrect, please use url,path,body,formValue", key, v.In))
		}

		if (val == nil || val == "") && v.Default != nil {
			val = v.Default
		}

		if v.Type != "" {
			val, err = types.Convert(v.Type, val)
			if err != nil {
				panic(fmt.Sprintf("params.%s types.Convert() error: %s", key, err.Error()))
			}
		}
		if v.Required && (val == "" || val == nil) {
			err = errors.New("The requested parameter %s cannot be empty", key)
			panic(err)
		}
		data[key] = val
	}
	return data
}

// GetParamsType
//
//	@Description: 获取参数类型定义
//	@receiver r
//	@param ictx
//	@return string  参数文件名
//	@return xtype.ParamsType  参数配置类型
func (s *ApiHandle) GetParamsType(ictx iris.Context) (string, common.ParamsType) {
	urlPars := s.GetUrlParams(ictx)
	return s.fun.GetParamsType(urlPars, s.service.FsOpts())
}

// GetUrlParams
//
//	@Description: 取URL地址中的参数
//	@receiver r
//	@param ctx
//	@return map[string]any
func (s *ApiHandle) GetUrlParams(ctx iris.Context) map[string]any {
	p := make(map[string]any)
	// 获取所有路径参数
	pathParams := ctx.Params()
	for _, param := range pathParams.Store {
		p[param.Key] = param.Value
	}

	// 获取所有查询参数
	queryParams := ctx.URLParams()
	for key, value := range queryParams {
		p[key] = value
	}

	return p
}
