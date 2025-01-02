package hserver

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Script struct {
	runCode     string
	runtime     *runtime.Pool
	logger      logrus.FieldLogger
	config      ScriptConfig
	fs          afero.Fs
	transType   runtime.TransformType //转换类型
	isBuildCode bool
}

type ScriptConfig struct {
	FuncName    string
	Code        string
	CodeType    string
	SrcFileName string
	UsePool     bool
	Alias       map[string]string
	TransType   runtime.TransformType //转换类型
}

type RunValues struct {
	Server     *Server
	Service    *Service
	WebContext *WebContext
	Request    *Request
	Alias      Alias
	WorkPath   string
	Self       any
}

type RunOptions = func(vm *goja.Runtime) error

type InitVM interface {
	InitVM(vm *goja.Runtime) error
}

func NewScript(config *ScriptConfig, logger logrus.FieldLogger, reader fs.Reader, pkg *types.CMap[any]) (*Script, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	script := &Script{
		runtime:   runtime.NewPool(config.UsePool, reader, pkg),
		config:    *config,
		logger:    logger,
		transType: config.TransType,
	}

	err := script.BuildCode()
	if err != nil {
		panic(err)
	}

	return script, nil
}

func ParseScript(parentEl *goquery.Selection, selector string, funcName string, srcFileName string) (*ScriptConfig, error) {
	var err error
	var cfg *ScriptConfig
	if parentEl == nil {
		return cfg, nil
	}
	scripts := parentEl.Find(selector).First()
	if scripts.Length() > 0 {
		cfg = &ScriptConfig{
			FuncName:    funcName,
			SrcFileName: srcFileName,
		}

		cfg.Code = scripts.Text()
		cfg.CodeType = scripts.AttrOr("type", "")
		alias := scripts.AttrOr("alias", "")
		cfg.Alias, err = ParseAlias(alias)
	}
	return cfg, err
}

func (r *Script) GetLogger() logrus.FieldLogger {
	return r.logger
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
