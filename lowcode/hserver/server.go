package hserver

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris/v12"
	common "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/sirupsen/logrus"
)

// Server 表示顶层结构
type Server struct {
	app       *iris.Application
	fs        *FsManager
	services  cmap.ConcurrentMap
	envConfig common.IEnvConfig
	tpl       *Template
	srcPath   string
	data      cmap.ConcurrentMap
	console   any
}

// NewServer 解析 HTML 并返回 Server 对象
func NewServer(app *iris.Application, mainFileName string, envConfig common.IEnvConfig, data map[string]any) (*Server, error) {
	fs, err := NewFsManger(envConfig)
	if err != nil {
		return nil, err
	}
	htmlData := fs.ReadFile(mainFileName)
	console := newConsole(logrus.New())
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
	dataMap := cmap.New()
	for k, v := range data {
		dataMap.Set(k, v)
	}
	server := &Server{
		app:       app,
		data:      dataMap,
		services:  cmap.New(),
		fs:        fs,
		srcPath:   srcPath,
		envConfig: envConfig,
		console:   console,
	}

	tpl := NewTemplate(envConfig, server)
	server.tpl = tpl

	// 遍历所有 <service> 节点
	doc.Find("server services service").Each(func(i int, s *goquery.Selection) {
		// 获取 url 属性
		url, exists := s.Attr("url")

		if exists {
			fileUrl := fileutils.AbsPath(url, server.srcPath)
			data := fs.ReadFile(fileUrl)
			service, err := NewService(server, data, fileUrl, nil)
			if err != nil {
				panic(err)
			}
			server.services.Set(service.Name, service)
		}
	})

	for _, key := range server.services.Keys() {
		item, ok := server.services.Get(key)
		if ok {
			service := item.(*Service)
			for _, key := range service.GetRequestKeys() {
				r := service.GetRequest(key)
				if r != nil {
					r.init(server, service)
					url := r.AbsURl
					if url == "" {
						url = service.URL + r.URL
					}
					app.Handle(r.Type, url, r.Handle)
				}
			}
		}
	}

	return server, nil
}
