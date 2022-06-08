package ddd_mongodb

import (
	"context"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
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

type Repository[T ddd.Entity] struct {
	entityBuilder *ddd_repository.EntityBuilder[T]
	collection    *mongo.Collection
	mongodb       *MongoDB
	emptyEntity   T
	newFun        func() T
}

func NewRepository[T ddd.Entity](newFun func() T, mongodb *MongoDB, collection *mongo.Collection) *Repository[T] {
	return &Repository[T]{
		newFun:     newFun,
		collection: collection,
		mongodb:    mongodb,
	}
}

func (r *Repository[T]) NewEntity() T {
	return r.newFun()
}

func (r *Repository[T]) NewEntityList() *[]T {
	return &[]T{}
}

func (r *Repository[T]) Insert(ctx context.Context, entity T, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult[T] {
	return r.DoSet(func() (T, error) {
		_, err := r.collection.InsertOne(ctx, entity, getInsertOneOptions(opts...))
		return entity, err
	})
}

func (r *Repository[T]) Update(ctx context.Context, entity T, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult[T] {
	return r.DoSet(func() (T, error) {
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

func (r *Repository[T]) Delete(ctx context.Context, entity ddd.Entity, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult[T] {
	return r.DeleteById(ctx, entity.GetTenantId(), entity.GetId(), opts...)
}

func (r *Repository[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult[T] {
	data := map[string]interface{}{
		IdField: id,
	}
	return r.DeleteByMap(ctx, tenantId, data)
}

func (r *Repository[T]) DeleteAll(ctx context.Context, tenantId string, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult[T] {
	data := map[string]interface{}{}
	return r.DeleteByMap(ctx, tenantId, data)
}

func (r *Repository[T]) DeleteByMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...*ddd_repository.SetOptions) *ddd_repository.SetResult[T] {
	if err := assert.NotNil(data, assert.NewOptions("data is nil")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return r.DoSet(func() (T, error) {
		filter := bson.D{
			{TenantIdField, tenantId},
		}
		for k, v := range data {
			e := bson.E{
				Key:   k,
				Value: v,
			}
			filter = append(filter, e)
		}
		deleteOptions := getDeleteOptions(opts...)
		_, err := r.collection.DeleteOne(ctx, filter, deleteOptions)
		var result T
		return result, err
	})
}

func (r *Repository[T]) NewFilter(tenantId string, filterMap map[string]interface{}) bson.D {
	filter := bson.D{
		{TenantIdField, tenantId},
	}
	if filterMap != nil {
		for k, v := range filterMap {
			e := bson.E{
				Key:   k,
				Value: v,
			}
			filter = append(filter, e)
		}
	}
	return filter
}

func (r *Repository[T]) FindById(ctx context.Context, tenantId string, id string, opts ...*ddd_repository.FindOptions) *ddd_repository.FindOneResult[T] {
	idMap := map[string]interface{}{
		IdField: id,
	}
	return r.FindOneByMap(ctx, tenantId, idMap, opts...)
}

func (r *Repository[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...*ddd_repository.FindOptions) *ddd_repository.FindOneResult[T] {
	return r.DoFindOne(func() (T, bool, error) {
		filter := r.NewFilter(tenantId, filterMap)
		findOneOptions := getFindOneOptions(opts...)
		data := r.NewEntity()
		result := r.collection.FindOne(ctx, filter, findOneOptions)
		if result.Err() != nil {
			return r.emptyEntity, false, result.Err()
		}
		if err := result.Decode(data); err != nil {
			return r.emptyEntity, false, err
		}
		return data, true, nil
	})
}

func (r *Repository[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...*ddd_repository.FindOptions) *ddd_repository.FindListResult[T] {
	return r.DoFindList(func() (*[]T, bool, error) {
		filter := r.NewFilter(tenantId, filterMap)
		data := r.NewEntityList()
		findOptions := getFindOptions(opts...)
		cursor, err := r.collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, data)
		return data, true, err
	})
}

func (r *Repository[T]) FindAll(ctx context.Context, tenantId string, opts ...*ddd_repository.FindOptions) *ddd_repository.FindListResult[T] {
	return r.FindListByMap(ctx, tenantId, nil, opts...)
}

func (r *Repository[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...*ddd_repository.FindOptions) *ddd_repository.FindPagingResult[T] {
	return r.DoFilter(query.GetTenantId(), query.GetFilter(), func(filter map[string]interface{}) (*ddd_repository.FindPagingResult[T], bool, error) {
		data := r.NewEntityList()

		findOptions := getFindOptions(opts...)
		if query.GetPageSize() > 0 {
			findOptions.SetLimit(query.GetPageSize())
			findOptions.SetSkip(query.GetPageSize() * query.GetPageNum())
		}
		if len(query.GetSort()) > 0 {
			sort, err := r.getSort(query.GetSort())
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
		totalRows, err := r.collection.CountDocuments(ctx, filter)
		findData := ddd_repository.NewFindPagingResult[T](data, totalRows, query, err)
		return findData, true, err
	})
}

func (r *Repository[T]) DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*ddd_repository.FindPagingResult[T], bool, error)) *ddd_repository.FindPagingResult[T] {
	p := NewMongoProcess()
	err := rsql.ParseProcess(filter, p)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	filterData := p.GetFilter(tenantId)
	data, _, err := fun(filterData)
	if err != nil {
		if ddd_errors.IsErrorMongoNoDocuments(err) {
			err = nil
		}
	}
	return data
}

func (r *Repository[T]) DoFindList(fun func() (*[]T, bool, error)) *ddd_repository.FindListResult[T] {
	data, isFound, err := fun()
	if err != nil {
		if ddd_errors.IsErrorMongoNoDocuments(err) {
			isFound = false
			err = nil
		}
	}
	return ddd_repository.NewFindListResult[T](data, isFound, err)
}

func (r *Repository[T]) DoFindOne(fun func() (T, bool, error)) *ddd_repository.FindOneResult[T] {
	data, isFound, err := fun()
	if err != nil {
		if ddd_errors.IsErrorMongoNoDocuments(err) {
			isFound = false
			err = nil
		}
	}
	return ddd_repository.NewFindOneResult[T](data, isFound, err)
}

func (r *Repository[T]) DoSet(fun func() (T, error)) *ddd_repository.SetResult[T] {
	data, err := fun()
	return ddd_repository.NewSetResult[T](data, err)
}

func (r *Repository[T]) getSort(sort string) (map[string]interface{}, error) {
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
