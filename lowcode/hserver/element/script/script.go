package script

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Script struct {
	runCode     string
	runtime     *runtime.Pool
	config      *element.ScriptConfig
	fs          afero.Fs
	transType   runtime.TransformType //转换类型
	isBuildCode bool
	logger      logrus.FieldLogger
}

type RunOptions = func(vm *goja.Runtime) error

type InitVM interface {
	InitVM(vm *goja.Runtime) error
}

func NewScript(config *element.ScriptConfig, logger logrus.FieldLogger, reader fs.Reader, pkg *types.CMap[any]) (element.Script, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	script := &Script{
		runtime:   runtime.NewPool(config.UsePool, reader, pkg),
		config:    config,
		logger:    logger,
		transType: config.TransType,
	}

	err := script.BuildCode()
	if err != nil {
		panic(err)
	}

	return script, nil
}

func (r *Script) Logger() logrus.FieldLogger {
	return r.logger
}

func (r *Script) Config() *element.ScriptConfig {
	return r.config
}

func (r *Script) BuildCode() error {
	if !r.isBuildCode {
		codeBytes, err := runtime.TransformCode(r.config.Code, r.config.SrcFileName, r.config.TransType)
		if err != nil {
			return err
		}
		r.isBuildCode = true
		r.runCode = string(codeBytes)
	}
	return nil
}

func (r *Script) Run(opts ...RunOptions) (res any, err error) {
	defer func() {
		err = utils.RecoverError(err, recover())
	}()

	if r != nil {
		opts = append(opts, func(vm *goja.Runtime) error {
			for oldName, newName := range r.config.Alias {
				value := vm.Get(oldName)
				if value != nil {
					_ = vm.Set(newName, value)
				}
			}
			return nil
		})
	}

	data, err := r.runtime.Run("\n"+r.runCode, opts...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("run error:%v in %s %s", err, r.config.FuncName, r.config.SrcFileName))
	}
	//r.logger.Printf("run %s return %v in %s", r.config.FuncName, data, r.config.SrcFileName)

	return data, err
}
