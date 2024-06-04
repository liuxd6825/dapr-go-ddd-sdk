package server

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/k6server/js/modules"
)

type Server struct {
	app      *iris.Application
	vu       modules.VU
	services map[string]*Service
}

func NewServer(app *iris.Application, vu modules.VU) *Server {
	return &Server{app: app, vu: vu, services: make(map[string]*Service)}
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

func (e *Server) CreateOpenApiService() {

}

type AddServiceOption struct {
	Name    string
	Desc    string
	Service *goja.Object
}

func (e *Server) AddService(opts *AddServiceOption) error {
	if opts == nil {
		return errors.New("opts is nil")
	}
	if opts.Service == nil {
		return errors.New("opts.Service is nil")
	}
	e.services[opts.Name] = &Service{
		Service: opts.Service,
		Name:    opts.Name,
		Desc:    opts.Desc,
		server:  e,
	}
	return nil
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
