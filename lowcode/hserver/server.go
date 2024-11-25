package hserver

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/tpl_pkg"
	common "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
)

// Server 表示顶层结构
type Server struct {
	app       *iris.Application
	fs        *fs_pkg.FsManager
	services  cmap.ConcurrentMap
	envConfig common.IEnvConfig
	tpl       *tpl_pkg.Template
	srcPath   string
	runValues cmap.ConcurrentMap
	pkg       cmap.ConcurrentMap
	code      string
	codeType  string
	logger    logrus.FieldLogger
	scripts   *ScriptManager
}

type NewServerOptions func(server *Server)

// NewServer 解析 HTML 并返回 Server 对象
func NewServer(app *iris.Application, mainFileName string, envConfig common.IEnvConfig, opts ...NewServerOptions) (*Server, error) {
	fs, err := fs_pkg.NewFsManger(envConfig)
	if err != nil {
		return nil, err
	}
	htmlData := fs.ReadFile(mainFileName)

	reader := bytes.NewReader(htmlData)
	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %w", err)
	}

	srcPath, err := getFilePath(mainFileName)
	if err != nil {
		return nil, err
	}

	logger := logrus.StandardLogger()
	server := &Server{
		app:       app,
		runValues: cmap.New(),
		services:  cmap.New(),
		fs:        fs,
		srcPath:   srcPath,
		envConfig: envConfig,
		logger:    logger,
		scripts:   NewScriptManager(logger),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(server)
		}
	}

	if err := server.parse(doc); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) Run() error {
	err := s.runInitScript()
	if err != nil {
		return err
	}

	for _, key := range s.services.Keys() {
		item, ok := s.services.Get(key)
		if ok {
			service := item.(*Service)
			service.Run()

		}
	}
	return nil
}

func RunInitScript(scripts *ScriptManager, code, codeType, fileName string, opt *RuntimeOption, logger logrus.FieldLogger) error {
	if code != "" {
		err := scripts.AddScript("init", code, codeType, fileName, false, logger)
		if err != nil {
			return err
		}
		_, err = scripts.RunScript("init", opt, nil)
		return nil
	}
	return nil
}

func (s *Server) runInitScript() error {
	if s.code != "" {
		opts := &RuntimeOption{
			Server: s,
		}
		err := RunInitScript(s.scripts, s.code, s.codeType, "server_init.js", opts, s.GetLogger())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) parse(doc *goquery.Document) error {
	serverEl := doc.Find("body server")
	if serverEl == nil {
		return errors.New("no server")
	}
	serverEl.Find("script").Each(func(i int, sel *goquery.Selection) {
		s.code += sel.Text()
		s.codeType = sel.AttrOr("type", "")
	})

	// 遍历所有 <service> 节点
	serverEl.Find("services service").Each(func(i int, sel *goquery.Selection) {
		// 获取 url 属性
		url, exists := sel.Attr("url")

		if exists {
			fileUrl := fileutils.AbsPath(url, s.srcPath)
			data := s.fs.ReadFile(fileUrl)
			service, err := NewService(s, data, fileUrl, nil)
			if err != nil {
				panic(err)
			}
			s.services.Set(service.Name, service)
		}
	})

	return nil
}

func (s *Server) GetFs() *fs_pkg.FsManager {
	return s.fs
}

func (s *Server) GetApp() *iris.Application {
	return s.app
}

func (s *Server) GetEnvConfig() common.IEnvConfig {
	return s.envConfig
}

func (s *Server) GetRunValues() map[string]any {
	return s.runValues.Items()
}

func (s *Server) GetSrcPath() string {
	return s.srcPath
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

func (s *Server) SetRunValues(data map[string]any) {
	s.runValues.MSet(data)
}

func (s *Server) InitVM(vm *goja.Runtime) error {
	values := s.runValues.Items()
	addRuntimeValues(vm, values)
	return nil
}
