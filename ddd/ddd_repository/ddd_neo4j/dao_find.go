package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"strings"
	"time"
)

func (d *Dao[T]) FindById(ctx context.Context, tenantId, id string, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	res := ddd_repository.NewFindOneResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.cypher.FindById(ctx, tenantId, id)
		if err != nil {
			return err
		}
		result, err := d.Query(ctx, cr.Cypher(), cr.Params())
		if err != nil {
			return err
		}
		entity := d.eb.NewEntity()
		_, err = result.GetOne(cr.ResultOneKey(), entity)
		res.SetData(entity)
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	res := ddd_repository.NewFindListResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.cypher.FindByIds(ctx, tenantId, ids)
		if err != nil {
			return err
		}
		result, err := d.Query(ctx, cr.Cypher(), cr.Params())
		if err != nil {
			return err
		}
		list, err := reflectutils.NewSlice[[]T]()
		if err != nil {
			return err
		}
		if err := result.GetList(ctx, cr.ResultOneKey(), &list); err != nil {
			return err
		}
		res.SetData(list)
		return nil
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
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
	list := d.NewEntityList()
	if err := result.GetList(ctx, cr.ResultOneKey(), &list); err != nil {
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
	return d.FindByRSQL(ctx, tenantId, filter)
}

func (d *Dao[T]) findPagingByCypher(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	res := ddd_repository.NewFindPagingResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.cypher.FindPaging(ctx, query)
		if err != nil {
			return err
		}

		err = assert.NotEmpty(query.GetTenantId(), assert.NewOptions("TenantId cannot be empty"))
		if err != nil {
			return err
		}

		cypher := cr.Cypher()
		result, err := d.Query(ctx, cypher, cr.Params())
		if err != nil {
			return err
		}

		list := d.NewEntityList()
		if err = result.GetList(ctx, cr.ResultOneKey(), &list); err != nil {
			return err
		}
		res.SetData(list)

		if query.GetIsTotalRows() {
			countCr, err := d.cypher.Count(ctx, query.GetTenantId(), query.GetFilter())
			result, err := d.Query(ctx, countCr.Cypher(), countCr.Params())
			total, err := result.GetInteger(countCr.ResultOneKey(), 0)
			if err != nil {
				return err
			}
			res.SetTotalRow(total)
		}

		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return d.findPagingByCypher(ctx, query, opts...)
}

func (d *Dao[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) FindByRSQL(ctx context.Context, tenantId, filter string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	res := ddd_repository.NewFindListResultEmpty[T]()
	gp.Try(func() error {
		if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
			return err
		}

		cr, err := d.cypher.GetFilter(ctx, tenantId, filter)
		if err != nil {
			return err
		}

		cypher := cr.Cypher()
		result, err := d.Query(ctx, cypher, cr.Params())
		if err != nil {
			return err
		}

		list := d.NewEntityList()
		if err = result.GetList(ctx, cr.ResultOneKey(), &list); err != nil {
			return err
		}
		res.SetData(list)
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
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
