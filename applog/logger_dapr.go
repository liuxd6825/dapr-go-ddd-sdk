package applog

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

type logger struct {
	daprClient daprclient.DaprDddClient
	level      Level
}

func NewLogger(daprClient daprclient.DaprDddClient) Logger {
	return &logger{
		daprClient: daprClient,
		level:      logs.ErrorLevel,
	}
}

func (l *logger) WriteEventLog(ctx context.Context, req *WriteEventLogRequest) (resp *WriteEventLogResponse, resErr error) {
	return l.daprClient.WriteEventLog(ctx, req)
}

func (l *logger) UpdateEventLog(ctx context.Context, req *UpdateEventLogRequest) (resp *UpdateEventLogResponse, resErr error) {
	return l.daprClient.UpdateEventLog(ctx, req)
}

func (l *logger) GetEventLogByCommandId(ctx context.Context, req *GetEventLogByCommandIdRequest) (resp *GetEventLogByCommandIdResponse, resErr error) {
	return l.daprClient.GetEventLogByCommandId(ctx, req)
}

func (l *logger) WriteAppLog(ctx context.Context, req *WriteAppLogRequest) (resp *WriteAppLogResponse, resErr error) {
	return l.daprClient.WriteAppLog(ctx, req)
}

func (l *logger) UpdateAppLog(ctx context.Context, req *UpdateAppLogRequest) (resp *UpdateAppLogResponse, resErr error) {
	return l.daprClient.UpdateAppLog(ctx, req)
}

func (l *logger) GetAppLogById(ctx context.Context, req *GetAppLogByIdRequest) (resp *GetAppLogByIdResponse, resErr error) {
	return l.daprClient.GetAppLogById(ctx, req)
}

func (l *logger) SetLevel(level Level) {
	l.level = level
}

func (l *logger) GetLevel() Level {
	return l.level
}
