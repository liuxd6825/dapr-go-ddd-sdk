package base

import (
	"context"
	"fmt"
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
	funcs       element.FuncManager
	logger      logrus.FieldLogger
	initFunc    element.Func
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
		funcs:       factory.NewFuncManager(logger, reader),
		logger:      logger,
		srcFileName: srcFileName,
		fsOpts:      fsOpts,
		runValues:   types.NewCMap[any](),
		pkg:         types.NewCMap[any](),
	}, nil
}

func (b *Base) RunInitScript(values *element.ApiRunValues, opts ...element.RunOptions) error {
	opts = append(opts, func(vm *goja.Runtime) error {
		_ = vm.Set("ctx", context.Background())
		return nil
	})
	_, err := b.funcs.Run("init", values, false, opts...)
	return err
}

func (b *Base) RunOnce(funcName, code, codeType, fileName string, runValues *element.ApiRunValues, logger logrus.FieldLogger, pkg *types.CMap[any], opts ...element.RunOptions) error {
	var err error
	if code != "" {
		_, err = b.funcs.Run(funcName, runValues, false, opts...)
	}
	return err
}

func (b *Base) ParseFunc(sel *goquery.Selection, server element.Server) error {
	funcConfigs := make([]*element.FuncConfig, 0)
	tags := make([]*element.FuncTag, 0)
	sel.Children().Each(func(i int, sel *goquery.Selection) {
		node := sel.Get(0) // 获取当前节点的 *html.Node
		if node.Type != 3 {
			return
		}

		switch node.Data {
		case "script":
			{
				funConfig, isFound, err := element.NewFuncConfig(sel, b.SrcFileName(), tags)
				tags = make([]*element.FuncTag, 0)
				if err != nil {
					panic(err)
				}
				if !isFound {
					panic("invalid script tag")
				}
				funcConfigs = append(funcConfigs, funConfig)
			}
		default:
			tag := element.NewFuncTag(node)
			tags = append(tags, tag)
		}
	})

	for _, cfg := range funcConfigs {
		fmt.Println("func = ", cfg.FuncName)
		fun, err := server.Factory().NewFunc(server, cfg, b.Logger(), server, b.Pkg())
		if err != nil {
			return err
		}
		if err = b.Funcs().Add(fun); err != nil {
			return err
		}
	}
	return nil
}

func (b *Base) SetPkg(pkg *types.CMap[any]) {
	b.pkg = pkg
}

func (b *Base) Pkg() *types.CMap[any] {
	return b.pkg
}

func (b *Base) AddPkg(key string, pkg any) {
	b.pkg.Set(key, pkg)
}

func (b *Base) FsOpts() *fsopts.Options {
	return b.fsOpts
}

func (b *Base) Funcs() element.FuncManager {
	return b.funcs
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
