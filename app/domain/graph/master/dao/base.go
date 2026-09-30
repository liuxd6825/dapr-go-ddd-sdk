package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type Base[T any] struct {
	DBSchema *dbschema.DBSchema
	idao.Dao[T]
}
