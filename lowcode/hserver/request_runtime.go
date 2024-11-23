package hserver

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
)

type RequestRuntime struct {
	runtime *runtime.Pool
	request *Request
	code    string
}

func NewRequestRuntime(request *Request, code string) *RequestRuntime {
	return &RequestRuntime{
		runtime: runtime.NewPool(5),
		request: request,
		code:    code,
	}
}

func (r *RequestRuntime) Execute(wctx *WebContext, requestParams map[string]any) (any, error) {
	data, err := r.runtime.Run(r.code, func(vm *goja.Runtime) error {
		return setRuntime(vm, &SetRuntimeOption{
			Service:    r.request.service,
			Server:     r.request.server,
			WebContext: wctx,
			Request:    r.request,
		}, map[string]any{"$params": requestParams})
	})
	return data, err
}
