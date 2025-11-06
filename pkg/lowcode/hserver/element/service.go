package element

import (
	"context"

	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type Attrs map[string]string

type ServiceConfig struct {
	Attrs Attrs `json:"attrs"`
}

type Service interface {
	Base
	Server() Server
	ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error)
	Config() *ServiceConfig
	Logger() logrus.FieldLogger
	SetSelfVMValue(name string, vm *goja.Runtime) error
	Initialize() error
	Close() error
}

type ApiRunValues struct {
	Server     Server
	Service    Service
	WebContext WebContext
	WorkPath   string
	Self       any
	Context    context.Context
}

func (c *ServiceConfig) Name() string {
	return c.Attrs.Get("name")
}

func (c *ServiceConfig) Url() string {
	return c.Attrs.Get("url")
}

func (c *ServiceConfig) AppId() string {
	return c.Attrs.Get("app-id")
}

func NewAttrs(node *html.Node) Attrs {
	attrs := make(Attrs)
	for _, a := range node.Attr {
		attrs[a.Key] = a.Val
	}
	return attrs
}

func (a Attrs) Set(key string, value string) {
	a[key] = value
}

func (a Attrs) Get(key string) string {
	if v, ok := a[key]; !ok {
		return ""
	} else {
		return v
	}
}
func (a Attrs) GetOr(key string, def string) string {
	v, ok := a[key]
	if !ok || v == "" {
		return def
	}
	return v
}
