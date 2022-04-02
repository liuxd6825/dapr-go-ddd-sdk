package ddd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/httpclient"
	"io"
	"net/http"
)

const (
	ApiEventStorageEventApply     = "/v1.0/event-storage/events/apply"
	ApiEventStorageSnapshotSave   = "/v1.0/event-storage/snapshot/save"
	ApiEventStorageExistAggregate = "/v1.0/event-storage/aggregates/%s/%s"
	ApiEventStorageLoadEvents     = "/v1.0/event-storage/events/%s/%s"
)

type daprEventStorage struct {
	httpClient *httpclient.HttpClient
	pubsubName string
	subscribes *[]Subscribe
}

func NewDaprEventStorage(httpClient *httpclient.HttpClient, options ...func(s EventStorage)) (EventStorage, error) {
	subscribes = make([]Subscribe, 0)
	res := &daprEventStorage{
		httpClient: httpClient,
		subscribes: &subscribes,
	}
	for _, option := range options {
		option(res)
	}
	return res, nil
}

/*func (s *daprEventStorage) GetHost() string {
	return s.host
}

func (s *daprEventStorage) GetPort() int {
	return s.port
}*/

func (s *daprEventStorage) GetPubsubName() string {
	return s.pubsubName
}

func (s *daprEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (res Aggregate, find bool, err error) {
	req := &LoadEventsRequest{
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
		sequenceNumber := int64(0)
		for _, record := range *resp.EventRecords {
			sequenceNumber = record.SequenceNumber
			if err = CallEventHandler(ctx, aggregate, &record); err != nil {
				return res, find, err
			}
		}

		if len(records) >= 3 {
			snapshot := &SaveSnapshotRequest{
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
			}
		}
	}
	res = aggregate
	find = true
	return res, find, err
}

func (s *daprEventStorage) LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error) {
	url := fmt.Sprintf(ApiEventStorageLoadEvents, req.TenantId, req.AggregateId)
	data, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	res := &LoadEventsResponse{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *daprEventStorage) ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error) {
	if len(req.PubsubName) == 0 {
		req.PubsubName = s.pubsubName
	}
	url := fmt.Sprintf(ApiEventStorageEventApply)
	if err := isEmpty(req.CommandId, "CommandId"); err != nil {
		return nil, err
	}
	if err := isEmpty(req.PubsubName, "PubsubName"); err != nil {
		return nil, err
	}
	if err := isEmpty(req.EventType, "EventType"); err != nil {
		return nil, err
	}
	if err := isEmpty(req.EventId, "EventId"); err != nil {
		return nil, err
	}
	if err := isEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := isEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := isEmpty(req.EventRevision, "EventRevision"); err != nil {
		return nil, err
	}
	if err := isEmpty(req.Topic, "Topic"); err != nil {
		return nil, err
	}
	if req.EventData == nil {
		return nil, errors.New("EventData cannot be null.")
	}

	data, err := s.httpClient.Post(url, req)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	res := &ApplyEventsResponse{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *daprEventStorage) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	url := fmt.Sprintf(ApiEventStorageSnapshotSave)
	data, err := s.httpClient.Post(url, req)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	res := &SaveSnapshotResponse{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *daprEventStorage) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error) {
	url := fmt.Sprintf(ApiEventStorageExistAggregate, tenantId, aggregateId)
	data, err := s.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	resp := &ExistAggregateResponse{}
	err = json.Unmarshal(data, resp)
	return resp.IsExist, err
}

func (s *daprEventStorage) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	return bytes, err
}

func isEmpty(v string, field string) error {
	if len(v) == 0 {
		return errors.New(fmt.Sprintf("%s  cannot be empty.", field))
	}
	return nil
}
