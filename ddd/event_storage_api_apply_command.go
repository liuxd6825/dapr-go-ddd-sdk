package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"reflect"
	"strings"
)

type CreateAggregateOptions struct {
	eventStorageKey *string
}

type ApplyCommandOptions struct {
	EventStorageKey string
}

func NewApplyCommandOptions() *ApplyCommandOptions {
	return &ApplyCommandOptions{
		EventStorageKey: strEmpty,
	}
}

func (o *ApplyCommandOptions) Merge(opts ...*ApplyCommandOptions) *ApplyCommandOptions {
	for _, item := range opts {
		if len(item.EventStorageKey) != 0 {
			o.EventStorageKey = item.EventStorageKey
		}
	}
	return o
}

func (o *ApplyCommandOptions) SetEventStorageKey(v string) *ApplyCommandOptions {
	o.EventStorageKey = v
	return o
}

//
// ApplyCommand
// @Description: 执行聚合命令
// @param ctx
// @param aggregate
// @param cmd
// @param opts
// @return error
//
func ApplyCommand(ctx context.Context, agg Aggregate, cmd Command, opts ...*ApplyCommandOptions) (err error) {
	if agg == nil {
		return errors.ErrorOf("ApplyCommand(ctx, agg, cmd) error: agg is nil")
	}
	if cmd == nil {
		return errors.ErrorOf("ApplyCommand(ctx, agg, cmd) error: cmd is nil")
	}
	opt := NewApplyCommandOptions()
	opt.Merge(opts...)

	if _, ok := cmd.(IsAggregateCreateCommand); ok {
		return callCommandHandler(ctx, agg, cmd)
	}
	if ok := isAggregateCreateCommand(ctx, agg, cmd); ok {
		return callCommandHandler(ctx, agg, cmd)
	}

	loadOpt := NewLoadAggregateOptions().SetEventStorageKey(opt.EventStorageKey)
	aggId := cmd.GetAggregateId().RootId()
	_, find, err := LoadAggregate(ctx, cmd.GetTenantId(), aggId, agg, loadOpt)
	if err != nil {
		return err
	}
	if !find {
		return errors.NewAggregateIdNotFondError(aggId)
	}
	return callCommandHandler(ctx, agg, cmd)
}

func isAggregateCreateCommand(ctx context.Context, aggregate Aggregate, cmd Command) bool {
	aggName := reflect.TypeOf(aggregate).Name()
	cmdName := reflect.TypeOf(cmd).Name()
	isAgg := strings.HasPrefix(aggName, cmdName)
	isCreate := strings.HasSuffix("CreateCommand", cmdName)
	if isAgg && isCreate {
		return true
	}
	return false
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
