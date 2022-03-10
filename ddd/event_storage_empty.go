package ddd

import (
	"context"
	"errors"
)

type emptyEventStorage struct {
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

func (s *emptyEventStorage) GetPussubName() string {
	return ""
}

func (s *emptyEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (res Aggregate, find bool, err error) {
	return nil, false, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error) {
	return false, errors.New("emptyEventStorage")
}
