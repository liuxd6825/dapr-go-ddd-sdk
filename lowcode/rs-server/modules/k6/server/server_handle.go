package server

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
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

type RequestParam struct {
	Required bool   `json:"required,omitempty"`
	Type     string `json:"type"`
}

type RequestData struct {
	Data       any            `json:"data,omitempty"`
	PathParams map[string]any `json:"pathParams,omitempty"`
	UrlParams  map[string]any `json:"urlParams,omitempty"`
}
type RequestOptions struct {
	ServiceName string                  `json:"serviceName"`
	Method      MethodType              `json:"method"`
	Path        string                  `json:"path"`
	Body        *schema.Schema          `json:"body"`
	PathParams  map[string]RequestParam `json:"pathParams"`
	UrlParams   map[string]RequestParam `json:"urlParams"`
	HandleName  string                  `json:"methodName"`
	Handle      func(cxt *WebContext, requestData *RequestData)
}

func (e *Server) Get(opt *RequestOptions) {
	opt.Method = GET
	e.Handle(opt)
}

func (e *Server) Post(opt *RequestOptions) {
	opt.Method = POST
	e.Handle(opt)
}

func (e *Server) Put(opt *RequestOptions) {
	opt.Method = PUT
	e.Handle(opt)
}

func (e *Server) Delete(opt *RequestOptions) {
	opt.Method = DELETE
	e.Handle(opt)
}

func (e *Server) Handle(opt *RequestOptions) {
	if opt.Handle == nil {
		e.app.Logger().Error("未在%s上设置处理函数", opt.Path)
		return
	}

	method := opt.Method.String()

	e.app.Handle(method, opt.Path, func(ictx iris.Context) {
		var err error
		defer func() {
			_ = catchError(ictx, err, recover())
		}()

		rctx := NewWebContext(ictx)
		var obj common.Object
		if opt.Body != nil {
			opt.Body.Init()
			obj, err = rctx.ReadObject(opt.Body)
			if err != nil {
				err = errors.NewErr(err, "数据验证失败")
			}
		}
		urlParams := map[string]any{}
		if val, err := e.GetUrlParams(ictx, opt.UrlParams); err != nil {
			return
		} else {
			urlParams = val
		}

		pathParams := map[string]any{}
		if val, err := e.GetPathParams(ictx, opt.PathParams); err != nil {
			return
		} else {
			pathParams = val
		}

		if err != nil {
			setError(ictx, err)
			return
		}
		requestData := &RequestData{
			Data:       e.newObject(obj),
			PathParams: pathParams,
			UrlParams:  urlParams,
		}
		opt.Handle(rctx, requestData)
	})
}

func (e *Server) GetUrlParams(ictx iris.Context, params map[string]RequestParam) (map[string]any, error) {
	var err error
	data := map[string]any{}
	if len(params) > 0 {
		for k, v := range params {
			val := ictx.URLParam(k)
			if v.Required && val == "" {
				err = errors.NewErr(err, "url 参数 %s 缺失", k)
			} else {
				if v, err := types.Convert(v.Type, val); err != nil {
					data[k] = v
				} else {
					err = errors.NewErr(err, "url 参数 %s 类型转换失败", k)
				}
			}
		}
	}
	if err != nil {
		setError(ictx, err)
	}
	return data, err
}

func (e *Server) GetPathParams(ictx iris.Context, params map[string]RequestParam) (map[string]any, error) {
	var err error
	data := map[string]any{}
	if len(params) > 0 {
		for k, v := range params {
			val := ictx.Params().Get(k)
			if v.Required && val == "" {
				err = errors.NewErr(err, "path 参数 %s 缺失", k)
			} else {
				if v, err := types.Convert(v.Type, val); err != nil {
					data[k] = v
				} else {
					err = errors.NewErr(err, "path 参数 %s 类型转换失败", k)
				}
			}
		}
	}
	if err != nil {
		setError(ictx, err)
	}
	return data, err
}

func (e *Server) Handles(handlers []*RequestOptions) error {
	for _, opt := range handlers {
		e.Handle(opt)
	}
	return nil
}
