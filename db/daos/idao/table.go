package idao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
)

type Table interface {
	GetTableName() string
	GetSchema() *dbschema.DBSchema
	AutoMigrate(ctx context.Context)
	Exist(ctx context.Context) bool
	Drop(ctx context.Context)
}
