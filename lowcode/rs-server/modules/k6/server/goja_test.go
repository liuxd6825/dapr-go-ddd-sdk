package server

import (
	"fmt"
	"github.com/dop251/goja"
	"testing"
)

type Console struct {
}

func NewConsole() *Console {
	return &Console{}
}

func (c *Console) Log(a ...interface{}) {
	fmt.Println(a...)
}

func (c *Console) Info(a ...interface{}) {
	fmt.Println(a...)
}

func run(t *testing.T, code string, values map[string]any) (*goja.Runtime, goja.Value, error) {
	vm := goja.New()
	for k, v := range values {
		if err := vm.Set(k, v); err != nil {
			return nil, nil, err
		}
	}
	value, err := vm.RunString(code)
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	if err != nil {
		t.Error(err)
		return nil, nil, err
	}
	return vm, value, nil
}

func call(t *testing.T, code string, funName string, fun any) error {
	vm, _, err := run(t, code, nil)
	if err != nil {
		return err
	}
	vm.Get(funName)
	err = vm.ExportTo(vm.Get(funName), fun)
	if err != nil {
		t.Error(err)
		return err
	}
	return nil
}

func runString(t *testing.T, code string, values map[string]any) (*goja.Runtime, goja.Value, error) {
	vm := goja.New()
	for k, v := range values {
		if err := vm.Set(k, v); err != nil {
			return nil, nil, err
		}
	}
	value, err := vm.RunString(code)
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	if err != nil {
		t.Error(err)
		return nil, nil, err
	}
	return vm, value, nil
}

func Test_HandleOptions(t *testing.T) {
	code := `
	var Service = /** @class */ (function () {
		function Service() {
			var _this = this;
			this.handle = function(request) {
				console.Log("name:", request.command.data.data.data.name);
			};
		}
		return Service;
	}());
	var service = new Service();
	console.Info("handle");
	setHandler({Handle:service.handle})
`

	type SetHandlerOptions struct {
		Handle func(request map[string]any)
	}
	setHandler := func(opts *SetHandlerOptions) {
		request := map[string]any{
			"tenantId": "0001",
			"command": map[string]any{
				"data": map[string]any{
					"name": "0000--0000",
					"data": map[string]any{
						"name": "1111--1111",
						"data": map[string]any{
							"name": "2222--2222",
						},
					},
				},
			},
		}
		opts.Handle(request)
	}
	data := map[string]interface{}{
		"console":    NewConsole(),
		"setHandler": setHandler,
	}
	_, _, err := run(t, code, data)
	if err != nil {
		t.Fatal(err)
	}
	/*
		v := vm.Get("service")
		service, ok := v.(*goja.Object)
		if ok {
			_ = &Request{
				Command: map[string]any{
					"commandId": "commandId",
					"tenantId":  "0001",
					"data": map[string]any{
						"name": "lxd",
						"size": "XXL",
					},
				},
				TenantId: "0001",
			}
			fmt.Println("return:", service)
		}*/
}
