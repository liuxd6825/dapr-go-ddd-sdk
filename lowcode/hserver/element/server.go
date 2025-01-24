package element

import (
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/definition"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type NewServerOptions func(server Server)

type Server interface {
	Base
	SchemaLoader() schema.URLLoader
	App() *iris.Application
	Restart() error
	FsPkg() FsPkg
	Definition() *definition.Definition
	RootPath() string
	CacheEnable() bool
	SrcFs() afero.Fs
	ReadSrcFile(fileName string, opts ...*fsopts.Options) ([]byte, error)
	ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error)
	Logs(level logrus.Level, format string, args ...interface{})
	NewSchemaCompiler() *jsonschema.Compiler
	Factory() Factory
	InitVM(vm *goja.Runtime) error
	SetSelfVMValue(name string, vm *goja.Runtime) error
	SetRunValues(data map[string]any)
	Start() error
}
