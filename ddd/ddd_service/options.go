package ddd_service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
)

type Options = store.Options

type ServiceOptions = store.RepositoryOptions

func NewOptions() Options {
	return &ServiceOptions{}
}

func MergeOptions(opts ...Options) Options {
	o := NewOptions()
	return o
}

func NewRepositoryOptions(opts ...Options) []store.Options {
	var list []store.Options
	for _, o := range opts {
		list = append(list, o)
	}
	return list
}
