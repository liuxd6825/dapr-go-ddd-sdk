package hserver

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
)

type Base struct {
	scripts     *ScriptManager
	logger      logrus.FieldLogger
	runValues   cmap.ConcurrentMap
	initScript  *Script
	srcFileName string
	fsOpts      *fsopts.Options
	reader      fs.Reader
	pkg         any
}

func NewBase(srcFileName string, logger logrus.FieldLogger, reader fs.Reader, fsOpts *fsopts.Options) (*Base, error) {
	//rootPath, err := getFilePath(srcFileName)
	//if err != nil {
	/*		return nil, err
	}*/
	return &Base{
		reader:      reader,
		scripts:     NewScriptManager(logger, reader),
		logger:      logger,
		runValues:   cmap.New(),
		srcFileName: srcFileName,
		fsOpts:      fsOpts,
	}, nil
}

func (b *Base) GetSrcFileName() string {
	return b.srcFileName
}

func (b *Base) GetRootPath() string {
	return b.fsOpts.RootPath
}

func (b *Base) GetWorkPath() string {
	return b.fsOpts.WorkPath
}

func (b *Base) ParseInitScript(parentEl *goquery.Selection) error {
	var err error
	opts, err := ParseScript(parentEl, "script", "init()", b.srcFileName)
	if err != nil {
		return err
	}
	if opts != nil {
		opts.UsePool = false
		return b.scripts.AddScript(opts, b.logger, b.pkg)
	}
	return err
}

func (b *Base) RunInitScript(values *RunValues, opts ...RunOptions) error {
	opts = append(opts, func(vm *goja.Runtime) error {
		_ = vm.Set("ctx", context.Background())
		return nil
	})
	_, err := b.scripts.RunScript("init()", values, false, opts...)
	return err
}

func (b *Base) RunOnce(funcName, code, codeType, fileName string, runValues *RunValues, logger logrus.FieldLogger, pkg any, opts ...RunOptions) error {
	if code != "" {
		addOpts := &ScriptConfig{
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

func (b *Base) SetPkg(pkg any) {
	b.pkg = pkg
}

func (b *Base) GetPkg() any {
	return b.pkg
}
