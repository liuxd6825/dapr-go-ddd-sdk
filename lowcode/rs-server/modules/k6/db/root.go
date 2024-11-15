package db

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/k6server/js/modules"
	"sync"
)

type RootModule struct {
	cfg common.IEnvConfig
}

// New returns a pointer to a new RootModule instance.
func New(cfg common.IEnvConfig) *RootModule {
	return &RootModule{
		cfg: cfg,
	}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (m *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Exports{
		vu:  vu,
		cfg: m.cfg,
	}
}

type Exports struct {
	vu  modules.VU
	cfg common.IEnvConfig
}

var _dbOnce sync.Once
var _db *DB

// Exports returns the exports of the k6 module.
func (e *Exports) Exports() modules.Exports {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			panic(err)
		}
	}()

	_dbOnce.Do(func() {
		_db = NewDB(e.cfg)
	})

	return modules.Exports{
		Named: map[string]interface{}{
			"db": e.vu.Runtime().ToValue(_db),
		},
	}
}
