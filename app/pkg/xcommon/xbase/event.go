package xbase

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
)

var _outboxDao events.EventPublisher
var _outboxOnce sync.Once

func newEventPublisher() events.EventPublisher {
	_outboxOnce.Do(func() {
		_outboxDao = dao.NewOutboxEventDao(config.DBKey)
	})
	return _outboxDao
}

func PublishEvent(ctx context.Context, data events.IEvent) error {
	publisher := newEventPublisher()
	return events.Publish(ctx, publisher, data)
}

func PublishEvents(ctx context.Context, data []events.IEvent) error {
	publisher := newEventPublisher()
	return events.PublishList(ctx, publisher, data)
}
