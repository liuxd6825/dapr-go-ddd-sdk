package db

import (
	"context"
	"github.com/liuxd6825/jsonschema/v6"
)

type Table interface {
	GetTableName() string
	GetSchema() *jsonschema.Schema
	AutoMigrate(ctx context.Context)
	Exist(ctx context.Context) bool
	Drop(ctx context.Context)
}
