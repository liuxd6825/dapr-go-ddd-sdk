package hserver

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
)

type ScriptManager struct {
	scriptMap *types.CMap[*Script]
	logger    logrus.FieldLogger
	reader    fs.Reader
	transType runtime.TransformType
}

func NewScriptManager(logger logrus.FieldLogger, reader fs.Reader) *ScriptManager {
	return &ScriptManager{
		scriptMap: types.NewCMap[*Script](),
		logger:    logger,
		reader:    reader,
	}
}

func (b *ScriptManager) AddScript(config *ScriptConfig, logger logrus.FieldLogger, pkg *types.CMap[any]) error {
	var err error
	if config == nil {
		return errors.New("AddScript() no config provided")
	}
	if config.FuncName == "" {
		return errors.New("AddScript() no config.FuncName provided")
	}
	if config.SrcFileName == "" {
		return errors.New("AddScript() no config.SrcFileName provided")
	}

	script, err := NewScript(config, logger, b.reader, pkg)
	if err != nil {
		return err
	}

	b.scriptMap.Set(config.FuncName, script)
	return nil
}

func (b *ScriptManager) DeleteScript(scriptName string) error {
	b.scriptMap.Remove(scriptName)
	return nil
}

// RunScript
//
//	@Description: 运行脚本
//	@receiver b
//	@param scriptName 脚步名称
//	@param runValues 运行的环境变量
//	@param checkHave 是否检查, false:当scriptName不存返回nil, nil
//	@param config
//	@return res
//	@return err
func (b *ScriptManager) RunScript(funcName string, runValues *RunValues, checkHave bool, opts ...RunOptions) (res any, err error) {
	srcFileName := ""

	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v in %s %s ", rec, funcName, srcFileName)
			b.logger.Error(err)
		}
	}()

	script, ok := b.scriptMap.Get(funcName)
	if !ok {
		if checkHave {
			return nil, fmt.Errorf("script %s not found", funcName)
		}
		return nil, nil
	}

	srcFileName = script.config.SrcFileName

	opts = append(opts, func(vm *goja.Runtime) error {
		return setRunValues(vm, runValues)
	})

	val, err := script.Run(opts...)
	return val, err
}

func setRunValues(vm *goja.Runtime, runValues *RunValues, data ...map[string]any) error {

	if runValues != nil {
		if runValues.Server != nil {
			if err := runValues.Server.SetSelfVMValue("server", vm); err != nil {
				return err
			}
		}
		if runValues.Service != nil {
			if err := runValues.Service.SetSelfVMValue("service", vm); err != nil {
				return err
			}
		}
		if runValues.Self != nil {
			self, ok := runValues.Self.(Self)
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

	err := addRuntimeValues(vm, data...)
	if err != nil {
		return err
	}
	if runValues != nil {

	}
	return nil
}

func addRuntimeValues(vm *goja.Runtime, data ...map[string]any) error {
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
