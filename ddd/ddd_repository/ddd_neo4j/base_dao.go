package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Entity interface {
	ddd.Entity
	GetTenantId() string
}

type BaseDao[T Entity] struct {
	driver neo4j.Driver
	cypher Cypher
}

func (d *BaseDao[T]) init(driver neo4j.Driver, cypher Cypher) {
	d.driver = driver
	d.cypher = cypher
}

func (d *BaseDao[T]) NewEntity() (res T, resErr error) {
	return reflectutils.NewStruct[T]()
}

func (d *BaseDao[T]) NewEntities() (res []T, resErr error) {
	return reflectutils.NewSlice[[]T]()
}

func (d *BaseDao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T]) {
	defer func() {
		if e := recover(); e != nil {
			if err := errors.GetRecoverError(e); err != nil {
				setResult = ddd_repository.NewSetResultError[T](err)
			}
		}
	}()

	cr, err := d.cypher.Insert(ctx, entity)
	if err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}

	res, err := d.doSet(ctx, entity.GetTenantId(), cr.Cypher(), cr.Params(), opts...)
	if _, err := res.GetOne("", entity); err != nil {
		return ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetResult(entity, err)
}

func (d *BaseDao[T]) doSet(ctx context.Context, tenantId string, cypher string, params map[string]interface{}, opts ...ddd_repository.Options) (*Neo4jResult, error) {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return nil, err
	}
	res, err := d.doSession(ctx, func(tx neo4j.Transaction) (*Neo4jResult, error) {
		r, err := tx.Run(cypher, params)
		return NewNeo4jResult(r), err
	})
	return res, err
}

func (d *BaseDao[T]) doSession(ctx context.Context, fun func(tx neo4j.Transaction) (*Neo4jResult, error), opts ...*SessionOptions) (*Neo4jResult, error) {
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
