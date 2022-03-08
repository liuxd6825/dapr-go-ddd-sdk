package ddd_errors

import (
	"errors"
	"fmt"
)

func NewNoEventHandledError(eventType string) error {
	return errors.New(fmt.Sprintf("“%s” no domain event type handled.", eventType))
}
