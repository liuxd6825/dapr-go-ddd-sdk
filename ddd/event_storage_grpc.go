package ddd

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"io"
	"net/http"
)

type grpcEventStorage struct {
	client     daprclient.DaprDddClient
	pubsubName string
	subscribes *[]Subscribe
}

func NewGrpcEventStorage(client daprclient.DaprDddClient, options ...func(s EventStorage)) (EventStorage, error) {
	subscribes = make([]Subscribe, 0)
	res := &grpcEventStorage{
		client:     client,
		subscribes: &subscribes,
	}
	for _, option := range options {
		option(res)
	}
	return res, nil
}

func (s *grpcEventStorage) GetPubsubName() string {
	return s.pubsubName
}

func (s *grpcEventStorage) LoadAggregate(ctx context.Context, tenantId string, aggregateId string, aggregate any) (agg Aggregate, isFound bool, resErr error) {
	a, ok := aggregate.(Aggregate)
	if !ok {
		return nil, false, errors.New("aggregate is not ddd.Aggregate interface")
	}

	defer func() {
		if e := errors.GetRecoverError(recover()); e != nil {
			resErr = e
		}
	}()

	if err := assert.NotNil(aggregate, assert.NewOptions("aggregate is nil")); err != nil {
		return nil, false, err
	}

	if err := assert.NotEmpty(aggregateId, assert.NewOptions("aggregateId is nil")); err != nil {
		return nil, false, err
	}

	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is nil")); err != nil {
		return nil, false, err
	}

	req := &daprclient.LoadEventsRequest{
		TenantId:      tenantId,
		AggregateType: a.GetAggregateType(),
		AggregateId:   aggregateId,
	}

	resp, err := s.LoadEvent(ctx, req)
	if err != nil {
		return nil, false, errors.New("grpcEventStorage.LoadEvent() error:%s", err.Error())
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
				return nil, false, errors.New("CallEventHandler(agg, record) eventType:%v, error:%v", record.EventType, err.Error())
			}
		}
	}
	return a, true, err
}

func (s *grpcEventStorage) LoadEvent(ctx context.Context, req *daprclient.LoadEventsRequest) (res *daprclient.LoadEventsResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()
	return s.client.LoadEvents(ctx, req)
}

func (s *grpcEventStorage) ApplyEvent(ctx context.Context, req *daprclient.ApplyEventRequest) (res *daprclient.ApplyEventResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()
	s.setEventsPubsubName(req.Events)
	return s.client.ApplyEvent(ctx, req)
}

func (s *grpcEventStorage) Commit(ctx context.Context, req *daprclient.CommitRequest) (res *daprclient.CommitResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()

	return s.client.Commit(ctx, req)
}

func (s *grpcEventStorage) Rollback(ctx context.Context, req *daprclient.RollbackRequest) (res *daprclient.RollbackResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()

	return s.client.Rollback(ctx, req)
}

/*func (s *grpcEventStorage) CreateEvent(ctx context.Context, req *daprclient.CreateEventRequest) (res *daprclient.CreateEventResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()
	s.setEventsPubsubName(req.Events)
	return s.client.CreateEvent(ctx, req)
}*/

/*func (s *grpcEventStorage) DeleteEvent(ctx context.Context, req *daprclient.DeleteEventRequest) (res *daprclient.DeleteEventResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()
	if req != nil && req.Event != nil && len(req.Event.PubsubName) == 0 {
		req.Event.PubsubName = s.pubsubName
	}
	return s.client.DeleteEvent(ctx, req)
}
*/
func (s *grpcEventStorage) SaveSnapshot(ctx context.Context, req *daprclient.SaveSnapshotRequest) (res *daprclient.SaveSnapshotResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()
	return s.client.SaveSnapshot(ctx, req)
}

func (s *grpcEventStorage) GetEvents(ctx context.Context, req *daprclient.GetEventsRequest) (res *daprclient.GetEventsResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()
	return s.client.GetEvents(ctx, req)
}

func (s *grpcEventStorage) GetRelations(ctx context.Context, req *daprclient.GetRelationsRequest) (res *daprclient.GetRelationsResponse, resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resErr = err
			}
		}
	}()
	return s.client.GetRelations(ctx, req)
}

func (s *grpcEventStorage) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	return bytes, err
}

func (s *grpcEventStorage) setEventsPubsubName(events []*daprclient.EventDto) {
	if events != nil {
		for _, event := range events {
			if len(event.PubsubName) == 0 {
				event.PubsubName = s.pubsubName
			}
		}
	}
}
