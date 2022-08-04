package ddd_repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

func NewIds[T ddd.Entity](ctx context.Context, list []T) ([]string, error) {
	var ids []string
	for i, e := range list {
		if stringutils.IsEmptyStr(e.GetId()) {
			return nil, errors.ErrorOf("params list index %v entity id is empty", i)
		}
		ids = append(ids, e.GetId())
	}
	return ids, nil
}
