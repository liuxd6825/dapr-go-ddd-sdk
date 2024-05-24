package gorm

import (
	"github.com/liuxd6825/k6server/js/modules"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Exports struct {
	vu modules.VU
}

// Exports returns the exports of the k6 module.
func (e *Exports) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"gorm": e.vu.Runtime().ToValue(&orm{}),
		},
	}
}

type orm struct {
}

func (e *orm) Open(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
	return gorm.Open(dialector, opts...)
}

// Expr returns clause.Expr, which can be used to pass SQL expression as params
func (e *orm) Expr(expr string, args ...interface{}) clause.Expr {
	return gorm.Expr(expr, args...)
}

func (e *orm) NewConfig() *gorm.Config {
	return &gorm.Config{}
}
