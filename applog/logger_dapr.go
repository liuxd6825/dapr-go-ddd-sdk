package applog

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/httpclient"
)

type logger struct {
	httpclient *httpclient.HttpClient
}

func NewLogger(httpclient *httpclient.HttpClient) Logger {
	return &logger{
		httpclient: httpclient,
	}
}

func (l *logger) WriteEventLog(ctx context.Context, req *WriteEventLogRequest) (resp *WriteEventLogResponse, resErr error) {
	resp = &WriteEventLogResponse{}
	l.httpclient.Post(ctx, ApiWriteEventLog, req).OnSuccess(resp, func() error {
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (l *logger) UpdateEventLog(ctx context.Context, req *UpdateEventLogRequest) (resp *UpdateEventLogResponse, resErr error) {
	data := &UpdateEventLogResponse{}
	l.httpclient.Post(ctx, ApiUpdateEventLog, req).OnSuccess(data, func() error {
		resp = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (l *logger) GetEventLogByCommandId(ctx context.Context, req *GetEventLogByCommandIdRequest) (resp *GetEventLogByCommandIdResponse, resErr error) {
	url := fmt.Sprintf(ApiGetEventLogByCommandId, req.TenantId, req.AppId, req.CommandId)
	data := &GetEventLogByCommandIdResponse{}
	l.httpclient.Get(ctx, url).OnSuccess(data, func() error {
		resp = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (l *logger) WriteAppLog(ctx context.Context, req *WriteAppLogRequest) (resp *WriteAppLogResponse, resErr error) {
	data := &WriteAppLogResponse{}
	l.httpclient.Post(ctx, ApiWriteAppLog, req).OnSuccess(data, func() error {
		resp = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (l *logger) UpdateAppLog(ctx context.Context, req *UpdateAppLogRequest) (resp *UpdateAppLogResponse, resErr error) {
	data := &UpdateAppLogResponse{}
	l.httpclient.Post(ctx, ApiUpdateAppLog, req).OnSuccess(resp, func() error {
		resp = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}

func (l *logger) GetAppLogById(ctx context.Context, req *GetAppLogByIdRequest) (resp *GetAppLogByIdResponse, resErr error) {
	data := &GetAppLogByIdResponse{}
	l.httpclient.Get(ctx, fmt.Sprintf(ApiGetAppLogById, req.TenantId, req.Id)).OnSuccess(data, func() error {
		resp = data
		return nil
	}).OnError(func(err error) {
		resErr = err
	})
	return
}
