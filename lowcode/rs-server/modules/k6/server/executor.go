package server

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type Executor struct {
	rctx        *WebContext
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

func (e *Executor) DoCommand(fun func(ctx context.Context)) *Executor {
	defer func() {
		e.error = catchError(e.rctx.ictx, e.error, recover())
	}()
	if !e.init() {
		return e
	}

	if fun != nil {
		fun(e.ctx)
	}
	return e
}

func (e *Executor) DoQuery(fun func(ctx context.Context) any) *Executor {
	defer func() {
		e.error = catchError(e.rctx.ictx, e.error, recover())
	}()
	if !e.init() {
		return e
	}
	if fun != nil {
		e.data = fun(e.ctx)
	}
	return e
}

// DoQueryOne
//
//	@Description: 如何数据库nil返回404错误。
//	@receiver e
//	@param fun
//	@return *Executor
func (e *Executor) DoQueryOne(fun func(ctx context.Context) any) *Executor {
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
		e.rctx.WriteJson(e.data)
	}
}
