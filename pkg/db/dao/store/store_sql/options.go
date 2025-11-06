package store_sql

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type Options[T any] struct {
	autoCreateCollection *bool // 自动建表
	autoCreateIndex      *bool // 自动建索引
	entityBuilder store.EntityBuilder[T]
}

func NewOptions[T any](opts ...*Options[T]) *Options[T] {
	o := &Options[T]{}
	for _, item := range opts {
		if item == nil {
			continue
		}
		if item.autoCreateCollection != nil {
			o.autoCreateCollection = item.autoCreateCollection
		}
		if item.autoCreateIndex != nil {
			o.autoCreateIndex = item.autoCreateIndex
		}
		if item.entityBuilder != nil {
			o.entityBuilder = item.entityBuilder
		}
	}
	return o
}

func (o *Options[T]) SetAutoCreateCollection(v bool) *Options[T] {
	o.autoCreateCollection = &v
	return o
}

func (o *Options[T]) SetAutoCreateIndex(v bool) *Options[T] {
	o.autoCreateIndex = &v
	return o
}

func (o *Options[T]) SetEntityBuilder(v store.EntityBuilder[T]) *Options[T] {
	o.entityBuilder = v
	return o
}

func (o *Options[T]) GetAutoCreateCollection() bool {
	if o == nil || o.autoCreateCollection == nil {
		return false
	}
	v := o.autoCreateCollection
	return *v
}

func (o *Options[T]) GetAutoCreateIndex() bool {
	if o == nil || o.autoCreateIndex == nil {
		return false
	}
	v := o.autoCreateIndex
	return *v
}
