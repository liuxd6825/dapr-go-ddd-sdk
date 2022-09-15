package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"strings"
	"time"
)

type Neo4jEntity interface {
	ddd.Entity
}

type Neo4jDao[T ElementEntity] struct {
	driver        neo4j.Driver
	cypherBuilder CypherBuilder
}

type SessionOptions struct {
	AccessMode *neo4j.AccessMode
}

func NewSessionOptions() *SessionOptions {
	return &SessionOptions{}
}

func NewNeo4jDao[T ElementEntity](driver neo4j.Driver, builder CypherBuilder) *Neo4jDao[T] {
	base := &Neo4jDao[T]{}
	return base.Init(driver, builder)
}

func (d *Neo4jDao[T]) Init(driver neo4j.Driver, builder CypherBuilder) *Neo4jDao[T] {
	d.driver = driver
	d.cypherBuilder = builder
	return d
}

func (d *Neo4jDao[T]) NewEntity() (res T, resErr error) {
	return reflectutils.NewStruct[T]()
}

func (d *Neo4jDao[T]) NewEntities() (res []T, resErr error) {
	return reflectutils.NewSlice[[]T]()
}

func (d *Neo4jDao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
	defer func() {
		if e := recover(); e != nil {
			if err := errors.GetRecoverError(e); err != nil {
				setResult = ddd_repository.NewSetResultError[T](err)
			}
		}
	}()

	cr, err := d.cypherBuilder.Insert(ctx, entity)
	if err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}

	res, err := d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...)
	if _, err := res.GetOne("", entity); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetResult(entity, err)
}

func (d *Neo4jDao[T]) InsertMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	for _, e := range entities {
		if err := d.Insert(ctx, e, opts...).GetError(); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}
	return ddd_repository.NewSetManyResult[T](entities, nil)
}

func (d *Neo4jDao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	cr, err := d.cypherBuilder.Update(ctx, entity)
	res, err := d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...)
	if _, err := res.GetOne("", entity); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetResult(entity, err)
}

func (d *Neo4jDao[T]) UpdateMany(ctx context.Context, list []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	for _, entity := range list {
		if cr, err := d.cypherBuilder.Update(ctx, entity); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		} else {
			if res, err := d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...); err != nil {
				return ddd_repository.NewSetManyResultError[T](err)
			} else if _, err := res.GetOne(cr.ResultKeys()[0], entity); err != nil {
				return ddd_repository.NewSetManyResultError[T](err)
			}
		}
	}
	return ddd_repository.NewSetManyResult(list, nil)
}

func (d *Neo4jDao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) error {
	cr, err := d.cypherBuilder.DeleteById(ctx, tenantId, id)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Neo4jDao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error {
	cr, err := d.cypherBuilder.DeleteByIds(ctx, tenantId, ids)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Neo4jDao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error {
	cr, err := d.cypherBuilder.DeleteAll(ctx, tenantId)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Neo4jDao[T]) DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...ddd_repository.Options) error {
	cr, err := d.cypherBuilder.DeleteByFilter(ctx, tenantId, filter)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Neo4jDao[T]) FindById(ctx context.Context, tenantId, id string, opts ...ddd_repository.Options) (T, bool, error) {
	var null T
	cr, err := d.cypherBuilder.FindById(ctx, tenantId, id)
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

func (d *Neo4jDao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) ([]T, bool, error) {
	var null []T
	cr, err := d.cypherBuilder.FindByIds(ctx, tenantId, ids)
	if err != nil {
		return null, false, err
	}
	result, err := d.Query(ctx, cr.Cypher(), cr.Params())
	if err != nil {
		return null, false, err
	}
	list, err := d.NewEntities()
	if err != nil {
		return null, false, err
	}
	if err := result.GetList(cr.ResultOneKey(), &list); err != nil {
		return null, false, err
	}
	return list, len(list) > 0, nil
}

func (d *Neo4jDao[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	cr, err := d.cypherBuilder.FindAll(ctx, tenantId)
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

func (d *Neo4jDao[T]) FindByGraphId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	cr, err := d.cypherBuilder.FindByGraphId(ctx, tenantId, graphId)
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
}

func (d *Neo4jDao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
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

func (d *Neo4jDao[T]) FindByFilter(ctx context.Context, tenantId, filter string) *ddd_repository.FindListResult[T] {
	return d.DoList(ctx, tenantId, filter, func() (*ddd_repository.FindListResult[T], bool, error) {
		if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
			return nil, false, err
		}

		cr, err := d.cypherBuilder.GetFilter(ctx, tenantId, filter)
		if err != nil {
			return ddd_repository.NewFindListResultError[T](err), false, err
		}

		cyhper := cr.Cypher()
		result, err := d.Query(ctx, cyhper, cr.Params())
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

func (d *Neo4jDao[T]) Write(ctx context.Context, cypher string) (*Neo4jResult, error) {
	return d.doSession(ctx, func(tx neo4j.Transaction) (*Neo4jResult, error) {
		result, err := tx.Run(cypher, nil)
		if err != nil {
			return nil, err
		}
		return NewNeo4jResult(result), err
	})
}

func (d *Neo4jDao[T]) Query(ctx context.Context, cypher string, params map[string]interface{}) (*Neo4jResult, error) {
	var resultData *Neo4jResult
	_, err := d.doSession(ctx, func(tx neo4j.Transaction) (*Neo4jResult, error) {
		result, err := tx.Run(cypher, params)
		if err != nil {
			log.Println("wirte to DB with error:", err)
			return nil, err
		}
		resultData = NewNeo4jResult(result)
		return nil, err
	})
	return resultData, err
}

func (d *Neo4jDao[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return d.DoFilter(query.GetTenantId(), query.GetFilter(), func() (*ddd_repository.FindPagingResult[T], bool, error) {
		if err := assert.NotEmpty(query.GetTenantId(), assert.NewOptions("tenantId is empty")); err != nil {
			return nil, false, err
		}

		cr, err := d.cypherBuilder.FindPaging(ctx, query)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err), false, err
		}

		cyhper := cr.Cypher()
		result, err := d.Query(ctx, cyhper, cr.Params())
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err), false, err
		}

		list, err := reflectutils.NewSlice[[]T]()
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err), false, err
		}

		listKey := cr.ResultKeys()[0]
		if err = result.GetList(listKey, &list); err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err), false, err
		}

		var totalRows *int64
		if query.GetIsTotalRows() {
			totalKey := cr.ResultKeys()[1]
			total, err := result.GetInteger(totalKey, 0)
			if err != nil {
				return ddd_repository.NewFindPagingResultWithError[T](err), false, err
			}
			totalRows = &total
		}

		return ddd_repository.NewFindPagingResult[T](list, totalRows, query, nil), true, err
	})
}

func (d *Neo4jDao[T]) DoFilter(tenantId, filter string, fun func() (*ddd_repository.FindPagingResult[T], bool, error), opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	p := NewRSqlProcess()
	if err := ParseProcess(filter, p); err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	data, _, err := fun()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	return data
}

func (d *Neo4jDao[T]) DoList(ctx context.Context, tenantId string, filter string, fun func() (*ddd_repository.FindListResult[T], bool, error), opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	p := NewRSqlProcess()
	if err := ParseProcess(filter, p); err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	data, _, err := fun()
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	return data
}

func (d *Neo4jDao[T]) doSession(ctx context.Context, fun func(tx neo4j.Transaction) (*Neo4jResult, error), opts ...*SessionOptions) (*Neo4jResult, error) {
	if fun == nil {
		return nil, errors.New("doSession(ctx, fun) fun is nil")
	}
	if sc, ok := GetSessionContext(ctx); ok {
		tx := sc.GetTransaction()
		_, err := fun(tx)
		return nil, err
	}

	opt := NewSessionOptions()
	opt.Merge(opts...)
	opt.setDefault()

	session := d.driver.NewSession(neo4j.SessionConfig{AccessMode: *opt.AccessMode})
	defer func() {
		_ = session.Close()
	}()

	var res interface{}
	var err error
	if *opt.AccessMode == neo4j.AccessModeRead {
		res, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return fun(tx)
		})
	} else if *opt.AccessMode == neo4j.AccessModeWrite {
		res, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return fun(tx)
		})
	}

	if result, ok := res.(*Neo4jResult); ok {
		return result, err
	}

	return nil, err
}

func (d *Neo4jDao[T]) doSet(ctx context.Context, tenantId string, cypher string, params map[string]interface{}, opts ...ddd_repository.Options) (*Neo4jResult, error) {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return nil, err
	}
	res, err := d.doSession(ctx, func(tx neo4j.Transaction) (*Neo4jResult, error) {
		r, err := tx.Run(cypher, params)
		return NewNeo4jResult(r), err
	})
	return res, err
}

func (d *Neo4jDao[T]) getLabels(entity ElementEntity) string {
	label := ""
	for _, l := range entity.GetLabels() {
		label = label + " :" + l
	}
	return label
}

func (d *Neo4jDao[T]) newSetManyResult(result *Neo4jResult, err error) *ddd_repository.SetManyResult[T] {
	if err != nil {
		return ddd_repository.NewSetManyResultError[T](err)
	}
	var data []T
	if err := result.GetList("n", &data); err != nil {
		ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetManyResult[T](data, err)
}

func (r *SessionOptions) SetAccessMode(accessMode neo4j.AccessMode) {
	r.AccessMode = &accessMode
}

func (r *SessionOptions) Merge(opts ...*SessionOptions) {
	for _, o := range opts {
		if o == nil {
			continue
		}
		if o.AccessMode != nil {
			r.SetAccessMode(*o.AccessMode)
		}
	}
}

func (r *SessionOptions) setDefault() {
	if r.AccessMode == nil {
		r.SetAccessMode(neo4j.AccessModeWrite)
	}
}
