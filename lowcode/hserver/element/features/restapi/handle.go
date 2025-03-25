package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/jsonschema_ext"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/convert"
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
	params := s.GetParams(wctx)
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

// GetParams
//
//	@Description: 获取请求的参数
//	@param wctx 请求的web上下文
//	@param aParams  通过对象定义的参数类型
//	@param cfgUrl  通过url定义的参数类型,
//	@return map[string]any 参数
func (s *ApiHandle) GetParams(wctx element.WebContext) any {
	var err error

	// 获取参数文件与参数类型
	sch := s.GetParamsSchema(wctx.ICtx())
	if sch == nil {
		return nil
	}

	// 要返回的值
	data := map[string]any{}
	ictx := wctx.ICtx()

	metaSch := jsonschema_ext.GetMetaExtension(sch)
	if metaSch != nil {
		if metaSch.Param != nil && metaSch.Param.Type == jsonschema_ext.HParamType_Body {
			res := wctx.ReadObject(sch)
			return res
		}
	}

	properties := sch.GetAllProperties()
	for key, prop := range properties {
		var val any
		meta := jsonschema_ext.GetMetaExtension(prop)
		if meta == nil {
			panic("no metadata")
		}
		if meta.Param == nil {
			panic("no param")
		}
		param := meta.Param
		paramName := param.Name
		required := sch.IsRequired(key)
		switch param.Type {
		case jsonschema_ext.HParamType_Path:
			val = ictx.URLParam(paramName)
		case jsonschema_ext.HParamType_Query:
			val = ictx.URLParam(paramName)
		case jsonschema_ext.HParamType_Body:
			val = wctx.ReadObject(prop)
		case jsonschema_ext.HParamTypee_FormValue:
			val = wctx.FormValue(paramName, required)
		case jsonschema_ext.HParamType_FormObject:
			val = wctx.FormObject(paramName, required, sch)
		case jsonschema_ext.HParamType_FormFile:
			val = wctx.FormFile(paramName)
		default:
			panic(fmt.Sprintf("The requested parameter [%s] type [%s] is incorrect, please use url,path,body,formValue", key, param.Name))
		}

		if (val == nil || val == "") && prop.Default != nil {
			val = prop.Default
			switch val.(type) {
			case json.Number:
				num := val.(json.Number)
				val, err = num.Int64()
			case string:
				val = val.(string)
			}
		}

		if prop.Types != nil {
			val, err = s.convert(prop.Name(), prop.Types, val)
			if err != nil {
				panic(fmt.Sprintf("params.%s types.Convert() error: %s", key, err.Error()))
			}
		}
		if required && (val == "" || val == nil) {
			err = errors.New("The requested parameter %s cannot be empty", key)
			panic(err)
		}
		data[key] = val
	}

	paramData, err := jsonschema_ext.DoConvert(sch, data)
	if err != nil {
		panic(err)
	}
	return paramData
}

func (s *ApiHandle) convert(propName string, schTypes *jsonschema.Types, val any) (res any, err error) {
	isNull := false
	if ok := schTypes.Contains(jsonschema.JsonType_NullType); ok {
		isNull = true
	}
	if val == nil {
		if isNull {
			return nil, errors.New("%s not null", propName)
		}
	}
	if ok := schTypes.Contains(jsonschema.JsonType_DateTimeType); ok {
		res, err = convert.Convert(convert.ConvertTypeDateTime, val)
	} else if ok := schTypes.Contains(jsonschema.JsonType_StringType); ok {
		res, err = convert.Convert(convert.ConvertTypeString, val)
	} else if ok := schTypes.Contains(jsonschema.JsonType_IntegerType); ok {
		res, err = convert.Convert(convert.ConvertTypeInt, val)
	} else if ok := schTypes.Contains(jsonschema.JsonType_BooleanType); ok {
		res, err = convert.Convert(convert.ConvertTypeBool, val)
	} else if ok := schTypes.Contains(jsonschema.JsonType_NumberType); ok {
		res, err = convert.Convert(convert.ConvertTypeNumber, val)
	} else if ok := schTypes.Contains(jsonschema.JsonType_ObjectType); ok {
		res, err = convert.Convert(convert.ConvertTypeString, val)
	} else if ok := schTypes.Contains(jsonschema.JsonType_ArrayType); ok {
		res, err = convert.Convert(convert.ConvertTypeString, val)
	} else if ok := schTypes.Contains(jsonschema.JsonType_DateType); ok {
		res, err = convert.Convert(convert.ConvertTypeDateTime, val)
	}
	if err != nil {
		err = errors.New("%s types.Convert() error: %s", propName, err.Error())
	}
	return res, err
}

// GetParamsValue
//
//	@Description: 获取请求的参数
//	@param wctx 请求的web上下文
//	@param aParams  通过对象定义的参数类型
//	@param cfgUrl  通过url定义的参数类型,
//	@return map[string]any 参数

/*
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

*/

// GetParamsSchema
//
//	@Description: 获取参数类型定义
//	@receiver r
//	@param ictx
//	@return string  参数文件名
//	@return xtype.ParamsType  参数配置类型
func (s *ApiHandle) GetParamsSchema(ictx iris.Context) *jsonschema.Schema {
	urlPars := s.GetUrlParams(ictx)
	return s.fun.GetParamsSchema(urlPars, s.service.FsOpts())
}

// GetParamsType
//
//	@Description: 获取参数类型定义
//	@receiver r
//	@param ictx
//	@return string  参数文件名
//	@return xtype.ParamsType  参数配置类型
/*
func (s *ApiHandle) GetParams(ictx iris.Context) (string, common.ParamsType) {
	urlPars := s.GetUrlParams(ictx)
	return s.fun.GetParamsType(urlPars, s.service.FsOpts())
}
*/

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
