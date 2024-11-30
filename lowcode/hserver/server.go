package hserver

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/definition"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/tpl_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/utils/schema_utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/xtype"
	common "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	iofs "io/fs"
)

// Server 表示顶层结构
type Server struct {
	*Base
	app          *iris.Application
	srcFs        afero.Fs
	fsm          *fs_pkg.FsManager
	services     *xtype.Map[*Service] // 服务Map
	envConfig    common.IEnvConfig    // 环境变量
	tpl          *tpl_pkg.Template    // 模板渲染服务
	definition   *definition.Definition
	cacheEnable  bool // 是否启用缓存
	schemaLoader schema.URLLoader
}

type NewServerOptions func(server *Server)

// NewServer 解析 HTML 并返回 Server 对象
func NewServer(app *iris.Application, srcFileName string, srcFs afero.Fs, env common.IEnvConfig, opts ...NewServerOptions) (*Server, error) {
	fsma, err := env.GetFsManager()
	if err != nil {
		return nil, err
	}
	fsm, err := fs_pkg.NewFsManger(env)
	if err != nil {
		return nil, err
	}
	//rootPath := env.GetRsServerSrcPath()
	//filename := rootPath + srcFileName
	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, "/")

	htmlData, err := fs.ReadFile(srcFs, srcFileName, fsOpts)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(htmlData)
	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %w", err)
	}

	logger := logrus.StandardLogger()
	server := &Server{
		app:          app,
		services:     xtype.NewMap[*Service](),
		fsm:          fsm,
		envConfig:    env,
		srcFs:        srcFs,
		cacheEnable:  true,
		schemaLoader: schema_utils.NewSchemaLoader(fsma),
	}
	server.definition, err = definition.NewDefinition(srcFs, "/definition/params")
	if err != nil {
		return nil, err
	}

	server.Base, err = NewBase(srcFileName, logger, server, fsOpts)
	if err != nil {
		return nil, err
	}

	if err = server.parse(doc); err != nil {
		return nil, err
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
			fmt.Println("File:", path)
		}
		return nil
	})

	return server, err
}

func (s *Server) GetCacheEnable() bool {
	return s.cacheEnable
}

func (s *Server) GetSrcFs() afero.Fs {
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
	data := s.fsm.ReadFile(filename, opts...)
	return data, nil
}

// Start
//
//	@Description: 启动服务
//	@receiver s
//	@return error
func (s *Server) Start() error {
	runValue := &RunValues{
		Server: s,
	}
	if err := s.RunInitScript(runValue); err != nil {
		return err
	}
	for _, val := range s.services.Items() {
		if service, ok := val.(*Service); ok {
			if err := service.Initialize(); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetService
//
//	@Description: 取得服务
//	@receiver s
//	@param serviceName
//	@return *Service
func (s *Server) GetService(serviceName string) *Service {
	service, ok := s.services.Get(serviceName)
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
	serverEl := doc.Find("body server")
	if serverEl == nil {
		return errors.New("no server")
	}

	err := s.ParseInitScript(serverEl)
	if err != nil {
		return err
	}

	// 遍历所有 <service> 节点
	serverEl.Find("services").Children().Each(func(i int, sel *goquery.Selection) {
		node := sel.Get(0)
		if node != nil && node.Type == 3 && node.Data == "link" {
			// 获取 url 属性
			fileUrl, exists := sel.Attr("href")
			if exists {
				data, err := s.ReadSrcFile(fileUrl, s.fsOpts)
				if err != nil {
					panic(err)
				}
				var service *Service
				service, err = NewService(s, data, fileUrl, nil)
				if err != nil {
					panic(err)
				}
				s.services.Set(service.config.Name, service)
			}
		}
	})

	return nil
}

// GetFsm
//
//	@Description: 取得文件管理器
//	@receiver s
//	@return *fs_pkg.FsManager
func (s *Server) GetFsm() *fs_pkg.FsManager {
	return s.fsm
}

// GetApp
//
//	@Description: 取iris.Application
//	@receiver s
//	@return *iris.Application
func (s *Server) GetApp() *iris.Application {
	return s.app
}

func (s *Server) GetEnvConfig() common.IEnvConfig {
	return s.envConfig
}

func (s *Server) GetRunValues() map[string]any {
	return s.runValues.Items()
}

func (s *Server) GetRootPath() string {
	return s.fsOpts.RootPath
}

func (s *Server) GetWorkPath() string {
	return s.fsOpts.WorkPath
}

func (s *Server) GetLogger() logrus.FieldLogger {
	return s.logger
}

// SetRunValue
//
//	@Description: 添加运行时变量
//	@receiver s
//	@param key
//	@param value
func (s *Server) SetRunValue(key string, value any) {
	s.runValues.Set(key, value)
}

// SetRunValues
//
//	@Description: 添加运行时变量
//	@receiver s
//	@param data
func (s *Server) SetRunValues(data map[string]any) {
	s.runValues.MSet(data)
}

func (s *Server) InitVM(vm *goja.Runtime) error {
	values := s.runValues.Items()
	if err := addRuntimeValues(vm, values); err != nil {
		return err
	}
	return nil
}
