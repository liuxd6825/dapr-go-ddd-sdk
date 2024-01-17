package ddd

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-sdk/actor"
	dapr "github.com/liuxd6825/dapr-go-sdk/client"
)

const AggregateSnapshotActorType = "ddd.AggregateSnapshotActorType"

type AggregateSnapshotActorService struct {
	actor.ServerImplBaseCtx
	daprClient dapr.Client
}

func NewAggregateSnapshotActorServer(daprClient dapr.Client) actor.ServerContext {
	s := newAggregateSnapshotActorService(daprClient)
	return s
}

func newAggregateSnapshotActorService(daprClient dapr.Client) *AggregateSnapshotActorService {
	return &AggregateSnapshotActorService{
		daprClient: daprClient,
	}
}

func (s *AggregateSnapshotActorService) Type() string {
	return AggregateSnapshotActorType
}

func (s *AggregateSnapshotActorService) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (resp *SaveSnapshotResponse, err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	logs.Debug(ctx, req.TenantId, logs.Fields{"aggregateType": req.AggregateType, "eventStoreKey": req.EventStoreKey, "request": func() any {
		bs, _ := json.Marshal(req)
		return string(bs)
	}})

	if err := SaveSnapshot(ctx, req.TenantId, req.AggregateType, req.AggregateId, req.EventStoreKey); err != nil {
		return nil, err
	}
	return &SaveSnapshotResponse{}, nil
}
