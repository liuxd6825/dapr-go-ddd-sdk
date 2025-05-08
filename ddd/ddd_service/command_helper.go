package ddd_service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
)

var validate = validator.New()

type CommandHelper[T any] struct {
	null        T
	aggTypeName string
}

func NewCommandHelper[T any](aggTypeName string) *CommandHelper[T] {
	return &CommandHelper[T]{
		aggTypeName: aggTypeName,
	}
}

func (s *CommandHelper[T]) ValidateCommand(cmd interface{}) error {
	if cmd == nil {
		return errors.New("command is nil")
	}
	if err := validate.Struct(cmd); err != nil {
		return err
	}
	return nil
}

func (s *CommandHelper[T]) NewAggregate() (T, error) {
	return reflectutils.NewStruct[T]()
}

// DoCommand
// @Description:
// @receiver s
// @param ctx
// @param cmd
// @return *model.UserAggregate
// @return error
func (s *CommandHelper[T]) DoCommand(ctx context.Context, cmd ddd.Command, validateFunc func() error, opts ...ddd.DoCommandOption) (res T, resErr error) {
	defer func() {
		resErr = errors.GetRecoverError(resErr, recover())
	}()

	option := ddd.NewDoCommandOptionMerges(opts...)

	// 进行业务检查
	if validateFunc != nil {
		if err := validateFunc(); err != nil {
			return s.null, err
		}
	} else if err := cmd.Validate(); err != nil {
		return s.null, err
	}

	// 如果只是业务检查，则不执行领域命令，
	validOnly := option.GetIsValidOnly()
	if (validOnly == nil && cmd.GetIsValidOnly()) || (validOnly != nil && *validOnly == true) {
		return s.null, nil
	}

	// 新建聚合根对象
	agg, err := s.NewAggregate()
	if err != nil {
		return s.null, err
	}

	// 如果领域命令执行时出错，则返回错误
	if err := ddd.ApplyCommand(ctx, agg, cmd); err != nil {
		return s.null, err
	}

	return agg, nil
}

// GetAggregateById
// @Description: 获取聚合对象
// @receiver s
// @param ctx 上下文
// @param tenantId 租户id
// @param id 主键id
// @return *user_model.UserCommandDomainService  聚合对象
// @return bool 是否找到聚合根对象
// @return error 错误对象
func (s *CommandHelper[T]) GetAggregateById(ctx context.Context, tenantId string, id string) (T, bool, error) {
	agg, err := s.NewAggregate()
	if err != nil {
		return s.null, false, err
	}
	_, ok, err := ddd.LoadAggregate(ctx, tenantId, id, agg)
	return agg, ok, err
}
