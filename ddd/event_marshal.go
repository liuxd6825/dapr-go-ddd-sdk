package ddd

import "github.com/dapr/dapr-go-ddd-sdk/daprclient"

type EventMarshal interface {
	Marshal(record *daprclient.EventRecord) error
}
