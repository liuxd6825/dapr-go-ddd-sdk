package ddd

import (
	"context"
	"encoding/json"
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

func (s *grpcEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (res Aggregate, find bool, err error) {
	aggregateType := aggregate.GetAggregateType()

	req := &daprclient.LoadEventsRequest{
		TenantId:    tenantId,
		AggregateId: aggregateId,
	}
	find = false
	resp, err := s.LoadEvents(ctx, req)
	if err != nil {
		return res, find, err
	}
	if resp.Snapshot == nil && (resp.EventRecords == nil || len(*resp.EventRecords) == 0) {
		return res, find, err
	}

	if resp.Snapshot != nil {
		bytes, err := json.Marshal(resp.Snapshot.AggregateData)
		if err != nil {
			return res, find, err
		}
		err = json.Unmarshal(bytes, aggregate)
		if err != nil {
			return res, find, err
		}
	}
	records := *resp.EventRecords
	if records != nil && len(records) > 0 {
		sequenceNumber := uint64(0)
		for _, record := range *resp.EventRecords {
			sequenceNumber = record.SequenceNumber
			if err = CallEventHandler(ctx, aggregate, &record); err != nil {
				return res, find, err
			}
		}

		if len(records) >= 3 {
			snapshot := &daprclient.SaveSnapshotRequest{
				TenantId:          tenantId,
				AggregateData:     aggregate,
				AggregateId:       aggregateId,
				AggregateType:     aggregateType,
				AggregateRevision: aggregate.GetAggregateRevision(),
				SequenceNumber:    sequenceNumber,
			}
			_, err := s.SaveSnapshot(ctx, snapshot)
			if err != nil {
				return res, find, err
			}
		}
	}
	res = aggregate
	find = true
	return res, find, err
}

func (s *grpcEventStorage) LoadEvents(ctx context.Context, req *daprclient.LoadEventsRequest) (res *daprclient.LoadEventsResponse, resErr error) {
	return s.client.LoadEvents(ctx, req)
}

func (s *grpcEventStorage) ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (*daprclient.ApplyEventsResponse, error) {
	if len(req.PubsubName) == 0 {
		req.PubsubName = s.pubsubName
	}
	return s.client.ApplyEvent(ctx, req)
}

func (s *grpcEventStorage) SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (*daprclient.SaveSnapshotResponse, error) {
	return s.client.SaveSnapshot(ctx, req)
}

func (s *grpcEventStorage) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (isFind bool, resErr error) {
	return s.client.ExistAggregate(ctx, tenantId, aggregateId)
}

func (s *grpcEventStorage) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	return bytes, err
}
