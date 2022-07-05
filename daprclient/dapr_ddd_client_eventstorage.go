package daprclient

import (
	"context"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_utils"
	pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
)

func (c *daprDddClient) LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := ddd_utils.IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}

	in := &pb.LoadEventRequest{
		TenantId:      req.TenantId,
		AggregateType: req.AggregateType,
		AggregateId:   req.AggregateId,
	}
	out, err := c.grpcClient.LoadEvents(ctx, in)
	if err != nil {
		return nil, err
	}

	resp := &LoadEventsResponse{
		TenantId:      out.TenantId,
		AggregateId:   out.AggregateId,
		AggregateType: out.AggregateType,
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
	resp.Headers = c.newResponseHeaders(out.Headers)
	return resp, nil
}

func (c *daprDddClient) ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "tenantId"); err != nil {
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
	out, err := c.grpcClient.ApplyEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &ApplyEventResponse{
		Headers: c.newResponseHeaders(out.Headers),
	}
	return resp, nil
}

func (c *daprDddClient) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "tenantId"); err != nil {
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
	out, err := c.grpcClient.CreateEvent(ctx, in)
	if err != nil {
		return nil, err
	}

	resp := &CreateEventResponse{
		Headers: c.newResponseHeaders(out.Headers),
	}
	return resp, nil
}

func (c *daprDddClient) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "tenantId"); err != nil {
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
	out, err := c.grpcClient.DeleteEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &DeleteEventResponse{
		Headers: c.newResponseHeaders(out.Headers),
	}
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
		Relations:    e.Relations,
	}
	return event, nil
}

func (c *daprDddClient) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	if err := ddd_utils.IsEmpty(req.TenantId, "tenantId"); err != nil {
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
	out, err := c.grpcClient.SaveSnapshot(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &SaveSnapshotResponse{
		Headers: c.newResponseHeaders(out.Headers),
	}
	return resp, nil
}

func (c *daprDddClient) GetRelations(ctx context.Context, req *GetRelationsRequest) (*GetRelationsResponse, error) {
	if req == nil {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req is nil")
	}
	if len(req.TenantId) == 0 {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req.TenantId is nil")
	}
	if len(req.AggregateType) == 0 {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req.AggregateType is nil")
	}

	in := &pb.GetRelationsRequest{
		TenantId:      req.TenantId,
		AggregateType: req.AggregateType,
		Filter:        req.Filter,
		Sort:          req.Sort,
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
	}
	out, err := c.grpcClient.GetRelations(ctx, in)
	if err != nil {
		return nil, err
	}

	var relations []*Relation
	if out != nil && len(out.Data) > 0 {
		for _, datum := range out.Data {
			items := datum.Items
			if items != nil {
				items = make(map[string]string)
			}
			relation := &Relation{
				Id:          datum.Id,
				TenantId:    datum.TenantId,
				AggregateId: datum.AggregateId,
				IsDeleted:   datum.IsDeleted,
				TableName:   datum.TableName,
				Items:       items,
			}
			relations = append(relations, relation)
		}
	}

	resp := &GetRelationsResponse{}
	resp.Sort = out.Sort
	resp.PageNum = out.PageNum
	resp.PageSize = out.PageSize
	resp.Filter = out.Filter
	resp.Error = out.Error
	resp.IsFound = out.IsFound
	resp.TotalRows = out.TotalRows
	resp.TotalPages = out.TotalPages
	resp.Data = relations
	resp.Headers = c.newResponseHeaders(out.Headers)

	return resp, nil
}

func (c *daprDddClient) newResponseHeaders(out *pb.ResponseHeaders) *ResponseHeaders {
	if out == nil {
		return &ResponseHeaders{
			Status:  ResponseStatusSuccess,
			Message: "",
			Values:  map[string]string{},
		}
	}
	values := out.Values
	if values == nil {
		values = map[string]string{}
	}
	res := &ResponseHeaders{
		Status:  ResponseStatus(out.Status),
		Message: out.Message,
		Values:  values,
	}
	return res
}
