package ddd_service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

type Options = ddd_repository.Options

type ServiceOptions = ddd_repository.RepositoryOptions

func NewOptions() Options {
	return &ServiceOptions{}
}

func MergeOptions(opts ...Options) Options {
	o := NewOptions()
	return o
}

func NewRepositoryOptions(opts ...Options) []ddd_repository.Options {
	var list []ddd_repository.Options
	for _, o := range opts {
		list = append(list, o)
	}
	return list
}
