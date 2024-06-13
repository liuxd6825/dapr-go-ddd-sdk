package server

import (
	"github.com/dop251/goja"
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

type InParamType string

const (
	InParamType_URL   InParamType = "url"
	InParamType_Param InParamType = "param"
	InParamType_Body  InParamType = "body"
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
	Handle      func(cxt *WebContext, requestData *RequestData)
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
		if err != nil {
			setError(ictx, err)
			return
		}

		params, err := e.GetParams(ictx, opt.Params)
		if err != nil {
			setError(ictx, err)
			return
		}
		data, err := e.newObject(obj)
		if err != nil {
			setError(ictx, err)
			return
		}
		requestData := &RequestData{
			Data:   data,
			Params: params,
		}

		opt.Handle(rctx, requestData)
	})
}

func (e *Server) GetParams(ictx iris.Context, params map[string]RequestParam) (map[string]any, error) {
	var err error
	data := map[string]any{}
	if len(params) > 0 {
		for key, v := range params {
			var val string
			switch v.In {
			case "url":
				val = ictx.URLParam(key)
			case "path":
				val = ictx.Params().Get(key)
			}

			if v.Required && val == "" {
				err = errors.NewErr(err, "参数 %s 缺失", key)
			} else {
				if v, err := types.Convert(v.Type, val); err != nil {
					data[key] = v
				} else {
					err = errors.NewErr(err, "参数 %s 类型转换失败", key)
				}
			}
		}
	}
	if err != nil {
		setError(ictx, err)
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
