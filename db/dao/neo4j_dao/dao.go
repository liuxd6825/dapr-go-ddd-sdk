package neo4j_dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Dao[T ddd_neo4j.Element] struct {
	dao *ddd_neo4j.Dao[T]
}

type RepositoryOptions struct {
	driver neo4j.DriverWithContext
}

var _driver neo4j.DriverWithContext

func NewRepositoryOptions() *RepositoryOptions {
	return &RepositoryOptions{}
}

func NewSession(isWrite bool) ddd_repository.Session {
	return ddd_neo4j.NewSession(isWrite, GetDB())
}

// NewNodeDao 创建节点DAO
func NewNodeDao[T ddd_neo4j.Node](labels []string, opts ...*RepositoryOptions) *Dao[T] {
	options := NewRepositoryOptions()
	options.Merge(opts...)
	return &Dao[T]{
		dao: ddd_neo4j.NewNodeDao[T](options.driver, ddd_neo4j.NewNodeCypher(labels...)),
	}
}

// NewRelationDao 创建关系DAO
func NewRelationDao[T ddd_neo4j.Relation](labels string, opts ...*RepositoryOptions) *Dao[T] {
	options := NewRepositoryOptions()
	options.Merge(opts...)
	return &Dao[T]{
		dao: ddd_neo4j.NewRelationDao[T](options.driver, ddd_neo4j.NewRelationCypher(labels)),
	}
}

func (u *Dao[T]) Save(ctx context.Context, data *ddd.SetData[T], opts ...ddd_repository.Options) error {
	return u.dao.Save(ctx, data, opts...).GetError()
}

func (u *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) error {
	return u.dao.Insert(ctx, entity, opts...).GetError()
}

func (u *Dao[T]) InsertMany(ctx context.Context, entity []T, opts ...ddd_repository.Options) error {
	return u.dao.InsertMany(ctx, entity, opts...).GetError()
}

func (u *Dao[T]) InsertOrUpdate(ctx context.Context, entity T, opts ...ddd_repository.Options) error {
	return u.dao.InsertOrUpdate(ctx, entity, opts...).GetError()
}

func (u *Dao[T]) InsertOrUpdateMany(ctx context.Context, entity []T, opts ...ddd_repository.Options) error {
	return u.dao.InsertOrUpdateMany(ctx, entity, opts...).GetError()
}

func (u *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) error {
	return u.dao.Update(ctx, entity, opts...).GetError()
}

func (u *Dao[T]) UpdateMany(ctx context.Context, entity []T, opts ...ddd_repository.Options) error {
	return u.dao.UpdateMany(ctx, entity, opts...).GetError()
}

func (u *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) error {
	return u.dao.DeleteById(ctx, tenantId, id, opts...)
}

func (u *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error {
	return u.dao.DeleteByIds(ctx, tenantId, ids)
}

func (u *Dao[T]) DeleteByCaseId(ctx context.Context, tenantId string, caseId string, opts ...ddd_repository.Options) error {
	return u.dao.DeleteByCaseId(ctx, tenantId, caseId)
}

func (u *Dao[T]) DeleteByTenantId(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error {
	return u.dao.DeleteByTenantId(ctx, tenantId, opts...)
}

func (u *Dao[T]) DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...ddd_repository.Options) error {
	return u.dao.DeleteByFilter(ctx, tenantId, filter, opts...)
}

func (u *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error {
	return u.dao.DeleteAll(ctx, tenantId, opts...)
}

func (u *Dao[T]) FindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) (T, bool, error) {
	return u.dao.FindById(ctx, tenantId, id, opts...)
}

func (u *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) ([]T, bool, error) {
	return u.dao.FindByIds(ctx, tenantId, ids, opts...)
}

/*
	func (u *Dao[T]) FindByCaseId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) ([]T, bool, error) {
		return u.dao.FindByCaseId(ctx, tenantId, graphId, opts...).Result()
	}
*/
func (u *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return u.dao.FindAll(ctx, tenantId, opts...)
}

func (u *Dao[T]) FindPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return u.dao.FindPaging(ctx, query, opts...)
}

func (u *Dao[T]) FindPagingByCypher(ctx context.Context, tenantId, cypher string, pageNum, pageSize int64, resultKey string, isTotalRows bool, params map[string]any, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return u.dao.FindPagingByCypher(ctx, tenantId, cypher, pageNum, pageSize, resultKey, isTotalRows, params, opts...)
}

func (u *Dao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	return u.dao.FindListByMap(ctx, tenantId, filterMap, opts...)
}

func (u *Dao[T]) Query(ctx context.Context, match string, pars map[string]interface{}) (*ddd_neo4j.Neo4jResult, error) {
	return u.dao.Query(ctx, match, pars)
}

func (u *Dao[T]) Run(ctx context.Context, cypher string, params map[string]any, isWriteMode bool, opts ...ddd_repository.Options) (*ddd_neo4j.Neo4jResult, error) {
	return u.dao.Run(ctx, cypher, params, isWriteMode, opts...)
}
func (u *Dao[T]) ImportNodeCsvFile(ctx context.Context, cmd ddd_neo4j.ImportCsvCmd, opts ...ddd_repository.Options) error {
	return u.dao.ImportNodeCsv(ctx, cmd, opts...)
}

func (u *Dao[T]) ImportRelationCsvFile(ctx context.Context, cmd ddd_neo4j.ImportCsvCmd, opts ...ddd_repository.Options) error {
	return u.dao.ImportRelationCsv(ctx, cmd, opts...)
}

func (u *Dao[T]) ImportJsonFile(ctx context.Context, cmd ddd_neo4j.ImportJsonCmd, opts ...ddd_repository.Options) error {
	return u.dao.ImportJson(ctx, cmd, opts...)
}

func (u *Dao[T]) FindAutoComplete(ctx context.Context, qry ddd_repository.FindAutoCompleteQuery, opts ...ddd_repository.Options) (*ddd_repository.FindPagingResult[T], bool, error) {
	//TODO implement me
	panic("AutoComplete implement me")
}

func (u *Dao[T]) FindDistinct(ctx context.Context, qry ddd_repository.FindDistinctQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	//TODO implement me
	panic("AutoComplete implement me")
}

func (o *RepositoryOptions) SetDriver(driver neo4j.DriverWithContext) *RepositoryOptions {
	o.driver = driver
	return o
}

func (o *RepositoryOptions) Merge(opts ...*RepositoryOptions) *RepositoryOptions {
	if opts != nil {
		for _, item := range opts {
			if item.driver != nil {
				o.driver = item.driver
			}
		}
	}
	o.driver = GetDB()
	return o
}

func GetDB() neo4j.DriverWithContext {
	if _driver != nil {
		return _driver
	}
	return restapp.GetNeo4j()
}

func SetDB(driver neo4j.DriverWithContext) {
	_driver = driver
}
