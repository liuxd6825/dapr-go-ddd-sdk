package store_neo4j

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
)

func (d *Dao[T]) FindById(ctx context.Context, tenantId, id string, opts ...store2.Options) *store2.FindOneResult[T] {
	res := store2.NewFindOneResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.FindById(ctx, tenantId, id)
		if err != nil {
			return err
		}
		result, err := d.Query(ctx, cr.Cypher(), cr.Params())
		if err != nil {
			return err
		}
		entity := d.config.EntityBuilder.NewEntity()
		_, err = result.GetOne(cr.ResultOneKey(), entity, d.config.DBSchema)
		res.SetData(entity)
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...store2.Options) *store2.FindListResult[T] {
	res := store2.NewFindListResultEmpty[T]()
	gp.Try(func() error {
		cr, err := d.Cypher.FindByIds(ctx, tenantId, ids)
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
		if err := result.GetList(ctx, cr.ResultOneKey(), &list, d.GetSchema()); err != nil {
			return err
		}
		res.SetData(list)
		return nil
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...store2.Options) *store2.FindListResult[T] {
	cr, err := d.Cypher.FindAll(ctx, tenantId)
	if err != nil {
		return store2.NewFindListResultError[T](err)
	}
	result, err := d.Query(ctx, cr.Cypher(), cr.Params())
	if err != nil {
		return store2.NewFindListResultError[T](err)
	}
	list := d.NewEntityList()
	if err := result.GetList(ctx, cr.ResultOneKey(), &list, d.GetSchema()); err != nil {
		return store2.NewFindListResultError[T](err)
	}
	return store2.NewFindListResult[T](list, len(list) > 0, nil)
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

func (d *Dao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...store2.Options) *store2.FindListResult[T] {
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

func (d *Dao[T]) findPagingByCypher(ctx context.Context, query store2.FindPagingQuery, opts ...store2.Options) store2.FindPagingResult[T] {
	res := store2.NewFindPagingResultEmpty[T]()
	gp.Try(func() error {
		tenantId := appctx.GetTenantId2(ctx)
		err := assert2.NotEmpty(tenantId, assert2.NewOptions("TenantId cannot be empty"))
		if err != nil {
			return err
		}

		cr, err := d.Cypher.FindPaging(ctx, query)
		if err != nil {
			return err
		}

		cypher := cr.Cypher()

		result, err := d.Query(ctx, cypher, cr.Params())
		if err != nil {
			return err
		}

		list := d.NewEntityList()
		if err = result.GetList(ctx, cr.ResultOneKey(), &list, d.GetSchema()); err != nil {
			return err
		}
		res.SetData(list)
		res.SetFilter(query.GetFilter())
		res.SetFields(query.GetFields())
		res.SetPageSize(query.GetPageSize())
		res.SetPageNum(query.GetPageNum())
		if len(list) > 0 {
			res.SetIsFound(true)
		}

		if query.GetIsTotalRows() {
			count, err := d.Cypher.Count(ctx, tenantId, query.GetFilter())
			if err != nil {
				return err
			}
			result, err := d.Query(ctx, count.Cypher(), count.Params())
			if err != nil {
				return err
			}
			total, err := result.GetInteger(count.ResultOneKey(), 0)
			if err != nil {
				return err
			}
			res.SetTotalRows(total)
		}

		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) FindPaging(ctx context.Context, query store2.FindPagingQuery, opts ...store2.Options) store2.FindPagingResult[T] {
	return d.findPagingByCypher(ctx, query, opts...)
}

func (d *Dao[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...store2.Options) *store2.FindOneResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) FindByRSQL(ctx context.Context, tenantId, filter string, opts ...store2.Options) *store2.FindListResult[T] {
	res := store2.NewFindListResultEmpty[T]()
	gp.Try(func() error {
		if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
			return err
		}

		cr, err := d.Cypher.GetRSQL(ctx, tenantId, filter)
		if err != nil {
			return err
		}

		cypher := cr.Cypher()
		result, err := d.Query(ctx, cypher, cr.Params())
		if err != nil {
			return err
		}

		list := d.NewEntityList()
		if err = result.GetList(ctx, cr.ResultOneKey(), &list, d.GetSchema()); err != nil {
			return err
		}
		res.SetData(list)
		return err
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) FindAutoComplete(ctx context.Context, qry store2.FindAutoCompleteQuery, opts ...store2.Options) store2.FindPagingResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) FindDistinct(ctx context.Context, qry store2.FindDistinctQuery, opts ...store2.Options) store2.FindPagingResult[T] {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...store2.Options) (T, error) {
	//TODO implement me
	panic("implement me")
}
