package server_pkg

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/definition"
	element2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/tpl_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/spf13/afero"
)

type IServer interface {
	App() *iris.Application
	SrcFs() afero.Fs
	FsPkg() fspkg.FsPkg
	Env() *env.Env
	Tpl() *tpl_pkg.Template             // 模板渲染服务
	Definition() *definition.Definition // 系统定义类
	CacheEnable() bool                  // 是否启用缓存
	GetSchemaLoader() schema.URLLoader  // schema加载器
	SrcFileName() string                // 源代码文件名
}

type Server struct {
	app          *iris.Application           // App应用实例
	srcFs        afero.Fs                    // 源代码文件系统
	fsPkg        fspkg.IFsPkg                // 文件系统管理器
	env          *env.Env                    // 环境变量
	tpl          *tpl_pkg.Template           // 模板渲染服务
	definition   *definition.Definition      // 系统定义类
	cacheEnable  bool                        // 是否启用缓存
	schemaLoader schema.URLLoader            // schema加载器
	srcFileName  string                      // 源代码文件名
	factory      element2.Factory            // 工厂类
	opts         []element2.NewServerOptions // 选项
}

func GetServer() *Server {
	return &Server{}
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) App() *iris.Application {
	return s.app
}

func (s *Server) SrcFs() afero.Fs {
	return s.srcFs
}

func (s *Server) FsPkg() fspkg.IFsPkg {
	return s.fsPkg
}

func (s *Server) Env() *env.Env {
	return s.env
}

func (s *Server) Tpl() *tpl_pkg.Template {
	return s.tpl
}

func (s *Server) Definition() *definition.Definition {
	return s.definition
}

func (s *Server) CacheEnable() bool {
	return s.cacheEnable
}

func (s *Server) SchemaLoader() schema.URLLoader {
	return s.schemaLoader
}

func (s *Server) SrcFileName() string {
	return s.srcFileName
}
