package ddd_mongodb

import (
	"github.com/dapr/dapr-go-ddd-sdk/ddd/ddd_repository"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getDeleteOptions(opts ...*ddd_repository.SetOptions) *options.DeleteOptions {
	deleteOptions := &options.DeleteOptions{}
	return deleteOptions
}

func getFindOptions(opts ...*ddd_repository.FindOptions) *options.FindOptions {
	opt := ddd_repository.MergeFindOptions(opts...)
	findOneOptions := &options.FindOptions{}
	findOneOptions.MaxTime = opt.MaxTime
	return findOneOptions
}

func getFindOneOptions(opts ...*ddd_repository.FindOptions) *options.FindOneOptions {
	opt := ddd_repository.MergeFindOptions(opts...)
	findOneOptions := &options.FindOneOptions{}
	findOneOptions.MaxTime = opt.MaxTime
	return findOneOptions
}

func getUpdateOptions(opts ...*ddd_repository.SetOptions) *options.UpdateOptions {
	updateOptions := &options.UpdateOptions{}
	return updateOptions
}

func getInsertOneOptions(opts ...*ddd_repository.SetOptions) *options.InsertOneOptions {
	return options.InsertOne()
}
