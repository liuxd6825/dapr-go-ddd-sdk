package service

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	element2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element/features/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element/features/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/sirupsen/logrus"
)

// Service 定义服务的基本结构

type Service struct {
	element2.Base
	selection *goquery.Selection
	server    element2.Server
	data      *types.CMap[any]
	config    *element2.ServiceConfig
}

func NewService(server element2.Server, sel *goquery.Selection, srcFileName string, data map[string]any) (element2.Service, error) {
	var err error
	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, server.RootPath())

	dataMap := types.NewCMap[any]()
	for k, v := range data {
		dataMap.Set(k, v)
	}

	// 解析 Service 节点
	service := &Service{
		server: server,
		data:   dataMap,
		config: &element2.ServiceConfig{},
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

func (s *Service) FindFunc(tagType string) []element2.Func {
	funcs := make([]element2.Func, 0)
	for _, fun := range s.Funcs().Items() {
		if fun.Config().ContainTag(tagType) {
			funcs = append(funcs, fun)
		}
	}
	return funcs
}

func (s *Service) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	data, err := s.server.ReadFile(filename, opts...)
	return data, err
}

func (s *Service) parse(serviceEl *goquery.Selection) error {
	s.selection = serviceEl
	s.config.Attrs = element2.NewAttrs(serviceEl.Get(0))
	if len(s.config.Name()) == 0 {
		panic(errors.New("service name is empty in %s", s.SrcFileName()))
	}

	return s.Base.ParseFunc(serviceEl, s.server)
}

func (s *Service) Initialize() error {
	runValues := &element2.ApiRunValues{
		Server:   s.server,
		Self:     s,
		WorkPath: s.WorkPath(),
	}
	ctx := context.Background()
	if err := s.RunInitScript(ctx, runValues); err != nil {
		return err
	}

	for _, fun := range s.Funcs().Items() {
		for _, tag := range fun.Config().Tags {
			s.register(tag, fun)
		}
	}
	return nil
}

// registerHandle 注册服务各类控制器
func (s *Service) register(tag *element2.FuncTag, fun element2.Func) {
	tagType := tag.GetType()
	switch tagType {
	case element2.FuncTag_DomainEvent:
		event.Register(s, fun, tag)
	case element2.FuncTag_RestApi:
		restapi.Register(s, fun, tag)
	}
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
	if errs.IsEmpty() {
		return nil
	}
	return errs
}

func (s *Service) Server() element2.Server {
	return s.server
}

func (s *Service) Config() *element2.ServiceConfig {
	return s.config
}

func (s *Service) InitVM(vm *goja.Runtime) error {
	return nil
}
