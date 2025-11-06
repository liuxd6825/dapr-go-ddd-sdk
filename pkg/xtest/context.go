package xtest

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
)

func NewContext() context.Context {
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		panic(err)
	}
	return ctx
}
