package common

import (
	"github.com/liuxd6825/k6server/js/modules"
)

type RootModule struct {
}

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &Exports{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (m *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Exports{
		vu: vu,
	}
}

type Exports struct {
	vu modules.VU
}

// Exports returns the exports of the k6 module.
func (e *Exports) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"context": e.vu.Runtime().ToValue(NewContextPkg()),
			"common":  e.vu.Runtime().ToValue(NewCommonPkg()),
		},
	}
}
