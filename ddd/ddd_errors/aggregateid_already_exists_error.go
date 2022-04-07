package ddd_errors

import "fmt"

type AggregateExistsError struct {
	AggregateId string
}

func NewAggregateIdExistsError(aggregateId string) *AggregateExistsError {
	return &AggregateExistsError{
		aggregateId,
	}
}

func (e *AggregateExistsError) Error() string {
	return fmt.Sprintf("aggregate root id %s  already exists.", e.AggregateId)
}

func IsErrorAggregateExists(err error) bool {
	switch err.(type) {
	case *AggregateExistsError:
		return true
	}
	return false
}
