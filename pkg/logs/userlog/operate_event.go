package userlog

import (
	"context"
	"time"

	appctx2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
)

type OperateEvent struct {
	CommandId   string       `json:"commandId"`
	TenantId    string       `json:"tenantId"`
	EventId     string       `json:"eventId"`
	EventType   string       `json:"eventType"`   // 事件类型
	EventVer    string       `json:"eventVer"`    // 事件版本号
	AggId       string       `json:"aggId"`       // 聚合根Id
	AggType     string       `json:"aggType"`     // 聚合类型
	CreatedTime time.Time    `json:"createdTime"` // 创建时间
	Data        *OperateData `json:"data"`
}

type OperateData struct {
	Id         string    `json:"id"`
	TenantId   string    `json:"tenantId"`
	AppId      string    `json:"appId"`
	AppName    string    `json:"appName"`
	UserId     string    `json:"userId"`
	UserName   string    `json:"userName"`
	Time       time.Time `json:"time"`
	ModelName  string    `json:"modelName"`
	Message    string    `json:"message"`
	ActionType string    `json:"actionType"`
}

const OperateEventType = "system.UserOperateLogEventType"
const OperateEventVersion = "v1.0"

type DomainEvent interface {
	GetData() any // 事件数据
}

func NewOperateEvent(ctx context.Context, commandId string, eventId string, oData *OperateData) *OperateEvent {
	event := &OperateEvent{
		CommandId:   commandId,
		TenantId:    oData.TenantId,
		EventId:     eventId,
		CreatedTime: time.Now(),
		AggId:       newAggregateId(oData.UserId),
		AggType:     AggregateType,
		EventType:   OperateEventType,
		EventVer:    OperateEventVersion,
		Data:        oData,
	}
	return event
}

func NewOperateData(ctx context.Context, modelName string, actionType string, data any) (res *OperateData, err error) {
	appId, appName := getAppIdName(ctx)
	msg, err := reflectutils.ParseDesc(data)
	if err != nil {
		return nil, err
	}
	operateData, err := newOperateData(ctx, appId, appName, actionType, modelName, msg)
	return operateData, err
}

func newOperateData(ctx context.Context, appId string, appName string, actionType string, modelName string, message string) (*OperateData, error) {
	authUser, ok := appctx2.GetAuthUser(ctx)
	if !ok {
		return nil, errors.ErrNotFoundLoginUser
	}

	tenantId, ok := appctx2.GetTenantId(ctx)
	if !ok && authUser != nil {
		tenantId = authUser.GetTenantId()
	}
	data := &OperateData{}
	data.Id = idutils.NewId()
	data.Time = time.Now()
	data.TenantId = tenantId
	data.ModelName = modelName
	data.ActionType = actionType
	data.AppId = appId
	data.AppName = appName
	data.Message = message
	if authUser != nil {
		data.UserId = authUser.GetId()
		data.UserName = authUser.GetName()
	}
	return data, nil
}

func (d *OperateData) GetId() string {
	return d.Id
}

func (d *OperateData) GetTime() time.Time {
	return d.Time
}

func (d *OperateData) GetUserId() string {
	return d.UserId
}

func (d *OperateData) GetUserName() string {
	return d.UserName
}

func (d *OperateData) GetMessage() string {
	return d.Message
}

func (d *OperateData) GetActionType() string {
	return d.ActionType
}

func (d *OperateData) GetAppId() string {
	return d.AppId
}

func (d *OperateData) GetTenantId() string {
	return d.TenantId
}

func (l *OperateEvent) GetTenantId() string {
	return l.TenantId
}

func (l *OperateEvent) GetCommandId() string {
	return l.CommandId
}

func (l *OperateEvent) GetEventId() string {
	return l.EventId
}

func (l *OperateEvent) GetEventType() string {
	return l.EventType
}

func (l *OperateEvent) GetEventVer() string {
	return l.EventVer
}

func (l *OperateEvent) GetAggId() string {
	return l.AggId
}

func (l *OperateEvent) GetAggType() string {
	return l.AggType
}

func (l *OperateEvent) GetCreatedTime() time.Time {
	return l.CreatedTime
}

func (l *OperateEvent) GetData() any {
	return l.Data
}

func (l *OperateEvent) GetIsSourcing() bool {
	return false
}

func getAppIdName(ctx context.Context) (string, string) {
	appId := DefaultAppId
	appName := DefaultAppName
	appInfo, ok := appctx2.GetAppInfo(ctx)
	if ok {
		appId = appInfo.GetAppId()
		appName = appInfo.GetAppName()
	}
	return appId, appName
}
