package api

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/subevents"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
	"strings"
)

// Service 定义服务的基本结构
type ApiService struct {
	element.Base
	server     element.Server
	requests   *types.CMap[element.ApiRequest]
	data       *types.CMap[any]
	config     *element.ApiServerConfig
	initScript *element.ScriptConfig
	subEvents  *types.CMap[*subevents.SubEvents]
}

// Call 定义调用信息结构
type Call struct {
	Ref    string            `json:"ref"`
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

func NewApiService(server element.Server, sel *goquery.Selection, srcFileName string, data map[string]any) (element.ApiService, error) {
	var err error
	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, server.RootPath())

	dataMap := types.NewCMap[any]()
	for k, v := range data {
		dataMap.Set(k, v)
	}

	// 解析 Service 节点
	service := &ApiService{
		server:    server,
		data:      dataMap,
		requests:  types.NewCMap[element.ApiRequest](),
		config:    &element.ApiServerConfig{},
		subEvents: types.NewCMap[*subevents.SubEvents](),
	}
	service.Base, err = server.Factory().NewBase(srcFileName, server.Logger(), service, server.Factory(), fsOpts)
	if err != nil {
		return nil, err
	}

	service.SetPkg(server.Pkg())

	if err := service.parse(sel); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ApiService) Config() *element.ApiServerConfig {
	return s.config
}

func (s *ApiService) InitVM(vm *goja.Runtime) error {
	return nil
}

func (s *ApiService) AddRequest(r element.ApiRequest) {
	s.requests.Set(r.Config().Name, r)
}

func (s *ApiService) GetRequest(name string) element.ApiRequest {
	r, ok := s.requests.Get(name)
	if !ok {
		return nil
	}
	return r
}

func (s *ApiService) GetRequestKeys() []string {
	return s.requests.Keys()
}

func (s *ApiService) GetRequestCount() int {
	return s.requests.Count()
}

func (s *ApiService) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	data, err := s.server.ReadFile(filename, opts...)
	return data, err
}

func (s *ApiService) parse(serviceEl *goquery.Selection) error {

	s.config.URL = serviceEl.AttrOr("url", "")
	s.config.Name = serviceEl.AttrOr("name", "")
	s.config.Description = serviceEl.AttrOr("description", "")
	if len(s.config.Name) == 0 {
		panic(errors.New("service name is empty in %s", s.SrcFileName()))
	}

	var reqConfigs []element.ApiRequestConfig

	err := s.ParseInitScript(serviceEl)
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
		case common.NodeType_Request:
			// 解析 Request 列表
			reqCfg := element.ApiRequestConfig{
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
				scriptConfig, _, er := element.GetScriptConfig(req, funcName, s.SrcFileName())
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
		request, err := s.server.Factory().NewApiRequest(s.server, s, s.SrcFileName(), &cfg)
		if err != nil {
			return err
		}
		s.AddRequest(request)
	}
	return err
}

func (s *ApiService) Initialize() error {
	runValues := &element.ApiRunValues{
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

func (s *ApiService) Logger() logrus.FieldLogger {
	return s.server.Logger()
}

func (s *ApiService) SetSelfVMValue(name string, vm *goja.Runtime) error {
	obj := vm.NewDynamicObject(NewServiceProxy(s, vm))
	_ = vm.Set(name, obj)
	if err := s.InitVM(vm); err != nil {
		return err
	}
	return nil
}

func (s *ApiService) Close() error {
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
