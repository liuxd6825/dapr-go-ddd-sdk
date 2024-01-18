package ddd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr/actor"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type AggregateSnapshotActorServer struct {
	actor.ServerImplBaseCtx
	daprClient dapr.Client
}

type AggregateSnapshotActorClient struct {
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

type SaveSnapshotResponse = actor.Response

const AggregateSnapshotActorType = "ddd.AggregateSnapshotActorType"

func NewAggregateSnapshotClient(client dapr.Client, aggregateType, aggregateId string) *AggregateSnapshotActorClient {
	actor := &AggregateSnapshotActorClient{
		aggregateType: aggregateType,
		aggregateId:   aggregateId,
	}
	client.ImplActorClientStub(actor)
	return actor
}

func NewAggregateSnapshotActorServer(daprClient dapr.Client) actor.ServerContext {
	return &AggregateSnapshotActorServer{
		daprClient: daprClient,
	}
}

func (s *AggregateSnapshotActorServer) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (resp *SaveSnapshotResponse, err error) {
	resp = actor.NewResponse(nil)
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			resp.SetError(err)
			err = nil
		}
	}()
	logs.Debug(ctx, req.TenantId, logs.Fields{"funcName": "AggregateSnapshotActorServer.SaveSnapshot()", "aggregateType": req.AggregateType, "eventStoreKey": req.EventStoreKey, "request": func() any {
		bs, _ := json.Marshal(req)
		return string(bs)
	}})

	err = SaveSnapshot(ctx, req.TenantId, req.AggregateType, req.AggregateId, req.EventStoreKey)
	return resp, err
}

func (s *AggregateSnapshotActorServer) Type() string {
	return AggregateSnapshotActorType
}

func (a *AggregateSnapshotActorClient) Type() string {
	return AggregateSnapshotActorType
}

func (a *AggregateSnapshotActorClient) ID() string {
	return fmt.Sprintf("(aggType:%s,aggId:%s)", a.aggregateType, a.aggregateId)
}
