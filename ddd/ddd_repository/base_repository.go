package ddd_repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type BaseRepository interface {
	BaseCreate(ctx context.Context, entity ddd.Entity, opts ...SetOption) *SetResult
	BaseUpdate(ctx context.Context, entity ddd.Entity, opts ...SetOption) *SetResult
	BaseFindById(ctx context.Context, tenantId string, id string, opts ...FindOption) *FindResult
}
