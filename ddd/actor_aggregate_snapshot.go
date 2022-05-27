package ddd

import (
	"context"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
)

type AggregateSnapshotActor struct {
	aggregateType string
	aggregateId   string
	SaveSnapshot  func(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
}

func (a *AggregateSnapshotActor) Type() string {
	return aggregateSnapshotActorType
}

func (a *AggregateSnapshotActor) ID() string {
	return fmt.Sprintf("aggType(%s),aggId(%s)", a.aggregateType, a.aggregateId)
}

type SaveSnapshotRequest struct {
	TenantId        string `json:"tenantId"`
	AggregateId     string `json:"aggregateId"`
	AggregateType   string `json:"aggregateType"`
	EventStorageKey string `json:"eventStorageKey"`
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
