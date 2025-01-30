package element

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/sirupsen/logrus"
)

type ApiServerConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type ApiService interface {
	Base
	ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error)
	Config() *ApiServerConfig
	Logger() logrus.FieldLogger
	SetSelfVMValue(name string, vm *goja.Runtime) error
	Initialize() error
	Close() error
}

type ApiRequestConfig struct {
	Type          string        `json:"type"`
	Name          string        `json:"name"`
	URL           string        `json:"url"`
	AbsURl        string        `json:"absUri"`
	Description   string        `json:"description"`
	ParamsType    string        `json:"paramsType"`
	LinkParamsUrl string        `json:"linkParamsUrl"`
	Script        *ScriptConfig `json:"script"`
}

type ApiRequest interface {
	GetParamsValue(wctx WebContext) map[string]any
	Config() *ApiRequestConfig
	Initialize() error
	Close() error
}

type ApiRunValues struct {
	Server     Server
	ApiService ApiService
	ApiRequest ApiRequest
	WebContext WebContext
	WorkPath   string
	Self       any
}
