package hserver

import "github.com/dop251/goja"

type SetRuntimeOption struct {
	Server     *Server
	Service    *Service
	WebContext *WebContext
	Request    *Request
}

func setRuntime(vm *goja.Runtime, option *SetRuntimeOption, data ...map[string]any) error {
	if option != nil {
		if option.Server != nil {
			_ = vm.Set("$server", vm.NewDynamicObject(NewServerProxy(option.Server, vm)))
			if option.Server.console != nil {
				_ = vm.Set("console", option.Server.console)
			}
		}
		if option.Service != nil {
			_ = vm.Set("$service", vm.NewDynamicObject(NewServiceProxy(option.Service, vm)))
		}
		if option.WebContext != nil {
			_ = vm.Set("$wctx", option.WebContext)
		}
	}
	for _, d := range data {
		for k, v := range d {
			if obj, ok := v.(goja.DynamicObject); ok {
				_ = vm.Set(k, vm.NewDynamicObject(obj))
			} else {
				_ = vm.Set(k, v)
			}
		}
	}
	return nil
}
