package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/k6server/js/modules"
)

type Exports struct {
	vu modules.VU
}

// NewExports returns a new instance of Exports.
func NewExports(vu modules.VU) *Exports {
	return &Exports{vu: vu}
}

// Exports returns the exports of the k6 module.
func (e *Exports) Exports() modules.Exports {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			panic(err)
		}
	}()
	return modules.Exports{
		Named: map[string]interface{}{
			"newSchema": e.vu.Runtime().ToValue(NewSchema),
		},
	}
}
