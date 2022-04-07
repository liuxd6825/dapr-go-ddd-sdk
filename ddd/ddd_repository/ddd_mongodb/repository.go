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

const (
	IdField       = "_id"
	TenantIdField = "tenantId"
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

func (r *Repository) Insert(ctx context.Context, entity ddd.Entity) *ddd_repository.SetResult {
	return r.DoSet(func() (interface{}, error) {
		_, err := r.collection.InsertOne(ctx, entity)
		return entity, err
	})
}

func (r *Repository) Update(ctx context.Context, entity ddd.Entity) *ddd_repository.SetResult {
	return r.DoSet(func() (interface{}, error) {
		objId, err := GetObjectID(entity.GetId())
		if err != nil {
			return entity, err
		}
		filter := bson.D{{IdField, objId}}
		_, err = r.collection.UpdateOne(ctx, filter, entity, options.Update())
		return entity, err
	})
}

func (r *Repository) DeleteById(ctx context.Context, tenantId string, id string) *ddd_repository.SetResult {
	return r.DoSet(func() (interface{}, error) {
		filter := bson.D{
			{IdField, id},
			{TenantIdField, tenantId},
		}
		_, err := r.collection.DeleteOne(ctx, filter)
		return nil, err
	})
}

func (r *Repository) FindById(ctx context.Context, tenantId string, id string) *ddd_repository.FindResult {
	return r.DoFind(func() (interface{}, bool, error) {
		filter := bson.M{
			IdField:       id,
			TenantIdField: tenantId,
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

func (r *Repository) FindAll(ctx context.Context, tenantId string) *ddd_repository.FindResult {
	return r.DoFind(func() (interface{}, bool, error) {
		filter := bson.D{{TenantIdField, tenantId}}
		data := r.NewEntityList()
		cursor, err := r.collection.Find(ctx, filter)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &data)
		return data, true, err
	})
}

func (r *Repository) FindPagingData(ctx context.Context, query *ddd_repository.PagingQuery) *ddd_repository.FindPagingDataResult {
	return r.DoFindPagingData(func() (*ddd_repository.PagingData, bool, error) {
		p := NewMongoProcess()
		err := rsql.ParseProcess(query.Filter, p)
		if err != nil {
			return nil, false, err
		}
		filter := p.GetFilter(query.TenantId)
		data := r.NewEntityList()

		var findOptions *options.FindOptions
		if query.Size > 0 {
			findOptions = &options.FindOptions{}
			findOptions.SetLimit(query.Size)
			findOptions.SetSkip(query.Size * query.Page)
		}
		if len(query.Sort) > 0 {
			if findOptions == nil {
				findOptions = &options.FindOptions{}
			}
			sort, oerr := r.getSort(query.Sort)
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

		count, err := r.collection.CountDocuments(ctx, filter)
		findData := &ddd_repository.PagingData{
			Data:      data,
			Count:     count,
			TotalPage: r.getTotalPage(count, query.Size),
			Filter:    query.Filter,
			Sort:      query.Sort,
			Size:      query.Size,
		}
		return findData, true, err
	})
}

func (r *Repository) DoFindPagingData(fun func() (*ddd_repository.PagingData, bool, error)) *ddd_repository.FindPagingDataResult {
	data, isFound, err := fun()
	if err != nil {
		if ddd_errors.IsErrorMongoNoDocuments(err) {
			isFound = false
			err = nil
		}
	}
	return ddd_repository.NewFindPagingDataResult(data, isFound, err)
}

func (r *Repository) DoFind(fun func() (interface{}, bool, error)) *ddd_repository.FindResult {
	data, isFound, err := fun()
	if err != nil {
		if ddd_errors.IsErrorMongoNoDocuments(err) {
			isFound = false
			err = nil
		}
	}
	return ddd_repository.NewFindResult(data, isFound, err)
}

func (r *Repository) DoSet(fun func() (interface{}, error)) *ddd_repository.SetResult {
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
		if name == "id" {
			name = IdField
		}
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

func (r *Repository) getTotalPage(count int64, size int64) int64 {
	totalPage := count / size
	if count%size > 1 {
		totalPage++
	}
	return totalPage
}
