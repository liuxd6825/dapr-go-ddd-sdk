package ddd

import (
	"errors"
	"fmt"
)

// Aggregate 聚合根接口类
type Aggregate interface {
	GetTenantId() string
	GetAggregateId() string
	GetAggregateType() string
	GetAggregateVersion() string
}

type NewAggregateFunc func() Aggregate

type AggregateTypes map[string]NewAggregateFunc

var aggregateTypes = AggregateTypes{}

func RegisterAggregateType(aggregateType string, fn NewAggregateFunc) {
	if aggregateType == "" {
		panic(errors.New("aggregateType is cannot be empty"))
	}
	if fn == nil {
		panic(errors.New("fn is cannot be nil"))
	}
	if t := aggregateTypes[aggregateType]; t != nil {
		panic(errors.New(fmt.Sprintf("aggregateType %s already exists", aggregateType)))
	}
	aggregateTypes[aggregateType] = fn
}

func NewAggregate(aggregateType string) (Aggregate, error) {
	if aggregateType == "" {
		return nil, errors.New("aggregateType is cannot be empty")
	}
	fn := aggregateTypes[aggregateType]
	if fn == nil {
		return nil, errors.New(fmt.Sprintf("aggregateType %s not registered ", aggregateType))
	}
	return fn(), nil
}
