package base

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	element2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/sirupsen/logrus"
)

type Base struct {
	srcFileName string
	funcs       element2.FuncManager
	logger      logrus.FieldLogger
	initFunc    element2.Func
	fsOpts      *fsopts.Options
	reader      fs.Reader
	runValues   *types.CMap[any]
	pkg         *types.CMap[any]
}

func NewBase(srcFileName string, logger logrus.FieldLogger, reader fs.Reader, factory element2.Factory, fsOpts *fsopts.Options) (element2.Base, error) {
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

func (b *Base) RunInitScript(ctx context.Context, values *element2.ApiRunValues, opts ...element2.RunOptions) error {
	opts = append(opts, func(vm *goja.Runtime) error {
		_ = vm.Set("ctx", context.Background())
		return nil
	})
	_, err := b.funcs.Run(ctx, "init", values, false, opts...)
	return err
}

func (b *Base) RunOnce(ctx context.Context, funcName, code, codeType, fileName string, runValues *element2.ApiRunValues, logger logrus.FieldLogger, pkg *types.CMap[any], opts ...element2.RunOptions) error {
	var err error
	if code != "" {
		_, err = b.funcs.Run(ctx, funcName, runValues, false, opts...)
	}
	return err
}

func (b *Base) ParseFunc(sel *goquery.Selection, server element2.Server) error {
	funcConfigs := make([]*element2.FuncConfig, 0)
	tags := make([]*element2.FuncTag, 0)
	sel.Children().Each(func(i int, sel *goquery.Selection) {
		node := sel.Get(0) // 获取当前节点的 *html.Node
		if node.Type != 3 {
			return
		}

		switch node.Data {
		case "script":
			{
				funConfig, isFound, err := element2.NewFuncConfig(sel, b.SrcFileName(), tags)
				tags = make([]*element2.FuncTag, 0)
				if err != nil {
					panic(err)
				}
				if !isFound {
					panic("invalid script tag")
				}
				funcConfigs = append(funcConfigs, funConfig)
			}
		default:
			tag := element2.NewFuncTag(node)
			tags = append(tags, tag)
		}
	})

	for _, cfg := range funcConfigs {
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

func (b *Base) Funcs() element2.FuncManager {
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
