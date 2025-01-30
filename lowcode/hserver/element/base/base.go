package base

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
)

type Base struct {
	srcFileName string
	scripts     element.ScriptManager
	logger      logrus.FieldLogger
	initScript  element.Script
	fsOpts      *fsopts.Options
	reader      fs.Reader
	runValues   *types.CMap[any]
	pkg         *types.CMap[any]
}

func NewBase(srcFileName string, logger logrus.FieldLogger, reader fs.Reader, factory element.Factory, fsOpts *fsopts.Options) (element.Base, error) {
	//rootPath, err := getFilePath(srcFileName)
	//if err != nil {
	/*		return nil, err
	}*/
	return &Base{
		reader:      reader,
		scripts:     factory.NewScriptManager(logger, reader),
		logger:      logger,
		runValues:   types.NewCMap[any](),
		srcFileName: srcFileName,
		fsOpts:      fsOpts,
		pkg:         types.NewCMap[any](),
	}, nil
}

func (b *Base) ParseInitScript(parentEl *goquery.Selection) error {
	var err error
	opts, _, err := element.GetScriptConfig(parentEl, "init()", b.srcFileName)
	if err != nil {
		return err
	}
	if opts != nil {
		opts.UsePool = false
		return b.scripts.AddScript(opts, b.logger, b.pkg)
	}
	return err
}

func (b *Base) RunInitScript(values *element.ApiRunValues, opts ...element.RunOptions) error {
	opts = append(opts, func(vm *goja.Runtime) error {
		_ = vm.Set("ctx", context.Background())
		return nil
	})
	_, err := b.scripts.RunScript("init()", values, false, opts...)
	return err
}

func (b *Base) RunOnce(funcName, code, codeType, fileName string, runValues *element.ApiRunValues, logger logrus.FieldLogger, pkg *types.CMap[any], opts ...element.RunOptions) error {
	if code != "" {
		addOpts := &element.ScriptConfig{
			FuncName:    funcName,
			Code:        code,
			CodeType:    codeType,
			SrcFileName: fileName,
			UsePool:     false,
		}
		err := b.scripts.AddScript(addOpts, logger, pkg)
		if err != nil {
			return err
		}
		_, err = b.scripts.RunScript(funcName, runValues, true, opts...)
		return nil
	}
	return nil
}

func (b *Base) SetPkg(pkg *types.CMap[any]) {
	b.pkg = pkg
}

func (b *Base) Pkg() *types.CMap[any] {
	return b.pkg
}

func (b *Base) FsOpts() *fsopts.Options {
	return b.fsOpts
}

func (b *Base) Scripts() element.ScriptManager {
	return b.scripts
}

func (b *Base) Logger() logrus.FieldLogger {
	return b.logger
}

func (b *Base) SrcFileName() string {
	return b.srcFileName
}
func (b *Base) RunValues() *types.CMap[any] {
	return b.runValues
}

func (b *Base) RunValuesMap() map[string]any {
	return b.runValues.Items()
}

func (b *Base) RootPath() string {
	return b.fsOpts.RootPath
}

func (b *Base) WorkPath() string {
	return b.fsOpts.WorkPath
}
