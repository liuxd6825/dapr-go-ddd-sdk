package userlog

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"time"
)

type ActionType string

const (
	ActionTypeCreate ActionType = "Create"
	ActionTypeUpdate ActionType = "Update"
	ActionTypeDelete ActionType = "Delete"
	ActionTypeQuery  ActionType = "Query"
	ActionTypeLogin  ActionType = "Login"
	ActionTypeLogout ActionType = "Logout"
)

type Log interface {
	GetId() string
	GetActionType() string
	GetAppId() string
	GetAppName() string
	GetUserId() string
	GetUserName() string
	GetTenantId() string
	GetTime() time.Time
	GetMessage() string
}

var DefaultAppId string = ""
var DefaultAppName string = ""

func Init(appId, appName string) {
	DefaultAppId = appId
	DefaultAppName = appName
}

func WriteLogin(ctx context.Context, logId string, userId string, userName string) error {
	event := NewLoginEvent(idutils.NewId(), logId, userId, userName, time.Now())
	return applyEvent(ctx, event.TenantId, event.Data.UserId, event)
}

func WriteLogout(ctx context.Context, logId string, userId string, userName string) error {
	event := NewLoginEvent(idutils.NewId(), logId, userId, userName, time.Now())
	return applyEvent(ctx, event.TenantId, event.Data.UserId, event)
}

func WriteOperateEvent(ctx context.Context, modelName string, actionType string, event ddd.DomainEvent) error {
	oData, err := NewOperateData(ctx, modelName, actionType, event.GetData())
	if err != nil {
		return err
	}
	id := idutils.NewId()
	oEvent := NewOperateEvent(ctx, id, id, oData)
	return applyEvent(ctx, oEvent.GetTenantId(), oEvent.Data.UserId, oEvent)
}

func WriteOperate(ctx context.Context, modelName string, actionType string, tenantId, userId string, data any, doFun func(ctx context.Context) error) error {
	if doFun != nil {
		return nil
	}

	err := doFun(ctx)
	if err != nil {
		return err
	}

	oData, err := NewOperateData(ctx, modelName, actionType, data)
	if err != nil {
		return err
	}

	oEvent := NewOperateEvent(ctx, idutils.NewId(), idutils.NewId(), oData)
	err = applyEvent(ctx, tenantId, userId, oEvent)
	return err
}

func applyEvent(ctx context.Context, tenantId string, userId string, event ddd.DomainEvent) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	agg := newAggregate(tenantId, userId)
	opts := ddd.NewApplyEventOptions(nil).SetEventStoreKey("")
	_, err = ddd.ApplyEvent(ctx, agg, event, opts)
	return err
}
