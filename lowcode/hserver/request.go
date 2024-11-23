package hserver

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

// Request 定义单个请求的结构
type Request struct {
	Type        string                  `json:"type"`
	Name        string                  `json:"name"`
	URL         string                  `json:"url"`
	AbsURl      string                  `json:"abs_uri"`
	Description string                  `json:"description"`
	Params      string                  `json:"params"`
	Code        string                  `json:"code,omitempty"`
	server      *Server                 `json:"-"`
	service     *Service                `json:"-"`
	runtime     *RequestRuntime         `json:"-"`
	params      map[string]RequestParam `json:"-"`
}

func (e *Request) init(server *Server, service *Service) {
	e.server = server
	e.service = service
	if e.Code != "" {
		e.runtime = NewRequestRuntime(e, e.Code)
	}

}

func (e *Request) Handle(ictx iris.Context) {

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
	params := e.GetParams(wctx, e.Params)
	res := e.handle(wctx, params)
	if err, ok := res.(error); ok {
		setError(ictx, err)
		return
	}
	if res != nil {
		wctx.WriteJson(res)
	}
}

func (e *Request) handle(wctx *WebContext, params map[string]any) any {
	if e.runtime != nil {
		data, err := e.runtime.Execute(wctx, params)
		if err != nil {
			panic(err)
		}
		return data
	}
	return nil
}

func (e *Request) GetUrlParams(ctx iris.Context) map[string]any {
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
func (e *Request) GetParams(wctx *WebContext, aParamsURL string) map[string]any {
	var err error
	var params = e.params

	if params == nil && aParamsURL != "" {
		schemaFileUrl := aParamsURL
		urlPars := e.GetUrlParams(wctx.ictx)
		if (urlPars != nil) && (len(urlPars) > 0) {
			schemaFileUrl = stringutils.ReplacePlaceholders(aParamsURL, urlPars)
		}
		fileUrl := fileutils.AbsPath(schemaFileUrl, e.service.srcPath)
		bytes := e.server.fs.ReadFile(fileUrl)
		if bytes != nil && len(bytes) > 0 {
			if err = jsonutils.Unmarshal(bytes, &params); err != nil {
				panic(err)
			}
			e.params = params
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
