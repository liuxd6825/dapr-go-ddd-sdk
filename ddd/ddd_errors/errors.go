package ddd_errors

import (
	"fmt"
)

type NoEventHandledError struct {
	eventType string
}

func NewNoEventHandledError(eventType string) error {
	return &NoEventHandledError{
		eventType: eventType,
	}
}

func (e *NoEventHandledError) Error() string {
	return fmt.Sprintf("“%s” no domain event type handled.", e.eventType)
}
