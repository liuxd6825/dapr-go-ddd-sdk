package userlog

import "time"

type LoginEvent struct {
	CommandId     string    `json:"commandId"`
	TenantId      string    `json:"tenantId"`
	EventId       string    `json:"eventId"`
	EventType     string    `json:"eventType"`     // 事件类型
	EventVersion  string    `json:"eventVersion"`  // 事件版本号
	AggregateId   string    `json:"aggregateId"`   // 聚合根Id
	AggregateType string    `json:"aggregateType"` // 聚合类型
	CreatedTime   time.Time `json:"createdTime"`   // 创建时间
	Data          LoginData `json:"data"`
}

type LoginData struct {
	Id       string    `json:"id"`
	UserId   string    `json:"userId"`
	UserName string    `json:"userName"`
	Date     time.Time `json:"date"`
}

const UserLoginEventType = "system.UserLoginEventType"
const UserLoginEventVersion = "v1.0"

func NewLoginEvent(commandId string, logId string, userId, userName string, logTime time.Time) *LoginEvent {
	return &LoginEvent{
		CommandId: commandId,
		TenantId:  SystemTenantId,

		EventType:    UserLoginEventType,
		EventVersion: UserLoginEventVersion,

		AggregateId:   newAggregateId(userId),
		AggregateType: AggregateType,

		CreatedTime: logTime,
		Data: LoginData{
			Id:       logId,
			UserId:   userId,
			UserName: userName,
			Date:     logTime,
		},
	}
}

func (l *LoginEvent) GetTenantId() string {
	return l.TenantId
}

func (l *LoginEvent) GetCommandId() string {
	return l.CommandId
}

func (l *LoginEvent) GetEventId() string {
	return l.EventId
}

func (l *LoginEvent) GetEventType() string {
	return UserLogoutEventType
}

func (l *LoginEvent) GetEventVersion() string {
	return UserLogoutEventEventVersion
}

func (l *LoginEvent) GetAggregateId() string {
	return l.AggregateId
}

func (l *LoginEvent) GetAggregateType() string {
	return l.AggregateType
}

func (l *LoginEvent) GetCreatedTime() time.Time {
	return l.CreatedTime
}

func (l *LoginEvent) GetData() any {
	return l.Data
}

func (l *LoginEvent) GetIsSourcing() bool {
	return false
}
