package ddd

import (
	"errors"
	"fmt"
	"sync"
)

// Aggregate 聚合根接口类
type Aggregate interface {
	GetTenantId() string
	GetAggId() string
	GetAggType() string
	GetAggVer() string
}

type agg struct {
	tenantId string
	aggId    string
	aggType  string
	aggVer   string
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

func NewAggregateEmpty(tenantId, aggId, aggType, aggVer string) Aggregate {
	return &agg{
		tenantId: tenantId,
		aggId:    aggId,
		aggType:  aggType,
		aggVer:   aggVer,
	}
}

func (a *agg) GetTenantId() string {
	return a.tenantId
}

func (a *agg) GetAggId() string {
	return a.aggId
}

func (a *agg) GetAggType() string {
	return a.aggType
}

func (a *agg) GetAggVer() string {
	return a.aggVer
}
