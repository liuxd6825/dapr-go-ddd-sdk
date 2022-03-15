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
	"reflect"
	"strings"
	"time"
)

const (
	DefaultMaxIdleConns        = 10
	DefaultMaxIdleConnsPerHost = 50
	DefaultIdleConnTimeout     = 5
)

type daprEventStorage struct {
	host       string
	port       int
	client     *http.Client
	pubsubName string
	subscribes *[]Subscribe
}

func NewDaprEventStorage(host string, port int, options ...func(s EventStorage)) (EventStorage, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        DefaultMaxIdleConns,
			MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
			IdleConnTimeout:     DefaultIdleConnTimeout * time.Second,
		},
	}
	subscribes = make([]Subscribe, 0)
	res := &daprEventStorage{
		host:       host,
		port:       port,
		client:     client,
		subscribes: &subscribes,
	}
	for _, option := range options {
		option(res)
	}
	return res, nil
}

func (s *daprEventStorage) GetHost() string {
	return s.host
}

func (s *daprEventStorage) GetPort() int {
	return s.port
}

func (s *daprEventStorage) GetPubsubName() string {
	return s.pubsubName
}

func CallEventHandler(ctx context.Context, handler interface{}, eventType string, eventRevision string) error {
	v := reflect.ValueOf(handler)
	domainEvent, err := NewDomainEvent(eventType, eventRevision)
	if err != nil {
		return errors.New(fmt.Sprintf("Method is not found or not exported."))
	}
	methodName := getMethodName(eventType, eventRevision)
	method := v.MethodByName(methodName)
	if method.IsValid() {
		p1 := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(domainEvent),
		}
		resValues := method.Call(p1)
		if len(resValues) == 1 {
			err, ok := resValues[0].Interface().(error)
			if ok {
				return err
			}
		}
	} else {
		return errors.New(fmt.Sprintf("Method %s is not found or not exported.", methodName))
	}
	return nil
}

func getMethodName(eventType string, revision string) string {
	names := strings.Split(eventType, ".")
	name := names[len(names)-1]
	ver := strings.Replace(revision, ".", "_", -1)
	return fmt.Sprintf("On%sV%s", name, ver)
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

			domainEvent, err := NewDomainEvent(record.EventType, record.EventRevision)
			if err != nil {
				return nil, true, errors.New(fmt.Sprintf("Method is not found or not exported."))
			}
			//domainEvent := aggregate.CreateDomainEvent(ctx, &record)

			if err := record.Marshal(domainEvent); err != nil {
				return res, find, err
			}
			if err = CallEventHandler(ctx, aggregate, record.EventType, record.EventRevision); err != nil {
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
		req.PubsubName = s.pubsubName
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
