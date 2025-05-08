package ddd

import "github.com/liuxd6825/dapr-go-ddd-sdk/dapr"

type EventMarshal interface {
	Marshal(record *dapr.EventRecord) error
}
