package null

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/k6server/js/modules"
)

type Exports struct {
	vu  modules.VU
	app *iris.Application
}

// NewExports returns a new instance of Exports.
func NewExports(vu modules.VU) *Exports {
	return &Exports{vu: vu}
}

// Exports returns the exports of the k6 module.
func (e *Exports) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{},
	}
}
