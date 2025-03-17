package ddd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/dapr"
)

type EventMarshal interface {
	Marshal(record *dapr.EventRecord) error
}
