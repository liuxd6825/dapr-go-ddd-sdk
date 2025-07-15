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

func (b *Service) ValidCommand(object any) error {
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

func (b *Service) DoCommand(ctx context.Context, cmd any, fun func(ctx context.Context) error) error {
	if err := b.ValidCommand(cmd); err != nil {
		return err
	}
	return fun(ctx)
}
