package server

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/k6/schema"
)

type Server struct {
	app *iris.Application
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
	Method MethodType                                 `json:"method"`
	Path   string                                     `json:"path"`
	Handle func(cxt *RContext, data ...common.Object) `json:"handle"`
	Schema *schema.Schema                             `json:"schema"`
}

func NewServer(app *iris.Application) *Server {
	return &Server{app: app}
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
	e.app.Handle(method, opt.Path, func(ictx iris.Context) {
		var object = common.NewObject()
		var err error
		if opt.Schema != nil {
			if err = ictx.ReadJSON(&object); err == nil {
				if err = opt.Schema.Validate(object); err == nil {
					object, err = opt.Schema.Convertor(object)
				}
			}
		}
		if err != nil {
			setError(ictx, err)
			return
		}
		defer ctxRecover(ictx, recover())
		opt.Handle(NewRContext(ictx), object)
	})
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

func setError(ctx iris.Context, err error) {
	ctx.SetErr(err)
	ctx.StatusCode(httptest.StatusInternalServerError)
	_, _ = ctx.WriteString(err.Error())
}

func ctxRecover(ctx iris.Context, recover any) {
	if recover != nil {
		if err := recover.(error); err != nil {
			setError(ctx, err)
		}
	}
}
