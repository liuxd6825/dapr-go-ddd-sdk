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
				EventId:      item.EventId,
				EventData:    eventData,
				EventVersion: item.EventVersion,
				EventType:    item.EventType,
			}
			events = append(events, event)
		}
	}
	resp.EventRecords = &events

	return resp, nil
}

func (c *daprDddClient) ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateType, "AggregateType"); err != nil {
		return nil, err
	}
	if req.Events == nil {
		return nil, errors.New("req.events cannot be nil")
	}
	events, err := c.newEvents(req.Events)
	if err != nil {
		return nil, err
	}

	in := &pb.ApplyEventRequest{
		TenantId:      req.TenantId,
		AggregateId:   req.AggregateId,
		AggregateType: req.AggregateType,
		Events:        events,
	}
	_, err = c.grpcDaprClient.ApplyEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &ApplyEventResponse{}
	return resp, nil
}

func (c *daprDddClient) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateType, "AggregateType"); err != nil {
		return nil, err
	}
	if req.Events == nil {
		return nil, errors.New("req.events cannot be nil")
	}
	events, err := c.newEvents(req.Events)
	if err != nil {
		return nil, err
	}

	in := &pb.CreateEventRequest{
		TenantId:      req.TenantId,
		AggregateId:   req.AggregateId,
		AggregateType: req.AggregateType,
		Events:        events,
	}
	_, err = c.grpcDaprClient.CreateEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &CreateEventResponse{}
	return resp, nil
}

func (c *daprDddClient) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateType, "AggregateType"); err != nil {
		return nil, err
	}
	if req.Event == nil {
		return nil, errors.New("req.event cannot be nil")
	}
	event, err := c.newEvent(req.Event)
	if err != nil {
		return nil, err
	}

	in := &pb.DeleteEventRequest{
		TenantId:      req.TenantId,
		AggregateId:   req.AggregateId,
		AggregateType: req.AggregateType,
		Event:         event,
	}
	_, err = c.grpcDaprClient.DeleteEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &DeleteEventResponse{}
	return resp, nil
}

func (c *daprDddClient) newEvents(events []*EventDto) ([]*pb.EventDto, error) {
	var resList []*pb.EventDto
	for _, e := range events {
		event, err := c.newEvent(e)
		if err != nil {
			return nil, err
		}
		resList = append(resList, event)
	}
	return resList, nil
}

func (c *daprDddClient) newEvent(e *EventDto) (*pb.EventDto, error) {
	if err := ddd_utils.IsEmpty(e.CommandId, "CommandId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(e.PubsubName, "PubsubName"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(e.EventType, "EventType"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(e.EventId, "EventId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(e.EventVersion, "EventVersion"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(e.Topic, "Topic"); err != nil {
		return nil, err
	}
	eventData, err := ddd_utils.ToJson(e.EventData)
	if err != nil {
		return nil, err
	}

	metadata, err := ddd_utils.ToJson(e.Metadata)
	if err != nil {
		return nil, err
	}

	event := &pb.EventDto{
		Metadata:     metadata,
		CommandId:    e.CommandId,
		EventId:      e.EventId,
		EventData:    eventData,
		EventType:    e.EventType,
		EventVersion: e.EventVersion,
		PubsubName:   e.PubsubName,
		Topic:        e.Topic,
	}
	return event, nil
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
	if err := ddd_utils.IsEmpty(req.AggregateVersion, "AggregateVersion"); err != nil {
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
		TenantId:         req.TenantId,
		AggregateId:      req.AggregateId,
		AggregateType:    req.AggregateType,
		AggregateData:    aggregateData,
		AggregateVersion: req.AggregateVersion,
		SequenceNumber:   req.SequenceNumber,
		Metadata:         metadata,
	}
	_, err = c.grpcDaprClient.SaveSnapshot(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &SaveSnapshotResponse{}
	return resp, nil
}
