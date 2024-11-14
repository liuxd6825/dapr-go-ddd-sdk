package server

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

type Handle func(cxt *WebContext, data any)
type MethodType string

const (
	GET    MethodType = "GET"
	POST   MethodType = "POST"
	PUT    MethodType = "PUT"
	DELETE MethodType = "DELETE"
)

func (m MethodType) String() string {
	return string(m)
}

type InParamType string

const (
	InParamTypeURL        InParamType = "url"        // URL中的参数
	InParamTypePath       InParamType = "path"       // 路径参数
	InParamTypeBody       InParamType = "body"       // 请求body参数
	InParamTypeFormValue  InParamType = "formValue"  // 从FormData中读取string
	InParamTypeFormObject InParamType = "formObject" // 从FormData中读取json转成对象
	InParamTypeFormFile   InParamType = "formFile"   // 从FormFile中读取文件
)

type RequestParam struct {
	In       string         `json:"in"` // InParamType
	Required bool           `json:"required"`
	Type     string         `json:"type"`
	Desc     string         `json:"desc"`
	Schema   *schema.Schema `json:"schema"`
	Example  any            `json:"example"`
}

type RequestData struct {
	Data   *goja.Object   `json:"data"` // 从body中读取的map数据
	Params map[string]any `json:"params"`
}

type HandleOptions struct {
	Method      MethodType              `json:"method"`
	Path        string                  `json:"path"`
	Description string                  `json:"description"`
	Body        *schema.Schema          `json:"body"`
	Params      map[string]RequestParam `json:"params"`
	ParamsUrl   string                  `json:"paramsUrl"`
	HandleName  string                  `json:"handleName"`
	Handle      func(cxt *WebContext, params map[string]any) any
}

func (h HandleOptions) GetMethod() MethodType {
	return h.Method
}

func (h *HandleOptions) GetPath() string {
	return h.Path
}

func (h *HandleOptions) GetDescription() string {
	return h.Description
}

func (h *HandleOptions) GetBody() *schema.Schema {
	return h.Body
}

func (h *HandleOptions) GetParams() map[string]RequestParam {
	return h.Params
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
	})
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
	var params map[string]RequestParam = aParams
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
				data[key] = bodyData
			case InParamTypeFormValue.String():
				val = wctx.FormValue(key, v.Required)
				data[key] = val
			case InParamTypeFormObject.String():
				val = wctx.FormObject(key, v.Required, v.Schema)
				data[key] = val
			case InParamTypeFormFile.String():
				val = wctx.FormFile(key)
				data[key] = val
			default:
				panic(fmt.Sprintf("The requested parameter [%s] type [%s] is incorrect, please use url,path,body,formValue", key, v.In))
			}

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
