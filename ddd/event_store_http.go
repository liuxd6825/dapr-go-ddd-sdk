package ddd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
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

type httpEventStore struct {
	client     dapr.DaprClient
	pubsubName string
	subscribes []*Subscribe
}

func NewHttpEventStore(httpClient dapr.DaprClient, options ...func(s EventStore)) (EventStore, error) {
	subscribes = make([]*Subscribe, 0)
	res := &httpEventStore{
		client:     httpClient,
		subscribes: subscribes,
	}
	for _, option := range options {
		option(res)
	}
	return res, nil
}

func (s *httpEventStore) Commit(ctx context.Context, req *dapr.CommitRequest) (res *dapr.CommitResponse, resErr error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) Rollback(ctx context.Context, req *dapr.RollbackRequest) (res *dapr.RollbackResponse, resErr error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) GetEvents(ctx context.Context, req *dapr.GetEventsRequest) (*dapr.GetEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) GetRelations(ctx context.Context, req *dapr.GetRelationsRequest) (*dapr.GetRelationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) LoadEvent(ctx context.Context, req *dapr.LoadEventsRequest) (*dapr.LoadEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) CreateEvent(ctx context.Context, req *dapr.CreateEventRequest) (*dapr.CreateEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) DeleteEvent(ctx context.Context, req *dapr.DeleteEventRequest) (*dapr.DeleteEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) GetPubsubName() string {
	return s.pubsubName
}

func (s *httpEventStore) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate any) (Aggregate, bool, error) {
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
		TenantId:    tenantId,
		AggregateId: aggregateId,
	}
	agg := aggregate.(Aggregate)
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
	return agg, true, nil
}

func (s *httpEventStore) LoadEvents(ctx context.Context, req *dapr.LoadEventsRequest) (res *dapr.LoadEventsResponse, resErr error) {
	url := fmt.Sprintf(ApiEventStorageLoadEvents, req.TenantId, req.AggregateId)
	data := &dapr.LoadEventsResponse{}
	s.client.HttpGet(ctx, url).OnSuccess(data, func() error {
		res = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStore) ApplyEvent(ctx context.Context, req *dapr.ApplyEventRequest) (res *dapr.ApplyEventResponse, resErr error) {
	url := fmt.Sprintf(ApiEventStorageEventApply)
	if err := ddd_utils.IsEmpty(req.TenantId, "tenantId"); err != nil {
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

	data := &dapr.ApplyEventResponse{}
	s.client.HttpPost(ctx, url, req).OnSuccess(data, func() error {
		res = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStore) SaveSnapshot(ctx context.Context, req *dapr.SaveSnapshotRequest) (res *dapr.SaveSnapshotResponse, resErr error) {
	url := fmt.Sprintf(ApiEventStorageSnapshotSave)
	data := &dapr.SaveSnapshotResponse{}
	s.client.HttpPost(ctx, url, req).OnSuccess(data, func() error {
		res = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStore) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (isFind bool, resErr error) {
	url := fmt.Sprintf(ApiEventStorageExistAggregate, tenantId, aggregateId)
	data := &dapr.ExistAggregateResponse{}
	isFind = false
	s.client.HttpGet(ctx, url).OnSuccess(data, func() error {
		isFind = data.IsExist
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStore) getBodyBytes(resp *http.Response) ([]byte, error) {
	bytes, err := io.ReadAll(resp.Body)
	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if err != nil {
			err = e
		}
	}(resp.Body)

	return bytes, err
}
