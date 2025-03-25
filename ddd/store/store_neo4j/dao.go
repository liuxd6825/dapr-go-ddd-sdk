package store_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql/rsql_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/errors/assert"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"strings"
)

type Dao[T any] struct {
	Driver neo4j.DriverWithContext
	Cypher Cypher[T]
	eb     store.EntityBuilder[T]
	Labels []string
	Schema *store.DBSchema
}

func (d *Dao[T]) NewEntity() T {
	return d.eb.NewEntity()
}

func (d *Dao[T]) NewEntityList() []T {
	return d.eb.NewEntityList()
}

func (d *Dao[T]) GetTenantId(entity T) string {
	return d.eb.GetTenantId(entity)
}

func (d *Dao[T]) SetTenantId(entity T, tenantId string) {
	d.eb.SetTenantId(entity, tenantId)
}

func (d *Dao[T]) GetId(entity T) string {
	return d.eb.GetId(entity)
}

func (d *Dao[T]) SetId(entity T, id string) {
	d.eb.SetId(entity, id)
}

func (d *Dao[T]) GetAggId(entity T) string {
	return d.eb.GetAggId(entity)
}

func (d *Dao[T]) GetSchema() *store.DBSchema {
	return d.Schema
}

func NewNodeDao[T any](driver neo4j.DriverWithContext, dbSch *store.DBSchema, labels []string, opts ...*Options[T]) store.IStore[T] {
	eb := NewNodeEntityBuilder[T](dbSch)
	cypher := NewNodeCypher[T](eb, dbSch, labels...)
	return NewDao(driver, labels, cypher, eb, dbSch, opts...)
}

func NewRelationDao[T any](driver neo4j.DriverWithContext, dbSch *store.DBSchema, labels []string, opts ...*Options[T]) store.IStore[T] {
	eb := NewRelationEntityBuilder[T](dbSch)
	cypher := NewRelationCypher[T](eb, dbSch, labels...)
	return NewDao(driver, labels, cypher, eb, dbSch, opts...)
}

func NewDao[T any](driver neo4j.DriverWithContext, labels []string, cypher Cypher[T], eb store.EntityBuilder[T], dbSch *store.DBSchema, opts ...*Options[T]) store.IStore[T] {
	return newDao(driver, labels, cypher, eb, dbSch, opts...)
}

func newNodeDao[T any](driver neo4j.DriverWithContext, labels []string, eb NodeEntityBuilder[T], dbSch *store.DBSchema, opts ...*Options[T]) *Dao[T] {
	cypher := NewNodeCypher(eb, dbSch, labels...)
	return newDao(driver, labels, cypher, eb, dbSch, opts...)
}

func newDao[T any](driver neo4j.DriverWithContext, labels []string, cypher Cypher[T], eb store.EntityBuilder[T], dbSch *store.DBSchema, opts ...*Options[T]) *Dao[T] {
	dao := &Dao[T]{
		Driver: driver,
		Cypher: cypher,
		eb:     eb,
		Labels: labels,
		Schema: dbSch,
	}
	return dao
}

func (d *Dao[T]) StartTx(ctx context.Context, fun store.TxFunc, options ...*store.SessionOptions) error {
	return nil
}

func (d *Dao[T]) DoFilter(ctx context.Context, tenantId string, fun func() (*store.FindPagingResult[T], bool, error), opts ...store.Options) *store.FindPagingResult[T] {
	data, _, err := fun()
	if err != nil {
		return store.NewFindPagingResultWithError[T](err)
	}
	return data
}

func (d *Dao[T]) DoList(ctx context.Context, tenantId string, fun func() *store.FindListResult[T], opts ...store.Options) *store.FindListResult[T] {
	data := fun()
	return data
}

func (d *Dao[T]) newSetManyResult(ctx context.Context, result *Neo4jResult, err error) *store.SetResult[T] {
	if err != nil {
		return store.NewSetResultError[T](err)
	}
	var data []T
	if err := result.GetList(ctx, "n", &data); err != nil {
		store.NewSetResultError[T](err)
	}
	return store.NewSetResultEmpty[T]()
}

func (d *Dao[T]) doSet(ctx context.Context, tenantId string, cypher string, params map[string]interface{}, opts ...store.Options) (*Neo4jResult, error) {
	if err := assert2.NotEmpty(tenantId, assert2.NewOptions("tenantId is empty")); err != nil {
		return nil, err
	}
	return d.Run(ctx, cypher, params, true, opts...)
}

func (d *Dao[T]) Run(ctx context.Context, cypher string, params map[string]any, isWriteMode bool, opts ...store.Options) (*Neo4jResult, error) {
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

func (d *Dao[T]) CreateIndex(ctx context.Context, index, label, property string) (err error) {
	cypher := fmt.Sprintf("CREATE INDEX %s IF NOT EXISTS FOR (n:%s) ON (n.%s) ", index, label, property)
	idxSession := d.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

	_, err = idxSession.Run(ctx, cypher, nil)
	defer func() {
		_ = idxSession.Close(ctx)
	}()
	return err
}

func (d *Dao[T]) query(ctx context.Context, query string, data map[string]any) (any, error) {
	result, err := neo4j.ExecuteQuery(ctx, d.Driver, query, data, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
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

	session := d.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: *opt.AccessMode})
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

func getLabels(labels ...string) string {
	var s string
	for _, l := range labels {
		if len(l) > 0 {
			s = fmt.Sprintf("%v:`%v`", s, l)
		}
	}
	return strings.ToLower(s)
}

func getNeo4jWhere(tenantId string, dataKey string, filter string) (string, error) {
	process := rsql_neo4j.NewProcess(tenantId, dataKey)
	if err := rsql.ParseProcess(filter, process); err != nil {
		return "", err
	}
	where := process.GetSQL()
	if len(where) > 0 {
		where = "WHERE " + where
	}
	return where, nil
}
