package idao

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type Table interface {
	GetTableName() string
	GetSchema() *store.DBSchema
	AutoMigrate(ctx context.Context) error
	Exist(ctx context.Context) bool
	Drop(ctx context.Context) error
}
