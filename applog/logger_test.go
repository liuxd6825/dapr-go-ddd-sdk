package applog

import (
	"context"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/httpclient"
	"testing"
	"time"
)

func TestLogger_WriteEventLog(t *testing.T) {
	logger := newLogger()
	ctx := context.Background()
	uid := uuid.New()
	writeReq := &WriteEventLogRequest{
		Id:        uid.String(),
		AppId:     "test_subAppId",
		Class:     "test",
		Func:      "TestLogger_WriteEventLog",
		Level:     "info",
		TenantId:  "test",
		Time:      newTime(),
		Status:    false,
		Message:   "test message-create",
		PubAppId:  "test_pubAppId",
		EventId:   uid.String(),
		CommandId: uid.String(),
	}
	createResp, err := logger.WriteEventLog(ctx, writeReq)
	if err != nil {
		t.Error(err)
	}

	println(createResp)

	updateReq := &UpdateEventLogRequest{
		Id:        uid.String(),
		AppId:     "test_subAppId",
		Class:     "test",
		Func:      "TestLogger_WriteEventLog",
		Level:     "info",
		TenantId:  "test",
		Time:      newTime(),
		Status:    false,
		Message:   "test message-update",
		PubAppId:  "test_pubAppId",
		EventId:   uid.String(),
		CommandId: uid.String(),
	}
	updateResp, err := logger.UpdateEventLog(ctx, updateReq)
	if err != nil {
		t.Error(err)
	}
	println(updateResp)

	getReq := &GetEventLogByCommandIdRequest{
		TenantId:  "test",
		AppId:     "test_subAppId",
		CommandId: uid.String(),
	}
	getResp, err := logger.GetEventLogByCommandId(ctx, getReq)
	if err != nil {
		t.Error(err)
	}
	println(getResp)
}

func newLogger() Logger {
	client, _ := httpclient.NewHttpClient("localhost", 9011)
	logger := NewLogger(client)
	return logger
}

func newTime() *time.Time {
	t := time.Now()
	return &t
}
