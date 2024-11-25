package hserver

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"strings"
)

type ScriptManager struct {
	scriptMap cmap.ConcurrentMap
	logger    logrus.FieldLogger
}

func NewScriptManager(logger logrus.FieldLogger) *ScriptManager {
	return &ScriptManager{
		scriptMap: cmap.New(),
		logger:    logger,
	}
}

func (b *ScriptManager) AddScript(scriptName string, code string, codeType string, srcFileName string, poolSize bool, logger logrus.FieldLogger) error {
	var err error
	if strings.ToLower(codeType) == "text/typescript" {
		code, err = runtime.TransformTs(b.logger, code, srcFileName)
		if err != nil {
			return err
		}
	}

	script, err := NewScript(code, codeType, srcFileName, poolSize, logger)
	if err != nil {
		return err
	}
	b.scriptMap.Set(scriptName, script)
	return nil
}

func (b *ScriptManager) DeleteScript(scriptName string) error {
	b.scriptMap.Remove(scriptName)
	return nil
}

func (b *ScriptManager) RunScript(scriptName string, opts *RuntimeOption, initValues func(vm *goja.Runtime) (map[string]any, error)) (res any, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	val, ok := b.scriptMap.Get(scriptName)
	if !ok {
		return nil, fmt.Errorf("script %s not found", scriptName)
	}

	script, ok := val.(*Script)
	if !ok {
		return nil, fmt.Errorf("script %s is not of type *Script", scriptName)
	}

	val, err = script.Run(opts, func(vm *goja.Runtime) error {
		var mapData map[string]any
		if initValues != nil {
			mapData, err = initValues(vm)
			if err != nil {
				return err
			}
		}
		err = setRuntime(vm, opts, mapData)
		return err
	})
	return val, err
}

type Script struct {
	runtime     *runtime.Pool
	code        string
	codeType    string
	srcFileName string
	logger      logrus.FieldLogger
}

func NewScript(code string, codeType string, srcFileName string, userPool bool, logger logrus.FieldLogger) (*Script, error) {
	code, err := runtime.TransformTs(logger, code, srcFileName)
	if err != nil {
		return nil, err
	}
	return &Script{
		runtime:     runtime.NewPool(userPool),
		code:        code,
		codeType:    codeType,
		srcFileName: srcFileName,
		logger:      logger,
	}, nil
}

func (r *Script) GetLogger() logrus.FieldLogger {
	return r.logger
}

func (r *Script) Run(opts *RuntimeOption, init func(vm *goja.Runtime) error) (any, error) {
	data, err := r.runtime.Run(r.code, init)
	return data, err
}

type RuntimeOption struct {
	Server     *Server
	Service    *Service
	WebContext *WebContext
	Request    *Request
}

type InitVM interface {
	InitVM(vm *goja.Runtime) error
}

func setRuntime(vm *goja.Runtime, option *RuntimeOption, data ...map[string]any) error {
	if option != nil {
		if option.Server != nil {
			obj := vm.NewDynamicObject(NewServerProxy(option.Server, vm))
			_ = vm.Set("server", obj)
			option.Server.InitVM(vm)
		}
		if option.Service != nil {
			obj := vm.NewDynamicObject(NewServiceProxy(option.Service, vm))
			_ = vm.Set("service", obj)
			option.Service.InitVM(vm)
		}

		if option.WebContext != nil {
			_ = vm.Set("wctx", option.WebContext)
		}
	}
	addRuntimeValues(vm, data...)
	return nil
}

func addRuntimeValues(vm *goja.Runtime, data ...map[string]any) error {
	for _, d := range data {
		for k, v := range d {
			if init, ok := v.(InitVM); ok {
				if init.InitVM != nil {
					init.InitVM(vm)
				}
			}
			if obj, ok := v.(goja.DynamicObject); ok {
				_ = vm.Set(k, vm.NewDynamicObject(obj))
			} else {
				_ = vm.Set(k, v)
			}
		}
	}
	return nil
}
