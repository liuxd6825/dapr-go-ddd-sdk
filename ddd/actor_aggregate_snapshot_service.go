package ddd

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"
	"log"
	"os"
)

const aggregateSnapshotActorType = "ddd.AggregateSnapshotActorType"

var logger = log.New(os.Stdout, "", 0)

type AggregateSnapshotActorService struct {
	actor.ServerImplBase
	daprClient dapr.Client
}

func (s *AggregateSnapshotActorService) Type() string {
	return aggregateSnapshotActorType
}

func (s *AggregateSnapshotActorService) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	logger.Println("AggregateSnapshotActorService.SaveSnapshot()")
	logger.Println(json.Marshal(req))

	if err := SaveSnapshot(ctx, req.TenantId, req.AggregateType, req.AggregateId, req.EventStorageKey); err != nil {
		return nil, err
	}
	return &SaveSnapshotResponse{}, nil
}

func NewAggregateSnapshotActorService(daprClient dapr.Client) *AggregateSnapshotActorService {
	return &AggregateSnapshotActorService{
		daprClient: daprClient,
	}
}
