package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	common "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/definition"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element/funcs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/console_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/tpl_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils/schema_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	iofs "io/fs"
)

// Server 表示顶层结构
type Server struct {
	element.Base
	app          *iris.Application            // App应用实例
	httpServer   *restapp.HttpServer          // Http服务实例
	srcFs        afero.Fs                     // 源代码文件系统
	fsPkg        element.FsPkg                // 文件系统管理器
	services     *types.CMap[element.Service] // 服务Map
	envCfg       common.IEnvConfig            // 环境变量
	tpl          *tpl_pkg.Template            // 模板渲染服务
	definition   *definition.Definition       // 系统定义类
	cacheEnable  bool                         // 是否启用缓存
	schemaLoader schema.URLLoader             // schema加载器
	srcFileName  string                       // 源代码文件名
	factory      element.Factory              // 工厂类
	opts         []element.NewServerOptions   // 选项
	pkgSetup     PkgSetup                     // 包安装类
	eventPrefix  string                       // 事件前缀
	isPubEvent   bool                         // DAO是否发布事件的默认值
}

// NewServer 解析 HTML 并返回 Server 对象
func NewServer(httpServer *restapp.HttpServer, srcFileName string, srcFs afero.Fs, factory element.Factory, env common.IEnvConfig, opts ...element.NewServerOptions) (element.Server, error) {
	fsma := env.GetFsManager()

	fsPkg, err := fs_pkg.NewFsPkg(env, "")
	if err != nil {
		return nil, err
	}

	app := httpServer.App()
	logger := logrus.StandardLogger()
	server := &Server{
		srcFileName:  srcFileName,
		opts:         opts,
		httpServer:   httpServer,
		app:          app,
		services:     types.NewCMap[element.Service](),
		fsPkg:        fsPkg,
		envCfg:       env,
		srcFs:        srcFs,
		cacheEnable:  true,
		factory:      factory,
		schemaLoader: schema_utils.NewSchemaLoader(fsma),
		eventPrefix:  env.GetAppId(),
	}
	server.SetPkgSetup(NewPkgSetup(server))

	if err := server.init(logger); err != nil {
		return nil, err
	}
	return server, err
}

func (s *Server) init(logger logrus.FieldLogger) error {
	server := s
	srcFs := server.srcFs
	opts := server.opts

	var err error
	fsOpts := fsopts.NewOptionsWidthFileName(server.srcFileName, "/")

	server.definition, err = definition.NewDefinition(srcFs, "/definition/params")
	if err != nil {
		return err
	}

	server.Base, err = s.factory.NewBase(server.srcFileName, logger, server, server.factory, fsOpts)
	if err != nil {
		return err
	}

	if server.Pkg().Count() == 0 {
		server.LoadPkg("all")
	}

	for _, opt := range opts {
		if opt != nil {
			opt(server)
		}
	}

	// 使用 afero.Walk 遍历
	err = afero.Walk(srcFs, "/definition/params", func(path string, info iofs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 打印文件路径（忽略目录）
		if !info.IsDir() {
			//fmt.Println("File:", path)
		}
		return nil
	})

	logs := server.Logger()
	server.SetRunValue("console", console_pkg.NewConsole(logs))
	server.SetRunValue("logs", logs)
	return nil
}

// GetSchemaLoader
//
//	@Description: 取得Schema加载器，实现对引用schema文件加载
//	@receiver s
//	@return schema.URLLoader

func (s *Server) SchemaLoader() schema.URLLoader {
	return s.schemaLoader
}

func (s *Server) App() *iris.Application {
	return s.app
}

func (s *Server) FsPkg() element.FsPkg {
	return s.fsPkg
}

func (s *Server) Definition() *definition.Definition {
	return s.definition
}

func (s *Server) HttpServer() *restapp.HttpServer {
	return s.httpServer
}

// CacheEnable
//
//	@Description: 是否启用缓存模式
//	@receiver s
//	@return bool
func (s *Server) CacheEnable() bool {
	return s.cacheEnable
}

func (s *Server) Factory() element.Factory {
	return s.factory
}

// SrcFs
//
//	@Description: 取得源代码文件FS
//	@receiver s
//	@return afero.Fs
func (s *Server) SrcFs() afero.Fs {
	return s.srcFs
}

func (s *Server) ReadSrcFile(fileName string, opts ...*fsopts.Options) ([]byte, error) {
	return fs.ReadFile(s.srcFs, fileName, opts...)
}

// ReadFile
//
//	@Description: 读取文件内容
//	@receiver s
//	@param filename
//	@param opts
//	@return []byte
//	@return error
func (s *Server) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	data := s.fsPkg.ReadFile(filename, opts...)
	return data, nil
}

func (s *Server) Start() (err error) {
	defer func() {
		err = utils.RecoverError(err, recover())
		if err != nil {
			err = errors.New("Start : %s", err.Error())
		}
	}()
	err = s.start()
	return err
}

// Start
//
//	@Description: 启动服务
//	@receiver s
//	@return error
func (s *Server) start() error {
	htmlData, err := fs.ReadFile(s.srcFs, s.srcFileName, s.FsOpts())
	if err != nil {
		return err
	}

	reader := bytes.NewReader(htmlData)
	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return fmt.Errorf("error loading HTML: %w", err)
	}

	if err = s.parse(doc); err != nil {
		return err
	}

	runValue := &element.ApiRunValues{
		WorkPath: s.WorkPath(),
		Self:     s,
		Server:   s,
	}
	ctx := context.Background()
	if err := s.RunInitScript(ctx, runValue); err != nil {
		return err
	}

	if s.Pkg().Count() == 0 {
		s.LoadPkg("all")
	}

	for _, service := range s.services.Items() {
		if err := service.Initialize(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Restart() error {
	if err := s.init(s.Logger()); err != nil {
		return err
	}
	errs := errors.NewErrors()
	for _, service := range s.services.Items() {
		if e := service.Close(); e != nil {
			errs.AddError(e)
		}
	}
	if !errs.IsEmpty() {
		return errs
	}
	s.services.Clear()
	return s.Start()
}

// GetService
//
//	@Description: 取得服务
//	@receiver s
//	@param serviceName
//	@return *Service
func (s *Server) GetService(apiServiceName string) element.Service {
	service, ok := s.services.Get(apiServiceName)
	if ok {
		return service
	}
	return nil
}

// parse
//
//	@Description: html内容解析
//	@receiver s
//	@param doc
//	@return error
func (s *Server) parse(doc *goquery.Document) error {
	rootEl := doc.Find("body server")
	if rootEl == nil {
		return errors.New("no server")
	}

	if err := s.ParseFunc(rootEl, s); err != nil {
		return err
	}

	linksNodes := rootEl.Find(common.NodeType_Links)
	linksNodes.Each(func(i int, linkSel *goquery.Selection) {
		err := s.parseLinks(linkSel)
		if err != nil {
			panic(err)
		}
	})
	return nil
}

func (s *Server) parseLinks(linkSel *goquery.Selection) error {
	link := linkSel.Get(0)
	if link == nil {
		return nil
	}
	if _, exists := linkSel.Attr("close"); exists {
		return nil
	}

	linkSel.Children().Each(func(i int, sel *goquery.Selection) {
		node := sel.Get(0)
		if node != nil && node.Type == 3 && node.Data == common.NodeType_Link {
			// 获取 url 属性
			fileUrl, exists := sel.Attr("href")
			if exists {
				// s.Logs(logrus.InfoLevel, "service.parse() %s ", fileUrl)
				fileData, err := s.ReadSrcFile(fileUrl, s.FsOpts())
				if err != nil {
					panic(err)
				}

				fileReader := bytes.NewReader(fileData)
				// 使用 goquery 解析 HTML
				fileDoc, err := goquery.NewDocumentFromReader(fileReader)
				if err != nil {
					panic(err)
				}
				fileDoc.Find(common.NodeType_Service).Each(func(i int, serviceSel *goquery.Selection) {
					if err = s.addService(serviceSel, fileUrl); err != nil {
						panic(err)
					}
				})
			}
		}
	})
	return nil
}

func (s *Server) addService(sel *goquery.Selection, fileUrl string) error {
	apiService, err := s.factory.NewService(s, sel, fileUrl, nil)
	if err != nil {
		panic(err)
	}

	var config = apiService.Config()
	var name = config.Name()
	isHas := s.services.Has(name)
	if isHas {
		msg := fmt.Sprintf("addService %s already exists", name)
		return errors.New(msg)
	}
	s.services.Set(name, apiService)
	s.Logs(logrus.InfoLevel, "addService name=%s; file=%s;", name, fileUrl)
	return nil
}

func (s *Server) Logs(level logrus.Level, format string, args ...interface{}) {
	logger := s.Logger()
	switch level {
	case logrus.DebugLevel:
		logger.Debugf(format, args...)
	case logrus.InfoLevel:
		logger.Infof(format, args...)
	case logrus.WarnLevel:
		logger.Warnf(format, args...)
	case logrus.ErrorLevel:
		logger.Errorf(format, args...)
	case logrus.FatalLevel:
		logger.Fatalf(format, args...)
	case logrus.PanicLevel:
		logger.Panicf(format, args...)
	default:
		panic("unhandled default case")
	}
}

func (s *Server) EnvConfig() common.IEnvConfig {
	return s.envCfg
}

func (s *Server) RootPath() string {
	return s.FsOpts().RootPath
}

func (s *Server) WorkPath() string {
	return s.FsOpts().WorkPath
}

func (s *Server) Logger() logrus.FieldLogger {
	return s.Base.Logger()
}

// SetRunValue
//
//	@Description: 添加运行时变量
//	@receiver s
//	@param key
//	@param value
func (s *Server) SetRunValue(key string, value any) {
	s.RunValues().Set(key, value)
}

// SetRunValues
//
//	@Description: 添加运行时变量
//	@receiver s
//	@param data
func (s *Server) SetRunValues(data map[string]any) {
	s.RunValues().MSet(data)
}

func (s *Server) InitVM(vm *goja.Runtime) error {
	values := s.RunValues().Items()
	if err := funcs.AddRuntimeValues(vm, values); err != nil {
		return err
	}
	return nil
}

func (s *Server) SetSelfVMValue(name string, vm *goja.Runtime) error {
	obj := vm.NewDynamicObject(NewProxy(s, vm))
	_ = vm.Set(name, obj)
	if err := s.InitVM(vm); err != nil {
		return err
	}
	return nil
}

func (s *Server) NewSchemaCompiler() *jsonschema.Compiler {
	compiler := jsonschema.NewCompiler()
	compiler.UseLoader(s.schemaLoader)
	return compiler
}

func (s *Server) LoadPkg(names ...string) {
	s.pkgSetup.Setup(names...)
}

func (s *Server) GetPkgSetup() PkgSetup {
	return s.pkgSetup
}

func (s *Server) SetPkgSetup(setup PkgSetup) {
	s.pkgSetup = setup
}

// GetEventPrefix 事件前缀
func (s *Server) GetEventPrefix() string {
	return s.eventPrefix
}

// GetIsPubEvent 取DAO是否发布事件的默认值
func (s *Server) GetIsPubEvent() bool {
	return s.isPubEvent
}

func (s *Server) GetEnvCfg() common.IEnvConfig {
	return s.envCfg
}

func (s *Server) Init(opts *element.ServerInitOptions) {
	if opts == nil {
		return
	}
	if opts.IsPubEvent == nil {
		s.isPubEvent = *opts.IsPubEvent
	}
	if opts.EventPrefix == nil {
		s.eventPrefix = *opts.EventPrefix
	}
}
