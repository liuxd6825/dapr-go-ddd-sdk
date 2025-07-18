package xbase

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type Service struct {
}

type ValidateObject interface {
	Validate() error
}

type DoCommandOptions struct {
	CacheKey string
}

type DoQueryOptions struct {
	CacheKey string
}

func ValidCommand(object any) error {
	if object == nil {
		return errors.New("Validate() object is null")
	}
	if vo, ok := object.(ValidateObject); ok {
		if err := vo.Validate(); err != nil {
			return err
		}
	}
	return restapi.Validate(object)
}

func DoCommand(ctx context.Context, cmd any, fun func(ctx context.Context) error, opts ...DoCommandOptions) error {
	if err := ValidCommand(cmd); err != nil {
		return err
	}
	if isValidOnly, ok := cmd.(IsValidOnly); ok && isValidOnly.GetIsValidOnly() {
		return nil
	}
	return fun(ctx)
}

func DoQuery[T any](ctx context.Context, qry any, fun func(ctx context.Context) (T, error), opts ...DoQueryOptions) (T, error) {
	return fun(ctx)
}
