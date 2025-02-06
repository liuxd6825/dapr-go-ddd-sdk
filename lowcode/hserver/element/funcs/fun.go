package funcs

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

type Func struct {
	runCode     string
	runtime     *runtime.Pool
	config      *element.FuncConfig
	fs          afero.Fs
	transType   runtime.TransformType //转换类型
	isBuildCode bool
	logger      logrus.FieldLogger
}

type RunOptions = func(vm *goja.Runtime) error

type InitVM interface {
	InitVM(vm *goja.Runtime) error
}

func NewFunc(config *element.FuncConfig, logger logrus.FieldLogger, reader fs.Reader, pkg *types.CMap[any]) (element.Func, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	fun := &Func{
		runtime:   runtime.NewPool(config.UsePool, reader, pkg),
		config:    config,
		logger:    logger,
		transType: config.TransType,
	}

	err := fun.BuildCode()
	if err != nil {
		panic(err)
	}

	return fun, nil
}

func (r *Func) Logger() logrus.FieldLogger {
	return r.logger
}

func (r *Func) Config() *element.FuncConfig {
	return r.config
}

func (r *Func) BuildCode() error {
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

func (r *Func) Run(opts ...RunOptions) (res any, err error) {
	defer func() {
		err = utils.RecoverError(err, recover())
	}()

	if r != nil {
		opts = append(opts, func(vm *goja.Runtime) error {
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
