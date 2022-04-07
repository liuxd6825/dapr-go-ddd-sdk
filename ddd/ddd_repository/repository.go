package ddd_repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type Repository interface {
	Insert(ctx context.Context, entity ddd.Entity) *SetResult
	Update(ctx context.Context, entity ddd.Entity) *SetResult
	DeleteById(ctx context.Context, tenantId string, id string) *SetResult
	FindById(ctx context.Context, tenantId string, id string) *FindResult
	FindAll(ctx context.Context, tenantId string) *FindResult
	FindPagingData(ctx context.Context, search *PagingQuery) *FindPagingResult

	DoFind(fun func() (interface{}, bool, error)) *FindResult
	DoSet(fun func() (interface{}, error)) *SetResult
	DoFindPagingData(fun func() (*PagingData, bool, error)) *FindPagingResult
}

type Pageable struct {
	PageNumber int
	PageSize   int
}

type PagingQuery struct {
	TenantId string
	Fields   string
	Filter   string // name=="lxd" and (id=="1" or id=="2")
	Sort     string // name:desc,id:asc
	Page     int64
	Size     int64
}

type Option struct {
}
