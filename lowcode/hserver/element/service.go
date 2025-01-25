package element

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Service interface {
	Base
	ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error)
	Config() *ServerConfig
	Logger() logrus.FieldLogger
	SetSelfVMValue(name string, vm *goja.Runtime) error
	Initialize() error
	Close() error
}
