package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql/rsql_mongo"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/mongoutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Dao[T]) FindPaging(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) (result *ddd_repository.FindPagingResult[T]) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				result = ddd_repository.NewFindPagingResultWithError[T](err)
			}
		}
	}()
	var findData *ddd_repository.FindPagingResult[T]
	var err error
	data := r.NewEntityList()
	queryGroup := NewQueryGroup(qry)
	//findOptions := getFindOptions(opts...)
	ctx = r.getSessionCtx(ctx)

	groupQueryResult := &findByGroupQueryOptions{
		results: &data,
	}
	err = r.findByGroupQuery(ctx, queryGroup, groupQueryResult)
	if err != nil {
		panic(err)
	}

	if data == nil {
		data = []T{}
	}
	findData = ddd_repository.NewFindPagingResult[T](data, groupQueryResult.totalRows, qry, err)

	// 进行汇总计算
	if len(qry.GetValueCols()) > 0 {
		sumData, _, err := r.SumEntity(ctx, qry, opts...)
		findData.SetSum(true, sumData, err)
	} else {
		sumData := []T{}
		findData.SetSum(false, sumData, err)
	}

	findData.IsTotalRows = qry.GetIsTotalRows()
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
		return ddd_repository.NewFindPagingResultWithError[T](err, errt)
	}

	if isGroup && !isLeaf {
		d := make([]struct {
			Data      []T   `json:"data" bson:"data"`
			TotalRows int64 `json:"totalRows" bson:"total_rows"`
		}, 0)
		err = cur.All(ctx, &d)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}
		if d != nil && len(d) == 1 {
			data = d[0].Data
			totalRows = d[0].TotalRows
		}
	}

	if data == nil {
		err = cur.All(ctx, &data)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}
	}

	if data == nil {
		data = []T{}
	}
	findData = ddd_repository.NewFindPagingResult[T](data, totalRows, qry, err)
	// 进行汇总计算
	if len(qry.GetValueCols()) > 0 {
		sumData, _, err := r.SumEntity(ctx, qry, opts...)
		findData.SetSum(true, sumData, err)
	} else {
		sumData := []T{}
		findData.SetSum(false, sumData, err)
	}
	findData.IsTotalRows = qry.GetIsTotalRows()
	return findData
}

func (r *Dao[T]) find(ctx context.Context, fOpt *findOption, opts ...ddd_repository.Options) error {
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

func (r *Dao[T]) findByGroupQuery(ctx context.Context, qry *QueryGroup, results *findByGroupQueryOptions, opts ...ddd_repository.Options) error {
	coll := r.getCollection(ctx)

	filter := qry.GetFilter()
	if !filter.IsAggregate() && !qry.IsGroup() {
		findOptions, err := r.newFindOptions(qry.Query, opts...)
		if err != nil {
			return err
		}
		match := filter.Match
		//mongoutils.GetMQL(match)

		cursor, err := coll.Find(ctx, match, findOptions)
		if err != nil {
			return err
		}
		err = cursor.All(ctx, results.results)
		defer cursor.Close(ctx)

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
