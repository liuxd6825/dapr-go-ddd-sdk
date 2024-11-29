package hserver

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/sirupsen/logrus"
)

// Request 定义单个请求的结构
type Request struct {
	*Base
	server  *Server
	service *Service
	params  map[string]RequestParam
	config  RequestConfig
	script  *Script
}

type RequestConfig struct {
	Type        string       `json:"type"`
	Name        string       `json:"name"`
	URL         string       `json:"url"`
	AbsURl      string       `json:"abs_uri"`
	Description string       `json:"description"`
	ParamsType  string       `json:"params_type"`
	ParamsUrl   string       `json:"params_url"`
	Script      ScriptConfig `json:"script"`
}

func NewRequest(server *Server, service *Service, srcFileName string, config RequestConfig) (*Request, error) {
	var err error
	logger := service.GetLogger()
	r := &Request{
		server:  server,
		service: service,
		config:  config,
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
		_ = RecoverError(err, recover())
	}()

	ctx, err = restapp.NewContext(ictx)
	if err != nil {
		return
	}
	wctx := NewWebContext(ctx, ictx)
	params := r.GetParams(wctx, r.config.ParamsUrl)
	r.runScript(wctx, params)
}

func (r *Request) runScript(wctx *WebContext, params any) {
	values := &RunValues{
		Server:     r.server,
		Service:    r.service,
		Request:    r,
		WebContext: wctx,
	}
	findPaging := wctx.GetFindPaging()
	fmt.Sprintf("%v", findPaging)
	tenantId := wctx.GetTenantId()
	val, err := r.scripts.RunScript(r.config.Script.FuncName, values, true, func(vm *goja.Runtime) error {
		_ = vm.Set("params", params)
		_ = vm.Set("ctx", wctx)
		_ = vm.Set("tenantId", tenantId)
		_ = vm.Set("tokenUser", wctx.GetTokenUser())
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

func (r *Request) GetUrlParams(ctx iris.Context) map[string]any {
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
//	@param cfgUrl  通过url定义的参数类型,
//	@return map[string]any 参数
func (r *Request) GetParams(wctx *WebContext, cfgUrl string) map[string]any {
	var err error
	var params = r.params
	fsOpts := &fsopts.Options{
		RootPath: r.server.GetRootPath(),
		WorkPath: r.service.fsOpts.WorkPath,
	}
	if params == nil && cfgUrl != "" {
		fileUrl := cfgUrl
		urlPars := r.GetUrlParams(wctx.ictx)
		if (urlPars != nil) && (len(urlPars) > 0) {
			fileUrl = stringutils.ReplacePlaceholders(cfgUrl, urlPars)
		}
		bytes, err := r.server.ReadFile(fileUrl, fsOpts)
		if err != nil {
			panic(err)
		}
		if bytes != nil && len(bytes) > 0 {
			if err = jsonutils.Unmarshal(bytes, &params); err != nil {
				panic(fmt.Sprintf(" loading %s  error: %s", fileUrl, err.Error()))
			}
			r.params = params
		}
	}

	data := map[string]any{}
	if params == nil {
		return data
	}

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
