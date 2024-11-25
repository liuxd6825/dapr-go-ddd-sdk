package hserver

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/sirupsen/logrus"
)

// Request 定义单个请求的结构
type Request struct {
	server  *Server
	service *Service
	params  map[string]RequestParam
	scripts *ScriptManager
	opts    RequestOptions
}

type RequestOptions struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	AbsURl      string `json:"abs_uri"`
	Description string `json:"description"`
	ParamsUrl   string `json:"params"`
	PoolSize    int    `json:"pool_size"`
	Code        string `json:"code"`
	CodeType    string `json:"code_type"`
}

func NewRequest(server *Server, service *Service, opts RequestOptions) (*Request, error) {
	logger := service.GetLogger()
	r := &Request{
		server:  server,
		service: service,
		opts:    opts,
		scripts: NewScriptManager(logger),
	}
	server.scripts = NewScriptManager(r.GetLogger())
	if r.opts.Code != "" {
		srcFileName := fmt.Sprintf("%s_%s.js", service.Name, r.opts.Name)
		err := r.scripts.AddScript("handler", r.opts.Code, r.opts.CodeType, srcFileName, true, logger)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (r *Request) Run() error {
	url := r.opts.AbsURl
	if url == "" {
		url = r.service.URL + r.opts.URL
	}
	r.server.app.Handle(r.opts.Type, url, r.Handle)
	return nil
}

func (r *Request) GetLogger() logrus.FieldLogger {
	return r.service.GetLogger()
}

func (r *Request) Handle(ictx iris.Context) {
	var ctx context.Context
	var err error

	defer func() {
		_ = catchError(ictx, err, recover())
	}()

	ctx, err = restapp.NewContext(ictx)
	if err != nil {
		return
	}
	wctx := NewWebContext(ctx, ictx)
	params := r.GetParams(wctx, r.opts.ParamsUrl)
	res := r.handle(wctx, params)
	if err, ok := res.(error); ok {
		setError(ictx, err)
		return
	}
	if res != nil {
		wctx.WriteJson(res)
	}
}

func (r *Request) handle(wctx *WebContext, params map[string]any) any {
	opts := &RuntimeOption{
		Server:     r.server,
		Service:    r.service,
		Request:    r,
		WebContext: wctx,
	}
	val, err := r.scripts.RunScript("handler", opts, func(vm *goja.Runtime) (map[string]any, error) {
		vars := map[string]any{
			"params": params,
			"ctx":    wctx.ctx,
		}
		return vars, nil
	})
	if err != nil {
		panic(err)
	}
	return val
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

	if params == nil && cfgUrl != "" {
		schemaFileUrl := cfgUrl
		urlPars := r.GetUrlParams(wctx.ictx)
		if (urlPars != nil) && (len(urlPars) > 0) {
			schemaFileUrl = stringutils.ReplacePlaceholders(cfgUrl, urlPars)
		}
		fileUrl := fileutils.AbsPath(schemaFileUrl, r.service.srcPath)
		bytes := r.server.fs.ReadFile(fileUrl)
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
