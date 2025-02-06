package element

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
)

type RunOptions = func(vm *goja.Runtime) error

// Base
// @Description: Element基础类
type Base interface {
	FsOpts() *fsopts.Options
	ParseFunc(sel *goquery.Selection, server Server) error
	Funcs() FuncManager
	Logger() logrus.FieldLogger
	SrcFileName() string
	RootPath() string
	WorkPath() string
	SetPkg(pkg *types.CMap[any])
	Pkg() *types.CMap[any]
	RunValues() *types.CMap[any]
	RunInitScript(values *ApiRunValues, opts ...RunOptions) error
	RunOnce(funcName, code, codeType, fileName string, runValues *ApiRunValues, logger logrus.FieldLogger, pkg *types.CMap[any], opts ...RunOptions) error
}
