package ddd_repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type Repository interface {
	DoCreate(ctx context.Context, entity ddd.Entity) *SetResult
	DoUpdate(ctx context.Context, entity ddd.Entity) *SetResult
	DoFindById(ctx context.Context, tenantId string, id string) *FindResult
	DoFindAll(ctx context.Context, tenantId string) *FindResult
	DoDeleteById(ctx context.Context, tenantId string, id string) *SetResult
	DoSearch(ctx context.Context, search *SearchQuery) *FindResult
}

type Pageable struct {
	PageNumber int
	PageSize   int
}

type SearchQuery struct {
	TenantId string
	Fields   string
	Filter   string
	Sort     string
	Page     int
	Size     int
}

type Option struct {
}
