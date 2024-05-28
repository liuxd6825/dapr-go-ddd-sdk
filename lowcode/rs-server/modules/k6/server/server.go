package server

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"github.com/liuxd6825/k6server/js/modules"
)

type Server struct {
	app *iris.Application
	vu  modules.VU
}

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

type HandleOptions struct {
	Method MethodType                    `json:"method"`
	Path   string                        `json:"path"`
	Handle func(cxt *RContext, data any) `json:"handle"`
	Schema *schema.Schema                `json:"schema"`
}

func NewServer(app *iris.Application, vu modules.VU) *Server {
	return &Server{app: app, vu: vu}
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

func (e *Server) Handle(opt *HandleOptions) {

	method := opt.Method.String()
	if opt.Handle == nil {
		return
	}

	e.app.Handle(method, opt.Path, func(ictx iris.Context) {
		var err error
		defer catchError(ictx, err, recover())

		rctx := NewRContext(ictx)
		var obj common.Object
		if opt.Schema != nil {
			obj, err = rctx.ReadObject(opt.Schema)
		}
		if err != nil {
			setError(ictx, err)
			return
		}
		opt.Handle(rctx, e.newObject(obj))
	})
}

func (e *Server) newObject(v map[string]any) *goja.Object {
	if v == nil {
		obj := e.vu.Runtime().NewObject()
		return obj
	}
	obj := e.vu.Runtime().NewObject()
	for k, v := range v {
		if m, ok := v.(map[string]any); ok {
			obj.Set(k, e.newObject(m))
			continue
		} else if m, ok := v.(common.Object); ok {
			obj.Set(k, e.newObject(m))
			continue
		}
		if err := obj.Set(k, v); err != nil {
			fmt.Println(err)
		}
	}
	return obj
}

func (e *Server) ReadJson(ictx iris.Context, data ...any) *common.Result[any] {
	var v any
	if len(data) > 0 {
		v = data[0]
	} else {
		v = make(map[string]any)
	}

	err := ictx.ReadJSON(&v)
	return common.NewResult[any](v, err)
}

func catchError(ctx iris.Context, e error, recover any) error {
	var err error
	if e != nil {
		err = e
	} else if recover != nil {
		if ve := recover.(error); ve != nil {
			err = ve
		} else {
			err = fmt.Errorf("unknown error %s")
		}
	}
	if err != nil && ctx != nil {
		setError(ctx, err)
	}
	return nil
}

func setError(ctx iris.Context, err error) {
	if err != nil && ctx != nil {
		ctx.SetErr(err)
		ctx.StatusCode(httptest.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
	}
}
