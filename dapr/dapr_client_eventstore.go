package dapr

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
)

func (c *daprClient) LoadEvents(ctx context.Context, req *LoadEventsRequest) (*LoadEventsResponse, error) {
	if err := IsEmpty(req.TenantId, "TenantId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}

	in := &pb.LoadDomainEventRequest{
		CompName:      req.CompName,
		TenantId:      req.TenantId,
		AggregateType: req.AggregateType,
		AggregateId:   req.AggregateId,
		Headers:       newRequstHeaders(&req.Headers),
	}
	out, err := c.grpcClient.LoadDomainEvents(ctx, in)
	if err != nil {
		return nil, err
	}

	resp := &LoadEventsResponse{
		TenantId:      out.TenantId,
		AggregateId:   out.AggregateId,
		AggregateType: out.AggregateType,
	}

	if out.Snapshot != nil {
		aggregateData, err := NewMapInterface(out.Snapshot.AggregateData)
		if err != nil {
			return nil, err
		}
		metadata, err := NewMapString(out.Snapshot.Metadata)
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
			eventData, err := NewMapInterface(item.EventData)
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

func (c *daprClient) ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyEventResponse, error) {
	errs := errors.NewErrors()
	if err := IsEmpty(req.TenantId, "tenantId"); err != nil {
		errs.AddError(err)
	}
	if err := IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		errs.AddError(err)
	}
	if err := IsEmpty(req.AggregateType, "AggregateType"); err != nil {
		errs.AddError(err)
	}
	if err := IsEmpty(req.CompName, "CompName"); err != nil {
		errs.AddError(err)
	}
	if req.Events == nil {
		errs.AddError(errors.New("req.events cannot be nil"))
	}
	if !errs.IsEmpty() {
		return nil, errs.NewError()
	}
	events, err := c.newEvents(req.Events)
	if err != nil {
		return nil, err
	}

	in := &pb.ApplyDomainEventRequest{
		CompName:      req.CompName,
		SessionId:     req.SessionId,
		TenantId:      req.TenantId,
		AggregateId:   req.AggregateId,
		AggregateType: req.AggregateType,
		Events:        events,
		Headers:       newRequstHeaders(&req.Headers),
	}

	out, err := c.grpcClient.ApplyDomainEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &ApplyEventResponse{
		Headers: c.newResponseHeaders(out.Headers),
	}
	return resp, nil
}

func newRequstHeaders(request *RequestHeader) *pb.RequestHeaders {
	var values map[string][]string
	if request == nil {
		values = request.Values
	}
	return &pb.RequestHeaders{
		Values: newPbMapStrings(values),
	}
}

/*func (c *daprClient) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	if err := IsEmpty(req.TenantId, "tenantId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(req.AggregateType, "AggregateType"); err != nil {
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
}*/

/*func (c *daprClient) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {
	if err := IsEmpty(req.TenantId, "tenantId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(req.AggregateType, "AggregateType"); err != nil {
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
*/

func (c *daprClient) newEvents(events []*EventDto) ([]*pb.EventDto, error) {
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

func (c *daprClient) newEvent(e *EventDto) (*pb.EventDto, error) {
	if err := IsEmpty(e.CommandId, "CommandId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(e.PubsubName, "PubsubName"); err != nil {
		return nil, err
	}
	if err := IsEmpty(e.EventType, "EventType"); err != nil {
		return nil, err
	}
	if err := IsEmpty(e.EventId, "EventId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(e.EventVersion, "EventVersion"); err != nil {
		return nil, err
	}
	if err := IsEmpty(e.Topic, "Topic"); err != nil {
		return nil, err
	}
	eventData, err := ToJson(e.EventData)
	if err != nil {
		return nil, err
	}

	metadata, err := ToJson(e.Metadata)
	if err != nil {
		return nil, err
	}

	event := &pb.EventDto{
		ApplyType:    e.ApplyType,
		Metadata:     metadata,
		CommandId:    e.CommandId,
		EventId:      e.EventId,
		EventData:    eventData,
		EventType:    e.EventType,
		EventVersion: e.EventVersion,
		PubsubName:   e.PubsubName,
		Topic:        e.Topic,
		Relations:    e.Relations,
		IsSourcing:   e.IsSourcing,
	}
	return event, nil
}

func (c *daprClient) SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error) {
	if err := IsEmpty(req.TenantId, "tenantId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(req.AggregateId, "AggregateId"); err != nil {
		return nil, err
	}
	if err := IsEmpty(req.AggregateType, "AggregateType"); err != nil {
		return nil, err
	}
	if err := IsEmpty(req.AggregateVersion, "AggregateVersion"); err != nil {
		return nil, err
	}

	aggregateData, err := ToJson(req.AggregateData)
	if err != nil {
		return nil, err
	}
	metadata, err := ToJson(req.Metadata)
	if err != nil {
		return nil, err
	}
	in := &pb.SaveDomainEventSnapshotRequest{
		CompName:         req.CompName,
		TenantId:         req.TenantId,
		AggregateId:      req.AggregateId,
		AggregateType:    req.AggregateType,
		AggregateData:    aggregateData,
		AggregateVersion: req.AggregateVersion,
		SequenceNumber:   req.SequenceNumber,
		Metadata:         metadata,
		Headers:          newRequstHeaders(&req.Headers),
	}

	out, err := c.grpcClient.SaveDomainEventSnapshot(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &SaveSnapshotResponse{
		Headers: c.newResponseHeaders(out.Headers),
	}
	return resp, nil
}

func (c *daprClient) GetRelations(ctx context.Context, req *GetRelationsRequest) (*GetRelationsResponse, error) {
	if req == nil {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req is nil")
	}
	if len(req.TenantId) == 0 {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req.TenantId is nil")
	}
	if len(req.AggregateType) == 0 {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req.AggregateType is nil")
	}

	in := &pb.GetDomainEventRelationsRequest{
		CompName:      req.CompName,
		TenantId:      req.TenantId,
		AggregateType: req.AggregateType,
		Filter:        req.Filter,
		Sort:          req.Sort,
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		Headers:       newRequstHeaders(&req.Headers),
	}

	out, err := c.grpcClient.GetDomainEventRelations(ctx, in)
	if err != nil {
		return nil, err
	}

	var relations []*Relation
	if out != nil && len(out.Data) > 0 {
		for _, datum := range out.Data {
			relation := &Relation{
				Id:          datum.Id,
				TenantId:    datum.TenantId,
				AggregateId: datum.AggregateId,
				IsDeleted:   datum.IsDeleted,
				TableName:   datum.TableName,
				RelName:     datum.RelName,
				RelValue:    datum.RelValue,
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

func (c *daprClient) GetEvents(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	if req == nil {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req is nil")
	}
	if len(req.TenantId) == 0 {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req.TenantId is nil")
	}
	if len(req.AggregateType) == 0 {
		return nil, errors.New("daprclient.GetRelations(ctx, req) error: req.AggregateType is nil")
	}

	in := &pb.GetDomainEventsRequest{
		CompName:      req.CompName,
		TenantId:      req.TenantId,
		AggregateType: req.AggregateType,
		Filter:        req.Filter,
		Sort:          req.Sort,
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		Headers:       newRequstHeaders(&req.Headers),
	}

	out, err := c.grpcClient.GetDomainEvents(ctx, in)
	if err != nil {
		return nil, err
	}

	var events []*GetEventsItem
	if out != nil && len(out.Data) > 0 {
		for _, datum := range out.Data {
			eventData := map[string]interface{}{}
			if err := json.Unmarshal([]byte(datum.EventData), &eventData); err != nil {
				return nil, err
			}
			metadata := map[string]string{}
			if err := json.Unmarshal([]byte(datum.Metadata), &metadata); err != nil {
				return nil, err
			}
			eventItem := &GetEventsItem{
				EventId:      datum.EventId,
				CommandId:    datum.CommandId,
				EventData:    eventData,
				EventType:    datum.EventType,
				EventVersion: datum.EventVersion,
				// EventTime:    datum.EventTime,
				PubsubName: datum.PubsubName,
				Topic:      datum.Topic,
				Metadata:   metadata,
			}
			events = append(events, eventItem)
		}
	}

	resp := &GetEventsResponse{}
	resp.Sort = out.Sort
	resp.PageNum = out.PageNum
	resp.PageSize = out.PageSize
	resp.Filter = out.Filter
	resp.Error = out.Error
	resp.IsFound = out.IsFound
	resp.TotalRows = out.TotalRows
	resp.TotalPages = out.TotalPages
	resp.Data = events
	resp.Headers = c.newResponseHeaders(out.Headers)

	return resp, nil
}

func (c *daprClient) newResponseHeaders(out *pb.ResponseHeaders) *ResponseHeaders {
	if out == nil {
		return &ResponseHeaders{
			Status:  ResponseStatusSuccess,
			Message: "",
			Values:  make(map[string][]string),
		}
	}
	values := out.Values
	res := &ResponseHeaders{
		Status:  ResponseStatus(out.Status),
		Message: out.Message,
		Values:  newMapStrings(values),
	}
	return res
}
