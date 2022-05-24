package daprclient

import (
	"context"
	"errors"
	pb "github.com/dapr/dapr/pkg/proto/runtime/v1"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
)

func (c *daprDddClient) LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}

	in := &pb.LoadEventRequest{
		TenantId:    req.TenantId,
		AggregateId: req.AggregateId,
	}
	out, err := c.grpcDaprClient.LoadEvents(ctx, in)
	if err != nil {
		return nil, err
	}

	resp := &LoadEventsResponse{
		TenantId:    out.TenantId,
		AggregateId: out.AggregateId,
	}

	if out.Snapshot != nil {
		aggregateData, err := ddd_utils.NewMapInterface(out.Snapshot.AggregateData)
		if err != nil {
			return nil, err
		}
		metadata, err := ddd_utils.NewMapString(out.Snapshot.Metadata)
		if err != nil {
			return nil, err
		}

		snapshot := &Snapshot{
			AggregateData:     aggregateData,
			AggregateRevision: out.Snapshot.AggregateRevision,
			SequenceNumber:    out.Snapshot.SequenceNumber,
			Metadata:          metadata,
		}
		resp.Snapshot = snapshot
	}

	events := make([]EventRecord, 0)
	if out.Events != nil {
		for _, item := range out.Events {
			eventData, err := ddd_utils.NewMapInterface(item.EventData)
			if err != nil {
				return nil, err
			}
			event := EventRecord{
				EventId:       item.EventId,
				EventData:     eventData,
				EventRevision: item.EventRevision,
				EventType:     item.EventType,
			}
			events = append(events, event)
		}
	}
	resp.EventRecords = &events

	return resp, nil
}

func (c *daprDddClient) ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventsResponse, error) {

	if err := ddd_utils.IsEmpty(req.CommandId, "CommandId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.PubsubName, "PubsubName"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.EventType, "EventType"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.EventId, "EventId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.EventRevision, "EventRevision"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.Topic, "Topic"); err != nil {
		return nil, err
	}
	if req.EventData == nil {
		return nil, errors.New("EventData cannot be nil.")
	}

	eventData, err := ddd_utils.ToJson(req.EventData)
	if err != nil {
		return nil, err
	}

	metadata, err := ddd_utils.ToJson(req.Metadata)
	if err != nil {
		return nil, err
	}

	in := &pb.ApplyEventRequest{
		TenantId:      req.TenantId,
		Metadata:      metadata,
		CommandId:     req.CommandId,
		EventId:       req.EventId,
		EventData:     eventData,
		EventType:     req.EventType,
		EventRevision: req.EventRevision,
		AggregateId:   req.AggregateId,
		AggregateType: req.AggregateType,
		PubsubName:    req.PubsubName,
		Topic:         req.Topic,
	}
	_, err = c.grpcDaprClient.ApplyEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &ApplyEventsResponse{}
	return resp, nil
}

func (c *daprDddClient) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateType, "AggregateType"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateRevision, "AggregateRevision"); err != nil {
		return nil, err
	}

	aggregateData, err := ddd_utils.ToJson(req.AggregateData)
	if err != nil {
		return nil, err
	}
	metadata, err := ddd_utils.ToJson(req.Metadata)
	if err != nil {
		return nil, err
	}
	in := &pb.SaveSnapshotRequest{
		TenantId:          req.TenantId,
		AggregateId:       req.AggregateId,
		AggregateType:     req.AggregateType,
		AggregateData:     aggregateData,
		AggregateRevision: req.AggregateRevision,
		SequenceNumber:    req.SequenceNumber,
		Metadata:          metadata,
	}
	_, err = c.grpcDaprClient.SaveSnapshot(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &SaveSnapshotResponse{}
	return resp, nil
}

func (c *daprDddClient) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error) {
	if err := ddd_utils.IsEmpty(tenantId, "TenantId"); err != nil {
		return false, err
	}
	if err := ddd_utils.IsEmpty(aggregateId, "AggregateId"); err != nil {
		return false, err
	}

	in := &pb.ExistAggregateRequest{
		TenantId:    tenantId,
		AggregateId: aggregateId,
	}

	out, err := c.grpcDaprClient.ExistAggregate(ctx, in)
	if err != nil {
		return false, err
	}
	return out.IsExist, nil
}
