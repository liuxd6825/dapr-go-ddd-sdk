package ddd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	"io"
	"net/http"
)

const (
	ApiEventStorageEventApply     = "/v1.0/event-storage/events/apply"
	ApiEventStorageSnapshotSave   = "/v1.0/event-storage/snapshot/save"
	ApiEventStorageExistAggregate = "/v1.0/event-storage/aggregates/%s/%s"
	ApiEventStorageLoadEvents     = "/v1.0/event-storage/events/%s/%s"
)

type httpEventStorage struct {
	client     daprclient.DaprDddClient
	pubsubName string
	subscribes *[]Subscribe
}

func NewHttpEventStorage(httpClient daprclient.DaprDddClient, options ...func(s EventStorage)) (EventStorage, error) {
	subscribes = make([]Subscribe, 0)
	res := &httpEventStorage{
		client:     httpClient,
		subscribes: &subscribes,
	}
	for _, option := range options {
		option(res)
	}
	return res, nil
}

func (s *httpEventStorage) GetPubsubName() string {
	return s.pubsubName
}

func (s *httpEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggreageId string, aggregate Aggregate) (res Aggregate, find bool, err error) {
	aggregateId := aggregate.GetAggregateId()
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

func (s *httpEventStorage) LoadEvents(ctx context.Context, req *daprclient.LoadEventsRequest) (res *daprclient.LoadEventsResponse, resErr error) {
	url := fmt.Sprintf(ApiEventStorageLoadEvents, req.TenantId, req.AggregateId)
	data := &daprclient.LoadEventsResponse{}
	s.client.HttpGet(ctx, url).OnSuccess(data, func() error {
		res = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStorage) ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (res *daprclient.ApplyEventsResponse, resErr error) {
	if len(req.PubsubName) == 0 {
		req.PubsubName = s.pubsubName
	}
	url := fmt.Sprintf(ApiEventStorageEventApply)
	if err := ddd_utils.IsEmpty(req.CommandId, "CommandId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.PubsubName, "PubsubName"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.EventType, "EventType"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.EventId, "EventId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.EventRevision, "EventRevision"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.Topic, "Topic"); err != nil {
		return nil, err
	}
	if req.EventData == nil {
		return nil, errors.New("EventData cannot be null.")
	}

	data := &daprclient.ApplyEventsResponse{}
	s.client.HttpPost(ctx, url, req).OnSuccess(data, func() error {
		res = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStorage) SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (res *daprclient.SaveSnapshotResponse, resErr error) {
	url := fmt.Sprintf(ApiEventStorageSnapshotSave)
	data := &daprclient.SaveSnapshotResponse{}
	s.client.HttpPost(ctx, url, req).OnSuccess(data, func() error {
		res = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStorage) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (isFind bool, resErr error) {
	url := fmt.Sprintf(ApiEventStorageExistAggregate, tenantId, aggregateId)
	data := &daprclient.ExistAggregateResponse{}
	isFind = false
	s.client.HttpGet(ctx, url).OnSuccess(data, func() error {
		isFind = data.IsExist
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStorage) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	return bytes, err
}
