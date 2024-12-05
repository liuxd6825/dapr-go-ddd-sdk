package hserver

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/runtime"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Script struct {
	runtime *runtime.Pool
	logger  logrus.FieldLogger
	config  ScriptConfig
	fs      afero.Fs
}

type ScriptConfig struct {
	FuncName    string
	Code        string
	CodeType    string
	SrcFileName string
	UsePool     bool
	Alias       map[string]string
}

type ScriptManager struct {
	scriptMap cmap.ConcurrentMap
	logger    logrus.FieldLogger
	reader    fs.Reader
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

func NewScriptManager(logger logrus.FieldLogger, reader fs.Reader) *ScriptManager {
	return &ScriptManager{
		scriptMap: cmap.New(),
		logger:    logger,
		reader:    reader,
	}
}

func NewScript(config *ScriptConfig, logger logrus.FieldLogger, reader fs.Reader) (*Script, error) {
	var err error
	if config == nil {
		return nil, errors.New("config is nil")
	}

	codes, err := runtime.TransformTSCodeToJS(config.Code)
	if err != nil {
		return nil, err
	}
	config.Code = string(codes)

	return &Script{
		runtime: runtime.NewPool(config.UsePool, reader),
		config:  *config,
		logger:  logger,
	}, nil
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

func (b *ScriptManager) AddScript(config *ScriptConfig, logger logrus.FieldLogger) error {
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

	code, err := runtime.TransformTSCodeToJS(config.Code)
	if err != nil {
		return err
	}

	config.Code = string(code)

	script, err := NewScript(config, logger, b.reader)
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

	val, ok := b.scriptMap.Get(funcName)
	if !ok {
		if checkHave {
			return nil, fmt.Errorf("script %s not found", funcName)
		}
		return nil, nil
	}

	script, ok := val.(*Script)
	if !ok {
		return nil, fmt.Errorf("script %s is not of type *Script", funcName)
	}

	srcFileName = script.config.SrcFileName

	opts = append(opts, func(vm *goja.Runtime) error {
		return setRunValues(vm, runValues)
	})

	val, err = script.Run(opts...)
	return val, err
}

func (r *Script) GetLogger() logrus.FieldLogger {
	return r.logger
}

func (r *Script) Run(opts ...RunOptions) (res any, err error) {
	defer func() {
		err = RecoverError(err, recover())
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

	data, err := r.runtime.Run("\n"+r.config.Code, opts...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("run error:%v in %s %s", err, r.config.FuncName, r.config.SrcFileName))
	}
	r.logger.Printf("run %s return %v in %s", r.config.FuncName, data, r.config.SrcFileName)

	return data, err
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
