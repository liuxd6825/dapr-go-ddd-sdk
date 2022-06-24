package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"reflect"
)

type CreateAggregateOptions struct {
	eventStorageKey *string
}

//
// CreateAggregate
// @Description: 创建聚合根
// @param ctx
// @param aggregate
// @param cmd
// @param opts
// @return error
//
func CreateAggregate(ctx context.Context, aggregate Aggregate, cmd Command, opts ...*CreateAggregateOptions) error {
	options := &CreateAggregateOptions{
		eventStorageKey: &strEmpty,
	}
	for _, item := range opts {
		if item.eventStorageKey != nil {
			options.eventStorageKey = item.eventStorageKey
		}
	}
	return callCommandHandler(ctx, aggregate, cmd)
}

//
// CommandAggregate
// @Description: 执行聚合命令
// @param ctx
// @param aggregate
// @param cmd
// @param opts
// @return error
//
func CommandAggregate(ctx context.Context, aggregate Aggregate, cmd Command, opts ...LoadAggregateOption) error {
	aggId := cmd.GetAggregateId().RootId()
	_, find, err := LoadAggregate(ctx, cmd.GetTenantId(), aggId, aggregate, opts...)
	if err != nil {
		return err
	}
	if !find {
		return ddd_errors.NewAggregateIdNotFondError(aggId)
	}
	return callCommandHandler(ctx, aggregate, cmd)
}

func (o *CreateAggregateOptions) SetEventStorageKey(eventStorageKey string) {
	o.eventStorageKey = &eventStorageKey
}

func callCommandHandler(ctx context.Context, aggregate Aggregate, cmd Command) error {
	cmdTypeName := reflect.ValueOf(cmd).Elem().Type().Name()
	methodName := fmt.Sprintf("%s", cmdTypeName)
	metadata := ddd_context.GetMetadataContext(ctx)
	return CallMethod(aggregate, methodName, ctx, cmd, metadata)
}
