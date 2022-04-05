package ddd_errors

import "fmt"

type AggregateIdNotFondError struct {
	AggregateId string
}

func NewAggregateIdNotFondError(aggregateId string) *AggregateIdNotFondError {
	return &AggregateIdNotFondError{
		aggregateId,
	}
}

func (e *AggregateIdNotFondError) Error() string {
	return fmt.Sprintf("aggregate root id %s not fond error", e.AggregateId)
}
