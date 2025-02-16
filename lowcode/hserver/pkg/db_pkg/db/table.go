package db

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
)

type Table interface {
	GetName() string
	GetSchema() *schema.Schema
	AutoMigrate(ctx context.Context)
	Exist(ctx context.Context) bool
	Drop(ctx context.Context)
}
