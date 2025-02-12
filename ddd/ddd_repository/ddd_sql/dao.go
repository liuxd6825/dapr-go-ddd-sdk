package ddd_sql

import (
	"context"
	"fmt"
	"github.com/dapr/components-contrib/liuxd/eventstore/impl/gorm_impl/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"gorm.io/gorm"
	"strings"
)

type MapEntity = map[string]any

func NewMapEntity() MapEntity {
	return make(map[string]any)
}

type Dao[T any] struct {
	options       *Options[T]
	entityBuilder ddd.EntityBuilder[T] // 实体构造器
	db            *gorm.DB
	entity        T
	tableName     string
	_table        *gorm.DB
}

const (
	TenantId = "tenant_id"
	Id       = "id"
)

func NewDao[T any](db *gorm.DB, entityBuilder ddd.EntityBuilder[T], tableName string) ddd_repository.Dao[T] {
	entity, err := entityBuilder.NewEntity()
	if err != nil {
		panic(err)
	}
	return &Dao[T]{
		entityBuilder: entityBuilder,
		db:            db,
		entity:        entity,
		tableName:     tableName,
	}
}

func (d *Dao[T]) ExecSql(sql string) error {
	return d.db.Exec(sql).Error
}

func (d *Dao[T]) NewEntity() (T, error) {
	return d.entityBuilder.NewEntity()
}

func (d *Dao[T]) NewEntityList() ([]T, error) {
	return d.entityBuilder.NewEntityList()
}

func (d *Dao[T]) GetTenantId(entity T) string {
	return d.entityBuilder.GetTenantId(entity)
}

func (d *Dao[T]) SetTenantId(entity T, tenantId string) {
	d.entityBuilder.SetTenantId(entity, tenantId)
}

func (d *Dao[T]) GetId(entity T) string {
	return d.entityBuilder.GetId(entity)
}

func (d *Dao[T]) SetId(entity T, id string) {
	d.entityBuilder.SetId(entity, id)
}

func (d *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	err := d.table(ctx).Model(entity).Create(entity).Error
	return ddd_repository.NewSetResult[T](entity, err)
}

func (d *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...ddd_repository.Options) error {
	err := d.table(ctx).Model(data).Create(data).Error
	return err
}

func (d *Dao[T]) InsertMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	v, err := d.NewEntity()
	if err != nil {
		panic(err)
	}
	err = d.table(ctx).Model(v).CreateInBatches(entities, len(entities)).Error
	return ddd_repository.NewSetManyResult(entities, err)
}

func (d *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	err := d.table(ctx).Model(entity).Updates(entity).Error
	return ddd_repository.NewSetResult[T](entity, err)
}

func (d *Dao[T]) UpdateManyByFilter(ctx context.Context, tenantId, filter string, data any, opts ...ddd_repository.Options) *ddd_repository.SetManyCountResult {
	v, err := d.NewEntity()
	if err != nil {
		panic(err)
	}
	res := d.table(ctx).Model(v).Updates(data)
	return ddd_repository.NewSetManyCountResult().SetError(res.Error).SetModifiedCount(res.RowsAffected)
}

func (d *Dao[T]) UpdateManyById(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	v, err := d.NewEntity()
	if err != nil {
		panic(err)
	}
	model := d.table(ctx).Model(v)
	for _, e := range entities {
		id := d.GetId(e)
		model.Where("id = ?", id).UpdateColumns(e)
	}
	return ddd_repository.NewSetManyResult(entities, err)
}

func (d *Dao[T]) getIds(entities []T) []string {
	var ids []string
	for _, entity := range entities {
		ids = append(ids, d.GetId(entity))
	}
	return ids
}
func (d *Dao[T]) UpdateManyMaskById(ctx context.Context, entities []T, mask []string, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	var err error
	model := d.table(ctx).Model(d.entity)
	for _, e := range entities {
		id := d.GetId(e)
		err = model.Where("id = ?", id).UpdateColumn(strings.Join(mask, ","), e).Error
		if err != nil {
			break
		}
	}
	return ddd_repository.NewSetManyResult(entities, err)
}

func (d *Dao[T]) UpdateMapById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) error {
	data[TenantId] = tenantId
	return d.table(ctx).Model(d.entity).Where("id", id).Updates(data).Error
}

func (d *Dao[T]) FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) (T, error) {
	var null T
	err := d.UpdateMapById(ctx, tenantId, id, data, opts...)
	if err != nil {
		return null, err
	}
	res := d.FindById(ctx, tenantId, id, opts...)
	return res.Data, res.Error
}

func (d *Dao[T]) UpdateMap(ctx context.Context, tenantId string, filter any, data any, opts ...ddd_repository.Options) error {
	table := d.table(ctx)
	err := table.Where(filter).Updates(data).Error
	return err
}

func (d *Dao[T]) UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...ddd_repository.Options) (int64, error) {
	table := d.table(ctx)
	var res *gorm.DB
	_ = d.asFilter(filter, func(data map[string]any) error {
		res = table.Where(filter).Updates(data)
		return res.Error
	}, func(sql string) error {
		sql = fmt.Sprintf("%s='%s' and (%s)", TenantId, tenantId, sql)
		res = table.Where(sql).Updates(data)
		return res.Error
	})

	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (d *Dao[T]) Delete(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	tenantId := d.GetTenantId(entity)
	id := d.GetId(entity)
	res := d.table(ctx).Where("tenant_id='%s' and id='%s'", tenantId, id).Delete(entity)
	return ddd_repository.NewSetResult[T](entity, res.Error)
}

func (d *Dao[T]) DeleteByFilter(ctx context.Context, tenantId, filter string, opts ...ddd_repository.Options) error {
	sql, err := d.getSql(filter)
	if err != nil {
		return err
	}
	table := d.table(ctx)
	res := table.Where("tenant_id=?", tenantId).Delete(sql)
	return res.Error
}

func (d *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	var null T
	table := d.table(ctx)
	res := table.Where("tenant_id=? and id=?", tenantId, id)
	return ddd_repository.NewSetResult[T](null, res.Error)
}

func (d *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error {
	table := d.table(ctx)
	res := table.Where("tenant_id=?", tenantId).Delete(ids)
	return res.Error
}

func (d *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	var null T
	table := d.table(ctx)
	res := table.Where("tenant_id=?", tenantId).Delete("")
	return ddd_repository.NewSetResult[T](null, res.Error)
}

func (d *Dao[T]) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	var null T
	sql := d.mapAsSql(tenantId, filterMap)
	table := d.table(ctx)
	res := table.Where("tenant_id=?", tenantId).Delete(sql)
	return ddd_repository.NewSetResult[T](null, res.Error)
}

func (d *Dao[T]) FindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	var data T
	table := d.table(ctx)
	res := table.Where("id=? and tenant_id=?", id, tenantId).Find(&data)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindOneResult[T](data, isFound, res.Error)
}

func (d *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	var data []T
	table := d.table(ctx)
	res := table.Where("id in [?] and tenant_id=?", ids, tenantId).Find(&data)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindListResult[T](data, isFound, res.Error)
}

func (d *Dao[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	var data T
	sql := d.mapAsSql(tenantId, filterMap)
	table := d.table(ctx)
	res := table.Where(sql).Find(&data)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindOneResult[T](data, isFound, res.Error)
}

func (d *Dao[T]) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	var data []T
	sql := d.mapAsSql(tenantId, filterMap)
	table := d.table(ctx)
	res := table.Where(sql).Find(&data)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindListResult[T](data, isFound, res.Error)
}

func (d *Dao[T]) FindByRSQL(ctx context.Context, tenantId string, rSql string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	sql, err := d.getSql(rSql)
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	var list []T
	table := d.table(ctx)
	res := table.Where("tenant_id=?", tenantId).Find(&list, sql)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindListResult[T](list, isFound, res.Error)
}

func (d *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	var list []T
	table := d.table(ctx)
	res := table.Where("tenant_id=?", tenantId).Find(&list)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindListResult[T](list, isFound, res.Error)
}

func (d *Dao[T]) FindPaging(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) (result *ddd_repository.FindPagingResult[T]) {
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

func (d *Dao[T]) SumEntity(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]T, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) SumMap(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) Sum(ctx context.Context, qry ddd_repository.FindPagingQuery, data any, opts ...ddd_repository.Options) (any, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) CountRows(ctx context.Context, tenantId string, filterData any, opts ...ddd_repository.Options) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) Count(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) table(ctx context.Context, opts ...ddd_repository.Options) *gorm.DB {
	tx := db.GetTransaction(ctx)
	if tx != nil {
		return tx.Table(d.tableName)
	}
	if d._table == nil {
		d._table = d.db.Table(d.tableName)
	}
	return d._table
}

func (d *Dao[T]) asFilter(filter any, mapFunc func(data map[string]any) error, sqlFunc func(sql string) error) error {
	if filter == nil {
		return errors.New("filter cannot be nil")
	}
	if mapData, ok := filter.(map[string]any); ok {
		return mapFunc(mapData)
	}
	if sqlData, ok := filter.(string); ok {
		return sqlFunc(sqlData)
	}
	return errors.New("filter data type is not string")
}

func (r *Dao[T]) mapAsSql(tenantId string, filterMap map[string]interface{}) string {
	if filterMap == nil || len(filterMap) == 0 {
		return fmt.Sprintf(`tenant_id='%s'`, tenantId)
	}

	ands := make([]string, 1)
	ands[0] = fmt.Sprintf(`tenant_id='%s'`, tenantId)

	for fieldName, fieldValue := range filterMap {
		var item string
		if fieldValue == nil {

		} else if v, ok := fieldValue.(int64); ok {
			item = fmt.Sprintf(`%s=%v`, fieldName, v)
		} else {
			item = fmt.Sprintf(`%s='%s'`, fieldName, fieldValue)
		}
		ands = append(ands, item)
	}
	return strings.Join(ands, " AND ")
}

func (d *Dao[T]) getSql(rSql string) (string, error) {
	res, err := rsql.SqlParseProcess(rSql)
	return res, err
}
