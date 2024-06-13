package server

import (
	"github.com/dop251/goja"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/k6server/js/modules"
	"sync"
)

type RootModule struct {
	app  *iris.Application
	data map[string]any
}

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &Exports{}
)

// New returns a pointer to a new RootModule instance.
func New(app *iris.Application, data map[string]any) *RootModule {
	return &RootModule{app: app, data: data}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (m *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return NewExports(vu, m.app, m.data)
}

type Exports struct {
	vu   modules.VU
	app  *iris.Application
	data map[string]any
}

// NewExports returns a new instance of Exports.
func NewExports(vu modules.VU, app *iris.Application, data map[string]any) *Exports {
	return &Exports{vu: vu, app: app, data: data}
}

var _serverOnce sync.Once
var _server *Server

func newServer(app *iris.Application, vu modules.VU) *Server {
	_serverOnce.Do(func() {
		_server = NewServer(app, vu)
	})
	return _server
}

func GetServer() *Server {
	return _server
}

// Exports returns the exports of the k6 module.
func (e *Exports) Exports() modules.Exports {
	server := newServer(e.app, e.vu)
	feign := NewFeign(server)
	ctxPkg := NewContextPkg()
	values := map[string]interface{}{
		"fmt":             e.vu.Runtime().ToValue(NewLogs()),
		"server":          e.vu.Runtime().ToValue(server),
		"object":          e.vu.Runtime().ToValue(server.newObject),
		"feign":           e.vu.Runtime().ToValue(feign),
		"context":         e.vu.Runtime().ToValue(ctxPkg),
		"regEventHandler": e.vu.Runtime().ToValue(RegisterSubscribeService),
	}
	for k, v := range e.data {
		values[k] = e.vu.Runtime().ToValue(v)
	}
	return modules.Exports{
		Named: values,
	}
}

func (e *Exports) RegEventHandler(subscribes []*Subscribe, serviceObj *goja.Object, options ...*RegisterSubscribeOptions) {
	RegisterSubscribeService(subscribes, e.vu.Runtime(), serviceObj, options...)
}
