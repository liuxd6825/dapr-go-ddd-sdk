package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const (
	IdField       = "_id"
	TenantIdField = "tenant_id"
)

type Repository[T ddd.Entity] struct {
	entityBuilder *ddd_repository.EntityBuilder[T]
	collection    *mongo.Collection
	mongodb       *MongoDB
	emptyEntity   T
	newFun        func() T
}

func NewRepository[T ddd.Entity](newFun func() T, mongodb *MongoDB, collection *mongo.Collection) *Repository[T] {
	r := &Repository[T]{}
	r.Init(newFun, mongodb, collection)
	return r
}

func (r *Repository[T]) Init(newFun func() T, mongodb *MongoDB, collection *mongo.Collection) {
	r.newFun = newFun
	r.collection = collection
	r.mongodb = mongodb
}

func (r *Repository[T]) NewEntity() T {
	return r.newFun()
}

func (r *Repository[T]) NewEntityList() *[]T {
	return &[]T{}
}

func (r *Repository[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	if err := assert.NotEmpty(entity.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return r.DoSet(func() (T, error) {
		_, err := r.collection.InsertOne(ctx, entity, getInsertOneOptions(opts...))
		return entity, err
	})
}

func (r *Repository[T]) InsertMany(ctx context.Context, entitits *[]T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	if entitits == nil || len(*entitits) == 0 {
		return ddd_repository.NewSetManyResultError[T](errors.New("entitits is nil"))
	}

	for _, e := range *entitits {
		if err := assert.NotEmpty(e.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}

	var docs []interface{}
	for _, e := range *entitits {
		docs = append(docs, e)
	}

	return r.DoSetMany(func() ([]T, error) {
		_, err := r.collection.InsertMany(ctx, docs, getInsertManyOptions(opts...))
		return *entitits, err
	})
}

func (r *Repository[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	if err := assert.NotEmpty(entity.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return r.DoSet(func() (T, error) {
		objId, err := GetObjectID(entity.GetId())
		if err != nil {
			return entity, err
		}
		updateOptions := getUpdateOptions(opts...)
		filter := bson.D{{IdField, objId}}
		setData := bson.M{"$set": entity}
		_, err = r.collection.UpdateOne(ctx, filter, setData, updateOptions)
		return entity, err
	})
}

func (r *Repository[T]) UpdateManyByFilter(ctx context.Context, tenantId, filter string, data interface{}, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	filterMap, err := r.getFilterMap(tenantId, filter)
	if err != nil {
		return ddd_repository.NewSetManyResultError[T](err)
	}
	return r.DoSetMany(func() ([]T, error) {
		setData := bson.M{"$set": data}
		updateOptions := getUpdateOptions(opts...)
		_, err = r.collection.UpdateMany(ctx, filterMap, setData, updateOptions)
		return nil, err
	})
}

func (r *Repository[T]) UpdateManyById(ctx context.Context, entities *[]T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	if entities == nil || len(*entities) == 0 {
		return ddd_repository.NewSetManyResultError[T](errors.New("entities is nil"))
	}

	for _, e := range *entities {
		if err := assert.NotEmpty(e.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}

	var docs []interface{}
	for _, e := range *entities {
		docs = append(docs, e)
	}

	return r.DoSetMany(func() ([]T, error) {
		for _, entity := range *entities {
			objId, err := GetObjectID(entity.GetId())
			if err != nil {
				return nil, err
			}
			updateOptions := getUpdateOptions(opts...)
			filter := bson.D{{IdField, objId}}
			setData := bson.M{"$set": entity}
			_, err = r.collection.UpdateOne(ctx, filter, setData, updateOptions)
			if err != nil {
				return nil, err
			}
		}
		return *entities, nil
	})
}

func (r *Repository[T]) UpdateManyMaskById(ctx context.Context, entities *[]T, mask []string, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	if entities == nil || len(*entities) == 0 {
		return ddd_repository.NewSetManyResultError[T](errors.New("entities is nil"))
	}

	for _, e := range *entities {
		if err := assert.NotEmpty(e.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}

	var docs []interface{}
	for _, e := range *entities {
		if len(mask) == 0 {
			docs = append(docs, e)
		} else {
			m := make(map[string]interface{})
			if err := types.MaskMapper(e, &m, mask); err != nil {
				return ddd_repository.NewSetManyResultError[T](err)
			}
			m[IdField] = e.GetId()
			doc := asDocument(m)
			docs = append(docs, doc)
		}
	}

	return r.DoSetMany(func() ([]T, error) {
		for _, doc := range docs {
			id, err := getDocumentId(doc)
			if err != nil {
				return nil, err
			}
			updateOptions := getUpdateOptions(opts...)
			filter := bson.D{{IdField, id}}

			setData := bson.M{"$set": doc}
			_, err = r.collection.UpdateOne(ctx, filter, setData, updateOptions)
			if err != nil {
				return nil, err
			}
		}
		return *entities, nil
	})
}

func (r *Repository[T]) UpdateMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, data map[string]interface{}, opts ...ddd_repository.Options) error {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return err
	}

	if err := assert.NotNil(filterMap, assert.NewOptions("filterMap is nil")); err != nil {
		return err
	}

	filterMap[TenantIdField] = tenantId
	updateOptions := getUpdateOptions(opts...)
	filter := r.NewFilter(tenantId, filterMap)
	setData := bson.M{"$set": data}
	_, err := r.collection.UpdateOne(ctx, filter, setData, updateOptions)
	return err
}

func (r *Repository[T]) Delete(ctx context.Context, entity ddd.Entity, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	return r.DeleteById(ctx, entity.GetTenantId(), entity.GetId(), opts...)
}

func (r *Repository[T]) DeleteByFilter(ctx context.Context, tenantId, filter string) error {
	if filterMap, err := r.getFilterMap(tenantId, filter); err != nil {
		return err
	} else if err := r.DeleteByMap(ctx, tenantId, filterMap).GetError(); err != nil {
		return err
	}
	return nil
}

func (r *Repository[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	data := map[string]interface{}{
		IdField: id,
	}
	return r.DeleteByMap(ctx, tenantId, data)
}

func (r *Repository[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	data := map[string]interface{}{}
	return r.DeleteByMap(ctx, tenantId, data)
}

func (r *Repository[T]) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	if err := assert.NotNil(filterMap, assert.NewOptions("data is nil")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return r.DoSet(func() (T, error) {
		filter := r.NewFilter(tenantId, filterMap)
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
		for fieldName, fieldValue := range filterMap {
			if fieldName != IdField {
				fieldName = AsFieldName(fieldName)
			}
			e := bson.E{
				Key:   fieldName,
				Value: fieldValue,
			}
			filter = append(filter, e)
		}
	}
	return filter
}
func (r *Repository[T]) FindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	idMap := map[string]interface{}{
		IdField: id,
	}
	return r.FindOneByMap(ctx, tenantId, idMap, opts...)
}

func (r *Repository[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
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

func (r *Repository[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.DoFindList(func() ([]T, bool, error) {
		filter := r.NewFilter(tenantId, filterMap)
		data := r.NewEntityList()
		findOptions := getFindOptions(opts...)
		cursor, err := r.collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, data)
		return *data, true, err
	})
}

func (r *Repository[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return r.FindListByMap(ctx, tenantId, nil, opts...)
}

func (r *Repository[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {

	return r.DoFilter(query.GetTenantId(), query.GetFilter(), func(filter map[string]interface{}) (*ddd_repository.FindPagingResult[T], bool, error) {
		if err := assert.NotEmpty(query.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
			return nil, false, err
		}

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
		findData := ddd_repository.NewFindPagingResult[T](*data, totalRows, query, err)
		return findData, true, err
	})
}

func (r *Repository[T]) DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*ddd_repository.FindPagingResult[T], bool, error)) *ddd_repository.FindPagingResult[T] {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	filterData, err := r.getFilterMap(tenantId, filter)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	data, _, err := fun(filterData)
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			err = nil
		}
	}
	return data
}

func (r *Repository[T]) getFilterMap(tenantId, rsqlstr string) (map[string]interface{}, error) {
	p := NewMongoProcess()
	if err := rsql.ParseProcess(rsqlstr, p); err != nil {
		return nil, err
	}
	filterMap, err := p.GetFilter(tenantId)
	return filterMap, err
}

func (r *Repository[T]) DoFindList(fun func() ([]T, bool, error)) *ddd_repository.FindListResult[T] {
	data, isFound, err := fun()
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			isFound = false
			err = nil
		}
	}
	return ddd_repository.NewFindListResult[T](data, isFound, err)
}

func (r *Repository[T]) DoFindOne(fun func() (T, bool, error)) *ddd_repository.FindOneResult[T] {
	data, isFound, err := fun()
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
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

func (r *Repository[T]) DoSetMany(fun func() ([]T, error)) *ddd_repository.SetManyResult[T] {
	data, err := fun()
	return ddd_repository.NewSetManyResult[T](data, err)
}

/*
func (r *Dao[T]) DoSetMap(fun func() (map[string]interface{}, error)) *ddd_repository.SetResult[map[string]interface{}] {
	data, err := fun()
	return ddd_repository.NewSetResult[T](data, err)
}
*/
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
