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
	eventStorageKey *string
	wheres          map[string]string
	pageSize        *int
	pageNum         *int
	sort            *string
}

func NewWhereOptions() *WhereOptions {
	return &WhereOptions{wheres: make(map[string]string)}
}

func HasRelations(ctx context.Context, tenantId, aggregateType string, opts ...*WhereOptions) (bool, uint64, *daprclient.GetRelationsResponse, error) {
	options := NewWhereOptions()
	options.Merge(opts...)
	req := &daprclient.GetRelationsRequest{
		TenantId:      tenantId,
		AggregateType: aggregateType,
		Filter:        options.GetFilter(),
		PageNum:       uint64(options.GetPageNum()),
		PageSize:      uint64(options.GetPageSize()),
		Sort:          options.GetSort(),
	}

	resp, err := GetRelations(ctx, req, NewApplyCommandOptions().SetEventStorageKey(options.GetEventStorageKey()))
	return resp.IsFound, resp.TotalRows, resp, err
}

func HasAggregate(ctx context.Context, tenantId, aggregateType, aggregateId string) (bool, error) {
	options := NewWhereOptions().AddWhere("AggregateId", aggregateId)
	ok, _, _, err := HasRelations(ctx, tenantId, aggregateType, options)
	return ok, err
}

func GetRelations(ctx context.Context, req *daprclient.GetRelationsRequest, opts ...*ApplyCommandOptions) (*daprclient.GetRelationsResponse, error) {
	opt := NewApplyCommandOptions().Merge(opts...)
	eventStorage, err := GetEventStore(opt.EventStorageKey)
	if err != nil {
		return nil, err
	}
	return eventStorage.GetRelations(ctx, req)
}

func (o *WhereOptions) GetFilter() string {
	filter := ""
	count := len(o.wheres)
	i := 0
	for idName, idValue := range o.wheres {
		filter = filter + fmt.Sprintf("%s==\"%s\"", idName, idValue)
		if i < count-1 {
			filter = filter + " and "
		}
		i++
	}
	return filter
}

func (o *WhereOptions) Merge(opts ...*WhereOptions) {
	if opts == nil {
		return
	}
	for _, item := range opts {
		if item == nil {
			continue
		}
		if item.wheres != nil {
			for k, v := range item.wheres {
				o.wheres[k] = v
			}
		}
		if item.pageSize != nil {
			o.pageSize = item.pageSize
		}
		if item.pageNum != nil {
			o.pageNum = item.pageNum
		}
		if item.sort != nil {
			o.sort = item.sort
		}
	}
}

func (o *WhereOptions) SetPageSize(v int) *WhereOptions {
	o.pageSize = &v
	return o
}

func (o *WhereOptions) GetPageSize() int {
	if o.pageSize == nil {
		return 20
	}
	return *o.pageSize
}

func (o *WhereOptions) SetPageNum(v int) *WhereOptions {
	o.pageNum = &v
	return o
}

func (o *WhereOptions) GetPageNum() int {
	if o.pageSize == nil {
		return 0
	}
	return *o.pageNum
}

func (o *WhereOptions) SetSort(v string) *WhereOptions {
	o.sort = &v
	return o
}

func (o *WhereOptions) GetSort() string {
	if o.sort == nil {
		return ""
	}
	return *o.sort
}

func (o *WhereOptions) SetEventStorageKey(v string) *WhereOptions {
	o.eventStorageKey = &v
	return o
}

func (o *WhereOptions) GetEventStorageKey() string {
	if o.eventStorageKey == nil {
		return ""
	}
	return *o.eventStorageKey
}

func (o *WhereOptions) AddWhere(idName, idValue string) *WhereOptions {
	o.wheres[idName] = idValue
	return o
}
