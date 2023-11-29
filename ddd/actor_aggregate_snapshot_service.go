package ddd

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-sdk/actor"
	dapr "github.com/liuxd6825/dapr-go-sdk/client"
	"log"
	"os"
)

const aggregateSnapshotActorType = "ddd.AggregateSnapshotActorType"

var logger = log.New(os.Stdout, "", 0)

type AggregateSnapshotActorService struct {
	actor.ServerImplBaseCtx
	daprClient dapr.Client
}

func NewAggregateSnapshotActorService(daprClient dapr.Client) *AggregateSnapshotActorService {
	return &AggregateSnapshotActorService{
		daprClient: daprClient,
	}
}

func (s *AggregateSnapshotActorService) Type() string {
	return aggregateSnapshotActorType
}

func (s *AggregateSnapshotActorService) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	logs.Debugf(ctx, "AggregateSnapshotActorService.SaveSnapshot() req=%v", func() any {
		bs, _ := json.Marshal(req)
		return string(bs)
	})

	if err := SaveSnapshot(ctx, req.TenantId, req.AggregateType, req.AggregateId, req.EventStorageKey); err != nil {
		return nil, err
	}
	return &SaveSnapshotResponse{}, nil
}
