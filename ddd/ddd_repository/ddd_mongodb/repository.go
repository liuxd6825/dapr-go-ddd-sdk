package ddd_mongodb

import (
	"context"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type Repository struct {
	entityBuilder ddd_repository.EntityBuilder
	collection    *mongo.Collection
	mongodb       *MongoDB
}

func NewRepository(entityBuilder ddd_repository.EntityBuilder, mongodb *MongoDB, collection *mongo.Collection) ddd_repository.Repository {
	return &Repository{
		entityBuilder: entityBuilder,
		collection:    collection,
		mongodb:       mongodb,
	}
}

func (r *Repository) NewEntity() interface{} {
	return r.entityBuilder.New()
}

func (r *Repository) NewEntityList() interface{} {
	return r.entityBuilder.NewList()
}

func (r *Repository) DoFindList(ctx context.Context, search *ddd_repository.ListQuery) *ddd_repository.FindResult {
	return r.Find(func() (interface{}, bool, error) {
		p := NewMongoProcess()
		err := rsql.ParseProcess(search.Filter, p)
		if err != nil {
			return nil, false, err
		}
		filter := p.GetFilter(search.TenantId)
		data := r.NewEntityList()

		var findOptions *options.FindOptions
		if search.Size > 0 {
			findOptions = &options.FindOptions{}
			findOptions.SetLimit(search.Size)
			findOptions.SetSkip(search.Size * search.Page)
		}
		if len(search.Sort) > 0 {
			if findOptions == nil {
				findOptions = &options.FindOptions{}
			}
			sort, oerr := r.getSort(search.Sort)
			if oerr != nil {
				return nil, false, oerr
			}
			findOptions.SetSort(sort)
		}

		cursor, err := r.collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, data)
		return data, true, err
	})
}

func (r *Repository) DoCreate(ctx context.Context, entity ddd.Entity) *ddd_repository.SetResult {
	return r.Set(func() (interface{}, error) {
		_, err := r.collection.InsertOne(ctx, entity)
		return entity, err
	})
}

func (r *Repository) DoUpdate(ctx context.Context, entity ddd.Entity) *ddd_repository.SetResult {
	return r.Set(func() (interface{}, error) {
		filter := bson.D{{"_id", id}}
		_, err := r.collection.UpdateOne(ctx, filter, entity, options.Update())
		return entity, err
	})
}

func (r *Repository) DoDeleteById(ctx context.Context, tenantId string, id string) *ddd_repository.SetResult {
	return r.Set(func() (interface{}, error) {
		filter := bson.D{{"_id", id}, {"tenantId", tenantId}}
		_, err := r.collection.DeleteOne(ctx, filter)
		return nil, err
	})
}

func (r *Repository) DoFindById(ctx context.Context, tenantId string, id string) *ddd_repository.FindResult {
	return r.Find(func() (interface{}, bool, error) {
		filter := bson.M{
			"tenantId": tenantId,
			"id":       id,
		}
		data := r.NewEntity()
		result := r.collection.FindOne(ctx, filter)
		if result.Err() != nil {
			return nil, false, result.Err()
		}
		if err := result.Decode(data); err != nil {
			return nil, false, err
		}
		return data, true, nil
	})
}

func (r *Repository) DoFindAll(ctx context.Context, tenantId string) *ddd_repository.FindResult {
	return r.Find(func() (interface{}, bool, error) {
		filter := bson.D{{"tenantId", tenantId}}
		data := r.NewEntityList()
		cursor, err := r.collection.Find(ctx, filter)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &data)
		return data, true, err
	})
}

func (r *Repository) Find(fun func() (interface{}, bool, error)) *ddd_repository.FindResult {
	data, isFound, err := fun()
	if err != nil {
		if ddd_errors.IsMongoNoDocumentsInResult(err) {
			isFound = false
			err = nil
		}
	}
	return ddd_repository.NewFindResult(data, isFound, err)
}

func (r *Repository) Set(fun func() (interface{}, error)) *ddd_repository.SetResult {
	data, err := fun()
	return ddd_repository.NewSetResult(data, err)
}

func (r *Repository) getSort(sort string) (map[string]interface{}, error) {
	if len(sort) == 0 {
		return nil, nil
	}
	//name:desc,id:asc
	res := map[string]interface{}{}
	list := strings.Split(sort, ",")
	for _, s := range list {
		sortItem := strings.Split(s, ":")
		name := sortItem[0]
		name = strings.Trim(name, " ")

		order := "asc"
		if len(sortItem) > 1 {
			order = sortItem[1]
			order = strings.ToLower(order)
			order = strings.Trim(order, " ")
		}

		// 其中 1 为升序排列，而-1是用于降序排列.
		orderVal := 1
		var oerr error
		switch order {
		case "asc":
			orderVal = 1
		case "desc":
			orderVal = -1
		default:
			oerr = errors.New("order " + order + " is error")
		}
		if oerr != nil {
			return nil, oerr
		}
		res[name] = orderVal
	}
	return res, nil
}
