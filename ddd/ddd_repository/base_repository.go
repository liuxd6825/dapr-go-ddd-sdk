package ddd_repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type BaseRepository interface {
	BaseCreate(ctx context.Context, entity ddd.Entity) *SetResult
	BaseUpdate(ctx context.Context, entity ddd.Entity) *SetResult
	BaseFindById(ctx context.Context, tenantId string, id string) *FindResult
	BaseFindAll(ctx context.Context, tenantId string) *FindResult
	BaseDeleteById(ctx context.Context, tenantId string, id string) *SetResult
}
