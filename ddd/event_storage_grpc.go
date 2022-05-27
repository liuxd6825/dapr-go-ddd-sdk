package ddd

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"io"
	"net/http"
)

type grpcEventStorage struct {
	client     daprclient.DaprDddClient
	pubsubName string
	subscribes *[]Subscribe
}

func NewGrpcEventStorage(client daprclient.DaprDddClient, options ...func(s EventStorage)) (EventStorage, error) {
	subscribes = make([]Subscribe, 0)
	res := &grpcEventStorage{
		client:     client,
		subscribes: &subscribes,
	}
	for _, option := range options {
		option(res)
	}
	return res, nil
}

func (s *grpcEventStorage) GetPubsubName() string {
	return s.pubsubName
}

func (s *grpcEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error) {
	if err := assert.NotNil(aggregate, assert.NewOptions("aggregate is nil")); err != nil {
		return nil, false, err
	}

	if err := assert.NotEmpty(aggregateId, assert.NewOptions("aggregateId is nil")); err != nil {
		return nil, false, err
	}

	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is nil")); err != nil {
		return nil, false, err
	}

	req := &daprclient.LoadEventsRequest{
		TenantId:    tenantId,
		AggregateId: aggregateId,
	}

	resp, err := s.LoadEvent(ctx, req)
	if err != nil {
		return nil, false, err
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
<<<<<<< HEAD
				return res, find, err
			}
		}

		if len(records) >= 10 {
			snapshot := &daprclient.SaveSnapshotRequest{
				TenantId:          tenantId,
				AggregateData:     aggregate,
				AggregateId:       aggregate.GetAggregateId(),
				AggregateType:     aggregate.GetAggregateType(),
				AggregateRevision: aggregate.GetAggregateRevision(),
				SequenceNumber:    sequenceNumber,
			}
			_, err := s.SaveSnapshot(ctx, snapshot)
			if err != nil {
				return res, find, err
=======
				return nil, false, err
>>>>>>> actor_event_storage
			}
		}
	}
	return aggregate, true, err
}

func (s *grpcEventStorage) LoadEvent(ctx context.Context, req *daprclient.LoadEventsRequest) (res *daprclient.LoadEventsResponse, resErr error) {
	return s.client.LoadEvents(ctx, req)
}

func (s *grpcEventStorage) ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (*daprclient.ApplyEventResponse, error) {
	s.setEventsPubsubName(req.Events)
	return s.client.ApplyEvent(ctx, req)
}

func (s *grpcEventStorage) CreateEvent(ctx context.Context, req *daprclient.CreateEventRequest) (*daprclient.CreateEventResponse, error) {
	s.setEventsPubsubName(req.Events)
	return s.client.CreateEvent(ctx, req)
}

func (s *grpcEventStorage) DeleteEvent(ctx context.Context, req *daprclient.DeleteEventRequest) (*daprclient.DeleteEventResponse, error) {
	if req != nil && req.Event != nil && len(req.Event.PubsubName) == 0 {
		req.Event.PubsubName = s.pubsubName
	}
	return s.client.DeleteEvent(ctx, req)
}

func (s *grpcEventStorage) SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (*daprclient.SaveSnapshotResponse, error) {
	return s.client.SaveSnapshot(ctx, req)
}

func (s *grpcEventStorage) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	return bytes, err
}

func (s *grpcEventStorage) setEventsPubsubName(events []*daprclient.EventDto) {
	if events != nil {
		for _, event := range events {
			if len(event.PubsubName) == 0 {
				event.PubsubName = s.pubsubName
			}
		}
	}
}
