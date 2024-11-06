package server

import (
	"context"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
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

type InParamType string

const (
	InParamType_URL  InParamType = "url"
	InParamType_Path InParamType = "path"
	InParamType_Body InParamType = "body"
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
		params, err := e.GetParams(wctx, opt.Params)
		if err != nil {
			setError(ictx, err)
			return
		}

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

func (e *Server) GetParams(wctx *WebContext, params map[string]RequestParam) (map[string]any, error) {
	var err error
	data := map[string]any{}
	var bodyData any = nil
	if len(params) > 0 {
		for key, v := range params {
			var val string
			switch v.In {
			case InParamType_URL.String():
				val = wctx.ictx.URLParam(key)
			case InParamType_Path.String():
				val = wctx.ictx.Params().Get(key)
			case InParamType_Body.String():
				if v.Schema != nil && bodyData == nil {
					v.Schema.Init()
					obj := wctx.ReadObject(v.Schema)
					bodyData = obj
				}
				data[key] = bodyData
			}

			if v.In == InParamType_URL.String() || v.In == InParamType_Path.String() {
				if v.Required && val == "" {
					err = errors.NewErr(err, "参数 %s 缺失", key)
				} else {
					if v, err := types.Convert(v.Type, val); err == nil {
						data[key] = v
					} else {
						err = errors.NewErr(err, "参数 %s 类型转换失败", key)
					}
				}
			}
		}
	}
	if err != nil {
		setError(wctx.ictx, err)
	}
	return data, err
}

func (e *Server) AddHandles(handlers ...*HandleOptions) error {
	for _, opt := range handlers {
		if opt != nil {
			e.Handle(opt)
		}
	}
	return nil
}
