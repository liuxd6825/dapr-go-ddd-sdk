package userlog

import "time"

type LoginEvent struct {
	CommandId   string    `json:"commandId"`
	TenantId    string    `json:"tenantId"`
	EventId     string    `json:"eventId"`
	EventType   string    `json:"eventType"`   // 事件类型
	EventVer    string    `json:"eventVer"`    // 事件版本号
	AggId       string    `json:"aggId"`       // 聚合根Id
	AggType     string    `json:"aggType"`     // 聚合类型
	CreatedTime time.Time `json:"createdTime"` // 创建时间
	Data        LoginData `json:"data"`
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
		CommandId:   commandId,
		TenantId:    SystemTenantId,
		EventType:   UserLoginEventType,
		EventVer:    UserLoginEventVersion,
		AggId:       newAggregateId(userId),
		AggType:     AggregateType,
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

func (l *LoginEvent) GetEventVer() string {
	return UserLogoutEventEventVersion
}

func (l *LoginEvent) GetAggId() string {
	return l.AggId
}

func (l *LoginEvent) GetAggType() string {
	return l.AggType
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
