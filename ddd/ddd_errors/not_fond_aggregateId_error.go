package ddd_errors

import "fmt"

type NotFondAggregateIdError struct {
	AggregateId string
}

func NewNotFondAggregateIdError(aggregateId string) *NotFondAggregateIdError {
	return &NotFondAggregateIdError{
		aggregateId,
	}
}

func (e *NotFondAggregateIdError) Error() string {
	return fmt.Sprintf("aggregate root id %s  already exists.", e.AggregateId)
}
