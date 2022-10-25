package ddd

import (
	"context"
)

type sessionKey struct {
}

type sessionValue struct {
	eventItems []*eventItem
}

type eventItem struct {
	aggregate     Aggregate
	event         DomainEvent
	ctx           context.Context
	callEventType CallEventType
	opts          []*ApplyEventOptions
}

func newSession(parent context.Context) context.Context {
	ctx, _ := context.WithCancel(parent)
	return context.WithValue(ctx, sessionKey{}, &sessionValue{})
}

func getSessionValue(ctx context.Context) (*sessionValue, bool) {
	value := ctx.Value(sessionKey{})
	events, ok := value.(*sessionValue)
	return events, ok
}

func StartSession(ctx context.Context, do func(sessionCtx context.Context) error) error {
	var sessionCtx context.Context
	_, ok := getSessionValue(ctx)
	if !ok {
		sessionCtx = newSession(ctx)
	}
	if err := do(sessionCtx); err != nil {
		if batch, ok := getSessionValue(sessionCtx); ok {
			return batch.Commit()
		}
	}
	return nil
}

func (b *sessionValue) AddEvent(ctx context.Context, aggregate Aggregate, event DomainEvent, callEventType CallEventType, opts ...*ApplyEventOptions) {
	if b == nil {
		return
	}
	item := &eventItem{
		ctx:           ctx,
		aggregate:     aggregate,
		event:         event,
		callEventType: callEventType,
		opts:          opts,
	}
	b.eventItems = append(b.eventItems, item)
}

func (b *sessionValue) Commit() error {
	for _, item := range b.eventItems {
		_, err := callDaprEventMethod(item.ctx, item.callEventType, item.aggregate, item.event, item.opts...)
		if err != nil {
			return err
		}
	}
	return nil
}
