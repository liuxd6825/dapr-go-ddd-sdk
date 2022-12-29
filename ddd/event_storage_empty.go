package ddd

import (
	"context"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

type emptyEventStorage struct {
}

func (s *emptyEventStorage) Commit(ctx context.Context, req *daprclient.CommitRequest) (res *daprclient.CommitResponse, resErr error) {
	//TODO implement me
	panic("implement me")
}

func (s *emptyEventStorage) Rollback(ctx context.Context, req *daprclient.RollbackRequest) (res *daprclient.RollbackResponse, resErr error) {
	//TODO implement me
	panic("implement me")
}

func (s *emptyEventStorage) GetEvents(ctx context.Context, req *daprclient.GetEventsRequest) (*daprclient.GetEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *emptyEventStorage) GetRelations(ctx context.Context, req *daprclient.GetRelationsRequest) (*daprclient.GetRelationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *emptyEventStorage) LoadEvent(ctx context.Context, req *daprclient.LoadEventsRequest) (*daprclient.LoadEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *emptyEventStorage) CreateEvent(ctx context.Context, req *daprclient.CreateEventRequest) (*daprclient.CreateEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *emptyEventStorage) DeleteEvent(ctx context.Context, req *daprclient.DeleteEventRequest) (*daprclient.DeleteEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *emptyEventStorage) GetPubsubName() string {
	return ""
}

func NewEmptyEventStorage() EventStorage {
	return &emptyEventStorage{}
}

func (s *emptyEventStorage) GetHost() string {
	return ""
}

func (s *emptyEventStorage) GetPort() int {
	return -1
}

func (s *emptyEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate any) (res Aggregate, find bool, err error) {
	return nil, false, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) LoadEvents(ctx context.Context, req *daprclient.LoadEventsRequest) (*daprclient.LoadEventsResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (*daprclient.ApplyEventResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (*daprclient.SaveSnapshotResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error) {
	return false, errors.New("emptyEventStorage")
}
