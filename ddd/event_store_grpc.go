package ddd

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"io"
	"net/http"
)

type grpcEventStore struct {
	client     dapr.DaprClient
	compName   string
	pubsubName string
	subscribes []*Subscribe
}

func NewGrpcEventStore(compName string, pubsubName string, client dapr.DaprClient, options ...func(s EventStore)) (EventStore, error) {
	subscribes = make([]*Subscribe, 0)
	res := &grpcEventStore{
		compName:   compName,
		pubsubName: pubsubName,
		client:     client,
		subscribes: subscribes,
	}
	for _, option := range options {
		option(res)
	}
	return res, nil
}

func (s *grpcEventStore) GetPubsubName() string {
	return s.pubsubName
}

func (s *grpcEventStore) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate any) (agg Aggregate, isFound bool, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()

	a, ok := aggregate.(Aggregate)
	if !ok {
		return nil, false, errors.New("agg is not ddd.Aggregate interface")
	}

	if err := assert.NotNil(aggregate, assert.NewOptions("agg is nil")); err != nil {
		return nil, false, err
	}

	if err := assert.NotEmpty(aggregateId, assert.NewOptions("aggregateId is nil")); err != nil {
		return nil, false, err
	}

	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is nil")); err != nil {
		return nil, false, err
	}

	req := &dapr.LoadEventsRequest{
		TenantId:      tenantId,
		AggregateType: a.GetAggregateType(),
		AggregateId:   aggregateId,
	}

	resp, err := s.LoadEvent(ctx, req)
	if err != nil {
		return nil, false, errors.New("grpcEventStore.LoadEvent() error:%s", err.Error())
	}
	if resp.Snapshot == nil && (resp.EventRecords == nil || len(*resp.EventRecords) == 0) {
		return nil, false, err
	}

	if resp.Snapshot != nil {
		bytes, err := json.Marshal(resp.Snapshot.AggregateData)
		if err != nil {
			return nil, false, err
		}
		err = json.Unmarshal(bytes, aggregate)
		if err != nil {
			return nil, false, err
		}
	}
	records := *resp.EventRecords
	if records != nil && len(records) > 0 {
		for _, record := range *resp.EventRecords {
			if err = CallEventHandler(ctx, aggregate, &record); err != nil {
				return nil, false, errors.New("CallEventHandler(agg, record) eventType:%v, error:%v", record.EventType, err.Error())
			}
		}
	}
	return a, true, err
}

func (s *grpcEventStore) LoadEvent(ctx context.Context, req *dapr.LoadEventsRequest) (res *dapr.LoadEventsResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	return s.client.LoadEvents(ctx, req)
}

func (s *grpcEventStore) ApplyEvent(ctx context.Context, req *dapr.ApplyEventRequest) (res *dapr.ApplyEventResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	req.CompName = s.compName
	s.setEventsPubsubName(req.Events)
	return s.client.ApplyEvent(ctx, req)
}

func (s *grpcEventStore) Commit(ctx context.Context, req *dapr.CommitRequest) (res *dapr.CommitResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	req.CompName = s.compName
	return s.client.CommitEvent(ctx, req)
}

func (s *grpcEventStore) Rollback(ctx context.Context, req *dapr.RollbackRequest) (res *dapr.RollbackResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	req.CompName = s.compName
	return s.client.RollbackEvent(ctx, req)
}

func (s *grpcEventStore) SaveSnapshot(ctx context.Context, req *dapr.SaveSnapshotRequest) (res *dapr.SaveSnapshotResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	req.CompName = s.compName
	return s.client.SaveSnapshot(ctx, req)
}

func (s *grpcEventStore) GetEvents(ctx context.Context, req *dapr.GetEventsRequest) (res *dapr.GetEventsResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	req.CompName = s.compName
	return s.client.GetEvents(ctx, req)
}

func (s *grpcEventStore) GetRelations(ctx context.Context, req *dapr.GetRelationsRequest) (res *dapr.GetRelationsResponse, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()
	req.CompName = s.compName
	return s.client.GetRelations(ctx, req)
}

func (s *grpcEventStore) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	return bytes, err
}

func (s *grpcEventStore) setEventsPubsubName(events []*dapr.EventDto) {
	if events != nil {
		for _, event := range events {
			if len(event.PubsubName) == 0 {
				event.PubsubName = s.pubsubName
			}
		}
	}
}
