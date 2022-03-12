package mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	entityBuilder ddd_repository.EntityBuilder
	collection    *mongo.Collection
}

func NewRepository(entityBuilder ddd_repository.EntityBuilder, collection *mongo.Collection) ddd_repository.Repository {
	return &Repository{
		entityBuilder: entityBuilder,
		collection:    collection,
	}
}

func (r *Repository) NewEntity() interface{} {
	return r.entityBuilder.New()
}

func (r *Repository) NewEntityList() interface{} {
	return r.entityBuilder.NewList()
}

func (r *Repository) BaseCreate(ctx context.Context, entity ddd.Entity, opts ...ddd_repository.SetOption) *ddd_repository.SetResult {
	return r.BaseSet(func() (interface{}, error) {
		_, err := r.collection.InsertOne(ctx, entity)
		return entity, err
	}, opts...)
}

func (r *Repository) BaseUpdate(ctx context.Context, entity ddd.Entity, opts ...ddd_repository.SetOption) *ddd_repository.SetResult {
	return r.BaseSet(func() (interface{}, error) {
		filter := bson.D{{"_id", id}}
		_, err := r.collection.UpdateOne(ctx, filter, entity, options.Update())
		return entity, err
	}, opts...)
}

func (r *Repository) BaseFindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.FindOption) *ddd_repository.FindResult {
	return r.BaseFind(func() (interface{}, error) {
		filter := bson.D{{"_id", id}, {"tenantId", tenantId}}
		data := r.NewEntity()
		err := r.collection.FindOne(ctx, filter).Decode(&data)
		return data, err
	}, opts...)
}

func (r *Repository) BaseFindAll(ctx context.Context, tenantId string, opts ...ddd_repository.FindOption) *ddd_repository.FindResult {
	return r.BaseFind(func() (interface{}, error) {
		filter := bson.D{{"tenantId", tenantId}}
		data := r.NewEntityList()
		cursor, err := r.collection.Find(ctx, filter)
		if err != nil {
			return nil, err
		}
		err = cursor.All(ctx, &data)
		return data, err
	}, opts...)
}

func (r *Repository) BaseFind(doFind func() (interface{}, error), opts ...ddd_repository.FindOption) *ddd_repository.FindResult {
	o := ddd_repository.NewFindOptions().Init(opts...)
	isFind := false
	data, err := doFind()
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			isFind = false
			err = nil
			err = o.OnNotFond()
		}
		err = o.OnError(err)
	}
	if err != nil {
		err = o.OnSuccess(data)
	}
	return ddd_repository.NewFindResult(data, isFind, err)
}

func (r *Repository) BaseSet(doFunc func() (interface{}, error), opts ...ddd_repository.SetOption) *ddd_repository.SetResult {
	o := ddd_repository.NewSetOptions().Init(opts...)
	data, err := doFunc()
	if err != nil {
		err = o.OnError(err)
	}
	if err != nil {
		err = o.OnSuccess(data)
	}

	return ddd_repository.NewSetResult(data, err)
}
