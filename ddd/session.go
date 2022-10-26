package ddd

import (
	"context"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

type sessionKey struct {
}

type sessionValue struct {
	sessionId  string
	eventItems []*eventItem
}

type eventItem struct {
	aggregate     Aggregate
	event         DomainEvent
	ctx           context.Context
	callEventType CallEventType
	opts          []*ApplyEventOptions
}

func newSession(parent context.Context) (context.Context, *sessionValue) {
	ctx, _ := context.WithCancel(parent)
	value := newSessionValue()
	return context.WithValue(ctx, sessionKey{}, value), value
}

func newSessionValue() *sessionValue {
	value := &sessionValue{
		sessionId:  uuid.New().String(),
		eventItems: make([]*eventItem, 0),
	}
	return value
}

func getSessionValue(ctx context.Context) (*sessionValue, bool) {
	value := ctx.Value(sessionKey{})
	events, ok := value.(*sessionValue)
	return events, ok
}

func StartSession(ctx context.Context, do func(sessionCtx context.Context) error) (resErr error) {
	defer func() {
		if err := errors.GetRecoverError(recover()); err != nil {
			resErr = err
		}
	}()
	sessValue, ok := getSessionValue(ctx)
	if !ok {
		ctx, sessValue = newSession(ctx)
	}
	err := do(ctx)
	if err != nil {
		return sessValue.Rollback()
	}
	return sessValue.Commit()
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

func (b *sessionValue) Rollback() error {
	for _, item := range b.eventItems {
		_, err := callDaprEventMethod(item.ctx, item.callEventType, item.aggregate, item.event, item.opts...)
		if err != nil {
			return err
		}
	}
	return nil
}
