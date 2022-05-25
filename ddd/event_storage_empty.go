package ddd

import (
	"context"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

type emptyEventStorage struct {
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

func (s *emptyEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate Aggregate) (res Aggregate, find bool, err error) {
	return nil, false, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) LoadEvents(ctx context.Context, req *daprclient.LoadEventsRequest) (*daprclient.LoadEventsResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (*daprclient.ApplyEventsResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (*daprclient.SaveSnapshotResponse, error) {
	return nil, errors.New("emptyEventStorage")
}

func (s *emptyEventStorage) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error) {
	return false, errors.New("emptyEventStorage")
}
