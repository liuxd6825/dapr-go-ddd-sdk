package server

import (
	"context"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/jsserver/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/k6server/js/modules"
)

type Server struct {
	app *iris.Application
	vu  modules.VU

	timerIDCounter uint64
}

// Exports returns the exports of the k6 module.
func (e *Server) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"server": e.vu.Runtime().ToValue(e),
			"time":   e.vu.Runtime().ToValue(NewTime()),
		},
	}
}

func (e *Server) nextID() uint64 {
	e.timerIDCounter++
	return e.timerIDCounter
}

func (e *Server) call(callback goja.Callable, args []goja.Value) error {
	// TODO: investigate, not sure GlobalObject() is always the correct value for `this`?
	_, err := callback(e.vu.Runtime().GlobalObject(), args...)
	return err
}

func (e *Server) Post(relativePath string, handle func(ictx iris.Context)) {
	e.Handle("POST", relativePath, handle)
}

func (e *Server) Get(relativePath string, handle func(ictx iris.Context)) {
	e.Handle("GET", relativePath, handle)
}

func (e *Server) Delete(relativePath string, handle func(ictx iris.Context)) {
	e.Handle("DELETE", relativePath, handle)
}

func (e *Server) Put(relativePath string, handle func(ictx iris.Context)) {
	e.Handle("PUT", relativePath, handle)
}

func (e *Server) Handle(method string, relativePath string, fun func(ictx iris.Context)) {
	e.app.Handle(method, relativePath, func(ictx iris.Context) {
		defer ctxRecover(ictx, recover())
		fun(ictx)
	})
}

func (e *Server) DoQuery(ictx iris.Context, tenantId string, fun restapp.QueryFunc, opts ...restapp.DoOptions) *common.Result[any] {
	data, _, err := restapp.DoQuery(ictx, tenantId, fun, opts...)
	return common.NewResult[any](data, err)
}

func (e *Server) DoRequest(ictx iris.Context, tenantId string, fun func(ctx context.Context) error, opts ...restapp.DoOptions) *common.Result[any] {
	err := restapp.Do(ictx, tenantId, fun, opts...)
	return common.NewResult[any](nil, err)
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

type QueryFunc = func(tx context.Context) *common.Result[any]

func (e *Server) DoQueryOne(ictx iris.Context, tenantId string, fun QueryFunc, opts ...restapp.DoOptions) *common.Result[any] {
	data, _, err := restapp.DoQueryOne(ictx, tenantId, func(ctx context.Context) (interface{}, bool, error) {
		return fun(ctx).GetResults()
	}, opts...)
	return common.NewResult[any](data, err)
}

func (e *Server) DoCmdAndQueryOne(ictx iris.Context, tenantId, queryAppId string, cmd restapp.Command, cmdFun restapp.CmdFunc, queryFun restapp.QueryFunc, opts ...restapp.CmdAndQueryOption) *common.Result[any] {
	data, _, err := restapp.DoCmdAndQueryOne(ictx, tenantId, queryAppId, cmd, cmdFun, queryFun, opts...)
	return common.NewResult[any](data, err)
}

func (e *Server) DoCmdAndQueryList(ictx iris.Context, tenantId string, queryAppId string, cmd restapp.Command, cmdFun restapp.CmdFunc, queryFun restapp.QueryFunc, opts ...restapp.CmdAndQueryOption) *common.Result[any] {
	data, _, err := restapp.DoCmdAndQueryList(ictx, tenantId, queryAppId, cmd, cmdFun, queryFun, opts...)
	return common.NewResult[any](data, err)
}

func (e *Server) DoCmd(ictx iris.Context, tenantId string, fun restapp.CmdFunc, opts ...restapp.DoOptions) *common.Result[any] {
	err := restapp.DoCmd(ictx, tenantId, fun, opts...)
	return common.NewResult[any](nil, err)
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
