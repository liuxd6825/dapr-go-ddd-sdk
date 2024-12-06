package hserver

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
	"strings"
)

// Service 定义服务的基本结构
type Service struct {
	*Base
	server     *Server
	requests   *types.CMap[*Request]
	data       *types.CMap[any]
	config     ServerConfig
	initScript *ScriptConfig
}

type ServerConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// Call 定义调用信息结构
type Call struct {
	Ref    string            `json:"ref"`
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

func NewService(server *Server, html []byte, srcFileName string, data map[string]any) (*Service, error) {
	var err error
	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, server.GetRootPath())

	dataMap := types.NewCMap[any]()
	for k, v := range data {
		dataMap.Set(k, v)
	}

	// 解析 Service 节点
	service := &Service{
		server:   server,
		data:     dataMap,
		requests: types.NewCMap[*Request](),
		config:   ServerConfig{},
	}
	service.Base, err = NewBase(srcFileName, server.logger, service, fsOpts)
	if err != nil {
		return nil, err
	}

	service.SetPkg(server.GetPkg())

	if err := service.parse(html); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *Service) InitVM(vm *goja.Runtime) error {
	return nil
}

func (s *Service) AddRequest(r *Request) {
	s.requests.Set(r.config.Name, r)
}

func (s *Service) GetRequest(name string) *Request {
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
		panic(errors.New("service name is empty in %s", s.srcFileName))
	}

	var reqConfigs []RequestConfig

	err = s.ParseInitScript(serviceEl)
	if err != nil {
		return err
	}

	serviceEl.Children().Each(func(i int, request *goquery.Selection) {
		if err != nil {
			return
		}

		node := request.Get(0) // 获取当前节点的 *html.Node
		if node.Type != 3 {
			return
		}
		switch node.Data {
		case "request":
			// 解析 Request 列表
			reqCfg := RequestConfig{
				Type:        strings.ToUpper(request.AttrOr("type", "")),
				Name:        request.AttrOr("name", ""),
				URL:         request.AttrOr("url", ""),
				AbsURl:      request.AttrOr("abs-url", ""),
				Description: request.AttrOr("description", ""),
				ParamsType:  request.AttrOr("params-type", ""),
			}
			request.Find("link.params").First().Each(func(i int, selection *goquery.Selection) {
				reqCfg.LinkParamsUrl = selection.AttrOr("href", "")
			})
			// 获取 <script> 子节点
			scriptEl := request.Find("script")
			if scriptEl != nil {
				funcName := fmt.Sprintf("%s.%s()", s.config.Name, reqCfg.Name)
				scriptConfig, er := ParseScript(request, "script", funcName, s.srcFileName)
				if er != nil {
					err = er
					return
				}
				scriptConfig.UsePool = true
				reqCfg.Script = *scriptConfig
			}
			reqConfigs = append(reqConfigs, reqCfg)
		}

	})

	for _, cfg := range reqConfigs {
		request, err := NewRequest(s.server, s, s.GetSrcFileName(), cfg)
		if err != nil {
			return err
		}
		s.AddRequest(request)
	}
	return err
}

func (s *Service) Initialize() error {
	runValues := &RunValues{
		Server:   s.server,
		Self:     s,
		WorkPath: s.GetWorkPath(),
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

func (s *Service) GetLogger() logrus.FieldLogger {
	return s.server.logger
}

func (s *Service) SetSelfVMValue(name string, vm *goja.Runtime) error {
	obj := vm.NewDynamicObject(NewServiceProxy(s, vm))
	_ = vm.Set(name, obj)
	if err := s.InitVM(vm); err != nil {
		return err
	}
	return nil
}
