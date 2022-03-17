package ddd_mongodb

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

func NewBaseRepository(entityBuilder ddd_repository.EntityBuilder, collection *mongo.Collection) ddd_repository.BaseRepository {
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

func (r *Repository) BaseCreate(ctx context.Context, entity ddd.Entity) *ddd_repository.SetResult {
	return r.BaseSet(func() (interface{}, error) {
		_, err := r.collection.InsertOne(ctx, entity)
		return entity, err
	})
}

func (r *Repository) BaseUpdate(ctx context.Context, entity ddd.Entity) *ddd_repository.SetResult {
	return r.BaseSet(func() (interface{}, error) {
		filter := bson.D{{"_id", id}}
		_, err := r.collection.UpdateOne(ctx, filter, entity, options.Update())
		return entity, err
	})
}

func (r *Repository) BaseDeleteById(ctx context.Context, tenantId string, id string) *ddd_repository.SetResult {
	return r.BaseSet(func() (interface{}, error) {
		filter := bson.D{{"_id", id}, {"tenantId", tenantId}}
		_, err := r.collection.DeleteOne(ctx, filter)
		return nil, err
	})
}

func (r *Repository) BaseFindById(ctx context.Context, tenantId string, id string) *ddd_repository.FindResult {
	return r.BaseFind(func() (interface{}, error) {
		filter := bson.D{{"_id", id}, {"tenantId", tenantId}}
		data := r.NewEntity()
		err := r.collection.FindOne(ctx, filter).Decode(&data)
		return data, err
	})
}

func (r *Repository) BaseFindAll(ctx context.Context, tenantId string) *ddd_repository.FindResult {
	return r.BaseFind(func() (interface{}, error) {
		filter := bson.D{{"tenantId", tenantId}}
		data := r.NewEntityList()
		cursor, err := r.collection.Find(ctx, filter)
		if err != nil {
			return nil, err
		}
		err = cursor.All(ctx, &data)
		return data, err
	})
}

func (r *Repository) BaseFind(doFind func() (interface{}, error)) *ddd_repository.FindResult {
	isFind := false
	data, err := doFind()
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			isFind = false
			err = nil
		}
	}
	return ddd_repository.NewFindResult(data, isFind, err)
}

func (r *Repository) BaseSet(doFunc func() (interface{}, error)) *ddd_repository.SetResult {
	data, err := doFunc()
	return ddd_repository.NewSetResult(data, err)
}
