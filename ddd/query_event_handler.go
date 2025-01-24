package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"strings"
)

// QueryEventHandler
// @Description: CQRS模型中的Query端事件处理器
type QueryEventHandler interface {
	CallEventHandler(ctx context.Context, handler any, eventType string, eventVersion string, event any, metadata Metadata) error
}

// QueryEventHandlerDefault
// @Description: 默认的事件处理器， 通过反射查找handler中的方法名称
type QueryEventHandlerDefault struct {
	handler any
}

func NewQueryEventHandlerDefault(handler any) QueryEventHandler {
	return &QueryEventHandlerDefault{handler: handler}
}

// getEventMethodName
// @Description: 根据事件类型名称获取接受事件方法名称
// @param eventType
// @param revision
// @return string
func (q *QueryEventHandlerDefault) getEventMethodName(eventType string, revision string) string {
	names := strings.Split(eventType, ".")
	name := names[len(names)-1]
	ver := strings.Replace(revision, ".", "s", -1)
	if strings.HasPrefix(ver, "v") || strings.HasPrefix(ver, "V") {
		ver = "V" + ver[1:]
	} else {
		ver = "V" + ver
	}
	return fmt.Sprintf("On%s%s", name, ver)
}

// CallEventHandler
//
//	@Description: 调用QueryHandler事件订阅处理器
//	@param ctx  上下文
//	@param queryHandler 事件处理器
//	@param eventType  事件类型
//	@param eventVersion 事件版本号
//	@param event 事件对象
//	@param Metadata  事件元数据
//	@return error 错误
func (q *QueryEventHandlerDefault) CallEventHandler(ctx context.Context, handler any, eventType string, eventVersion string, event any, metadata Metadata) error {
	methodName := getEventMethodName(eventType, eventVersion)
	return reflectutils.CallMethod(handler, methodName, ctx, event, metadata)
}
