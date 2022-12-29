package ddd

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
)

var validate = validator.New()

type CommandDomainService[T Aggregate] struct {
	newAggregate func() T
}

func NewCommandDomainService[T Aggregate]() *CommandDomainService[T] {
	return &CommandDomainService[T]{}
}

func (s *CommandDomainService[T]) ValidateCommand(cmd interface{}) error {
	if cmd == nil {
		return errors.New("command is nil")
	}
	if err := validate.Struct(cmd); err != nil {
		return err
	}
	return nil
}

//
// DoCommand
// @Description: fcj
// @receiver s
// @param ctx
// @param cmd
// @return *model.SolutionAggregate
// @return error
//
func (s *CommandDomainService[T]) DoCommand(ctx context.Context, cmd Command, opts ...DoCommandOption) (T, error) {
	var null T
	option := NewDoCommandOptionMerges(opts...)

	// 进行业务检查
	if err := cmd.Validate(); err != nil {
		return null, err
	}

	// 如果只是业务检查，则不执行领域命令，
	validOnly := option.GetIsValidOnly()
	if (validOnly == nil && cmd.GetIsValidOnly()) || (validOnly != nil && *validOnly == true) {
		return null, nil
	}

	// 新建聚合根对象
	agg, err := s.NewAggregate()
	if err != nil {
		return null, err
	}

	// 如果领域命令执行时出错，则返回错误
	if err := ApplyCommand(ctx, agg, cmd); err != nil {
		return null, err
	}

	return agg, nil
}

//
// GetAggregateById
// @Description: 获取聚合对象
// @receiver s
// @param ctx 上下文
// @param tenantId 租户id
// @param id 主键id
// @return *graph_model.SolutionCommandDomainService  聚合对象
// @return bool 是否找到聚合根对象
// @return error 错误对象
//
func (s *CommandDomainService[T]) GetAggregateById(ctx context.Context, tenantId string, id string) (T, bool, error) {
	var null T
	agg, err := s.NewAggregate()
	if err != nil {
		return null, false, err
	}
	_, ok, err := LoadAggregate(ctx, tenantId, id, agg)
	return agg, ok, err
}

func (s *CommandDomainService[T]) NewAggregate() (T, error) {
	return reflectutils.NewStruct[T]()
}
