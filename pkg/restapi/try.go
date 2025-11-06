package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

const gotry_rethrow = "----> Founded an Exception!!!\n"

// Handler object
type Handler struct {
	ictx    iris.Context
	ctx     context.Context
	catch   func(ctx context.Context, err error)
	finally func(ctx context.Context)
	error   error
}

// Throw function (return or rethrow an exception)
func Throw(e error) {
	if e == nil {
		panic(gotry_rethrow)
	} else {
		panic(e)
	}
}

func Try(ictx iris.Context, funcToTry func(ctx context.Context) error) (h *Handler) {
	var err error
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if h != nil {
			h.setError(err)
		} else {
			SetError(ictx, err)
		}
	}()
	h = newHandler(ictx)
	uri := h.ictx.Request().RequestURI
	method := h.ictx.Request().Method
	logs.Info(h.ctx, logs.Fields{"uri": uri, method: method})
	err = funcToTry(h.ctx)
	return h
}

func newHandler(ictx iris.Context) *Handler {
	ctx := context.Background()
	ctx, _ = restapp.NewTestContext(context.Background())
	goTry := &Handler{ictx: ictx, ctx: ctx, catch: nil, finally: nil, error: nil}
	return goTry
}

func (h *Handler) setError(err error) {
	if err != nil {
		h.error = err
		uri := h.ictx.Request().RequestURI
		logs.Error(h.ctx, logs.Fields{"uri": uri}, h.error)
		SetError(h.ictx, err)
	} else {
		h.ictx.StatusCode(iris.StatusOK)
	}
}

func (h *Handler) newCtx(ictx iris.Context) context.Context {
	ctx := context.Background()
	return ctx
}

func (h *Handler) Catch(funcCaught func(ctx context.Context, err error)) *Handler {
	h.catch = funcCaught
	if h.error != nil {
		defer func() {
			// call finally
			if h.finally != nil {
				h.finally(h.ctx)
			}

			if err := recover(); err != nil {
				if err == gotry_rethrow {
					err = h.error
				}
				panic(err)
			}
		}()
		h.catch(h.ctx, h.error)
	} else if h.finally != nil {
		h.finally(h.ctx)
	}
	return h
}

func (h *Handler) Finally(finallyFunc func(ctx context.Context)) *Handler {
	if h.finally != nil {
		panic("Finally Function by default !!")
	} else {
		h.finally = finallyFunc
	}
	defer h.finally(h.ctx)
	return h
}
