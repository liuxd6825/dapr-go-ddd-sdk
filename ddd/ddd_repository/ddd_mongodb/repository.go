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
	return r.entityBuilder.NewOne()
}

func (r *Repository) NewEntityList() interface{} {
	return r.entityBuilder.NewList()
}

func (r *Repository) Insert(ctx context.Context, entity ddd.Entity, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult {
	return r.DoSet(func() (interface{}, error) {
		_, err := r.collection.InsertOne(ctx, entity, getInsertOneOptions(opts...))
		return entity, err
	})
}

func (r *Repository) Update(ctx context.Context, entity ddd.Entity, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult {
	return r.DoSet(func() (interface{}, error) {
		objId, err := GetObjectID(entity.GetId())
		if err != nil {
			return entity, err
		}
		updateOptions := getUpdateOptions(opts...)
		filter := bson.D{{IdField, objId}}
		_, err = r.collection.UpdateOne(ctx, filter, entity, updateOptions)
		return entity, err
	})
}

func (r *Repository) Delete(ctx context.Context, entity ddd.Entity, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult {
	return r.DeleteById(ctx, entity.GetTenantId(), entity.GetId(), opts...)
}

func (r *Repository) DeleteById(ctx context.Context, tenantId string, id string, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult {
	return r.DoSet(func() (interface{}, error) {
		filter := bson.D{
			{IdField, id},
			{TenantIdField, tenantId},
		}
		deleteOptions := getDeleteOptions(opts...)
		_, err := r.collection.DeleteOne(ctx, filter, deleteOptions)
		return nil, err
	})
}

func (r *Repository) FindById(ctx context.Context, tenantId string, id string, opts ...*ddd_repository.FindOptions) *ddd_repository.FindResult {
	return r.DoFind(func() (interface{}, bool, error) {
		filter := bson.M{
			IdField:       id,
			TenantIdField: tenantId,
		}
		findOneOptions := getFindOneOptions(opts...)
		data := r.NewEntity()
		result := r.collection.FindOne(ctx, filter, findOneOptions)
		if result.Err() != nil {
			return nil, false, result.Err()
		}
		if err := result.Decode(data); err != nil {
			return nil, false, err
		}
		return data, true, nil
	})
}

func (r *Repository) FindAll(ctx context.Context, tenantId string, opts ...*ddd_repository.FindOptions) *ddd_repository.FindResult {
	return r.DoFind(func() (interface{}, bool, error) {
		filter := bson.D{{TenantIdField, tenantId}}
		data := r.NewEntityList()
		findOptions := getFindOptions(opts...)
		cursor, err := r.collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &data)
		return data, true, err
	})
}

func (r *Repository) FindPaging(ctx context.Context, query *ddd_repository.PagingQuery, opts ...*ddd_repository.FindOptions) *ddd_repository.FindPagingResult {
	return r.DoFilter(query.TenantId, query.Filter, func(filter map[string]interface{}) (*ddd_repository.PagingData, bool, error) {
		data := r.NewEntityList()

		findOptions := getFindOptions(opts...)
		if query.Size > 0 {
			findOptions.SetLimit(query.Size)
			findOptions.SetSkip(query.Size * query.Page)
		}
		if len(query.Sort) > 0 {
			sort, err := r.getSort(query.Sort)
			if err != nil {
				return nil, false, err
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

func (r *Repository) DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*ddd_repository.PagingData, bool, error)) *ddd_repository.FindPagingResult {
	p := NewMongoProcess()
	err := rsql.ParseProcess(filter, p)
	if err != nil {
		return ddd_repository.NewFindPagingDataResult(nil, false, err)
	}
	filterData := p.GetFilter(tenantId)
	data, isFound, err := fun(filterData)
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
