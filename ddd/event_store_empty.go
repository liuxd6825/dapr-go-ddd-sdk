package ddd

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

type emptyEventStore struct {
}

const NotImplemented = "emptyEventStore not implemented %v"

func NewEmptyEventStore() EventStore {
	return &emptyEventStore{}
}

func (s *emptyEventStore) panicMessage(funcName string) string {
	return fmt.Sprintf(NotImplemented, funcName)
}

func (s *emptyEventStore) Commit(ctx context.Context, req *daprclient.CommitRequest) (res *daprclient.CommitResponse, resErr error) {
	panic(s.panicMessage("Commit()"))
}

func (s *emptyEventStore) Rollback(ctx context.Context, req *daprclient.RollbackRequest) (res *daprclient.RollbackResponse, resErr error) {
	panic(s.panicMessage("Rollback()"))
}

func (s *emptyEventStore) GetEvents(ctx context.Context, req *daprclient.GetEventsRequest) (*daprclient.GetEventsResponse, error) {
	panic(s.panicMessage("GetEvents()"))
}

func (s *emptyEventStore) GetRelations(ctx context.Context, req *daprclient.GetRelationsRequest) (*daprclient.GetRelationsResponse, error) {
	panic(s.panicMessage("GetRelations()"))
}

func (s *emptyEventStore) LoadEvent(ctx context.Context, req *daprclient.LoadEventsRequest) (*daprclient.LoadEventsResponse, error) {
	panic(s.panicMessage("LoadEvent()"))
}

func (s *emptyEventStore) CreateEvent(ctx context.Context, req *daprclient.CreateEventRequest) (*daprclient.CreateEventResponse, error) {
	panic(s.panicMessage("CreateEvent()"))
}

func (s *emptyEventStore) DeleteEvent(ctx context.Context, req *daprclient.DeleteEventRequest) (*daprclient.DeleteEventResponse, error) {
	panic(s.panicMessage("DeleteEvent()"))
}

func (s *emptyEventStore) GetPubsubName() string {
	panic(s.panicMessage("GetPubsubName()"))
}

func (s *emptyEventStore) GetHost() string {
	panic(s.panicMessage("GetHost()"))
}

func (s *emptyEventStore) GetPort() int {
	panic(s.panicMessage("GetPort()"))
}

func (s *emptyEventStore) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate any) (res Aggregate, find bool, err error) {
	panic(s.panicMessage("LoadAggregate()"))
}

func (s *emptyEventStore) LoadEvents(ctx context.Context, req *daprclient.LoadEventsRequest) (*daprclient.LoadEventsResponse, error) {
	panic(s.panicMessage("LoadEvents()"))
}

func (s *emptyEventStore) ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (*daprclient.ApplyEventResponse, error) {
	panic(s.panicMessage("ApplyEvent()"))
}

func (s *emptyEventStore) SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (*daprclient.SaveSnapshotResponse, error) {
	return nil, errors.New("emptyEventStore")
}

func (s *emptyEventStore) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error) {
	return false, errors.New("emptyEventStore")
}
