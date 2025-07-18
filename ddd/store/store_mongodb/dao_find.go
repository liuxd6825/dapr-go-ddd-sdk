package store_mongodb

import (
	"context"
	"encoding/json"
	"github.com/dapr/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql/rsql_mongo"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/mongoutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

func (r *Dao[T]) FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...store.Options) (T, error) {
	var null T
	if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
		return null, err
	}

	if err := assert2.NotEmpty(id, assert2.NewOptions("id is empty")); err != nil {
		return null, err
	}
	filter := bson.M{"tenant_id": tenantId, "id": id}
	udpate := r.getDbMap(data)
	sCtx := r.getSessionCtx(ctx)
	_, err := r.getCollection(ctx).UpdateOne(sCtx, filter, udpate)
	if err != nil {
		return null, err
	}
	find := r.FindById(ctx, tenantId, id, opts...)
	return find.GetData(), find.GetError()
}

func (r *Dao[T]) FindById(ctx context.Context, tenantId string, id string, opts ...store.Options) *store.FindOneResult[T] {
	idMap := map[string]interface{}{
		ConstIdField: id,
	}
	return r.FindOneByMap(ctx, tenantId, idMap, opts...)
}

func (r *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...store.Options) *store.FindListResult[T] {
	return r.DoFindList(func() ([]T, bool, error) {
		var list []T
		// 构建查询条件
		filter := bson.D{
			{Key: ConstTenantIdField, Value: tenantId},
			{Key: ConstIdField, Value: bson.D{{Key: "$in", Value: ids}}},
		}
		findOptions := getFindOptions(opts...)
		list, count, err := r.mFindList(ctx, filter, findOptions)
		return list, count, err
	})

}

func (r *Dao[T]) mFindList(ctx context.Context, filter any, opts ...*mongo_options.FindOptions) ([]T, bool, error) {
	var list []T
	sCtx := r.getSessionCtx(ctx)
	cursor, err := r.getCollection(sCtx).Find(sCtx, filter, opts...)
	if err != nil {
		return nil, false, err
	}
	err = cursor.All(ctx, &list)
	if err != nil {
		return nil, false, err
	}
	if len(list) == 0 {
		list = []T{}
		return list, false, nil
	}

	if r.eb.GetConfig().IsMap {
		for i, item := range list {
			if eMap, ok := any(item).(map[string]any); ok {
				e := r.db2entity(eMap)
				list[i] = e
			}
		}
	}

	return list, len(list) > 0, err
}

func (r *Dao[T]) mFindOne(ctx context.Context, filter any, opts ...*mongo_options.FindOneOptions) (T, bool, error) {
	var null T
	var entity T
	var err error

	sCtx := r.getSessionCtx(ctx)
	result := r.getCollection(sCtx).FindOne(sCtx, filter, opts...)
	err = result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return null, false, nil
		}
		return null, false, err
	}

	err = result.Decode(&entity)
	if err != nil {
		return null, false, err
	}

	if r.eb.GetConfig().IsMap {
		if eMap, ok := any(entity).(map[string]any); ok {
			r.db2entity(eMap)
		}
	}
	return entity, any(entity) != nil, err
}

func (r *Dao[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...store.Options) *store.FindOneResult[T] {
	return r.DoFindOne(func() (T, bool, error) {
		filter := r.NewFilter(tenantId, filterMap)
		logs.Info(ctx, logs.Fields{"filter": func() any {
			bytes, _ := json.Marshal(filter)
			return string(bytes)
		}})
		findOneOptions := getFindOneOptions(opts...)
		data, isFound, err := r.mFindOne(ctx, filter, findOneOptions)
		return data, isFound, err
	})
}

func (r *Dao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...store.Options) *store.FindListResult[T] {
	return r.DoFindList(func() ([]T, bool, error) {
		var list []T
		filter := r.NewFilter(tenantId, filterMap)
		findOptions := getFindOptions(opts...)
		list, count, err := r.mFindList(ctx, filter, findOptions)
		return list, count, err
	})
}

func (r *Dao[T]) FindListByBsonM(ctx context.Context, tenantId string, filter bson.M, opts ...store.Options) *store.FindListResult[T] {
	return r.DoFindList(func() ([]T, bool, error) {
		findOptions := getFindOptions(opts...)
		list, count, err := r.mFindList(ctx, filter, findOptions)
		return list, count, err
	})
}

func (r *Dao[T]) FindByRSQL(ctx context.Context, tenantId string, rsql string, opts ...store.Options) *store.FindListResult[T] {
	return r.doList(tenantId, rsql, func(filter *rsql_mongo.Filter) ([]T, bool, error) {
		list := r.NewEntityList()
		ctx = r.getSessionCtx(ctx)
		findOpts := &findByFilterOptions{
			resultsData: &list,
		}
		err := r.findByFilter(ctx, filter, findOpts)
		if err != nil {
			return nil, false, err
		}
		return list, len(list) > 0, err
	})
}

func (r *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...store.Options) *store.FindListResult[T] {
	return r.FindListByMap(ctx, tenantId, nil, opts...)
}

func (r *Dao[T]) findPaging(ctx context.Context, query store.FindPagingQuery, opts ...store.Options) store.FindPagingResult[T] {
	return r.doFilter(query.GetTenantId(), query.GetFilter(), func(filter *rsql_mongo.Filter) (store.FindPagingResult[T], bool, error) {
		if err := assert2.NotEmpty(query.GetTenantId(), assert2.NewOptions("tenantId is empty")); err != nil {
			return nil, false, err
		}
		ctx = r.getSessionCtx(ctx)
		data := r.NewEntityList()

		inOut := newFindOption(filter, query, &data)
		err := r.find(ctx, inOut)
		if err != nil {
			return nil, false, err
		}

		findData := store.NewFindPagingResult[T](data, inOut.totalRows, query, err)
		return findData, findData.GetIsFound(), err
	})
}

type findOption struct {
	filter    *rsql_mongo.Filter    // rsql的查询条件
	query     store.FindPagingQuery // 分页查询条件
	results   any                   // 返回数据
	totalRows int64                 // 返回记录数
}

func newFindOption(filter *rsql_mongo.Filter, query store.FindPagingQuery, results any) *findOption {
	return &findOption{
		filter:    filter,
		query:     query,
		results:   results,
		totalRows: 0,
	}
}

func (r *Dao[T]) newFindOptions(query store.FindPagingQuery, opts ...store.Options) (*mongo_options.FindOptions, error) {
	findOptions := getFindOptions(opts...)
	if query == nil {
		return findOptions, nil
	}

	if query.GetPageSize() > 0 {
		findOptions.SetLimit(query.GetPageSize())
		findOptions.SetSkip(query.GetPageSize() * query.GetPageNum())
	}
	if len(query.GetSort()) > 0 {
		sort, err := r.getSort(query.GetSort())
		if err != nil {
			return nil, err
		}
		findOptions.SetSort(sort)
	}

	if projection := r.getFindOptionsProjection(query); projection != nil {
		findOptions.SetProjection(projection)
	}
	return findOptions, nil
}

func (r *Dao[T]) getFindOptionsProjection(query store.FindPagingQuery) bson.D {
	var projection bson.D
	if len(query.GetFields()) > 0 {
		fields := strings.Split(query.GetFields(), ",")
		projection = bson.D{}
		for _, f := range fields {
			key := stringutils.SnakeString(strings.Trim(f, " "))
			projection = append(projection, bson.E{Key: key, Value: 1})
		}
	}
	return projection
}

/*
	func (r *Dao[T]) FindPaging2(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
		var err error
		findOptions := getFindOptions(opts...)
		queryGroup := NewQueryGroup(query)

		g, err := queryGroup.GetGroup()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		gpFilter, err := queryGroup.GetGroupPagingBsonFilter()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		gFilter, err := queryGroup.GetGroupExpandGroupNoPagingBsonFilter()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		gFilter1, err := queryGroup.GetGroupNoPagingFilter()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		filter, err := queryGroup.GetFilter()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		gSort, err := queryGroup.GetBsonFilterSort()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		sort, err := queryGroup.GetFilterSort()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		data, err := r.NewEntityList()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		coll := r.getCollection(ctx)
		var findData *ddd_repository.FindPagingResult[T]
		var cur *mongo.Cursor
		///var curt *mongo.Cursor
		var errt error
		var totalRows int64

		isGroup := queryGroup.IsGroup()
		isPaging := queryGroup.IsPaging()
		isLeaf := queryGroup.IsLeaf()

		if isGroup {
			if isPaging {
				pipeline := mongo.Pipeline{}
				if gpFilter != nil && len(gpFilter) > 0 {
					pipeline = append(pipeline, gpFilter)
				}
				if g != nil && len(g) > 0 {
					pipeline = append(pipeline, g)
				}
				if sort != nil && len(sort) > 0 {
					pipeline = append(pipeline, sort)
				}

				skip := query.GetPageSize() * query.GetPageNum()
				pipeline = append(pipeline, bson.D{{"$skip", skip}})

				limit := query.GetPageSize()
				pipeline = append(pipeline, bson.D{{"$limit", limit}})

				cur, err = coll.Aggregate(ctx, pipeline)
				if err == nil {
					totalRows = int64(cur.RemainingBatchLength())
				}
			} else if !isLeaf {
				pipeline := mongo.Pipeline{}
				if gFilter != nil && len(gFilter) > 0 {
					pipeline = append(pipeline, gFilter)
				}
				if g != nil && len(g) > 0 {
					pipeline = append(pipeline, g)
				}
				if gSort != nil && len(gSort) > 0 {
					pipeline = append(pipeline, gSort)
				}
				if cur, err = coll.Aggregate(ctx, pipeline); err == nil {
					totalRows = int64(cur.RemainingBatchLength())
				}
			} else if isLeaf {
				findOptions.SetSort(sort)
				cur, err = coll.Find(ctx, gFilter1, findOptions)
				if err == nil {
					totalRows, errt = coll.CountDocuments(ctx, gFilter1)
				}
			}
		} else if !isGroup {
			findOptions.SetSort(sort)
			if query.GetPageSize() > 0 {
				findOptions.SetLimit(query.GetPageSize())
				findOptions.SetSkip(query.GetPageSize() * query.GetPageNum())
			}
			cur, err = coll.Find(ctx, filter, findOptions)
			if query.GetIsTotalRows() {
				totalRows, errt = coll.CountDocuments(ctx, filter)
			}
		}
		if err != nil || errt != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err, errt)
		}

		err = cur.All(ctx, &data)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}

		findData = ddd_repository.NewFindPagingResult[T](data, &totalRows, query, err)
		// 进行汇总计算
		if len(query.GetValueCols()) > 0 {
			sumData, _, err := r.Sum(ctx, query, opts...)
			findData.SetSum(true, sumData, err)
		}
		return findData
	}
*/

func (r *Dao[T]) FindAutoComplete(ctx context.Context, qry store.FindAutoCompleteQuery, opts ...store.Options) store.FindPagingResult[T] {
	f := store.NewFindPagingQuery()
	groupCols := []*store.GroupCol{
		{Field: qry.GetField(), DataType: types.DataTypeString},
	}

	f.SetGroupCols(groupCols)
	f.SetTenantId(qry.GetTenantId())
	f.SetFields(qry.GetFields())
	f.SetFilter(qry.GetFilter())
	f.SetMustFilter(qry.GetMustWhere())

	f.SetPageNum(qry.GetPageNum())
	f.SetPageSize(qry.GetPageSize())
	f.SetSort(qry.GetSort())
	f.SetIsTotalRows(false)

	return r.FindPaging(ctx, f, opts...)
}

func (r *Dao[T]) FindDistinct(ctx context.Context, qry store.FindDistinctQuery, opts ...store.Options) store.FindPagingResult[T] {
	f := store.NewFindPagingQuery()

	f.SetGroupCols(qry.GetGroupCols())
	f.SetTenantId(qry.GetTenantId())
	f.SetFields(qry.GetFields())
	f.SetFilter(qry.GetFilter())
	f.SetMustFilter(qry.GetMustWhere())

	f.SetPageNum(qry.GetPageNum())
	f.SetPageSize(qry.GetPageSize())
	f.SetSort(qry.GetSort())
	f.SetIsTotalRows(false)

	return r.FindPaging(ctx, f, opts...)
}

func (r *Dao[T]) AggregateByPipeline(ctx context.Context, pipeline mongo.Pipeline, data interface{}, opts ...store.Options) error {
	sCtx := r.getSessionCtx(ctx)
	options := getAggregateOptions(opts...)
	cur, err := r.getCollection(ctx).Aggregate(sCtx, pipeline, options)
	if err != nil {
		return err
	}
	err = cur.All(ctx, data)
	return err
}

func (r *Dao[T]) CopyTo(ctx context.Context, tenantId string, rsql string, toCollectionName string, opts ...store.Options) error {
	options := getAggregateOptions(opts...)
	//db.record.aggregate([{$match:{opp_bank_name:"工商银行"}},{$out:"record1"}])
	filter, err := r.getFilter(tenantId, rsql)
	if err != nil {
		return err
	}
	pipeline := mongo.Pipeline{}
	if filter != nil && len(filter.Match) > 0 {
		pipeline = append(pipeline, bson.D{{"$match", filter.Match}})
	}
	pipeline = append(pipeline, bson.D{{"$out", toCollectionName}})
	sCtx := r.getSessionCtx(ctx)
	_, err = r.getCollection(ctx).Aggregate(sCtx, pipeline, options)
	return err
}

func (r *Dao[T]) SumEntity(ctx context.Context, qry store.FindPagingQuery, opts ...store.Options) ([]T, bool, error) {
	data := r.NewEntityList()
	sCtx := r.getSessionCtx(ctx)
	_, found, err := r.SumByQuery(sCtx, qry, &data, opts...)
	return data, found, err
}

func (r *Dao[T]) SumMap(ctx context.Context, qry store.FindPagingQuery, opts ...store.Options) ([]map[string]any, bool, error) {
	data := make([]map[string]any, 0)
	_, found, err := r.SumByQuery(r.getSessionCtx(ctx), qry, &data, opts...)
	return data, found, err
}

func (r *Dao[T]) SumByQuery(ctx context.Context, qry store.FindPagingQuery, data any, opts ...store.Options) (any, bool, error) {
	if len(qry.GetValueCols()) == 0 {
		return nil, false, nil
	}

	var err error
	process := rsql_mongo.NewProcess(qry.GetTenantId())

	f1 := qry.GetFilter()
	f2 := qry.GetMustFilter()
	f3 := ""
	mustWhere, ok := qry.(store.FindPagingQueryMustWhere)
	if ok {
		f3, err = mustWhere.GetMustWhere()
		if err != nil {
			return nil, false, err
		}
	}
	filterRSQL := getRsqlAnds(f1, f2, f3)

	if err := rsql.ParseProcess(filterRSQL, process); err != nil {
		return nil, false, err
	}
	filter, ok := process.GetFilter().(*rsql_mongo.Filter)
	if !ok {
		return nil, false, errors.New("filter does not implement Filter")
	}

	_, found, err := r.sum(ctx, filter, qry.GetValueCols(), data, opts...)
	return data, found, err
}

func (r *Dao[T]) SumByRSQL(ctx context.Context, tenantId string, rSql string, valueCols []*store.ValueCol, opts ...store.Options) map[string]any {
	filter, err := r.getFilter(tenantId, rSql)
	if err != nil {
		panic(err)
	}
	var list []map[string]any
	_, _, err = r.sum(ctx, filter, valueCols, &list, opts...)
	if err != nil {
		panic(err)
	}
	return list[0]
}

func (r *Dao[T]) sum(ctx context.Context, filter *rsql_mongo.Filter, valueCols []*store.ValueCol, list any, opts ...store.Options) (any, bool, error) {
	coll := r.getCollection(ctx)
	/*
		var match map[string]any
		if filter, ok := filterMap.(*rsql_mongo.Filter); ok {
			match = filter.Match
		} else if m, ok := filterMap.(map[string]any); ok {
			match = m
		} else {
			return nil, false, errors.New("sum() filterMap is not map[string]any")
		}
	*/
	var cur *mongo.Cursor
	summaryMap := make(map[string]interface{})
	summaryMap["_id"] = "total"
	for _, col := range valueCols {
		field := utils.SnakeString(col.Field)
		aggFunc := utils.SnakeString(col.AggFunc.Name())

		var fieldValue any = "$" + field
		if aggFunc == store.AggFuncCount.Name() {
			field = "count_" + field
			aggFunc = "$sum"
			fieldValue = 1
		} else {
			aggFunc = "$" + aggFunc
		}

		summaryMap[field] = map[string]interface{}{aggFunc: fieldValue}
	}

	pipeline := mongo.Pipeline{}
	filter.AddPipelines(pipeline)

	if summaryMap != nil {
		pipeline = append(pipeline, bson.D{{"$group", summaryMap}})
	}

	cur, err := coll.Aggregate(r.getSessionCtx(ctx), pipeline)
	if err != nil {
		result := store.NewFindPagingResultWithError[T](err)
		return result, false, err
	}
	err = cur.All(ctx, list)
	if err != nil {
		result := store.NewFindPagingResultWithError[T](err)
		return result, false, err
	}
	return list, true, nil

}

func (r *Dao[T]) CountByMap(ctx context.Context, tenantId string, filterData any, opts ...store.Options) (int64, error) {
	total, err := r.getCollection(ctx).CountDocuments(r.getSessionCtx(ctx), filterData)
	if err != nil {
		return 0, err
	}
	return total, err
}

func (r *Dao[T]) CountByRSQL(ctx context.Context, tenantId string, rsql string, opts ...store.Options) (int64, error) {
	f, err := r.getFilter(tenantId, rsql)
	if err != nil {
		return 0, err
	}
	ctx = r.getSessionCtx(ctx)
	total, err := r.getCollection(ctx).CountDocuments(ctx, f.Match)
	if err != nil {
		return 0, err
	}
	return total, err
}

func (r *Dao[T]) doList(tenantId, rsql string, fun func(filter *rsql_mongo.Filter) ([]T, bool, error)) *store.FindListResult[T] {
	if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
		return store.NewFindListResultError[T](err)
	}
	filterData, err := r.getFilter(tenantId, rsql)
	if err != nil {
		return store.NewFindListResultError[T](err)
	}
	data, ok, err := fun(filterData)
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			err = nil
		}
	}
	return store.NewFindListResult(data, ok, err)
}

func (r *Dao[T]) doFilter(tenantId, rsql string, fun func(filter *rsql_mongo.Filter) (store.FindPagingResult[T], bool, error)) store.FindPagingResult[T] {
	if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
		return store.NewFindPagingResultWithError[T](err)
	}
	filterData, err := r.getFilter(tenantId, rsql)
	if err != nil {
		return store.NewFindPagingResultWithError[T](err)
	}
	data, _, err := fun(filterData)
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			err = nil
		}
	}
	return data
}

func (r *Dao[T]) GetFilterMap(tenantId, rsql string) *rsql_mongo.Filter {
	data, err := r.getFilter(tenantId, rsql)
	if err != nil {
		panic(err)
	}
	return data
}

func (r *Dao[T]) getFilter(tenantId, rSql string) (*rsql_mongo.Filter, error) {
	process := rsql_mongo.NewProcess(tenantId)
	if err := rsql.ParseProcess(rSql, process); err != nil {
		return nil, err
	}
	filter := process.GetFilter().(*rsql_mongo.Filter)
	return filter, nil
}

func (r *Dao[T]) DoFindList(fun func() ([]T, bool, error)) *store.FindListResult[T] {
	data, isFound, err := fun()
	if err != nil {
		if errors.IsErrorMongoNoDocuments(err) {
			isFound = false
			err = nil
		}
	}
	return store.NewFindListResult[T](data, isFound, err)
}

func (r *Dao[T]) DoFindOne(fun func() (T, bool, error)) *store.FindOneResult[T] {
	data, isFound, err := fun()
	return store.NewFindOneResult[T](data, isFound, err)
}

func (r *Dao[T]) DoSet(fun func() (T, error)) *store.SetResult[T] {
	data, err := fun()
	return store.NewSetResult[T]().SetError(err).SetData(data)
}

func (r *Dao[T]) StartTx(ctx context.Context, txFun store.TxFunc, options ...*store.SessionOptions) (err error) {
	return StartTx(ctx, r.mongodb, r.mongodb.GetDatabase().Name(), txFun, options...)
}

/*
func (r *Dao[T]) DoSetMap(fun func() (map[string]interface{}, error)) *ddd_repository.SetResult[map[string]interface{}] {
	data, err := fun()
	return ddd_repository.NewSetResult[T](data, err)
}
*/

// getSort
// @Description: 返回排序bson.D
// @receiver r
// @param sort  排序语句 "name:desc,id:asc"
// @return bson.D
// @return error
func (r *Dao[T]) getSort(sort string) (bson.D, error) {
	if len(sort) == 0 {
		return bson.D{}, nil
	}
	// 输入
	// name:desc,id:asc
	// 输出
	/*	sort := bson.D{
		bson.E{"update_time", -1},
		bson.E{"goods_id", -1},
	}*/
	res := bson.D{}
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
			oerr = errors.New("order %s is error", order)
		}
		if oerr != nil {
			return nil, oerr
		}
		item := bson.E{name, orderVal}
		res = append(res, item)
	}
	return res, nil
}

/*
	func (r *Dao[T]) getDocuments(entities []T) []any {
		var list []any
		for _, item := range entities {
			list = append(list, r.getDocument(item))
		}
		return list
	}

	func (r *Dao[T]) getDocument(entity any) any {
		if dataMap, ok := entity.(DataMap); ok {
			return dataMap.GetDataMap()
		}
		return entity
	}
*/
func (r *Dao[T]) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...store.Options) (result store.FindPagingResult[T]) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				result = store.NewFindPagingResultWithError[T](err)
			}
		}
	}()
	var findData store.FindPagingResult[T]
	var err error
	//data := r.NewEntityList()
	queryGroup := NewQueryGroup(qry)
	//findOptions := getFindOptions(opts...)
	ctx = r.getSessionCtx(ctx)

	groupQueryResult := &findByGroupQueryOptions{}
	err = r.findByGroupQuery(ctx, queryGroup, groupQueryResult)
	if err != nil {
		panic(err)
	}
	data := groupQueryResult.results.([]T)
	if data == nil {
		data = []T{}
	}
	findData = store.NewFindPagingResult[T](data, groupQueryResult.totalRows, qry, err)

	// 进行汇总计算
	if len(qry.GetValueCols()) > 0 {
		sumData, _, err := r.SumEntity(ctx, qry, opts...)
		findData.SetSumData(sumData)
		findData.SetError(err)
		findData.SetIsSum(true)
	} else {
		sumData := []T{}
		findData.SetSumData(sumData)
		findData.SetError(err)
		findData.SetIsSum(false)
	}

	findData.SetIsTotalRows(qry.GetIsTotalRows())
	return findData

	g := queryGroup.GetGroup()

	totalGroup := queryGroup.GetTotalGroup()
	filter := queryGroup.GetFilter()
	//expandFilter := queryGroup.GetGroupExpandFilter()
	gSort := queryGroup.GetBsonFilterSort()
	//sort := queryGroup.GetFilterSort()

	coll := r.getCollection(ctx)

	var cur *mongo.Cursor
	var errt error
	var totalRows int64

	isGroup := queryGroup.IsGroup()
	isLeaf := queryGroup.IsLeaf()

	// 是分组查询
	if isGroup {
		// 不是树型叶子查询
		if !isLeaf {
			pipeline := filter.NewPipeline()
			if g != nil && len(g) > 0 {
				pipeline = append(pipeline, g)
			}
			if gSort != nil && len(gSort) > 0 {
				pipeline = append(pipeline, gSort)
			}
			pipeline = append(pipeline, totalGroup)
			pipeline = append(pipeline, bson.D{{
				"$project", map[string]interface{}{
					"_id":        "$_id",
					"data":       map[string]interface{}{"$slice": []interface{}{"$data", qry.GetPageSize() * qry.GetPageNum(), qry.GetPageSize()}},
					"total_rows": "$total_rows",
				},
			}})
			cur, err = coll.Aggregate(ctx, pipeline)
		}
	}
	// 标准查询
	if !isGroup || isLeaf {
		/*
			var f map[string]any
			if isLeaf {
				f = expandFilter
			} else {
				f = filter.Match
			}

			findOptions.SetSort(sort)
			if qry.GetPageSize() > 0 {
				findOptions.SetLimit(qry.GetPageSize())
				findOptions.SetSkip(qry.GetPageSize() * qry.GetPageNum())
			}
			if projection := r.getFindOptionsProjection(qry); projection != nil {
				findOptions.SetProjection(projection)
			}

			cur, err = coll.Find(sCtx, f, findOptions)
			if qry.GetIsTotalRows() {
				totalRows, errt = coll.CountDocuments(sCtx, f)
			}
		*/
		fOpt := newFindOption(filter, qry, &data)
		err = r.find(ctx, fOpt, opts...)
		totalRows = fOpt.totalRows

	}
	if err != nil || errt != nil {
		return store.NewFindPagingResultWithError[T](err, errt)
	}

	if isGroup && !isLeaf {
		d := make([]struct {
			Data      []T   `json:"data" bson:"data"`
			TotalRows int64 `json:"totalRows" bson:"total_rows"`
		}, 0)
		err = cur.All(ctx, &d)
		if err != nil {
			return store.NewFindPagingResultWithError[T](err)
		}
		if d != nil && len(d) == 1 {
			data = d[0].Data
			totalRows = d[0].TotalRows
		}
	}

	if data == nil {
		err = cur.All(ctx, &data)
		if err != nil {
			return store.NewFindPagingResultWithError[T](err)
		}
	}

	if data == nil {
		data = []T{}
	}
	findData = store.NewFindPagingResult[T](data, totalRows, qry, err)
	// 进行汇总计算
	if len(qry.GetValueCols()) > 0 {
		sumData, _, err := r.SumEntity(ctx, qry, opts...)
		findData.SetSumData(sumData)
		findData.SetError(err)
		findData.SetIsSum(true)
	} else {
		sumData := []T{}
		findData.SetSumData(sumData)
		findData.SetError(err)
		findData.SetIsSum(false)
	}
	findData.SetIsTotalRows(qry.GetIsTotalRows())
	return findData
}

func (r *Dao[T]) find(ctx context.Context, fOpt *findOption, opts ...store.Options) error {
	filter := fOpt.filter

	var qry *QueryGroup
	if fOpt.query != nil {
		qry = NewQueryGroup(fOpt.query)
	}

	if qry == nil {
		coll := r.getCollection(ctx)
		if !filter.IsAggregate() {
			findOptions, err := r.newFindOptions(fOpt.query, opts...)
			if err != nil {
				return err
			}
			cursor, err := coll.Find(ctx, filter.Match, findOptions)
			if err != nil {
				return err
			}
			err = cursor.All(ctx, fOpt.results)
			if fOpt.query.GetIsTotalRows() {
				total, err := coll.CountDocuments(ctx, filter.Match)
				if err != nil {
					return err
				}
				fOpt.totalRows = total
			}
			return err
		} else {
			pipeline := filter.NewPipeline()
			cur, err := coll.Aggregate(ctx, pipeline)
			if err != nil {
				return err
			}
			err = cur.All(ctx, fOpt.results)
			if err != nil {
				return err
			}
		}
	} else {
		results := &findByGroupQueryOptions{
			results: fOpt.results,
		}
		err := r.findByGroupQuery(ctx, qry, results, opts...)
		if err != nil {
			return err
		}
		fOpt.totalRows = results.totalRows
	}
	return nil
}

type findByFilterOptions struct {
	resultsData      any
	resultsTotalRows int64
	isTotalRows      bool
}

func (r *Dao[T]) findByFilter(ctx context.Context, filter *rsql_mongo.Filter, opt *findByFilterOptions) error {
	coll := r.getCollection(ctx)
	if !filter.IsAggregate() {
		cursor, err := coll.Find(ctx, filter.Match)
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)
		err = cursor.All(ctx, opt.resultsData)

		if opt.isTotalRows {
			total, err := coll.CountDocuments(ctx, filter.Match)
			if err != nil {
				return err
			}
			opt.resultsTotalRows = total
		}

		return err

	} else {
		pipeline := filter.NewPipeline()
		cur, err := coll.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		mongoutils.PrintPipeline(pipeline)

		defer cur.Close(ctx)

		err = cur.All(ctx, opt.resultsData)
		if err != nil {
			return err
		}
		if opt.isTotalRows {
			opt.resultsTotalRows, err = r.aggregateTotal(ctx, coll, pipeline)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type findByGroupQueryOptions struct {
	results   any
	totalRows int64
}

func (r *Dao[T]) findByGroupQuery(ctx context.Context, qry *QueryGroup, results *findByGroupQueryOptions, opts ...store.Options) error {
	coll := r.getCollection(ctx)

	filter := qry.GetFilter()
	if !filter.IsAggregate() && !qry.IsGroup() {
		findOptions, err := r.newFindOptions(qry.Query, opts...)
		if err != nil {
			return err
		}
		match := filter.Match
		list, _, err := r.mFindList(ctx, match, findOptions)
		results.results = list

		if qry.Query.GetIsTotalRows() {
			total, err := coll.CountDocuments(ctx, match)
			if err != nil {
				return err
			}
			results.totalRows = total
		}
		return err
	}

	pipeline := filter.NewPipeline()
	gSort := qry.GetBsonFilterSort()
	if gSort != nil && len(gSort) > 0 {
		pipeline = append(pipeline, gSort)
	}

	if qry.IsGroup() {
		g := qry.GetGroup()
		if g != nil && len(g) > 0 {
			pipeline = append(pipeline, g)
		}
		totalGroup := qry.GetTotalGroup()
		pipeline = append(pipeline, totalGroup)
		project := bson.D{{
			"$project", bson.M{
				"_id": "$_id",
				"data": bson.M{
					"$slice": bson.A{
						"$data",
						qry.GetPageSize() * qry.GetPageNum(),
						qry.GetPageSize(),
					},
				},
				"total_rows": "$total_rows",
			},
		}}
		pipeline = append(pipeline, project)
		/*
					pipeline = append(pipeline, bson.M{
			,			"$project": bson.M{
							"_id": "$_id",
							"data": map[string]any{
								"$slice": []any{
									"$data",
									qry.GetPageSize() * qry.GetPageNum(),
									qry.GetPageSize(),
								},
							},
							"total_rows": "$total_rows",
						}
					}) */
	}
	mongoutils.PrintPipeline(pipeline)
	cur, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, results.results)
	if err != nil {
		return err
	}

	if qry.Query.GetIsTotalRows() {
		total, err := r.aggregateTotal(ctx, coll, pipeline)
		if err != nil {
			return err
		}
		results.totalRows = total
	}

	return err

}

// aggregateTotal 使用aggregate方法进行count计算
func (r *Dao[T]) aggregateTotal(ctx context.Context, coll *mongo.Collection, pipeline mongo.Pipeline) (int64, error) {
	pipeline = append(pipeline, bson.D{{"$count", "totalCount"}})

	cur, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	for cur.Next(ctx) {
		var totalRes bson.M
		if err := cur.Decode(&totalRes); err != nil {
			return 0, err
		}
		if total, ok := totalRes["totalCount"].(int32); ok {
			return int64(total), nil
		}
	}

	return 0, nil
}
