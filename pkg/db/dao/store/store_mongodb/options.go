package store_mongodb

import (
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Options[T any] struct {
	autoCreateCollection *bool // 自动建表
	autoCreateIndex      *bool // 自动建索引
	entityBuilder store2.EntityBuilder[T]
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

func (o *Options[T]) SetEntityBuilder(v store2.EntityBuilder[T]) *Options[T] {
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

func getDeleteOptions(opts ...store2.Options) *options.DeleteOptions {
	deleteOptions := &options.DeleteOptions{}
	return deleteOptions
}

func getFindOptions(opts ...store2.Options) *options.FindOptions {
	opt := store2.NewOptions().Merge(opts...)
	findOneOptions := &options.FindOptions{}
	findOneOptions.MaxTime = opt.GetTimeout()
	return findOneOptions
}

func getAggregateOptions(opts ...store2.Options) *options.AggregateOptions {
	opt := store2.NewOptions().Merge(opts...)
	options := &options.AggregateOptions{}
	options.MaxTime = opt.GetTimeout()
	return options
}

func getFindOneOptions(opts ...store2.Options) *options.FindOneOptions {
	opt := store2.NewOptions().Merge(opts...)
	findOneOptions := &options.FindOneOptions{}
	findOneOptions.MaxTime = opt.GetTimeout()
	return findOneOptions
}

func getUpdateOptions(opts ...store2.Options) *options.UpdateOptions {
	opt := &options.UpdateOptions{}
	for _, o := range opts {
		if o.GetUpsert() != nil {
			opt.Upsert = o.GetUpsert()
		}
	}
	return opt
}

func getInsertOneOptions(opts ...store2.Options) *options.InsertOneOptions {
	return options.InsertOne()
}

func getInsertManyOptions(opts ...store2.Options) *options.InsertManyOptions {
	return options.InsertMany()
}
