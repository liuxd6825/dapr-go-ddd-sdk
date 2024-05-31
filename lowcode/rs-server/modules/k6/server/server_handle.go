package server

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
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

type HandleOpt struct {
	Method MethodType                      `json:"method"`
	Path   string                          `json:"path"`
	Handle func(cxt *WebContext, data any) `json:"handle"`
	Schema *schema.Schema                  `json:"schema"`
}

func (e *Server) Get(opt *HandleOpt) {
	opt.Method = GET
	e.Handle(opt)
}

func (e *Server) Post(opt *HandleOpt) {
	opt.Method = POST
	e.Handle(opt)
}

func (e *Server) Put(opt *HandleOpt) {
	opt.Method = PUT
	e.Handle(opt)
}

func (e *Server) Delete(opt *HandleOpt) {
	opt.Method = DELETE
	e.Handle(opt)
}

func (e *Server) Handle(opt *HandleOpt) {
	method := opt.Method.String()
	if opt.Handle == nil {
		return
	}

	e.app.Handle(method, opt.Path, func(ictx iris.Context) {
		var err error
		defer catchError(ictx, err, recover())

		rctx := NewWebContext(ictx)
		var obj common.Object
		if opt.Schema != nil {
			obj, err = rctx.ReadObject(opt.Schema)
			if err != nil {
				err = errors.NewErr(err, "数据验证失败")
			}
		}
		if err != nil {
			setError(ictx, err)
			return
		}
		opt.Handle(rctx, e.newObject(obj))
	})
}

func (e *Server) Handles(handlers []*HandleOpt) error {
	for _, opt := range handlers {
		e.Handle(opt)
	}
	return nil
}
