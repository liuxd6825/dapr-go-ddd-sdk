package ddd_repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"strings"
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

func RSqlKeyValueToMap(s string) map[string]string {
	res := make(map[string]string, 0)
	if len(s) == 0 {
		return res
	}
	list := strings.Split(s, ",")
	for _, item := range list {
		kv := strings.Split(item, ":")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			res[key] = val
		}
	}
	return res
}

type KeyValue struct {
	Key   string
	Value string
}

func RSqlKeyValueToList(s string) []KeyValue {
	res := make([]KeyValue, 0)
	if len(s) == 0 {
		return res
	}
	list := strings.Split(s, ",")
	for _, item := range list {
		kv := strings.Split(item, ":")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			res = append(res, KeyValue{Key: key, Value: val})
		}
	}
	return res
}
