package ddd

import (
	"errors"
	"fmt"
	"sync"
)

// Aggregate 聚合根接口类
type Aggregate interface {
	GetTenantId() string
	GetAggregateId() string
	GetAggregateType() string
	GetAggregateVersion() string
}

type agg struct {
	tenantId         string
	aggregateId      string
	aggregateType    string
	aggregateVersion string
}

type AggregateFactory func() Aggregate

type AggregateTypes map[string]AggregateFactory

var aggregateTypes = AggregateTypes{}
var rw sync.RWMutex

func RegisterAggregateType(aggregateType string, fn AggregateFactory) {
	if aggregateType == "" {
		panic(errors.New("aggregateType is cannot be empty"))
	}
	if fn == nil {
		panic(errors.New("fn is cannot be nil"))
	}

	rw.Lock()
	defer func() {
		rw.Unlock()
	}()
	aggregateTypes[aggregateType] = fn

}

func NewAggregateByType(aggregateType string) (Aggregate, error) {
	if aggregateType == "" {
		return nil, errors.New("aggregateType is cannot be empty")
	}
	rw.RLocker()
	defer rw.RUnlock()

	fn := aggregateTypes[aggregateType]
	if fn == nil {
		return nil, errors.New(fmt.Sprintf("aggregateType %s not registered ", aggregateType))
	}
	return fn(), nil
}

func NewAggregateEmpty(tenantId, aggId, aggType, aggVersion string) Aggregate {
	return &agg{
		tenantId:         tenantId,
		aggregateId:      aggId,
		aggregateType:    aggType,
		aggregateVersion: aggVersion,
	}
}

func (a *agg) GetTenantId() string {
	return a.tenantId
}

func (a *agg) GetAggregateId() string {
	return a.aggregateId
}

func (a *agg) GetAggregateType() string {
	return a.aggregateType
}

func (a *agg) GetAggregateVersion() string {
	return a.aggregateVersion
}
