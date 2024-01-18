package dapr

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
	"time"
)

type WriteEventLogRequest struct {
	Headers  RequestHeader `json:"headers"`
	CompName string        `json:"compName"`
	Id       string        `json:"id"`
	TenantId string        `json:"tenantId"`
	AppId    string        `json:"appId"`
	Class    string        `json:"class"`
	Func     string        `json:"func"`
	Level    string        `json:"level"`
	Time     *time.Time    `json:"time"`
	Status   bool          `json:"status"`
	Message  string        `json:"message"`

	PubAppId string `json:"pubAppId"`
	EventId  string `json:"eventId"`
	// EventType string `json:"eventType"`
	CommandId string `json:"commandId"`
}

type WriteEventLogResponse struct {
}

type UpdateEventLogRequest struct {
	Headers  RequestHeader `json:"headers"`
	CompName string        `json:"compName"`
	Id       string        `json:"id"`
	TenantId string        `json:"tenantId"`
	AppId    string        `json:"appId"`
	Class    string        `json:"class"`
	Func     string        `json:"func"`
	Level    string        `json:"level"`
	Time     *time.Time    `json:"time"`
	Status   bool          `json:"status"`
	Message  string        `json:"message"`

	PubAppId  string `json:"pubAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
}

type UpdateEventLogResponse struct {
}

// GetLogByCommandId

type GetEventLogByCommandIdRequest struct {
	Headers   RequestHeader `json:"headers"`
	CompName  string        `json:"compName"`
	TenantId  string        `json:"tenantId"`
	AppId     string        `json:"appId"`
	CommandId string        `json:"commandId"`
}

type GetEventLogByCommandIdResponse struct {
	Data []*EventLogDto `json:"data"`
}

type EventLogDto struct {
	Id       string     `json:"id"`
	TenantId string     `json:"tenantId"`
	AppId    string     `json:"appId"`
	Class    string     `json:"class"`
	Func     string     `json:"func"`
	Level    string     `json:"level"`
	Time     *time.Time `json:"time"`
	Status   bool       `json:"status"`
	Message  string     `json:"message"`

	PubAppId  string `json:"pubAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
}

//

type WriteAppLogRequest struct {
	Headers  RequestHeader `json:"headers"`
	CompName string        `json:"compName"`
	Id       string        `json:"id"`
	TenantId string        `json:"tenantId"`
	AppId    string        `json:"appId"`
	Class    string        `json:"class"`
	Func     string        `json:"func"`
	Level    string        `json:"level"`
	Time     *time.Time    `json:"time"`
	Status   bool          `json:"status"`
	Message  string        `json:"message"`
}

type WriteAppLogResponse struct {
}

type UpdateAppLogRequest struct {
	Headers  RequestHeader `json:"headers"`
	CompName string        `json:"compName"`
	Id       string        `json:"id"`
	TenantId string        `json:"tenantId"`
	AppId    string        `json:"appId"`
	Class    string        `json:"class"`
	Func     string        `json:"func"`
	Level    string        `json:"level"`
	Time     *time.Time    `json:"time"`
	Status   bool          `json:"status"`
	Message  string        `json:"message"`
}

type UpdateAppLogResponse struct {
}

// GetLogByCommandId

type GetAppLogByIdRequest struct {
	Headers  RequestHeader `json:"headers"`
	CompName string        `json:"compName"`
	TenantId string        `json:"tenantId"`
	Id       string        `json:"id"`
}

type GetAppLogByIdResponse struct {
	Id       string     `json:"id"`
	TenantId string     `json:"tenantId"`
	AppId    string     `json:"appId"`
	Class    string     `json:"class"`
	Func     string     `json:"func"`
	Level    string     `json:"level"`
	Time     *time.Time `json:"time"`
	Status   bool       `json:"status"`
	Message  string     `json:"message"`
}

func (c *daprClient) WriteEventLog(ctx context.Context, req *WriteEventLogRequest) (resp *WriteEventLogResponse, resErr error) {
	request := &pb.WriteAppEventLogRequest{
		Headers:  newRequstHeaders(&req.Headers),
		TenantId: req.TenantId,
		Id:       req.Id,
		AppId:    req.AppId,
		Class:    req.Class,
		Func:     req.Func,
		Level:    req.Level,
		Time:     timeutils.AsTimestamp(req.Time),
		Status:   req.Status,
		Message:  req.Message,

		PubAppId:  req.PubAppId,
		EventId:   req.EventId,
		CommandId: req.CommandId,
	}
	_, err := c.grpcClient.WriteAppEventLog(ctx, request)
	if err != nil {
		return nil, err
	}
	resp = &WriteEventLogResponse{}
	return resp, nil
}

func (c *daprClient) UpdateEventLog(ctx context.Context, req *UpdateEventLogRequest) (resp *UpdateEventLogResponse, resErr error) {
	request := &pb.UpdateAppEventLogRequest{
		Headers:  newRequstHeaders(&req.Headers),
		TenantId: req.TenantId,
		Id:       req.Id,
		AppId:    req.AppId,
		Class:    req.Class,
		Func:     req.Func,
		Level:    req.Level,
		Time:     timeutils.AsTimestamp(req.Time),
		Status:   req.Status,
		Message:  req.Message,

		PubAppId:  req.PubAppId,
		EventId:   req.EventId,
		CommandId: req.CommandId,
	}

	_, err := c.grpcClient.UpdateAppEventLog(ctx, request)
	if err != nil {
		return nil, err
	}

	response := &UpdateEventLogResponse{}
	return response, nil
}

func (c *daprClient) GetEventLogByCommandId(ctx context.Context, req *GetEventLogByCommandIdRequest) (resp *GetEventLogByCommandIdResponse, resErr error) {
	request := &pb.GetAppEventLogByCommandIdRequest{
		Headers:   newRequstHeaders(&req.Headers),
		TenantId:  req.TenantId,
		AppId:     req.AppId,
		CommandId: req.CommandId,
	}

	data, err := c.grpcClient.GetAppEventLogByCommandId(ctx, request)
	if err != nil {
		return nil, err
	}

	var list []*EventLogDto
	for _, item := range data.Data {
		dto := &EventLogDto{
			Id:        item.Id,
			TenantId:  item.TenantId,
			AppId:     item.AppId,
			Class:     item.Class,
			Func:      item.Func,
			Level:     item.Level,
			Time:      timeutils.ToPTime(item.Time.AsTime()),
			Status:    item.Status,
			Message:   item.Message,
			PubAppId:  item.PubAppId,
			EventId:   item.EventId,
			CommandId: item.CommandId,
		}
		list = append(list, dto)
	}
	response := &GetEventLogByCommandIdResponse{
		Data: list,
	}
	return response, nil
}

func (c *daprClient) WriteAppLog(ctx context.Context, req *WriteAppLogRequest) (resp *WriteAppLogResponse, resErr error) {
	request := &pb.WriteAppLogRequest{
		Headers:  newRequstHeaders(&req.Headers),
		TenantId: req.TenantId,
		Id:       req.Id,
		AppId:    req.AppId,
		Class:    req.Class,
		Func:     req.Func,
		Level:    req.Level,
		Time:     timeutils.AsTimestamp(req.Time),
		Status:   req.Status,
		Message:  req.Message,
	}
	_, err := c.grpcClient.WriteAppLog(ctx, request)
	if err != nil {
		return nil, err
	}
	resp = &WriteAppLogResponse{}
	return resp, nil
}

func (c *daprClient) UpdateAppLog(ctx context.Context, req *UpdateAppLogRequest) (resp *UpdateAppLogResponse, resErr error) {
	request := &pb.UpdateAppLogRequest{
		Headers:  newRequstHeaders(&req.Headers),
		TenantId: req.TenantId,
		Id:       req.Id,
		AppId:    req.AppId,
		Class:    req.Class,
		Func:     req.Func,
		Level:    req.Level,
		Time:     timeutils.AsTimestamp(req.Time),
		Status:   req.Status,
		Message:  req.Message,
	}
	_, err := c.grpcClient.UpdateAppLog(ctx, request)
	if err != nil {
		return nil, err
	}
	resp = &UpdateAppLogResponse{}
	return resp, nil
}

func (c *daprClient) GetAppLogById(ctx context.Context, req *GetAppLogByIdRequest) (resp *GetAppLogByIdResponse, resErr error) {
	request := &pb.GetAppLogByIdRequest{
		Headers:  newRequstHeaders(&req.Headers),
		TenantId: req.TenantId,
		Id:       req.Id,
	}
	data, err := c.grpcClient.GetAppLogById(ctx, request)
	if err != nil {
		return nil, err
	}
	resp = &GetAppLogByIdResponse{
		Id:       data.Id,
		TenantId: data.TenantId,
		AppId:    data.AppId,
		Class:    data.Class,
		Func:     data.Func,
		Level:    data.Level,
		Time:     timeutils.ToPTime(data.Time.AsTime()),
		Status:   data.Status,
		Message:  data.Message,
	}
	return resp, nil
}
