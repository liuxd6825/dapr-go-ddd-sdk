package server

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/k6server/js/modules"
	"sync"
)

type RootModule struct {
	app  *iris.Application
	data map[string]any
	cfg  common.IEnvConfig
}

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &Exports{}
)

// New returns a pointer to a new RootModule instance.
func New(app *iris.Application, data map[string]any, cfg common.IEnvConfig) *RootModule {
	return &RootModule{app: app, data: data, cfg: cfg}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (m *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return NewExports(vu, m.app, m.data, m.cfg)
}

type Exports struct {
	vu   modules.VU
	app  *iris.Application
	data map[string]any
	cfg  common.IEnvConfig
}

// NewExports returns a new instance of Exports.
func NewExports(vu modules.VU, app *iris.Application, data map[string]any, cfg common.IEnvConfig) *Exports {
	return &Exports{vu: vu, app: app, data: data, cfg: cfg}
}

var _serverOnce sync.Once
var _server *Server

func newServer(app *iris.Application, vu modules.VU, cfg common.IEnvConfig) *Server {
	_serverOnce.Do(func() {
		_server = NewServer(app, vu, cfg)
	})
	return _server
}

func GetServer() *Server {
	return _server
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
	server := newServer(e.app, e.vu, e.cfg)
	feign := NewFeign(server)
	ctxPkg := NewContextPkg()
	runtime := e.vu.Runtime()
	events := NewEventPkg(runtime, e.cfg)
	values := map[string]interface{}{
		"fmt":     runtime.ToValue(NewLogs()),
		"server":  runtime.ToValue(server),
		"feign":   runtime.ToValue(feign),
		"context": runtime.ToValue(ctxPkg),
		"events":  runtime.ToValue(events),
	}
	for k, v := range e.data {
		values[k] = e.vu.Runtime().ToValue(v)
	}
	return modules.Exports{
		Named: values,
	}
}
