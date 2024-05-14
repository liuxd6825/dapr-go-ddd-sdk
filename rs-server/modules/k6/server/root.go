package server

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/k6server/js/modules"
)

type RootModule struct {
	app *iris.Application
}

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &Exports{}
)

// New returns a pointer to a new RootModule instance.
func New(app *iris.Application) *RootModule {
	return &RootModule{app: app}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (m *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return NewExports(vu, m.app)
}
