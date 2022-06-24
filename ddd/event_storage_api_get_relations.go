package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

type RelationsOptions struct {
	TenantId      string `json:"tenantId"`
	AggregateType string `json:"aggregateType"`
	Filter        string `json:"filter"`
	Sort          string `json:"sort"`
	PageNum       uint64 `json:"pageNum"`
	PageSize      uint64 `json:"pageSize"`
}

type WhereOptions struct {
	Wheres   map[string]string
	PageSize *int
	PageNum  *int
	Sort     *string
}

func GetRelationsWhere(ctx context.Context, tenantId, aggregateType string, opts ...WhereOptions) (*daprclient.GetRelationsResponse, error) {
	options := &WhereOptions{}
	options.Merge(opts...)
	req := &daprclient.GetRelationsRequest{
		TenantId:      tenantId,
		AggregateType: aggregateType,
		Filter:        options.GetFilter(),
		PageNum:       uint64(options.GetPageNum()),
		PageSize:      uint64(options.GetPageSize()),
		Sort:          options.GetSort(),
	}
	return GetRelations(ctx, "", req)
}

func GetRelations(ctx context.Context, eventStorageKey string, req *daprclient.GetRelationsRequest) (*daprclient.GetRelationsResponse, error) {
	eventStorage, err := GetEventStorage(eventStorageKey)
	if err != nil {
		return nil, err
	}
	return eventStorage.GetRelations(ctx, req)
}

func (o *WhereOptions) GetFilter() string {
	filter := ""
	count := len(o.Wheres)
	i := 0
	for idName, idValue := range o.Wheres {
		filter = filter + fmt.Sprintf("%s==\"%s\"", idName, idValue)
		if i < count-1 {
			filter = filter + " and "
		}
		i++
	}
	return filter
}

func (o *WhereOptions) Merge(opts ...WhereOptions) {
	if opts == nil {
		return
	}
	for _, item := range opts {
		if item.Wheres != nil {
			for k, v := range item.Wheres {
				o.Wheres[k] = v
			}
		}
		if item.PageSize != nil {
			o.PageSize = item.PageSize
		}
		if item.PageNum != nil {
			o.PageNum = item.PageNum
		}
		if item.Sort != nil {
			o.Sort = item.Sort
		}
	}
}

func (o *WhereOptions) SetPageSize(v int) *WhereOptions {
	o.PageSize = &v
	return o
}

func (o *WhereOptions) GetPageSize() int {
	if o.PageSize == nil {
		return 20
	}
	return *o.PageSize
}

func (o *WhereOptions) SetPageNum(v int) *WhereOptions {
	o.PageNum = &v
	return o
}

func (o *WhereOptions) GetPageNum() int {
	if o.PageSize == nil {
		return 0
	}
	return *o.PageNum
}

func (o *WhereOptions) SetSort(v string) *WhereOptions {
	o.Sort = &v
	return o
}

func (o *WhereOptions) GetSort() string {
	if o.Sort == nil {
		return ""
	}
	return *o.Sort
}
