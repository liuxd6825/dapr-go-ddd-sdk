package ddd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dapr/dapr-go-ddd-sdk/assert"
	"github.com/dapr/dapr-go-ddd-sdk/daprclient"
	"github.com/dapr/dapr-go-ddd-sdk/ddd/ddd_utils"
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

func (s *httpEventStorage) LoadEvent(ctx context.Context, req *daprclient.LoadEventsRequest) (*daprclient.LoadEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStorage) CreateEvent(ctx context.Context, req *daprclient.CreateEventRequest) (*daprclient.CreateEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStorage) DeleteEvent(ctx context.Context, req *daprclient.DeleteEventRequest) (*daprclient.DeleteEventResponse, error) {
	//TODO implement me
	panic("implement me")
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

func (s *httpEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error) {
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

	resp, err := s.LoadEvents(ctx, req)
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
				return nil, false, err
			}
		}
	}
	return aggregate, true, nil
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

func (s *httpEventStorage) ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (res *daprclient.ApplyEventResponse, resErr error) {
	url := fmt.Sprintf(ApiEventStorageEventApply)
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if req.Events == nil {
		return nil, errors.New("EventData cannot be null.")
	}
	for _, e := range req.Events {
		if err := ddd_utils.IsEmpty(e.CommandId, "CommandId"); err != nil {
			return nil, err
		}
		if err := ddd_utils.IsEmpty(e.PubsubName, "PubsubName"); err != nil {
			return nil, err
		}
		if err := ddd_utils.IsEmpty(e.EventType, "EventType"); err != nil {
			return nil, err
		}
		if err := ddd_utils.IsEmpty(e.EventId, "EventId"); err != nil {
			return nil, err
		}
		if err := ddd_utils.IsEmpty(e.EventVersion, "EventVersion"); err != nil {
			return nil, err
		}
		if err := ddd_utils.IsEmpty(e.Topic, "Topic"); err != nil {
			return nil, err
		}
		if len(e.PubsubName) == 0 {
			e.PubsubName = s.pubsubName
		}
	}

	data := &daprclient.ApplyEventResponse{}
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
