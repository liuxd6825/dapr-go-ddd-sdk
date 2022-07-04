package validateutils

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Struct(ctx context.Context, data interface{}) error {
	err := validate.StructCtx(getContext(ctx), data)
	return getError(err)
}

func Variable(ctx context.Context, field interface{}, tag string) error {
	err := validate.VarCtx(getContext(ctx), field, tag)
	return getError(err)
}

func Map(ctx context.Context, data map[string]interface{}, rule map[string]interface{}) map[string]interface{} {
	return validate.ValidateMapCtx(getContext(ctx), data, rule)
}

func getContext(ctx context.Context) context.Context {
	c := ctx
	if c == nil {
		c = context.Background()
	}
	return c
}

func getError(err error) error {
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}
		verifyError := ddd_errors.NewVerifyError()
		for _, err := range err.(validator.ValidationErrors) {
			verifyError.AppendField(err.Field(), err.Error())
		}
		return verifyError
	}
	return err
}
