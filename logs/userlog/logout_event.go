package userlog

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"time"
)

type LogoutEvent struct {
	CommandId     string     `json:"commandId"`
	TenantId      string     `json:"tenantId"`
	EventId       string     `json:"eventId"`
	EventType     string     `json:"eventType"`     // 事件类型
	EventVersion  string     `json:"eventVersion"`  // 事件版本号
	AggregateId   string     `json:"aggregateId"`   // 聚合根Id
	AggregateType string     `json:"aggregateType"` // 聚合类型
	CreatedTime   time.Time  `json:"createdTime"`   // 创建时间
	Data          LogoutData `json:"data"`
}

type LogoutData struct {
	Id       string    `json:"id"`
	UserId   string    `json:"userId"`
	UserName string    `json:"userName"`
	Date     time.Time `json:"date"`
}

const UserLogoutEventType = "system.UserLogoutEventType"
const UserLogoutEventEventVersion = "v1.0"

func NewLogoutEvent(commandId string, logId string, userId, userName string, logTime time.Time) *LogoutEvent {
	return &LogoutEvent{
		CommandId: commandId,

		TenantId:     SystemTenantId,
		EventId:      idutils.NewId(),
		EventType:    UserLoginEventType,
		EventVersion: UserLoginEventVersion,

		AggregateId:   newAggregateId(userId),
		AggregateType: AggregateType,

		CreatedTime: logTime,
		Data: LogoutData{
			Id:       logId,
			UserId:   userId,
			UserName: userName,
			Date:     logTime,
		},
	}
}

func (l *LogoutEvent) GetTenantId() string {
	return l.TenantId
}

func (l *LogoutEvent) GetCommandId() string {
	return l.CommandId
}

func (l *LogoutEvent) GetEventId() string {
	return l.EventId
}

func (l *LogoutEvent) GetEventType() string {
	return UserLogoutEventType
}

func (l *LogoutEvent) GetEventVersion() string {
	return UserLogoutEventEventVersion
}

func (l *LogoutEvent) GetAggregateId() string {
	return l.AggregateId
}

func (l *LogoutEvent) GetAggregateType() string {
	return l.AggregateType
}

func (l *LogoutEvent) GetCreatedTime() time.Time {
	return l.CreatedTime
}

func (l *LogoutEvent) GetData() any {
	return l.Data
}

func (l *LogoutEvent) GetIsSourcing() bool {
	return false
}
