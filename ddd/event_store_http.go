package ddd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dapr2 "github.com/liuxd6825/dapr-go-ddd-sdk/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
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
	client     dapr2.DaprClient
	pubsubName string
	subscribes []*Subscribe
}

func NewHttpEventStore(httpClient dapr2.DaprClient, options ...func(s EventStore)) (EventStore, error) {
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

func (s *httpEventStore) Commit(ctx context.Context, req *dapr2.CommitRequest) (res *dapr2.CommitResponse, resErr error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) Rollback(ctx context.Context, req *dapr2.RollbackRequest) (res *dapr2.RollbackResponse, resErr error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) GetEvents(ctx context.Context, req *dapr2.GetEventsRequest) (*dapr2.GetEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) GetRelations(ctx context.Context, req *dapr2.GetRelationsRequest) (*dapr2.GetRelationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) LoadEvent(ctx context.Context, req *dapr2.LoadEventsRequest) (*dapr2.LoadEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) CreateEvent(ctx context.Context, req *dapr2.CreateEventRequest) (*dapr2.CreateEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) DeleteEvent(ctx context.Context, req *dapr2.DeleteEventRequest) (*dapr2.DeleteEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *httpEventStore) GetPubsubName() string {
	return s.pubsubName
}

func (s *httpEventStore) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate any) (Aggregate, bool, error) {
	if err := assert2.NotNil(aggregate, assert2.NewOptions("agg is nil")); err != nil {
		return nil, false, err
	}
	if err := assert2.NotEmpty(aggregateId, assert2.NewOptions("aggregateId is nil")); err != nil {
		return nil, false, err
	}
	if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is nil")); err != nil {
		return nil, false, err
	}

	req := &dapr2.LoadEventsRequest{
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

func (s *httpEventStore) LoadEvents(ctx context.Context, req *dapr2.LoadEventsRequest) (res *dapr2.LoadEventsResponse, resErr error) {
	url := fmt.Sprintf(ApiEventStorageLoadEvents, req.TenantId, req.AggregateId)
	data := &dapr2.LoadEventsResponse{}
	s.client.HttpGet(ctx, url).OnSuccess(data, func() error {
		res = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStore) ApplyEvent(ctx context.Context, req *dapr2.ApplyEventRequest) (res *dapr2.ApplyEventResponse, resErr error) {
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
		if err := ddd_utils.IsEmpty(e.CommandId, "EventId"); err != nil {
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

	data := &dapr2.ApplyEventResponse{}
	s.client.HttpPost(ctx, url, req).OnSuccess(data, func() error {
		res = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (s *httpEventStore) SaveSnapshot(ctx context.Context, req *dapr2.SaveSnapshotRequest) (res *dapr2.SaveSnapshotResponse, resErr error) {
	url := fmt.Sprintf(ApiEventStorageSnapshotSave)
	data := &dapr2.SaveSnapshotResponse{}
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
	data := &dapr2.ExistAggregateResponse{}
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
