package element

import (
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/definition"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type NewServerOptions func(server Server)

type ServerInitOptions struct {
	IsPubEvent  *bool
	EventPrefix *string
}

type Server interface {
	Base
	SchemaLoader() schema.URLLoader
	App() *iris.Application
	Restart() error
	FsPkg() fspkg.IFsPkg
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
	SetRunValues(data map[string]any) //设置运行时变量
	Start() error                     // 启动服务
	HttpServer() *restapp.HttpServer  // 取得HTTP服务实例
	GetEventPrefix() string           // 取得事件前缀
	GetIsPubEvent() bool              // 取得DAO更新数据时，是否发布事件的默认值
	GetEnvCfg() *env.Env              // 取当前配置环境变量
	Init(opts *ServerInitOptions)     // 初始化参数
}
