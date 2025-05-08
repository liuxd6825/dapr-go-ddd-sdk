package hserver

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"

	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
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
	paramsTypeCache *types.CMap[xtype.ParamsType]
	schemaCache     *types.CMap[*jsonschema.Schema]
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
		paramsTypeCache: types.NewCMap[xtype.ParamsType](),
		schemaCache:     types.NewCMap[*jsonschema.Schema](),
	}
	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, server.GetRootPath())
	r.Base, err = NewBase(srcFileName, logger, r, fsOpts)
	if err != nil {
		return nil, err
	}
	r.SetPkg(server.GetPkg())

	if r.config.Script.Code != "" {
		config.Script.FuncName = config.Name
		config.Script.SrcFileName = srcFileName
		err := r.scripts.AddScript(&r.config.Script, logger, r.pkg)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (s *Request) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	data, err := s.service.ReadFile(filename, opts...)
	return data, err
}

func (s *Request) Initialize() error {
	url := s.config.AbsURl
	if url == "" {
		url = s.service.config.URL + s.config.URL
	}
	s.server.app.Handle(s.config.Type, url, s.Handle)
	return nil
}

func (s *Request) GetLogger() logrus.FieldLogger {
	return s.service.GetLogger()
}

func (s *Request) Handle(ictx iris.Context) {
	var ctx context.Context
	var err error

	defer func() {
		err = utils.RecoverError(err, recover())
		if err != nil {
			utils.SetError(ictx, err)
		}
	}()

	ctx, err = restapp.NewContext(ictx)
	if err != nil {
		return
	}
	wctx := NewWebContext(ctx, ictx, s)
	s.runScript(wctx, s.GetParamsValue(wctx))
}

func (s *Request) runScript(wctx *WebContext, params any) {
	values := &RunValues{
		Server:     s.server,
		Self:       s.service,
		WebContext: wctx,
		WorkPath:   s.GetWorkPath(),
	}
	tenantId := wctx.GetTenantId()
	val, err := s.scripts.RunScript(s.config.Script.FuncName, values, true, func(vm *goja.Runtime) error {
		_ = vm.Set("params", params)
		_ = vm.Set("ctx", wctx)
		_ = vm.Set("tenantId", tenantId)
		return nil
	})
	if err != nil {
		wctx.SetError(err)
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

// GetUrlParams
//
//	@Description: 取URL地址中的参数
//	@receiver r
//	@param ctx
//	@return map[string]any
func (s *Request) GetUrlParams(ctx iris.Context) map[string]any {
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

// GetParamsType
//
//	@Description: 获取参数类型定义
//	@receiver r
//	@param ictx
//	@return string  参数文件名
//	@return xtype.ParamsType  参数配置类型
func (s *Request) GetParamsType(ictx iris.Context) (string, xtype.ParamsType) {
	fsOpts := &fsopts.Options{
		RootPath: s.server.GetRootPath(),
		WorkPath: s.service.fsOpts.WorkPath,
	}
	var paramsTypeFile string
	var paramsType xtype.ParamsType

	// 引用Schema文件
	if s.config.LinkParamsUrl != "" {
		fileUrl := s.config.LinkParamsUrl
		urlPars := s.GetUrlParams(ictx)
		if (urlPars != nil) && (len(urlPars) > 0) {
			fileUrl = stringutils.ReplacePlaceholders(fileUrl, urlPars)
		}
		paramsTypeFile = s.service.GetWorkPath() + fileUrl
		if pType, ok := s.paramsTypeCache.Get(paramsTypeFile); ok {
			return paramsTypeFile, pType
		}
		bytes, err := s.server.ReadFile(fileUrl, fsOpts)
		if err != nil {
			panic(err)
		}
		if bytes != nil && len(bytes) > 0 {
			if err = jsonutils.Unmarshal(bytes, &paramsType); err != nil {
				panic(fmt.Sprintf(" loading %s  error: %s", fileUrl, err.Error()))
			}
			s.paramsTypeCache.Set(paramsTypeFile, paramsType)
		}

	} else if s.config.ParamsType != "" {
		// 引用系统中的schema定义文件
		paramsTypeFile = fmt.Sprintf("/definition/params/%s.json", s.config.ParamsType)
		if pType, ok := s.paramsTypeCache.Get(paramsTypeFile); ok {
			return paramsTypeFile, pType
		}
		paramsType = s.server.definition.GetParamsType(s.config.ParamsType + ".json")
		s.paramsTypeCache.Set(paramsTypeFile, paramsType)
	}

	return paramsTypeFile, paramsType
}

// GetParamsValue
//
//	@Description: 获取请求的参数
//	@param wctx 请求的web上下文
//	@param aParams  通过对象定义的参数类型
//	@param cfgUrl  通过url定义的参数类型,
//	@return map[string]any 参数
func (s *Request) GetParamsValue(wctx *WebContext) map[string]any {
	var err error
	paramsTypeFileName, paramsType := s.GetParamsType(wctx.ictx)
	if paramsType == nil {
		return nil
	}

	data := map[string]any{}
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
				schema := v.Schema.Init(paramsTypeFileName, s.server.schemaLoader)
				bodyData = wctx.ReadObject(schema)
			}
			val = bodyData
		case InParamTypeFormValue.String():
			val = wctx.FormValue(key, v.Required)
		case InParamTypeFormObject.String():
			if v.Schema != nil {
				schema := v.Schema.Init(paramsTypeFileName, s.server.schemaLoader)
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

		val, err = types.Convert(v.Type, val)
		if err != nil {
			panic(fmt.Sprintf("params.%s types.Convert() error: %s", key, err.Error()))
		}

		if v.Required && (val == "" || val == nil) {
			err = errors.New("The requested parameter %s cannot be empty", key)
			panic(err)
		}
		data[key] = val
	}
	return data
}

func (s *Request) SetSelfVMValue(name string, vm *goja.Runtime) error {
	return nil
}
