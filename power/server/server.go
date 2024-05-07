package server

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/power/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/power/service/values"
	"github.com/liuxd6825/dapr-go-ddd-sdk/power/service/watcher"
)

type Server struct {
	path     string
	watcher  watcher.Watcher
	app      *iris.Application
	settings []Setting
}

type Setting service.ServiceOption

var Settings []Setting

func AddServerHandler(app *iris.Application, path string, setting ...Setting) error {
	server := NewServer(app, path, setting...)
	return server.Start()
}

func NewServer(app *iris.Application, path string, setting ...Setting) *Server {
	apiServer := &Server{app: app, path: path, settings: setting}
	return apiServer
}

func (s *Server) Start() error {
	manager := service.NewManager(s.path, func(service *service.Service) error {
		if s.settings != nil {
			for _, set := range s.settings {
				if err := set(service); err != nil {
					return err
				}
			}
		}
		if Settings != nil {
			for _, set := range Settings {
				if err := set(service); err != nil {
					return err
				}
			}
		}
		service.SetValue("console", values.Console())
		return service.SetValue("time", values.NewTime())
	})

	if err := manager.Start(); err != nil {
		return err
	}
	s.watcher = watcher.NewFileWatcher(s.path)
	err := s.watcher.Start(func(rootPath, fileName string, eventType watcher.EventType) error {
		reload := false
		switch eventType {
		case watcher.WriteEvent:
			reload = true
		}
		if reload {
			s.app.Logger().Info("reload: " + fileName)
			if service := manager.GetService(fileName); service != nil {
				if err := service.Restart(); err != nil {
					s.app.Logger().Error(err)
				}
			}
		}
		if reload {
			if err := s.app.RefreshRouter(); err != nil {
				s.app.Logger().Error(err)
			}
		}
		return nil
	})
	return err
}
