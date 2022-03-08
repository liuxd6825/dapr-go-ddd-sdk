package ddd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type daprEventStorage struct {
	host              string
	port              int
	client            *http.Client
	defaultPubsubName string
}

func NewDaprEventStorage(host string, port int, pubsubName string) EventStorage {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     IdleConnTimeout * time.Second,
		},
	}
	return &daprEventStorage{
		host:              host,
		port:              port,
		client:            client,
		defaultPubsubName: pubsubName,
	}
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
			domainEvent := aggregate.CreateDomainEvent(ctx, &record)
			sequenceNumber = record.SequenceNumber
			if err := record.Marshal(domainEvent); err != nil {
				return res, find, err
			}
			if err := aggregate.OnSourceEvent(ctx, domainEvent); err != nil {
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
	url := fmt.Sprintf("/v1.0/event-sourcing/events/%s/%s", req.TenantId, req.AggregateId)
	data, err := s.httpGet(url)
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
		req.PubsubName = s.defaultPubsubName
	}
	url := fmt.Sprintf("/v1.0/event-sourcing/events/apply")
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

	data, err := s.httpPost(url, req)
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
	url := fmt.Sprintf("/v1.0/event-sourcing/snapshot/save")
	data, err := s.httpPost(url, req)
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
	url := fmt.Sprintf("/v1.0/event-sourcing/aggregates/%s/%s", tenantId, aggregateId)
	data, err := s.httpGet(url)
	if err != nil {
		return false, err
	}
	resp := &ExistAggregateResponse{}
	err = json.Unmarshal(data, resp)
	return resp.IsExist, err
}

func (s *daprEventStorage) httpGet(url string) ([]byte, error) {
	resp, err := s.client.Get(fmt.Sprintf("http://%s:%d/%s", s.host, s.port, url))
	if err != nil {
		return nil, err
	}
	bytes, err := s.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(bytes))
	}
	return bytes, err
}

func (s *daprEventStorage) httpPost(url string, reqData interface{}) ([]byte, error) {
	httpUrl := fmt.Sprintf("http://%s:%d/%s", s.host, s.port, url)
	jsonData, err := json.Marshal(reqData)
	resp, err := s.client.Post(httpUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	bytes, err := s.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(bytes))
	}
	return bytes, err
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
