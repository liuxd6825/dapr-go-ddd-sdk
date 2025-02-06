package funcs

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	funcs     *types.CMap[element.Func]
	logger    logrus.FieldLogger
	reader    fs.Reader
	transType runtime.TransformType
}

func NewFuncManager(logger logrus.FieldLogger, reader fs.Reader) *Manager {
	return &Manager{
		funcs:  types.NewCMap[element.Func](),
		logger: logger,
		reader: reader,
	}
}

func (b *Manager) Items() map[string]element.Func {
	return b.funcs.Items()
}

func (b *Manager) Add(fun element.Func) error {
	b.funcs.Set(fun.Config().FuncName, fun)
	return nil
}

func (b *Manager) Delete(scriptName string) error {
	b.funcs.Remove(scriptName)
	return nil
}

// Run
//
//	@Description: 运行脚本
//	@receiver b
//	@param scriptName 脚步名称
//	@param runValues 运行的环境变量
//	@param checkHave 是否检查, false:当scriptName不存返回nil, nil
//	@param config
//	@return res
//	@return err
func (b *Manager) Run(funcName string, runValues *element.ApiRunValues, checkHave bool, opts ...RunOptions) (res any, err error) {
	srcFileName := ""

	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v in %s %s ", rec, funcName, srcFileName)
			b.logger.Error(err)
		}
	}()

	script, ok := b.funcs.Get(funcName)
	if !ok {
		if checkHave {
			return nil, fmt.Errorf("script %s not found", funcName)
		}
		return nil, nil
	}

	srcFileName = script.Config().SrcFileName

	opts = append(opts, func(vm *goja.Runtime) error {
		return SetRunValues(vm, runValues)
	})

	val, err := script.Run(opts...)
	return val, err
}

func SetRunValues(vm *goja.Runtime, runValues *element.ApiRunValues, data ...map[string]any) error {
	if runValues != nil {
		if runValues.Server != nil {
			if err := runValues.Server.SetSelfVMValue("server", vm); err != nil {
				return err
			}
		}
		if runValues.ApiService != nil {
			if err := runValues.ApiService.SetSelfVMValue("service", vm); err != nil {
				return err
			}
		}
		if runValues.Self != nil {
			self, ok := runValues.Self.(element.Self)
			if ok {
				if err := self.SetSelfVMValue("self", vm); err != nil {
					return err
				}
			}
		}

		if runValues.WebContext != nil {
			_ = vm.Set("ctx", runValues.WebContext)
		}

		if runValues.WorkPath != "" {
			_ = vm.Set("workPath", runValues.WorkPath)
		}
	}

	err := AddRuntimeValues(vm, data...)
	if err != nil {
		return err
	}
	if runValues != nil {

	}
	return nil
}

func AddRuntimeValues(vm *goja.Runtime, data ...map[string]any) error {
	for _, d := range data {
		for k, v := range d {
			if init, ok := v.(InitVM); ok {
				if init.InitVM != nil {
					if err := init.InitVM(vm); err != nil {
						return err
					}
				}
			}
			if obj, ok := v.(goja.DynamicObject); ok {
				_ = vm.Set(k, vm.NewDynamicObject(obj))
			} else if pkgVal, ok := v.(pkg.Package); ok {
				_ = vm.Set(k, vm.NewDynamicObject(pkgVal.NewProxy(vm)))
			} else {
				_ = vm.Set(k, v)
			}
		}
	}
	return nil
}
