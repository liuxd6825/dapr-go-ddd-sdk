package ddd

import "context"

type QueryEventHandler interface {
	OnEvent(ctx context.Context, record *EventRecord) error
}
