package server

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/k6server/js/modules"
)

type Exports struct {
	vu  modules.VU
	app *iris.Application
}

// NewExports returns a new instance of Exports.
func NewExports(vu modules.VU, app *iris.Application) *Exports {
	return &Exports{vu: vu, app: app}
}

// Exports returns the exports of the k6 module.
func (e *Exports) Exports() modules.Exports {
	server := NewServer(e.app, e.vu)
	return modules.Exports{
		Named: map[string]interface{}{
			"server": e.vu.Runtime().ToValue(server),
			"object": e.vu.Runtime().ToValue(server.newObject),
			"time":   e.vu.Runtime().ToValue(NewTime()),
		},
	}
}
