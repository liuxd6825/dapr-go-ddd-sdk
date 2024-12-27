package template

import (
	"errors"
	"github.com/dop251/goja"
)

// ScriptRuntime represents a script that can be executed on the server-side.
type ScriptRuntime struct {
	vm *goja.Runtime
}

// NewScriptRuntime creates a new server script.
func NewScriptRuntime() *ScriptRuntime {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	return &ScriptRuntime{vm: vm}
}

func (s *ScriptRuntime) Run(data any, script string) (any, error) {
	s.vm.Set("data", data)
	_, err := s.vm.RunString(script)
	if err != nil {
		return nil, err
	}

	init, ok := goja.AssertFunction(s.vm.Get("init"))
	if !ok {
		return nil, errors.New("init function is not defined")
	}

	res, err := init(goja.Undefined(), nil)
	if err != nil {
		return nil, err
	}
	data = res.Export()
	return data, err
}
