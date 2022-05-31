package ddd

import "github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"

type EventMarshal interface {
	Marshal(record *daprclient.EventRecord) error
}
