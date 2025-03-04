package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"strings"
	"time"
)

func (d *Dao[T]) FindById(ctx context.Context, tenantId, id string, opts ...ddd_repository.Options) (T, bool, error) {
	var null T
	cr, err := d.cypher.FindById(ctx, tenantId, id)
	if err != nil {
		return null, false, err
	}
	result, err := d.Query(ctx, cr.Cypher(), cr.Params())
	if err != nil {
		return null, false, err
	}
	entity, err := reflectutils.NewStruct[T]()
	if err != nil {
		return null, false, err
	}
	if ok, err := result.GetOne("", entity); err != nil {
		return null, false, err
	} else if !ok {
		return null, false, nil
	}
	return entity, true, nil
}

func (d *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) ([]T, bool, error) {
	var null []T
	cr, err := d.cypher.FindByIds(ctx, tenantId, ids)
	if err != nil {
		return null, false, err
	}
	result, err := d.Query(ctx, cr.Cypher(), cr.Params())
	if err != nil {
		return null, false, err
	}
	list, err := reflectutils.NewSlice[[]T]()
	if err != nil {
		return null, false, err
	}
	if err := result.GetList(cr.ResultOneKey(), &list); err != nil {
		return null, false, err
	}
	return list, len(list) > 0, nil
}

func (d *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	cr, err := d.cypher.FindAll(ctx, tenantId)
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	result, err := d.Query(ctx, cr.Cypher(), cr.Params())
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	list, err := reflectutils.NewSlice[[]T]()
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	if err := result.GetList(cr.ResultOneKey(), &list); err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	return ddd_repository.NewFindListResult[T](list, len(list) > 0, nil)
}

/*func (d *Dao[T]) FindByGraphId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	cr, err := d.cypher.FindByGraphId(ctx, tenantId, graphId)
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	result, err := d.Query(ctx, cr.Cypher(), cr.Params())
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	list, err := reflectutils.NewSlice[[]T]()
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	if err := result.GetLists(cr.ResultKeys(), &list); err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	return ddd_repository.NewFindListResult[T](list, len(list) > 0, err)
}*/

func (d *Dao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	sb := strings.Builder{}
	for k, v := range filterMap {
		switch v.(type) {
		case string:
			sb.WriteString(fmt.Sprintf("%v=='%v'", k, v))
		case time.Time:
			sb.WriteString(fmt.Sprintf("%v=='%v'", k, v))
		case *time.Time:
			sb.WriteString(fmt.Sprintf("%v=='%v'", k, v))
		default:
			sb.WriteString(fmt.Sprintf("%v==%v", k, v))
		}
		sb.WriteString(" and ")
	}
	filter := sb.String()
	if strings.HasSuffix(filter, " and ") {
		filter = filter[0 : len(filter)-5]
	}
	return d.FindByFilter(ctx, tenantId, filter)
}

func (d *Dao[T]) FindByFilter(ctx context.Context, tenantId, filter string) *ddd_repository.FindListResult[T] {
	return d.DoList(ctx, tenantId, func() (*ddd_repository.FindListResult[T], bool, error) {
		if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
			return nil, false, err
		}

		cr, err := d.cypher.GetFilter(ctx, tenantId, filter)
		if err != nil {
			return ddd_repository.NewFindListResultError[T](err), false, err
		}

		cypher := cr.Cypher()
		result, err := d.Query(ctx, cypher, cr.Params())
		if err != nil {
			return ddd_repository.NewFindListResultError[T](err), false, err
		}

		list, err := reflectutils.NewSlice[[]T]()
		if err != nil {
			return ddd_repository.NewFindListResultError[T](err), false, err
		}

		if err = result.GetList(cr.ResultOneKey(), &list); err != nil {
			return ddd_repository.NewFindListResultError[T](err), false, err
		}
		return ddd_repository.NewFindListResult[T](list, len(list) > 0, err), false, nil
	})
}

func (d *Dao[T]) FindPagingByCypher(ctx context.Context, tenantId, cypher string, pageNum, pageSize int64, resultKey string, isTotalRows bool, params map[string]any, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return d.DoFilter(ctx, tenantId, func() (*ddd_repository.FindPagingResult[T], bool, error) {
		return d.findPagingByCypher(ctx, tenantId, cypher, pageNum, pageSize, resultKey, isTotalRows, params, opts...)
	})
}

func (d *Dao[T]) findPagingByCypher(ctx context.Context, tenantId, cypher string, pageNum, pageSize int64, resultKey string, isTotalRows bool, params map[string]any, opts ...ddd_repository.Options) (*ddd_repository.FindPagingResult[T], bool, error) {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("TenantId cannot be empty")); err != nil {
		return nil, false, err
	}
	result, err := d.Query(ctx, cypher+" RETURN "+resultKey, params)
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err), false, err
	}

	list, err := reflectutils.NewSlice[[]T]()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err), false, err
	}

	if err = result.GetList(resultKey, &list); err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err), false, err
	}

	var totalRows int64
	if isTotalRows {
		totalKey := "count"
		countCypher := cypher + fmt.Sprintf(" RETURN count(%s) as %s ", resultKey, totalKey)
		result, err := d.Query(ctx, countCypher, params)
		total, err := result.GetInteger(totalKey, 0)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err), false, err
		}
		totalRows = total
	}

	res := ddd_repository.NewFindPagingResult[T](list, totalRows, nil, nil)
	return res, true, err
}

func (d *Dao[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return d.DoFilter(ctx, query.GetTenantId(), func() (*ddd_repository.FindPagingResult[T], bool, error) {
		cr, err := d.cypher.FindPaging(ctx, query)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err), false, err
		}
		return d.findPagingByCypher(ctx, query.GetTenantId(), cr.Cypher(), query.GetPageNum(), query.GetPageSize(), cr.ResultKeys()[0], query.GetIsTotalRows(), cr.Params(), opts...)
	})
}

func (d *Dao[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) FindByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) FindAutoComplete(ctx context.Context, qry ddd_repository.FindAutoCompleteQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) FindDistinct(ctx context.Context, qry ddd_repository.FindDistinctQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) (T, error) {
	//TODO implement me
	panic("implement me")
}
