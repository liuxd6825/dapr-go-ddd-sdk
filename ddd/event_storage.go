package ddd

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

type EventStorage interface {
	LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error)
	LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error)
	ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error)
	SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
	ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error)
	GetPussubName() string
	GetHost() string
	GetPort() int
}

type EventStorageOption func(EventStorage)

func PubsubName(pubsubName string) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		s.pubsubName = pubsubName
	}
}

func IdleConnTimeout(idleConnTimeout time.Duration) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		t, _ := s.client.Transport.(*http.Transport)
		t.IdleConnTimeout = idleConnTimeout
	}
}

func MaxIdleConns(maxIdleConns int) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		t, _ := s.client.Transport.(*http.Transport)
		t.MaxIdleConns = maxIdleConns
	}
}

func MaxIdleConnsPerHost(maxIdleConnsPerHost int) EventStorageOption {
	return func(es EventStorage) {
		s, _ := es.(*daprEventStorage)
		t, _ := s.client.Transport.(*http.Transport)
		t.MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

func NewEventStorage(host string, port int, options ...func(s EventStorage)) (EventStorage, error) {
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
	res := &daprEventStorage{
		host:   host,
		port:   port,
		client: client,
	}
	for _, option := range options {
		option(res)
	}
	return res, nil
}

func LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (Aggregate, bool, error) {
	return eventStorage.LoadAggregate(ctx, tenantId, aggregateId, aggregate)
}

func LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error) {
	return eventStorage.LoadEvents(ctx, req)
}

type ApplyOptions struct {
	pubsubName string
	metadata   map[string]string
}
type ApplyOption func(*ApplyOptions)

func ApplyPubsubName(pubsubName string) ApplyOption {
	return func(o *ApplyOptions) {
		o.pubsubName = pubsubName
	}
}

func ApplyMetadata(metadata map[string]string) ApplyOption {
	return func(o *ApplyOptions) {
		o.metadata = metadata
	}
}

func Apply(ctx context.Context, aggregate Aggregate, event DomainEvent, options ...ApplyOption) error {
	appOptions := &ApplyOptions{
		pubsubName: "",
		metadata:   map[string]string{},
	}
	for _, option := range options {
		option(appOptions)
	}
	req := &ApplyEventRequest{
		TenantId:      event.GetTenantId(),
		CommandId:     event.GetCommandId(),
		EventId:       event.GetEventId(),
		EventRevision: event.GetEventRevision(),
		EventType:     event.GetEventType(),
		AggregateId:   event.GetAggregateId(),
		AggregateType: aggregate.GetAggregateType(),
		Metadata:      appOptions.metadata,
		PubsubName:    appOptions.pubsubName,
		EventData:     event,
		Topic:         event.GetEventType(),
	}
	if _, err := eventStorage.ApplyEvent(ctx, req); err != nil {
		return err
	}
	if err := aggregate.OnSourceEvent(ctx, event); err != nil {
		return err
	}
	return nil
}

func CreateAggregate(ctx context.Context, aggregate Aggregate, cmd DomainCommand) error {
	ok, err := eventStorage.ExistAggregate(ctx, cmd.GetTenantId(), cmd.GetAggregateId())
	if err != nil {
		return err
	}
	if ok {
		return errors.New(cmd.GetAggregateId() + " aggregate root already exists.")
	}
	return aggregate.OnCommand(ctx, cmd)
}

func CommandAggregate(ctx context.Context, aggregate Aggregate, cmd DomainCommand) error {
	_, find, err := LoadAggregate(ctx, cmd.GetTenantId(), cmd.GetAggregateId(), aggregate)
	if err != nil {
		return err
	}
	if !find {
		return errors.New(cmd.GetAggregateId() + " aggregate root not fond.")
	}
	return aggregate.OnCommand(ctx, cmd)
}

func applyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error) {
	return eventStorage.ApplyEvent(ctx, req)
}

func saveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	return eventStorage.SaveSnapshot(ctx, req)
}
