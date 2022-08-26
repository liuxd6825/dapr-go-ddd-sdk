package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

type GetEventsOptions struct {
	TenantId      string `json:"tenantId"`
	AggregateType string `json:"aggregateType"`
	Filter        string `json:"filter"`
	Sort          string `json:"sort"`
	PageNum       uint64 `json:"pageNum"`
	PageSize      uint64 `json:"pageSize"`
}

type GetEventsWhereOptions struct {
	eventStorageKey *string
	wheres          map[string]string
	pageSize        *int
	pageNum         *int
	sort            *string
}

func NewGetEventsWhereOptions() *GetEventsWhereOptions {
	return &GetEventsWhereOptions{wheres: make(map[string]string)}
}

func GetEvents(ctx context.Context, req *daprclient.GetEventsRequest, opts ...*ApplyCommandOptions) (*daprclient.GetEventsResponse, error) {
	opt := NewApplyCommandOptions().Merge(opts...)
	eventStorage, err := GetEventStorage(opt.EventStorageKey)
	if err != nil {
		return nil, err
	}
	return eventStorage.GetEvents(ctx, req)
}

func (o *GetEventsWhereOptions) GetFilter() string {
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

func (o *GetEventsWhereOptions) Merge(opts ...*GetEventsWhereOptions) {
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

func (o *GetEventsWhereOptions) SetPageSize(v int) *GetEventsWhereOptions {
	o.pageSize = &v
	return o
}

func (o *GetEventsWhereOptions) GetPageSize() int {
	if o.pageSize == nil {
		return 20
	}
	return *o.pageSize
}

func (o *GetEventsWhereOptions) SetPageNum(v int) *GetEventsWhereOptions {
	o.pageNum = &v
	return o
}

func (o *GetEventsWhereOptions) GetPageNum() int {
	if o.pageSize == nil {
		return 0
	}
	return *o.pageNum
}

func (o *GetEventsWhereOptions) SetSort(v string) *GetEventsWhereOptions {
	o.sort = &v
	return o
}

func (o *GetEventsWhereOptions) GetSort() string {
	if o.sort == nil {
		return ""
	}
	return *o.sort
}

func (o *GetEventsWhereOptions) SetEventStorageKey(v string) *GetEventsWhereOptions {
	o.eventStorageKey = &v
	return o
}

func (o *GetEventsWhereOptions) GetEventStorageKey() string {
	if o.eventStorageKey == nil {
		return ""
	}
	return *o.eventStorageKey
}

func (o *GetEventsWhereOptions) AddWhere(idName, idValue string) *GetEventsWhereOptions {
	o.wheres[idName] = idValue
	return o
}
