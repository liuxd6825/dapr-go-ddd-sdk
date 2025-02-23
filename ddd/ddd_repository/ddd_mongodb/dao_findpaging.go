package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
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

	var err error
	//findOptions := getFindOptions(opts...)

	queryGroup := NewQueryGroup(qry)
	g := queryGroup.GetGroup()

	totalGroup := queryGroup.GetTotalGroup()
	filter := queryGroup.GetFilter()
	//expandFilter := queryGroup.GetGroupExpandFilter()
	gSort := queryGroup.GetBsonFilterSort()
	//sort := queryGroup.GetFilterSort()
	data := r.NewEntityList()

	coll := r.getCollection(ctx)
	var findData *ddd_repository.FindPagingResult[T]
	var cur *mongo.Cursor
	var errt error
	var totalRows int64

	isGroup := queryGroup.IsGroup()
	isLeaf := queryGroup.IsLeaf()

	sCtx := r.getSessionCtx(ctx)

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
			cur, err = coll.Aggregate(sCtx, pipeline)
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
	findData = ddd_repository.NewFindPagingResult[T](data, &totalRows, qry, err)
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
	qry := fOpt.query
	filter := fOpt.filter
	coll := r.getCollection(ctx)
	if !filter.IsAggregate() {
		findOptions, err := r.newFindOptions(qry, opts...)
		if err != nil {
			return err
		}
		cursor, err := coll.Find(ctx, filter.Match, findOptions)
		if err != nil {
			return err
		}
		err = cursor.All(ctx, fOpt.results)
		if qry.GetIsTotalRows() {
			total, err := coll.CountDocuments(ctx, filter.Match)
			if err != nil {
				return err
			}
			fOpt.totalRows = total
		}
	} else {
		queryGroup := NewQueryGroup(qry)
		g := queryGroup.GetGroup()

		totalGroup := queryGroup.GetTotalGroup()
		filter := queryGroup.GetFilter()
		//expandFilter := queryGroup.GetGroupExpandFilter()
		gSort := queryGroup.GetBsonFilterSort()

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
		cur, err := coll.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}
		err = cur.All(ctx, fOpt.results)
		if err != nil {
			return err
		}
	}
	return nil
}
