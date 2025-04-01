package idao

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"

type OutboxDao interface {
	Dao[*dbevent.Outbox]
}
