package server

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/common"
)

type Executor struct {
	rctx        *RContext
	error       error
	data        any
	errorStatus int
	ctx         context.Context
}

func (e *Executor) init() bool {
	if e.error != nil {
		return false
	}
	if e.ctx == nil {
		e.ctx, e.error = restapp.NewCtx(e.rctx.ictx)
	}
	return e.error == nil
}

func (e *Executor) DoCommand(fun func(ctx context.Context) error) *Executor {
	defer func() {
		e.error = catchError(e.rctx.ictx, e.error, recover())
	}()
	if !e.init() {
		return e
	}

	if fun != nil {
		e.error = fun(e.ctx)
	}
	return e
}

func (e *Executor) DoQuery(fun func(ctx context.Context) *common.Result[any]) *Executor {
	defer func() {
		e.error = catchError(e.rctx.ictx, e.error, recover())
	}()
	if !e.init() {
		return e
	}

	if fun != nil {
		r := fun(e.ctx)
		e.data = r.Data
		e.error = r.Error
	}
	return e
}

// DoQueryOne
//
//	@Description: 如何数据库nil返回404错误。
//	@receiver e
//	@param fun
//	@return *Executor
func (e *Executor) DoQueryOne(fun func(ctx context.Context) *common.Result[any]) *Executor {
	defer func() {
		e.error = catchError(e.rctx.ictx, e.error, recover())
	}()
	e.DoQuery(fun)
	if e.error == nil && e.data == nil {
		e.error = notFoundError
	}
	return e
}

func (e *Executor) DoError(fun func(err error)) *Executor {
	defer func() {
		e.error = catchError(e.rctx.ictx, e.error, recover())
	}()
	if fun != nil && e.error != nil {
		fun(e.error)
	}
	return e
}

func (e *Executor) GetError() error {
	return e.error
}

func (e *Executor) GetData() any {
	return e.data
}

func (e *Executor) GetResult() *common.Result[any] {
	return common.NewResult[any](e.data, e.error)
}

func (e *Executor) SetResponse() {
	if e.error != nil {
		if errors.Is(e.error, notFoundError) {
			e.rctx.SetError(e.error, iris.StatusNotFound)
		} else {
			e.rctx.SetError(e.error)
		}
	}
	if e.data != nil {
		_ = e.rctx.WriteJson(e.data)
	}
}
