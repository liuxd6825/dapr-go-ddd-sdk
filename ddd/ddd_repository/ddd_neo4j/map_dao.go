package ddd_neo4j

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type MapDao struct {
	dao ddd_repository.Dao[map[string]any]
}

func NewMapNodeDao(driver neo4j.DriverWithContext, tableName string, opts ...*Options[map[string]any]) ddd_repository.Dao[map[string]any] {
	eb := ddd.NewMapEntityBuilder()
	dao := NewNodeDao[map[string]any](driver, tableName, eb, opts...)
	return newMapDao(dao)
}

func NewMapRelationDao(driver neo4j.DriverWithContext, tableName string, opts ...*Options[map[string]any]) ddd_repository.Dao[map[string]any] {
	eb := ddd.NewMapEntityBuilder()
	dao := NewRelationDao[map[string]any](driver, tableName, eb, opts...)
	return newMapDao(dao)
}

func newMapDao(dao ddd_repository.Dao[map[string]any]) ddd_repository.Dao[map[string]any] {
	return &MapDao{
		dao: dao,
	}
}

func (m *MapDao) NewEntity() map[string]any {
	return m.dao.NewEntity()
}

func (m *MapDao) NewEntityList() []map[string]any {
	return m.dao.NewEntityList()
}

func (m *MapDao) GetTenantId(entity map[string]any) string {
	return m.dao.GetTenantId(entity)
}

func (m *MapDao) SetTenantId(entity map[string]any, tenantId string) {
	m.dao.SetTenantId(entity, tenantId)
}

func (m *MapDao) GetId(entity map[string]any) string {
	return m.dao.GetId(entity)
}

func (m *MapDao) SetId(entity map[string]any, id string) {
	m.dao.SetId(entity, id)
}

func (m *MapDao) Insert(ctx context.Context, entity map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.Insert(ctx, entity, opts...)
}

func (m *MapDao) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...ddd_repository.Options) (res *ddd_repository.SetResult[map[string]any]) {
	return m.dao.InsertMap(ctx, tenantId, data, opts...)
}

func (m *MapDao) InsertMany(ctx context.Context, entities []map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[map[string]any] {
	return m.dao.InsertMany(ctx, entities, opts...)
}

func (m *MapDao) Update(ctx context.Context, entity map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.Update(ctx, entity, opts...)
}

func (m *MapDao) UpdateByRSQL(ctx context.Context, tenantId, filterRSQL string, data map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetManyCountResult {
	return m.dao.UpdateByRSQL(ctx, tenantId, filterRSQL, data, opts...)
}

func (m *MapDao) UpdateMany(ctx context.Context, entities []map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[map[string]any] {
	return m.dao.UpdateMany(ctx, entities, opts...)
}

func (m *MapDao) UpdateMap(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.UpdateMap(ctx, tenantId, id, data, opts...)
}

func (m *MapDao) UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.UpdateMapAndGetCount(ctx, tenantId, filter, data, opts...)
}

func (m *MapDao) Delete(ctx context.Context, entity map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.Delete(ctx, entity, opts...)
}

func (m *MapDao) DeleteByRSQL(ctx context.Context, tenantId, filter string, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.DeleteByRSQL(ctx, tenantId, filter, opts...)
}

func (m *MapDao) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.DeleteById(ctx, tenantId, id, opts...)
}

func (m *MapDao) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.DeleteByIds(ctx, tenantId, ids, opts...)
}

func (m *MapDao) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.DeleteAll(ctx, tenantId, opts...)
}

func (m *MapDao) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[map[string]any] {
	return m.dao.DeleteByMap(ctx, tenantId, filterMap, opts...)
}

func (m *MapDao) FindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[map[string]any] {
	return m.dao.FindById(ctx, tenantId, id, opts...)
}

func (m *MapDao) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[map[string]any] {
	return m.dao.FindByIds(ctx, tenantId, ids, opts...)
}

func (m *MapDao) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[map[string]any] {
	return m.dao.FindOneByMap(ctx, tenantId, filterMap, opts...)
}

func (m *MapDao) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.FindListResult[map[string]any] {
	return m.dao.FindListByMap(ctx, tenantId, filterMap, opts...)
}

func (m *MapDao) FindByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[map[string]any] {
	return m.dao.FindByRSQL(ctx, tenantId, rsql, opts...)
}

func (m *MapDao) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[map[string]any] {
	return m.dao.FindAll(ctx, tenantId, opts...)
}

func (m *MapDao) FindPaging(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) (result *ddd_repository.FindPagingResult[map[string]any]) {
	return m.dao.FindPaging(ctx, qry, opts...)
}

func (m *MapDao) FindAutoComplete(ctx context.Context, qry ddd_repository.FindAutoCompleteQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[map[string]any] {
	return m.dao.FindAutoComplete(ctx, qry, opts...)
}

func (m *MapDao) FindDistinct(ctx context.Context, qry ddd_repository.FindDistinctQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[map[string]any] {
	return m.dao.FindDistinct(ctx, qry, opts...)
}

func (m *MapDao) FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) (map[string]any, error) {
	return m.dao.FindOneAndUpdateById(ctx, tenantId, id, data, opts...)
}

func (m *MapDao) SumEntity(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error) {
	return m.dao.SumEntity(ctx, qry, opts...)
}

func (m *MapDao) SumMap(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error) {
	return m.dao.SumMap(ctx, qry, opts...)
}

func (m *MapDao) Sum(ctx context.Context, qry ddd_repository.FindPagingQuery, resData any, opts ...ddd_repository.Options) (any, bool, error) {
	return m.dao.Sum(ctx, qry, resData, opts...)
}

func (m *MapDao) SumByRSQL(ctx context.Context, tenantId, rSql string, valueCols []*ddd_repository.ValueCol, opts ...ddd_repository.Options) map[string]any {
	return m.dao.SumByRSQL(ctx, tenantId, rSql, valueCols, opts...)
}

func (m *MapDao) CountByMap(ctx context.Context, tenantId string, filterData any, opts ...ddd_repository.Options) (int64, error) {
	return m.dao.CountByMap(ctx, tenantId, filterData, opts...)
}

func (m *MapDao) CountByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) (int64, error) {
	return m.dao.CountByRSQL(ctx, tenantId, rsql, opts...)
}

func (m *MapDao) StartTx(ctx context.Context, fun ddd_repository.TxFunc, options ...*ddd_repository.SessionOptions) error {
	return m.dao.StartTx(ctx, fun, options...)
}

func (m *MapDao) SetMetadata(metadata map[string]any) {
	m.dao.SetMetadata(metadata)
}

func (m *MapDao) GetMetadata() map[string]any {
	return m.dao.GetMetadata()
}

func (m *MapDao) AddMetadata(key string, val any) {
	m.dao.AddMetadata(key, val)
}
