package sub

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
	"strings"
)

type SubService struct {
	element.Base
	server element.Server
	events *types.CMap[element.SubEvent]
	data   *types.CMap[any]
	config *element.SubServiceConfig
}

func (s *SubService) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	data, err := s.server.ReadFile(filename, opts...)
	return data, err
}

func NewSubService(server element.Server, sel *goquery.Selection, srcFileName string, data map[string]any) (element.SubService, error) {
	var err error
	dataMap := types.NewCMap[any]()
	dataMap.MSet(data)

	service := &SubService{
		server: server,
		events: types.NewCMap[element.SubEvent](),
		data:   dataMap,
		config: &element.SubServiceConfig{},
	}
	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, server.RootPath())
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

func (s *SubService) Config() *element.SubServiceConfig {
	return s.config
}

func (s *SubService) parse(serviceEl *goquery.Selection) error {
	s.config.Name = serviceEl.AttrOr("name", "")
	s.config.Description = serviceEl.AttrOr("description", "")
	s.config.AppId = serviceEl.AttrOr("app-id", "")
	s.config.Pubsub = serviceEl.AttrOr("pubsub", "pubsub")

	if len(s.config.Name) == 0 {
		panic(errors.New("SubService name is empty in %s", s.SrcFileName()))
	}

	if len(s.config.AppId) == 0 {
		panic(errors.New("SubService appId is empty in %s", s.SrcFileName()))
	}

	var configs []element.SubEventConfig

	err := s.ParseInitScript(serviceEl)
	if err != nil {
		return err
	}

	serviceEl.Children().Each(func(i int, ev *goquery.Selection) {
		if err != nil {
			return
		}

		node := ev.Get(0) // 获取当前节点的 *html.Node
		if node.Type != 3 {
			return
		}
		switch node.Data {
		case "event":
			// 解析 Request 列表
			eventCfg := &element.SubEventConfig{
				Name:        ev.AttrOr("name", ""),
				Version:     strings.ToLower(ev.AttrOr("version", "v1.0")),
				Description: ev.AttrOr("description", ""),
				ParamsType:  ev.AttrOr("params-type", ""),
			}

			s.server.Logs(logrus.InfoLevel, "event service.parse() event %s.%s()", s.config.Name, eventCfg.Name)
			ev.Find("a.params").First().Each(func(i int, selection *goquery.Selection) {
				eventCfg.LinkParamsUrl = selection.AttrOr("href", "")
			})

			// 获取 <script> 子节点
			funcName := fmt.Sprintf("%s.%s()", s.config.Name, eventCfg.Name)
			script, _, er := element.GetScriptConfig(ev, funcName, s.SrcFileName())
			if er != nil {
				err = er
				return
			}
			if script != nil {
				script.UsePool = true
				eventCfg.Script = script
				configs = append(configs, *eventCfg)
			} else {
				//panic(errors.New(fmt.Sprintf("script %s is empty in %s", script, s.SrcFileName())))
			}

		}
	})

	for _, cfg := range configs {
		request, err := s.server.Factory().NewSubEvent(s.server, s, s.SrcFileName(), &cfg)
		if err != nil {
			return err
		}
		if err = s.AddEvent(request); err != nil {
			return err
		}
	}
	return err
}

func (s *SubService) AddEvent(event element.SubEvent) error {
	s.events.Set(event.Config().Name, event)
	return nil
}

func (s *SubService) Initialize() error {
	runValues := &element.ApiRunValues{
		Server:   s.server,
		Self:     s,
		WorkPath: s.WorkPath(),
	}
	if err := s.RunInitScript(runValues); err != nil {
		return err
	}

	for _, key := range s.events.Keys() {
		r := s.GetEvent(key)
		if r != nil {
			err := r.Initialize()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SubService) GetEvent(name string) element.SubEvent {
	r, ok := s.events.Get(name)
	if !ok {
		return nil
	}
	return r
}
