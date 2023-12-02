package ddd

import (
	"context"
	"fmt"
	dapr "github.com/liuxd6825/dapr-go-sdk/client"
)

type AggregateSnapshotActor struct {
	aggregateType string
	aggregateId   string
	SaveSnapshot  func(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
}

type SaveSnapshotRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateId   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`
	EventStoreKey string `json:"eventStoreKey"`
}

type SaveSnapshotResponse struct {
}

func NewAggregateSnapshotClient(client dapr.Client, aggregateType, aggregateId string) *AggregateSnapshotActor {
	actor := &AggregateSnapshotActor{
		aggregateType: aggregateType,
		aggregateId:   aggregateId,
	}
	client.ImplActorClientStub(actor)
	return actor
}

func (a *AggregateSnapshotActor) Type() string {
	return AggregateSnapshotActorType
}

func (a *AggregateSnapshotActor) ID() string {
	return fmt.Sprintf("aggType(%s),aggId(%s)", a.aggregateType, a.aggregateId)
}
