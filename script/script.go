package script

import (
	"github.com/dop251/goja"
)

func NewRuntime() (*goja.Runtime, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	return vm, nil
}

func RunScript[T any](vm *goja.Runtime, values map[string]T, script *string) (goja.Value, error) {
	if script == nil {
		return nil, nil
	}
	code := *script
	if len(code) == 0 {
		return nil, nil
	}
	for k, v := range values {
		if err := vm.Set(k, v); err != nil {
			return nil, err
		}
	}
	value, err := vm.RunString(code)
	return value, err

}
