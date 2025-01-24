package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

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
func (s *Service) CallEventHandler(ctx context.Context, eventHandler any, eventType string, eventVersion string, event any, metadata ddd.Metadata) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	/*
		funcName := fmt.Sprintf(eventType, ".", eventVersion)
		subEvent, ok := s.subEvents.Get(funcName)
		if !ok {
			return fmt.Errorf("event handler for event %s not found", eventType)
		}*/

	return nil
}
