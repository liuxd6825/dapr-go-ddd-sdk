package server

import (
	"context"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/k6server/js/modules"
)

type Server struct {
	app *iris.Application
	vu  modules.VU

	timerIDCounter uint64
}

type ResultDoQuery struct {
	Data    any
	IsFound bool
	Error   error
}

type ResultData struct {
	Data  any
	Error error
}

type QueryFunc func(ctx context.Context) (map[string]any, bool, error)

// Exports returns the exports of the k6 module.
func (e *Server) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"server":            e.vu.Runtime().ToValue(e),
			"readJson":          e.vu.Runtime().ToValue(ReadJson),
			"doRequest":         e.vu.Runtime().ToValue(DoRequest),
			"doQuery":           e.vu.Runtime().ToValue(DoQuery),
			"doQueryOne":        e.vu.Runtime().ToValue(DoQueryOne),
			"doCmdAndQueryOne":  e.vu.Runtime().ToValue(DoCmdAndQueryOne),
			"doCmdAndQueryList": e.vu.Runtime().ToValue(DoCmdAndQueryList),
			"doCmd":             e.vu.Runtime().ToValue(DoCmd),
			"time":              e.vu.Runtime().ToValue(NewTime()),
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

func DoQuery(ictx iris.Context, tenantId string, fun restapp.QueryFunc, opts ...restapp.DoOptions) (any, bool, error) {
	return restapp.DoQuery(ictx, tenantId, fun, opts...)
}

func DoRequest(ictx iris.Context, tenantId string, fun func(ctx context.Context) error, opts ...restapp.DoOptions) (err error) {
	return restapp.Do(ictx, tenantId, fun, opts...)
}

func ReadJson(ictx iris.Context) (map[string]any, error) {
	var data map[string]any
	err := ictx.ReadJSON(&data)
	return data, err
}

func DoQueryOne(ictx iris.Context, tenantId string, fun restapp.QueryFunc, opts ...restapp.DoOptions) *ResultDoQuery {
	data, isFound, err := restapp.DoQueryOne(ictx, tenantId, fun, opts...)
	return &ResultDoQuery{Data: data, IsFound: isFound, Error: err}
}

func DoCmdAndQueryOne(ictx iris.Context, tenantId, queryAppId string, cmd restapp.Command, cmdFun restapp.CmdFunc, queryFun restapp.QueryFunc, opts ...restapp.CmdAndQueryOption) *ResultDoQuery {
	data, isFound, err := restapp.DoCmdAndQueryOne(ictx, tenantId, queryAppId, cmd, cmdFun, queryFun, opts...)
	return &ResultDoQuery{Data: data, IsFound: isFound, Error: err}
}

func DoCmdAndQueryList(ictx iris.Context, tenantId string, queryAppId string, cmd restapp.Command, cmdFun restapp.CmdFunc, queryFun restapp.QueryFunc, opts ...restapp.CmdAndQueryOption) *ResultDoQuery {
	data, isFound, err := restapp.DoCmdAndQueryList(ictx, tenantId, queryAppId, cmd, cmdFun, queryFun, opts...)
	return &ResultDoQuery{Data: data, IsFound: isFound, Error: err}
}

func DoCmd(ictx iris.Context, tenantId string, fun restapp.CmdFunc, opts ...restapp.DoOptions) (err error) {
	return restapp.DoCmd(ictx, tenantId, fun, opts...)
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
