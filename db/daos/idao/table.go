package idao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
)

type Table interface {
	GetTableName() string
	GetSchema() *store.DBSchema
	AutoMigrate(ctx context.Context)
	Exist(ctx context.Context) bool
	Drop(ctx context.Context)
}
