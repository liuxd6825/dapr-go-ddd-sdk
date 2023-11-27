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

/*type Dao[T any] interface {
	Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T])
	InsertMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T]

	Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T]
	UpdateMany(ctx context.Context, list []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T]

	DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) error
	DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error
	DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error
	DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...ddd_repository.Options) error
	DeleteByGraphId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) error

	FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T]
	FindByFilter(ctx context.Context, tenantId, filter string) *ddd_repository.FindListResult[T]
	FindById(ctx context.Context, tenantId, id string, opts ...ddd_repository.Options) (T, bool, error)
	FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) ([]T, bool, error)
	FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T]
	FindByGraphId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T]
}*/

type Element interface {
	GetTenantId() string
	SetTenantId(v string)
	GetId() string
	SetId(v string)
}

type Dao[T Element] struct {
	driver  neo4j.Driver
	cypher  Cypher
	newOne  func() T
	newList func() []T
}

type Options[T interface{}] struct {
	newOne  func() T
	newList func() []T
}

func NewOptions[T interface{}](opts ...*Options[T]) *Options[T] {
	n := &Options[T]{}
	for _, o := range opts {
		if o.newList != nil {
			n.newList = o.newList
		}
		if o.newOne != nil {
			n.newOne = o.newOne
		}
	}
	return n
}

func (d *Dao[T]) init(driver neo4j.Driver, cypher Cypher, opts ...*Options[T]) {
	o := NewOptions[T](opts...)
	d.driver = driver
	d.cypher = cypher
	if o.newList != nil {
		d.newList = o.newList
	}
	if o.newOne != nil {
		d.newOne = o.newOne
	}
}

func (d *Dao[T]) NewEntity() (res T, resErr error) {
	if d.newOne != nil {
		return d.newOne(), nil
	}
	return reflectutils.NewStruct[T]()
}

func (d *Dao[T]) Save(ctx context.Context, data *ddd.SetData[T], opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			setResult = ddd_repository.NewSetResultError[T](err)
		}
	}()

	for _, item := range data.Items() {
		statue := item.Statue()
		entity := item.Data().(T)
		switch statue {
		case ddd.DataStatueCreate:
			err = d.Insert(ctx, entity, opts...).GetError()
		case ddd.DataStatueUpdate:
			err = d.Update(ctx, entity, opts...).GetError()
		case ddd.DataStatueDelete:
			err = d.DeleteById(ctx, entity.GetTenantId(), entity.GetId(), opts...)
		case ddd.DataStatueCreateOrUpdate:
			err = d.InsertOrUpdate(ctx, entity, opts...).GetError()
		}
		if err != nil {
			return ddd_repository.NewSetResultError[T](err)
		}
	}
	return ddd_repository.NewSetResultError[T](nil)
}

func (d *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			if err = errors.GetRecoverError(err, e); err != nil {
				setResult = ddd_repository.NewSetResultError[T](err)
			}
		}
	}()

	cr, err1 := d.cypher.Insert(ctx, entity)
	if err1 != nil {
		err = err1
		return ddd_repository.NewSetResultError[T](err)
	}

	_, err = d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...)
	if err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetResult(entity, err)
}

func (d *Dao[T]) InsertMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	for _, e := range entities {
		if err := d.Insert(ctx, e, opts...).GetError(); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}
	return ddd_repository.NewSetManyResult[T](entities, nil)
}

func (d *Dao[T]) InsertOrUpdate(ctx context.Context, entity T, opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
	var err error
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			setResult = ddd_repository.NewSetResultError[T](err)
		}
	}()

	cr, err1 := d.cypher.InsertOrUpdate(ctx, entity)
	if err1 != nil {
		err = err1
		return ddd_repository.NewSetResultError[T](err)
	}

	_, err = d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...)
	if err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetResult(entity, err)
}

func (d *Dao[T]) InsertOrUpdateMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	for _, e := range entities {
		if err := d.InsertOrUpdate(ctx, e, opts...).GetError(); err != nil {
			return ddd_repository.NewSetManyResultError[T](err)
		}
	}
	return ddd_repository.NewSetManyResult[T](entities, nil)
}

func (d *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	cr, err := d.cypher.Update(ctx, entity)
	res, err := d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...)
	if err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	if _, err := res.GetOne("", entity); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetResult(entity, err)
}

func (d *Dao[T]) UpdateMany(ctx context.Context, list []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	for _, entity := range list {
		if cr, err := d.cypher.Update(ctx, entity); err != nil {
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

func (d *Dao[T]) UpdateLabelById(ctx context.Context, tenantId string, id string, label string) error {
	cr, err := d.cypher.UpdateLabelById(ctx, tenantId, id, label)
	if err != nil {
		return err
	}
	if cr == nil {
		return nil
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) UpdateLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) error {
	cr, err := d.cypher.UpdateLabelByFilter(ctx, tenantId, filter, labels...)
	if err != nil {
		return err
	}
	if cr == nil {
		return nil
	}

	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) DeleteLabelById(ctx context.Context, tenantId string, id string, label string) error {
	cr, err := d.cypher.DeleteLabelById(ctx, tenantId, id, label)
	if err != nil || cr == nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) DeleteLabelByFilter(ctx context.Context, tenantId string, filter string, labels ...string) error {
	cr, err := d.cypher.DeleteLabelByFilter(ctx, tenantId, filter, labels...)
	if err != nil || cr == nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) error {
	cr, err := d.cypher.DeleteById(ctx, tenantId, id)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error {
	cr, err := d.cypher.DeleteByIds(ctx, tenantId, ids)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error {
	cr, err := d.cypher.DeleteAll(ctx, tenantId)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Dao[T]) DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...ddd_repository.Options) error {
	cr, err := d.cypher.DeleteByFilter(ctx, tenantId, filter)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

func (d *Dao[T]) DeleteByGraphId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) error {
	return d.DeleteByFilter(ctx, tenantId, fmt.Sprintf("graphId=='%v'", graphId))
}

func (d *Dao[T]) DeleteByCaseId(ctx context.Context, tenantId string, caseId string, opts ...ddd_repository.Options) error {
	return d.DeleteByFilter(ctx, tenantId, fmt.Sprintf("caseId=='%v'", caseId))
}

func (d *Dao[T]) DeleteByTenantId(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error {
	cr, err := d.cypher.DeleteByTenantId(ctx, tenantId)
	if err != nil {
		return err
	}
	_, err = d.doSet(ctx, tenantId, cr.Cypher(), cr.Params(), opts...)
	return err
}

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

func (d *Dao[T]) Write(ctx context.Context, cypher string) (*Neo4jResult, error) {
	return d.doSession(ctx, func(tx neo4j.Transaction) (*Neo4jResult, error) {
		result, err := tx.Run(cypher, nil)
		if err != nil {
			return nil, err
		}
		return NewNeo4jResult(result), err
	})
}

func (d *Dao[T]) Query(ctx context.Context, cypher string, params map[string]interface{}) (*Neo4jResult, error) {
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

	var totalRows *int64
	if isTotalRows {
		totalKey := "count"
		countCypher := cypher + fmt.Sprintf(" RETURN count(%s) as %s ", resultKey, totalKey)
		result, err := d.Query(ctx, countCypher, params)
		total, err := result.GetInteger(totalKey, 0)
		if err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err), false, err
		}
		totalRows = &total
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

func (d *Dao[T]) DoFilter(ctx context.Context, tenantId string, fun func() (*ddd_repository.FindPagingResult[T], bool, error), opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	data, _, err := fun()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	return data
}

func (d *Dao[T]) DoList(ctx context.Context, tenantId string, fun func() (*ddd_repository.FindListResult[T], bool, error), opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	data, _, err := fun()
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	return data
}

func (d *Dao[T]) newSetManyResult(result *Neo4jResult, err error) *ddd_repository.SetManyResult[T] {
	if err != nil {
		return ddd_repository.NewSetManyResultError[T](err)
	}
	var data []T
	if err := result.GetList("n", &data); err != nil {
		ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetManyResult[T](data, err)
}

func (d *Dao[T]) doSet(ctx context.Context, tenantId string, cypher string, params map[string]interface{}, opts ...ddd_repository.Options) (*Neo4jResult, error) {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return nil, err
	}
	res, err := d.doSession(ctx, func(tx neo4j.Transaction) (*Neo4jResult, error) {
		r, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return NewNeo4jResult(r), nil
	})
	return res, err
}

func (d *Dao[T]) doSession(ctx context.Context, fun func(tx neo4j.Transaction) (*Neo4jResult, error), opts ...*SessionOptions) (result *Neo4jResult, err error) {

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
		if e1 := errors.GetError(recover()); e1 != nil {
			err = e1
		}
	}()

	var res interface{}
	if *opt.AccessMode == neo4j.AccessModeRead {
		res, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return fun(tx)
		})
	} else if *opt.AccessMode == neo4j.AccessModeWrite {
		res, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return fun(tx)
		})
	}
	if err != nil {
		return nil, err
	}

	if result, ok := res.(*Neo4jResult); ok {
		return result, nil
	}

	return nil, err
}

func getLabels(labels ...string) string {
	var s string
	for _, l := range labels {
		if len(l) > 0 {
			s = fmt.Sprintf("%v :`%v`", s, l)
		}
	}
	return strings.ToLower(s)
}
