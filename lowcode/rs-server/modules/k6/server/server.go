package server

import (
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	swagger3 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/swagger/v3"
	"github.com/liuxd6825/k6server/js/modules"
)

type Server struct {
	app      *iris.Application
	vu       modules.VU
	services map[string]*Service
	swagger  *swagger3.Swagger
	cfg      common.IEnvConfig
}

type AddServiceOption struct {
	Name    string
	Desc    string
	Service *goja.Object
}

func NewServer(app *iris.Application, vu modules.VU, cfg common.IEnvConfig) *Server {
	return &Server{app: app, vu: vu, services: make(map[string]*Service), swagger: swagger3.NewSwagger(), cfg: cfg}
}

func (e *Server) Run() error {
	if err := e.initSwagger(); err != nil {
		return err
	}
	return nil
}

func (e *Server) initSwagger() error {
	sb := NewSwaggerBuilder()
	swagger, err := sb.Build(e)
	if err != nil {
		return err
	}
	e.swagger = swagger
	return nil
}

func (e *Server) Swagger() *swagger3.Swagger {
	return e.swagger
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

func (e *Server) AddService(opts *AddServiceOption) error {
	if opts == nil {
		return errors.New("opts is nil")
	}
	if opts.Service == nil {
		return errors.New("opts.Service is nil")
	}
	service := &Service{
		Service: opts.Service,
		Name:    opts.Name,
		Desc:    opts.Desc,
		server:  e,
	}
	e.services[opts.Name] = service
	return nil
}

func (e *Server) newObject(m map[string]any) (*goja.Object, error) {
	return NewObject(e.vu.Runtime(), m)
}
