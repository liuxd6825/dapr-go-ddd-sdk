package null

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
	return NewExports(vu)
}
