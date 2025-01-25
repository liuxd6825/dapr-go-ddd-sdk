package service

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/subevents"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
	"strings"
)

// Service 定义服务的基本结构
type Service struct {
	element.Base
	server     element.Server
	requests   *types.CMap[element.Request]
	data       *types.CMap[any]
	config     *element.ServerConfig
	initScript *element.ScriptConfig
	subEvents  *types.CMap[*subevents.SubEvents]
}

// Call 定义调用信息结构
type Call struct {
	Ref    string            `json:"ref"`
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

func NewService(server element.Server, html []byte, srcFileName string, data map[string]any) (element.Service, error) {
	var err error
	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, server.RootPath())

	dataMap := types.NewCMap[any]()
	for k, v := range data {
		dataMap.Set(k, v)
	}

	// 解析 Service 节点
	service := &Service{
		server:    server,
		data:      dataMap,
		requests:  types.NewCMap[element.Request](),
		config:    &element.ServerConfig{},
		subEvents: types.NewCMap[*subevents.SubEvents](),
	}
	service.Base, err = server.Factory().NewBase(srcFileName, server.Logger(), service, server.Factory(), fsOpts)
	if err != nil {
		return nil, err
	}

	service.SetPkg(server.Pkg())

	if err := service.parse(html); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *Service) Config() *element.ServerConfig {
	return s.config
}

func (s *Service) InitVM(vm *goja.Runtime) error {
	return nil
}

func (s *Service) AddRequest(r element.Request) {
	s.requests.Set(r.Config().Name, r)
}

func (s *Service) GetRequest(name string) element.Request {
	r, ok := s.requests.Get(name)
	if !ok {
		return nil
	}
	return r
}

func (s *Service) GetRequestKeys() []string {
	return s.requests.Keys()
}

func (s *Service) GetRequestCount() int {
	return s.requests.Count()
}

func (s *Service) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	data, err := s.server.ReadFile(filename, opts...)
	return data, err
}

func (s *Service) parse(html []byte) error {
	reader := bytes.NewReader(html)
	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return err
	}
	serviceEl := doc.Find("body service")
	if serviceEl == nil {
		return errors.New("no service found")
	}

	s.config.URL = serviceEl.AttrOr("url", "")
	s.config.Name = serviceEl.AttrOr("name", "")
	s.config.Description = serviceEl.AttrOr("description", "")
	if len(s.config.Name) == 0 {
		panic(errors.New("service name is empty in %s", s.SrcFileName()))
	}

	var reqConfigs []element.RequestConfig

	err = s.ParseInitScript(serviceEl)
	if err != nil {
		return err
	}

	serviceEl.Children().Each(func(i int, req *goquery.Selection) {
		if err != nil {
			return
		}

		node := req.Get(0) // 获取当前节点的 *html.Node
		if node.Type != 3 {
			return
		}
		switch node.Data {
		case "request":
			// 解析 Request 列表
			reqCfg := element.RequestConfig{
				Type:        strings.ToUpper(req.AttrOr("type", "")),
				Name:        req.AttrOr("name", ""),
				URL:         req.AttrOr("url", ""),
				AbsURl:      req.AttrOr("abs-url", ""),
				Description: req.AttrOr("description", ""),
				ParamsType:  req.AttrOr("params-type", ""),
			}

			s.server.Logs(logrus.InfoLevel, "service.parse() request %s.%s()", s.config.Name, reqCfg.Name)

			req.Find("a.params").First().Each(func(i int, selection *goquery.Selection) {
				reqCfg.LinkParamsUrl = selection.AttrOr("href", "")
			})

			// 获取 <script> 子节点
			scriptEl := req.Find("script")
			if scriptEl != nil {
				funcName := fmt.Sprintf("%s.%s()", s.config.Name, reqCfg.Name)
				scriptConfig, er := element.ParseScriptConfig(req, "script", funcName, s.SrcFileName())
				if er != nil {
					err = er
					return
				}
				if scriptConfig != nil {
					scriptConfig.UsePool = true
					reqCfg.Script = scriptConfig
					reqConfigs = append(reqConfigs, reqCfg)
				}
			}
		}
	})

	for _, cfg := range reqConfigs {
		request, err := s.server.Factory().NewRequest(s.server, s, s.SrcFileName(), &cfg)
		if err != nil {
			return err
		}
		s.AddRequest(request)
	}
	return err
}

func (s *Service) Initialize() error {
	runValues := &element.RunValues{
		Server:   s.server,
		Self:     s,
		WorkPath: s.WorkPath(),
	}
	if err := s.RunInitScript(runValues); err != nil {
		return err
	}

	for _, key := range s.GetRequestKeys() {
		r := s.GetRequest(key)
		if r != nil {
			err := r.Initialize()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) Logger() logrus.FieldLogger {
	return s.server.Logger()
}

func (s *Service) SetSelfVMValue(name string, vm *goja.Runtime) error {
	obj := vm.NewDynamicObject(NewServiceProxy(s, vm))
	_ = vm.Set(name, obj)
	if err := s.InitVM(vm); err != nil {
		return err
	}
	return nil
}

func (s *Service) Close() error {
	errs := errors.NewErrors()
	for _, req := range s.requests.Items() {
		if err := req.Close(); err != nil {
			errs.AddError(err)
		}
	}
	if errs.IsEmpty() {
		return nil
	}
	return errs
}
