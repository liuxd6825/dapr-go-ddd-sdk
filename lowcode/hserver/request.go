package hserver

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils/schema_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/xtype"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/sirupsen/logrus"
)

// Request 定义单个请求的结构
type Request struct {
	*Base
	server          *Server
	service         *Service
	config          RequestConfig
	script          *Script
	paramsTypeCache *xtype.Map[xtype.ParamsType]
	schemaCache     *xtype.Map[*jsonschema.Schema]
}

type RequestConfig struct {
	Type          string       `json:"type"`
	Name          string       `json:"name"`
	URL           string       `json:"url"`
	AbsURl        string       `json:"abs_uri"`
	Description   string       `json:"description"`
	ParamsType    string       `json:"params_type"`
	LinkParamsUrl string       `json:"link_params_url"`
	Script        ScriptConfig `json:"script"`
}

func NewRequest(server *Server, service *Service, srcFileName string, config RequestConfig) (*Request, error) {
	var err error
	logger := service.GetLogger()
	r := &Request{
		server:          server,
		service:         service,
		config:          config,
		paramsTypeCache: xtype.NewMap[xtype.ParamsType](),
		schemaCache:     xtype.NewMap[*jsonschema.Schema](),
	}
	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, server.GetRootPath())
	r.Base, err = NewBase(srcFileName, logger, r, fsOpts)
	if err != nil {
		return nil, err
	}

	if r.config.Script.Code != "" {
		config.Script.FuncName = config.Name
		config.Script.SrcFileName = srcFileName
		err := r.scripts.AddScript(&r.config.Script, logger)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (r *Request) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	data, err := r.service.ReadFile(filename, opts...)
	return data, err
}

func (r *Request) Initialize() error {
	url := r.config.AbsURl
	if url == "" {
		url = r.service.config.URL + r.config.URL
	}
	r.server.app.Handle(r.config.Type, url, r.Handle)
	return nil
}

func (r *Request) GetLogger() logrus.FieldLogger {
	return r.service.GetLogger()
}

func (r *Request) Handle(ictx iris.Context) {
	var ctx context.Context
	var err error

	defer func() {
		err = RecoverError(err, recover())
		if err != nil {
			setError(ictx, err)
		}
	}()

	ctx, err = restapp.NewContext(ictx)
	if err != nil {
		return
	}
	wctx := NewWebContext(ctx, ictx, r)
	r.runScript(wctx, r.GetParamsValue(wctx))
}

func (r *Request) runScript(wctx *WebContext, params any) {
	values := &RunValues{
		Server:     r.server,
		Service:    r.service,
		Request:    r,
		WebContext: wctx,
	}
	tenantId := wctx.GetTenantId()
	val, err := r.scripts.RunScript(r.config.Script.FuncName, values, true, func(vm *goja.Runtime) error {
		_ = vm.Set("params", params)
		_ = vm.Set("ctx", wctx)
		_ = vm.Set("tenantId", tenantId)
		return nil
	})
	if err != nil {
		logs.Error(wctx, tenantId, logs.Fields{"code": r.config.Script.Code})
		wctx.SetError(err)
	}
	if val != nil {
		if err, ok := val.(error); ok {
			wctx.SetError(err)
		} else if v, ok := val.(goja.Value); ok {
			wctx.SetData(v.Export())
		} else {
			wctx.SetData(val)
		}
	}
}

// GetUrlParams
//
//	@Description: 取URL地址中的参数
//	@receiver r
//	@param ctx
//	@return map[string]any
func (r *Request) GetUrlParams(ctx iris.Context) map[string]any {
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

func (r *Request) GetParamsType(ictx iris.Context) (string, xtype.ParamsType) {
	fsOpts := &fsopts.Options{
		RootPath: r.server.GetRootPath(),
		WorkPath: r.service.fsOpts.WorkPath,
	}
	var paramsTypeFile string
	var paramsType xtype.ParamsType

	if r.config.LinkParamsUrl != "" {
		fileUrl := r.config.LinkParamsUrl
		urlPars := r.GetUrlParams(ictx)
		if (urlPars != nil) && (len(urlPars) > 0) {
			fileUrl = stringutils.ReplacePlaceholders(fileUrl, urlPars)
		}
		paramsTypeFile = r.service.GetWorkPath() + fileUrl
		if pType, ok := r.paramsTypeCache.Get(paramsTypeFile); ok {
			return paramsTypeFile, pType
		}
		bytes, err := r.server.ReadFile(fileUrl, fsOpts)
		if err != nil {
			panic(err)
		}
		if bytes != nil && len(bytes) > 0 {
			if err = jsonutils.Unmarshal(bytes, &paramsType); err != nil {
				panic(fmt.Sprintf(" loading %s  error: %s", fileUrl, err.Error()))
			}
			r.paramsTypeCache.Set(paramsTypeFile, paramsType)
		}

	} else if r.config.ParamsType != "" {
		paramsTypeFile = "/definition/params/" + r.config.ParamsType
		if pType, ok := r.paramsTypeCache.Get(paramsTypeFile); ok {
			return paramsTypeFile, pType
		}
		paramsType = r.server.definition.GetParamsType(r.config.ParamsType)
		r.paramsTypeCache.Set(paramsTypeFile, paramsType)
	}

	return paramsTypeFile, paramsType
}

// GetParamsValue
//
//	@Description: 获取请求的参数，aParamsURL的优先级最高，当为空时aParams参数生效。
//	@param wctx 请求的web上下文
//	@param aParams  通过对象定义的参数类型
//	@param cfgUrl  通过url定义的参数类型,
//	@return map[string]any 参数
func (r *Request) GetParamsValue(wctx *WebContext) map[string]any {
	var err error
	data := map[string]any{}

	paramsTypeFileName, paramsType := r.GetParamsType(wctx.ictx)
	if paramsType == nil {
		return data
	}

	ictx := wctx.ictx

	var bodyData any = nil

	for key, v := range paramsType {
		var val any
		switch v.In {
		case InParamTypeURL.String():
			val = ictx.URLParam(key)
		case InParamTypePath.String():
			val = ictx.Params().Get(key)
		case InParamTypeBody.String():
			if v.Schema != nil && bodyData == nil {
				var schema *jsonschema.Schema
				var ok bool
				if schema, ok = r.schemaCache.Get(paramsTypeFileName); !ok {
					schema, err = schema_utils.Compile(paramsTypeFileName, v.Schema, func(c *jsonschema.Compiler) error {
						c.UseLoader(r.server.schemaLoader)
						return nil
					})
					if err != nil {
						panic(err)
					}
				}
				obj := wctx.ReadObject(schema)
				bodyData = obj
			}
			val = bodyData
		case InParamTypeFormValue.String():
			val = wctx.FormValue(key, v.Required)
		case InParamTypeFormObject.String():
			schema, err := schema_utils.Compile(paramsTypeFileName, v.Schema, func(c *jsonschema.Compiler) error {
				c.UseLoader(r.server.schemaLoader)
				return nil
			})
			if err != nil {
				panic(err)
			}
			val = wctx.FormObject(key, v.Required, schema)
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

	return data
}
