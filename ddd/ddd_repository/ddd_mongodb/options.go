package ddd_mongodb

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getDeleteOptions(opts ...ddd_repository.Options) *options.DeleteOptions {
	deleteOptions := &options.DeleteOptions{}
	return deleteOptions
}

func getFindOptions(opts ...ddd_repository.Options) *options.FindOptions {
	opt := ddd_repository.NewOptions().Merge(opts...)
	findOneOptions := &options.FindOptions{}
	findOneOptions.MaxTime = opt.GetTimeout()
	return findOneOptions
}

func getFindOneOptions(opts ...ddd_repository.Options) *options.FindOneOptions {
	opt := ddd_repository.NewOptions().Merge(opts...)
	findOneOptions := &options.FindOneOptions{}
	findOneOptions.MaxTime = opt.GetTimeout()
	return findOneOptions
}

func getUpdateOptions(opts ...ddd_repository.Options) *options.UpdateOptions {
	opt := &options.UpdateOptions{}
	for _, o := range opts {
		if o.GetUpsert() != nil {
			opt.Upsert = o.GetUpsert()
		}
	}
	return opt
}

func getInsertOneOptions(opts ...ddd_repository.Options) *options.InsertOneOptions {
	return options.InsertOne()
}

func getInsertManyOptions(opts ...ddd_repository.Options) *options.InsertManyOptions {
	return options.InsertMany()
}
