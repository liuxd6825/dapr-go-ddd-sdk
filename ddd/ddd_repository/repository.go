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
	DoFindList(ctx context.Context, search *ListQuery) *FindResult
}

type Pageable struct {
	PageNumber int
	PageSize   int
}

type ListQuery struct {
	TenantId string
	Fields   string
	Filter   string // name=="lxd" and (id=="1" or id=="2")
	Sort     string // name:desc,id:asc
	Page     int64
	Size     int64
}

type Option struct {
}
