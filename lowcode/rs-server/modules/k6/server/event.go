package server

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"time"
)

// EventHandler 事件处理器
type EventHandler struct {
	VM         *goja.Runtime
	Object     *goja.Object
	subscribes []*Subscribe // 处理器所订阅的消息列表
}

// CallEventHandler
//
//	@Description: 执行消息处理函数
//	@receiver h
//	@param ctx
//	@param handler 事件处理器
//	@param eventType    事件类型
//	@param eventVersion 事件版本
//	@param event 事件数据
//	@param metadata 事件元数据
//	@return err
func (h *EventHandler) CallEventHandler(ctx context.Context, handler any, eventType string, eventVersion string, event any, metadata ddd.Metadata) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	methodName := h.GetMethodName(eventType)
	v := h.Object.Get(methodName)
	if goja.IsUndefined(v) || goja.IsNaN(v) || goja.IsNull(v) {
		return errors.New("method: %s is not defined", methodName)
	}
	fun, ok := goja.AssertFunction(v)
	if !ok {
		return errors.New("rs-server subscribe handler not fount %s()", methodName)
	}
	c := h.VM.ToValue(ctx)
	e := h.VM.ToValue(event)
	m := h.VM.ToValue(metadata)
	_, vErr := fun(h.Object, c, e, m)
	if vErr != nil {
		return vErr
	}
	return nil
}

func (h *EventHandler) GetMethodName(eventType string) string {
	for _, s := range h.subscribes {
		if s.Topic == eventType {
			return s.FuncName
		}
	}
	return ""
}

// Event 领域事件
type Event map[string]any

// GetTenantId 租户Id
func (e Event) GetTenantId() string {
	return common.Object(e).GetString(common.TenantId)
}

// GetCommandId 命令Id
func (e Event) GetCommandId() string {
	return common.Object(e).GetString("commandId")
}

// GetEventId 事件Id
func (e Event) GetEventId() string {
	return common.Object(e).GetString("eventId")
}

// GetEventType 事件类型
func (e Event) GetEventType() string {
	return common.Object(e).GetString("eventType")
}

// GetEventVersion 事件版本号
func (e Event) GetEventVersion() string {
	return common.Object(e).GetString("eventVersion")
}

// GetAggregateId 聚合根Id
func (e Event) GetAggregateId() string {
	return common.Object(e).GetString("aggregateId")
}

// GetCreatedTime 创建时间
func (e Event) GetCreatedTime() time.Time {
	t, _ := common.Object(e).GetDateTime("createdTime")
	return t
}

// GetData 事件数据
func (e Event) GetData() interface{} {
	return common.Object(e).Get("data")
}
