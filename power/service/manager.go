package service

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Manager struct {
	path     string
	services map[string]*Service
	opts     []ServiceOption
}

type Option struct {
}

func NewManager(path string, opts ...ServiceOption) *Manager {
	return &Manager{path: path, services: map[string]*Service{}, opts: opts}
}

func (s *Manager) Start() error {
	path, err := filepath.Abs(s.path)
	if err != nil {
		return err
	}
	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range fileList {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), ".service.js") {
				fileName := s.path + "/" + f.Name()
				service := NewService(fileName, s.opts...)
				if err := service.Run(); err != nil {
					return err
				}
				s.services[f.Name()] = service
			}
		}
	}
	return err
}

func (s *Manager) GetService(name string) *Service {
	return s.services[name]
}

func (s *Manager) SetService(name string, service *Service) {
	s.services[name] = service
}
