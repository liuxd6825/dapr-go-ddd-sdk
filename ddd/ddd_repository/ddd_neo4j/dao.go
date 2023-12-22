package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/gocsv"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
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
	driver  neo4j.DriverWithContext
	cypher  Cypher
	newOne  func() T
	newList func() []T
}

type Options[T interface{}] struct {
	newOne  func() T
	newList func() []T
}

type ImportCsvCmd struct {
	TenantId         string                 `json:"tenantId" desc:"租户ID"`
	CaseId           string                 `json:"caseId" desc:""`
	ImportFile       string                 `json:"importFile"`
	Labels           []string               `json:"label"`
	Fields           []string               `json:"fields"`
	ImportType       ImportType             `json:"importType"`
	Data             ImportCsvCmdData       `json:"data"`
	SaveFileCallback ImportSaveFileCallback `json:"-"`
}

type ImportCsvCmdData interface {
	Data() any
	List() any
	Item(index int) any
	Append(item any)
	Length() int
}

type ImportType int

type ImportJsonCmd struct {
	TenantId   string     `json:"tenantId" desc:"租户ID"`
	CaseId     string     `json:"caseId" desc:""`
	Neo4jPath  string     `json:"neo4JPath"`
	ImportFile string     `json:"importFile"`
	Nodes      []Node     `json:"nodes"`
	Relations  []Relation `json:"relations"`
}

type importCsvCmdData struct {
	list []any
	data any
}

func (i *importCsvCmdData) List() any {
	return i.list
}

func (i *importCsvCmdData) Data() any {
	return i.data
}

func (i *importCsvCmdData) Item(index int) any {
	return i.list[index]
}

func (i *importCsvCmdData) Append(item any) {
	i.list = append(i.list, item)
}

func (i *importCsvCmdData) Length() int {
	return len(i.list)
}

func NewImportCsvCmdData(data any) ImportCsvCmdData {
	return &importCsvCmdData{data: data}
}

const (
	ImportTypeNode = iota
	ImportTypeRelation
)

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

func (d *Dao[T]) init(driver neo4j.DriverWithContext, cypher Cypher, opts ...*Options[T]) {
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

func (d *Dao[T]) query(ctx context.Context, query string, data map[string]any) (any, error) {
	result, err := neo4j.ExecuteQuery(ctx, d.driver, query, data, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
	return result, err
}

func (d *Dao[T]) doSession(ctx context.Context, fun func(tx neo4j.ManagedTransaction) (*Neo4jResult, error), opts ...*SessionOptions) (result *Neo4jResult, err error) {
	if fun == nil {
		return nil, errors.New("doSession(ctx, fun) fun is nil")
	}
	if sc, ok := GetSessionContext(ctx); ok {
		tx := sc.GetTransaction()
		_, err := fun(tx)
		return nil, err
	}

	opt := NewSessionOptions(opts...)
	opt.setDefault()

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: *opt.AccessMode})
	defer func() {
		_ = session.Close(ctx)
		if e1 := errors.GetError(recover()); e1 != nil {
			err = e1
		}
	}()
	/*
		ex, err := session.BeginTransaction(ctx, func(config *neo4j.TransactionConfig) {
			config.Timeout = 50 * time.Second
		})
		if err != nil {
			return nil, err
		}
	*/

	var res any
	if *opt.AccessMode == neo4j.AccessModeRead {
		res, err = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			return fun(tx)
		})
	} else if *opt.AccessMode == neo4j.AccessModeWrite {
		res, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
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

func (d *Dao[T]) Write(ctx context.Context, cypher string) (*Neo4jResult, error) {
	return d.doSession(ctx, func(tx neo4j.ManagedTransaction) (*Neo4jResult, error) {
		result, err := tx.Run(ctx, cypher, nil)
		if err != nil {
			return nil, err
		}
		return NewNeo4jResult(ctx, result), err
	})
}

func (d *Dao[T]) Query(ctx context.Context, cypher string, params map[string]interface{}) (*Neo4jResult, error) {
	var resultData *Neo4jResult
	_, err := d.doSession(ctx, func(tx neo4j.ManagedTransaction) (*Neo4jResult, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			log.Println("wirte to DB with error:", err)
			return nil, err
		}
		resultData = NewNeo4jResult(ctx, result)
		return nil, err
	})
	return resultData, err
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
	return d.Run(ctx, cypher, params, true, opts...)
}

func (d *Dao[T]) Run(ctx context.Context, cypher string, params map[string]any, isWriteMode bool, opts ...ddd_repository.Options) (*Neo4jResult, error) {
	sOptionsBuilder := NewSessionOptionsBuilder().SetAccessMode(neo4j.AccessModeRead)
	if isWriteMode {
		sOptionsBuilder.SetAccessMode(neo4j.AccessModeWrite)
	}

	res, err := d.doSession(ctx, func(tx neo4j.ManagedTransaction) (*Neo4jResult, error) {
		r, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return NewNeo4jResult(ctx, r), nil
	}, sOptionsBuilder.Build())

	return res, err
}

func (d *Dao[T]) ImportFile2(ctx context.Context, cmd ImportCsvCmd, data []T, opts ...ddd_repository.Options) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	labels := ""
	for _, item := range cmd.Labels {
		labels += ":" + item
	}

	field := strings.Builder{}
	fs := []string{
		"id", "caseId", "tenantId", "graphId", "nodeType", "isDeleted", "remarks",
		"createTime", "creatorId", "creatorName", "deletedTime", "deleterId", "deleterName",
		"updatedTime", "updaterId", "updaterName",
	}
	fs = append(fs, cmd.Fields...)
	fsLen := len(fs) - 1
	for i, name := range fs {
		field.WriteString(name + ":a." + stringutils.FirstUpper(name))
		if i < fsLen {
			field.WriteString(",")
		}
	}

	cypher := strings.Builder{}
	//cypher.WriteString(fmt.Sprintf(`USING PERIODIC COMMIT 500;"`))
	switch cmd.ImportType {
	case ImportTypeNode:
		batch, err := jsonutils.Marshal(data)
		if err != nil {
			return err
		}
		cypher.WriteString(fmt.Sprintf(`
		UNWIND {batch:[%s]} as a
		CREATE (n%s{%s})
        `, batch, labels, field.String()))
		break
	case ImportTypeRelation:
		var list []any
		for _, item := range data {
			var obj any = item
			rel, ok := obj.(Relation)
			if ok {
				rel.SetProperties(map[string]any{})
				d := struct {
					Id      string `json:"id"`
					StartId string `json:"startId"`
					EndId   string `json:"endId"`
					RelType string `json:"relType"`
					Props   any    `json:"props"`
				}{
					Id:      rel.GetId(),
					StartId: rel.GetStartId(),
					EndId:   rel.GetEndId(),
					RelType: rel.GetRelType(),
					Props:   nil,
				}

				list = append(list, d)
			}
		}
		batch, err := jsonutils.Marshal(list)
		if err != nil {
			return err
		}

		cypher.WriteString(fmt.Sprintf(`
		UNWIND {batch:%s} as row
		MATCH (from%s{id:row.startId})
		MATCH (to%s{id:row.endId})
		CALL apoc.create.relationship(from, row.relType, row.props, to) yield rel RETURN count(*)
		`, batch, labels, labels))
		break
	}
	//cypher.WriteString(fmt.Sprintf(`"),{batchSize:1000, parallel:true})`))

	fmt.Println("***********")
	logs.Debug(ctx, cypher.String())
	fmt.Println("***********")

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher.String(), nil)
	if err == nil {
		summary, _ := result.Consume(ctx)
		fmt.Println("Query updated the database?", summary.Counters().ContainsUpdates())
	}

	return err
}

func (d *Dao[T]) CreateIndex(ctx context.Context, index, label, property string) (err error) {
	cypher := fmt.Sprintf("CREATE INDEX %s IF NOT EXISTS FOR (n:%s) ON (n.%s) ", index, label, property)
	idxSession := d.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

	_, err = idxSession.Run(ctx, cypher, nil)
	defer func() {
		_ = idxSession.Close(ctx)
	}()
	return err
}

// SaveFileFun
// fileName 文件名称
// data 保存的数据
// return string 存储文件的URI，可以是file:///或http://
type ImportSaveFileCallback func(ctx context.Context, tenantId, fileName string, data any) (string, ImportSaveCompleteCallback, error)
type ImportSaveCompleteCallback func() error

func (d *Dao[T]) ImportCsv(ctx context.Context, cmd ImportCsvCmd, opts ...ddd_repository.Options) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	var complete ImportSaveCompleteCallback
	defer func() {
		if complete != nil {
			//三次重试
			for i := 1; i < 4; i++ {
				if err := complete(); err != nil {
					time.Sleep(time.Duration(3*i) * time.Second)
					continue
				}
			}
		}
	}()

	fileUri := ""
	if cmd.SaveFileCallback == nil {
		neo4jPath, err := restapp.GetConfigAppValue("neo4jPath")
		if err != nil {
			return err
		}
		fileUri = fmt.Sprintf("file:///%s", cmd.ImportFile)
		fileName := fmt.Sprintf("file:///%s/import/%s", neo4jPath, cmd.ImportFile)
		var csvFile *os.File
		csvFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}

		defer func() {
			_ = csvFile.Close()
			os.Remove(fileName)
		}()

		if err = gocsv.MarshalFile(cmd.Data.Data(), csvFile); err != nil {
			return err
		}
	} else {
		fileUri, complete, err = cmd.SaveFileCallback(ctx, cmd.TenantId, cmd.ImportFile, cmd.Data.Data())
		if err != nil {
			return err
		}
	}

	labels := ""
	for _, item := range cmd.Labels {
		labels += ":" + item
	}

	fs := []string{
		"id", "caseId", "tenantId", "graphId", "nodeType", "isDeleted", "remarks",
		"createTime", "creatorId", "creatorName", "deletedTime", "deleterId", "deleterName",
		"updatedTime", "updaterId", "updaterName",
	}
	fs = append(fs, cmd.Fields...)
	fsLen := len(fs) - 1
	props := strings.Builder{}
	setFields := strings.Builder{}
	for i, name := range fs {
		props.WriteString(name + ":a." + stringutils.FirstUpper(name))
		setFields.WriteString("n." + name + "=a." + stringutils.FirstUpper(name))
		if i < fsLen {
			props.WriteString(",")
			setFields.WriteString(",")
		}
	}

	//创建索引
	_ = d.CreateIndex(ctx, "case_"+cmd.CaseId+"_id", "case_"+cmd.CaseId, "id")

	cypher := strings.Builder{}

	switch cmd.ImportType {
	case ImportTypeNode:
		cypher.WriteString(fmt.Sprintf(`
		LOAD CSV WITH HEADERS FROM '%s' AS a
		CALL {
			WITH a
			MERGE (n%s{id:a.Id}) ON CREATE SET %v ON MATCH SET %v 
		}  IN TRANSACTIONS OF 100 ROWS;
        `, fileUri, labels, setFields.String(), setFields.String()))
		break

	case ImportTypeRelation:
		relTypes := make(map[string]int) // 存储去重后的关系
		cypher.WriteString(fmt.Sprintf(`
		LOAD CSV WITH HEADERS FROM '%s' AS a
		CALL {
			WITH a
			MATCH (s%s{id:a.StartId}),(e%s{id:a.EndId}) 
		`, fileUri, labels, labels))

		// 去重后的关系
		for i := 0; i < cmd.Data.Length(); i++ {
			item := cmd.Data.Item(i)
			if rel, ok := item.(Relation); ok {
				relType := rel.GetRelType()
				if _, find := relTypes[relType]; !find {
					relTypes[relType] = 0
					// 按关系类型生成创建关系语句
					cypher.WriteString(fmt.Sprintf(`
					FOREACH (_ IN case when a.RelType = '%s' then[1] else[] end|
						MERGE (s)-[n:%s{id:a.Id}]->(e)
						ON CREATE SET %s
						ON MATCH  SET %s
					)`, relType, relType, setFields.String(), setFields.String()))
				}
			}
		}
		cypher.WriteString(`} IN TRANSACTIONS OF 100 ROWS;`)
		break
	}

	logs.Debugf(ctx, `
	******
	%s
	******
	`, func() any {
		return cypher.String()
	})

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		_ = session.Close(ctx)
	}()

	_, err = session.Run(ctx, cypher.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

type ImportJsonRelation struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Label      string                  `json:"label"`
	Properties any                     `json:"properties"`
	Start      ImportJsonRelationStart `json:"start"`
	End        ImportJsonRelationEnd   `json:"end"`
}
type ImportJsonRelationStart struct {
	Id         string   `json:"id"`
	Labels     []string `json:"labels"`
	Properties any      `json:"properties"`
}
type ImportJsonRelationEnd struct {
	Id         string   `json:"id"`
	Labels     []string `json:"labels"`
	Properties any      `json:"properties"`
}

type ImportJsonNode struct {
	Id         string   `json:"id"`
	Type       string   `json:"type"`
	Labels     []string `json:"labels"`
	Properties any      `json:"properties"`
}
type Null struct {
}

func (d *Dao[T]) ImportJson(ctx context.Context, cmd ImportJsonCmd, opts ...ddd_repository.Options) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	fileName := cmd.Neo4jPath + "/import/" + cmd.ImportFile

	var jsonFile *os.File
	jsonFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer func() {
		_ = jsonFile.Close()
	}()

	//labelTenant := fmt.Sprintf("tenant_%s", cmd.TenantId)
	labelCase := fmt.Sprintf("case_%s", cmd.CaseId)
	lables := []string{
		//labelTenant,
		labelCase,
	}
	ids := map[string]int{}
	for _, item := range cmd.Nodes {
		id := item.GetId()
		//nodeLables := append(lables, "human")
		if _, ok := ids[id]; ok {
			continue
		}
		ids[id] = 0
		props := map[string]any{"id": item.GetId()}
		node := ImportJsonNode{
			Id:         item.GetId(),
			Type:       "node",
			Labels:     lables,
			Properties: props,
		}
		if item, err := jsonutils.Marshal(node); err != nil {
			return err
		} else {
			jsonFile.WriteString(item)
			jsonFile.WriteString("\r\n")
		}
	}

	for _, item := range cmd.Relations {
		rel := ImportJsonRelation{
			Id:         item.GetId(),
			Type:       "relationship",
			Label:      item.GetRelType(),
			Properties: item.GetProperties(),
			Start: ImportJsonRelationStart{
				Id:         item.GetStartId(),
				Labels:     lables,
				Properties: Null{},
			},
			End: ImportJsonRelationEnd{
				Id:         item.GetEndId(),
				Labels:     lables,
				Properties: Null{},
			},
		}
		if item, err := jsonutils.Marshal(rel); err != nil {
			return err
		} else {
			jsonFile.WriteString(item)
			jsonFile.WriteString("\r\n")
		}
	}

	cypher := fmt.Sprintf(`CALL apoc.import.json("file:///%s",{cleanup:false, importIdName:"id"} )`, cmd.ImportFile)

	fmt.Println("***********")
	logs.Debug(ctx, cypher)
	fmt.Println("***********")

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err = session.Run(ctx, cypher, nil)
	if err != nil {
		fmt.Println(err)

	} else {
		//summary, _ := result.Consume(ctx)
		//fmt.Println("Query updated the database?", summary.Counters().ContainsUpdates())
	}

	return err
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
