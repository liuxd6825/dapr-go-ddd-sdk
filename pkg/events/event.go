package events

import (
	"context"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

// IEvent
// @Description: 领域事件接口
type IEvent interface {
	GetId() string
	GetTenantId() string
	GetAppId() string
	GetData() any
	GetMeta() map[string]any
	GetCreatedTime() *time.Time
	GetVersion() string
	GetEventType() string
	SetFields(id, tenantId, appId, eventType, version string, createTime *time.Time, data any, meta map[string]any)
	NewData() any
}

// Event
// @Description: 领域事件
type Event[T any] struct {
	Id          string         `json:"id" gorm:"id;primaryKey" bson:"id" `                                    // ID
	TenantId    string         `json:"tenantId" gorm:"tenant_id" bson:"tenant_id"`                            // 租户ID
	AppId       string         `json:"appId" gorm:"app_id" bson:"app_id"`                                     // 应用ID
	EventType   string         `json:"eventType" gorm:"event_type" bson:"event_type" `                        // 事件类型
	Data        T              `json:"data" gorm:"data;json" bson:"data"`                                     // 事件数据
	Meta        map[string]any `json:"meta" gorm:"meta;json" bson:"meta"`                                     // 上下文数据
	CreatedTime *time.Time     `json:"createdTime" gorm:"created_time;index:,sort:desc," bson:"created_time"` // 发生事件
	Version     string         `json:"version" gorm:"version" bson:"version"  `                               // 事件数据版本
}

type EventOptions struct {
	TenantId    string         `json:"tenantId"`
	AppId       string         `json:"appId"`
	EventType   string         `json:"eventType"`
	Meta        map[string]any `json:"meta"`
	CreatedTime *time.Time     `json:"createdTime"`
	Version     string         `json:"version"`
}

func (e *Event[T]) SetData(ctx context.Context, appId string, eventData T, options ...*EventOptions) (event *Event[T]) {
	e.setOptions(options...)
	if e.Id == "" {
		e.Id = idutils.NewId()
	}
	if e.TenantId == "" {
		tenantId, err := appctx.GetTenantId3(ctx)
		if err != nil {
			panic(err)
		}
		e.TenantId = tenantId
	}

	if e.Meta == nil {
		meta, err := newMeta(ctx, e.TenantId, e.Meta)
		if err != nil {
			panic(err)
		}
		e.Meta = meta
	}

	if e.EventType == "" {
		e.EventType = GetEventType(eventData)
	}

	e.EventType = stringutils.MidlineString(e.EventType)

	if e.CreatedTime == nil {
		createTime := time.Now()
		e.CreatedTime = &createTime
	}

	e.AppId = appId
	e.Data = eventData

	return event
}

func (e *Event[T]) setOptions(opts ...*EventOptions) {
	for _, o := range opts {
		if o == nil {
			continue
		}
		if o.TenantId != "" {
			e.TenantId = o.TenantId
		}
		if o.AppId != "" {
			e.AppId = o.AppId
		}
		if o.EventType != "" {
			e.EventType = o.EventType
		}
		if o.CreatedTime != nil {
			e.CreatedTime = o.CreatedTime
		}
		if o.Version != "" {
			e.Version = o.Version
		}
		if o.Meta != nil {
			e.Meta = o.Meta
		}
	}
}

func (e *Event[T]) GetData() any {
	return e.Data
}

func (e *Event[T]) GetMeta() map[string]any {
	return e.Meta
}

func (e *Event[T]) GetCreatedTime() *time.Time {
	return e.CreatedTime
}

func (e *Event[T]) GetVersion() string {
	return e.Version
}

func (e *Event[T]) GetEventType() string {
	return e.EventType
}

func (e *Event[T]) GetId() string {
	return e.Id
}

func (e *Event[T]) GetEventId() string {
	return e.Id
}

func (e *Event[T]) GetTenantId() string {
	return e.TenantId
}

func (e *Event[T]) GetAppId() string {
	return e.AppId
}

func (e *Event[T]) SetFields(id, tenantId, appId, eventType, version string, createTime *time.Time, data any, meta map[string]any) {
	e.Id = id
	e.TenantId = tenantId
	e.AppId = appId
	e.EventType = eventType
	e.Version = version
	e.Meta = meta
	e.CreatedTime = createTime
	if vData, ok := data.(T); ok {
		e.Data = vData
	} else if vData, ok := data.(*T); ok {
		e.Data = *vData
	} else {
		panic("event.Data is not T")
	}
}

func (e *Event[T]) NewData() any {
	return reflectutils.NewInstance[T]()
}
