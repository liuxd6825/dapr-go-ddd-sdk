package server

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/k6server/js/modules"
	"sync"
)

type RootModule struct {
	app    *iris.Application
	data   map[string]any
	envCfg common.IEnvConfig
}

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &Exports{}
)

// New returns a pointer to a new RootModule instance.
func New(app *iris.Application, data map[string]any, envCfg common.IEnvConfig) *RootModule {
	return &RootModule{app: app, data: data, envCfg: envCfg}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (m *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return NewExports(vu, m.app, m.data, m.envCfg)
}

type Exports struct {
	vu     modules.VU
	app    *iris.Application
	data   map[string]any
	envCfg common.IEnvConfig
}

// NewExports returns a new instance of Exports.
func NewExports(vu modules.VU, app *iris.Application, data map[string]any, envCfg common.IEnvConfig) *Exports {
	return &Exports{vu: vu, app: app, data: data, envCfg: envCfg}
}

var _serverOnce sync.Once
var _server *Server
var _fs *FsManager
var _template *Template
var _feign *Feign
var _ctxPkg *ContextPkg

func newServer(app *iris.Application, vu modules.VU, cfg common.IEnvConfig) *Server {
	_serverOnce.Do(func() {
		_server = NewServer(app, vu, cfg)
		fsManager, err := NewFsManger(cfg)
		if err != nil {
			panic(err)
		}
		_fs = fsManager
		_template = NewTemplate(cfg)
		_feign = NewFeign(_server)
		_ctxPkg = NewContextPkg()
	})
	return _server
}

func GetServer() *Server {
	return _server
}

func GetFsManger() *FsManager {
	return _fs
}

func GetTemplate() *Template {
	return _template
}

// Exports returns the exports of the k6 module.
func (e *Exports) Exports() modules.Exports {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			err = errors.NewErr(err, "rs-server/modelds/server/Exports")
			panic(err)
		}
	}()

	newServer(e.app, e.vu, e.envCfg)
	runtime := e.vu.Runtime()
	events := NewEventPkg(runtime, e.envCfg)

	values := map[string]interface{}{
		"fmt":     runtime.ToValue(NewLogs()),
		"server":  runtime.ToValue(_server),
		"feign":   runtime.ToValue(_feign),
		"context": runtime.ToValue(_ctxPkg),
		"events":  runtime.ToValue(events),
		"fs":      runtime.ToValue(_fs),
		"tpl":     runtime.ToValue(_template),
	}

	for k, v := range e.data {
		values[k] = e.vu.Runtime().ToValue(v)
	}
	return modules.Exports{
		Named: values,
	}
}
