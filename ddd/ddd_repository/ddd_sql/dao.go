package ddd_sql

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql/rsql_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"gorm.io/gorm"
	"strings"
)

type Dao[T any] struct {
	options       *Options[T]
	entityBuilder ddd.EntityBuilder[T] // 实体构造器
	db            *gorm.DB
	dbKey         string
	entity        T
	tableName     string
	metadata      map[string]any
}

const (
	TenantId = "tenant_id"
	Id       = "id"
)

var initializePlugin bool = false

func NewDaoWithDbKey[T any](dbKey string, eb ddd.EntityBuilder[T], tableName string) ddd_repository.Dao[T] {
	item := restapp.GetDb(dbKey)
	if item == nil {
		panic(errors.New(fmt.Sprintf("db key %s not found", dbKey)))
	}
	var db *gorm.DB
	switch item.GetDBType() {
	case restapp.DbType_Postgres:
		db = item.GetPostgres()
	case restapp.DbType_MySQL:
		db = item.GetMySQL()
	case restapp.DbType_Sqlite:
		db = item.GetSqlite()
	case restapp.DbType_MsSQL:
		db = item.GetMsSQL()
	case restapp.DbType_Oracle:
		db = item.GetOracle()
	default:
		panic(errors.New(fmt.Sprintf("db type %s not supported", item.GetDBType())))
	}
	return NewDao[T](db, dbKey, eb, tableName)
}

func NewDao[T any](db *gorm.DB, dbKey string, entityBuilder ddd.EntityBuilder[T], tableName string) *Dao[T] {
	entity := entityBuilder.NewEntity()

	if !initializePlugin {
		initializePlugin = true
		if err := db.Use(NewGormFieldPlugin()); err != nil {
			panic(err)
		}
	}

	return &Dao[T]{
		entityBuilder: entityBuilder,
		db:            db,
		dbKey:         dbKey,
		entity:        entity,
		tableName:     tableName,
		metadata:      make(map[string]any),
	}
}

func newDao[T any](db *gorm.DB, dbKey string, entityBuilder ddd.EntityBuilder[T], tableName string) *Dao[T] {
	entity := entityBuilder.NewEntity()

	return &Dao[T]{
		entityBuilder: entityBuilder,
		db:            db,
		dbKey:         dbKey,
		entity:        entity,
		tableName:     tableName,
	}
}

func (d *Dao[T]) SetMetadata(metadata map[string]any) {
	d.metadata = metadata
}

func (d *Dao[T]) GetMetadata() map[string]any {
	return d.metadata
}

func (d *Dao[T]) ExecSql(sql string) error {
	return d.db.Exec(sql).Error
}

func (d *Dao[T]) NewEntity() T {
	return d.entityBuilder.NewEntity()
}

func (d *Dao[T]) NewEntityList() []T {
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

func (d *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) (res *ddd_repository.SetResult[T]) {
	err := gp.Try(func() error {
		db := d.table(ctx)
		return db.Create(entity).Error
	}).Error
	return ddd_repository.NewSetResult[T](entity, err)
}

func (d *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]interface{}, opts ...ddd_repository.Options) error {
	err := gp.Try(func() error {
		return d.table(ctx).Model(data).Create(data).Error
	}).Error
	return err
}

func (d *Dao[T]) InsertMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	v := d.NewEntity()
	err := gp.Try(func() error {
		return d.table(ctx).Model(v).CreateInBatches(entities, len(entities)).Error
	}).Error
	return ddd_repository.NewSetManyResult(entities, err)
}

func (d *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	err := gp.Try(func() error {
		return d.table(ctx).Model(entity).Updates(entity).Error
	}).Error
	return ddd_repository.NewSetResult[T](entity, err)
}

func (d *Dao[T]) UpdateManyByFilter(ctx context.Context, tenantId, filter string, data any, opts ...ddd_repository.Options) *ddd_repository.SetManyCountResult {
	v := d.NewEntity()
	var res *gorm.DB
	err := gp.Try(func() error {
		res = d.table(ctx).Model(v).Updates(data)
		return res.Error
	}).Error
	if res != nil {
		return ddd_repository.NewSetManyCountResult().SetError(res.Error).SetModifiedCount(res.RowsAffected)
	}
	return ddd_repository.NewSetManyCountResult().SetError(err)
}

func (d *Dao[T]) UpdateManyById(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T] {
	var err error
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	v := d.NewEntity()
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
	sql, err := d.getSql(tenantId, filter)
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

func (d *Dao[T]) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
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
	res := table.Where("id in ? and tenant_id=?", ids, tenantId).Find(&data)
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
	sql, err := d.getSql(tenantId, rSql)
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
	return d.findPaging(ctx, qry, opts...)
}

func (d *Dao[T]) findPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	return d.DoFilter(query.GetTenantId(), query.GetFilter(), func(sqlWhere string) (*ddd_repository.FindPagingResult[T], bool, error) {
		var err error
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()

		list := d.entityBuilder.NewEntityList()

		tx := d.table(ctx)

		if len(query.GetFields()) > 0 {
			fields := strings.Split(query.GetFields(), ",")
			for i, field := range fields {
				fields[i] = stringutils.AsFieldName(field)
			}
			tx = tx.Select(strings.Join(fields, ","))
		}

		if len(sqlWhere) > 0 {
			tx = tx.Where(sqlWhere)
		}

		if query.GetPageSize() > 0 {
			tx = tx.Limit(int(query.GetPageSize()))
		}

		if query.GetPageNum() > 0 {
			tx = tx.Offset(int(query.GetPageSize() * query.GetPageNum()))
		}

		if len(query.GetSort()) > 0 {
			tx = tx.Order(query.GetSort())
		}

		if err = tx.Find(&list).Error; err != nil {
			return nil, false, err
		}

		var totalRows int64 = -1
		if query.GetIsTotalRows() {
			tx := d.table(ctx)
			if len(sqlWhere) > 0 {
				tx = tx.Where(sqlWhere)
			}
			if err := tx.Count(&totalRows).Error; err != nil {
				return nil, false, err
			}
		}

		findData := ddd_repository.NewFindPagingResult[T](list, &totalRows, query, err)
		return findData, true, err
	})

}

func (d *Dao[T]) FindAutoComplete(ctx context.Context, qry ddd_repository.FindAutoCompleteQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	f := ddd_repository.NewFindPagingQuery()
	groupCols := []*ddd_repository.GroupCol{
		{Field: qry.GetField(), DataType: types.DataTypeString},
	}

	f.SetGroupCols(groupCols)
	f.SetTenantId(qry.GetTenantId())
	f.SetFields(qry.GetFields())
	f.SetFilter(qry.GetFilter())
	f.SetMustFilter(qry.GetMustWhere())

	f.SetPageNum(qry.GetPageNum())
	f.SetPageSize(qry.GetPageSize())
	f.SetSort(qry.GetSort())
	f.SetIsTotalRows(false)

	return d.FindPaging(ctx, f, opts...)
}

func (d *Dao[T]) FindDistinct(ctx context.Context, qry ddd_repository.FindDistinctQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	f := ddd_repository.NewFindPagingQuery()

	f.SetGroupCols(qry.GetGroupCols())
	f.SetTenantId(qry.GetTenantId())
	f.SetFields(qry.GetFields())
	f.SetFilter(qry.GetFilter())
	f.SetMustFilter(qry.GetMustWhere())

	f.SetPageNum(qry.GetPageNum())
	f.SetPageSize(qry.GetPageSize())
	f.SetSort(qry.GetSort())
	f.SetIsTotalRows(false)

	return d.FindPaging(ctx, f, opts...)
}

func (d *Dao[T]) SumEntity(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]T, bool, error) {
	data := d.NewEntityList()
	_, found, err := d.Sum(ctx, qry, &data, opts...)
	return data, found, err
}

func (d *Dao[T]) SumMap(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error) {
	data := make([]map[string]any, 0)
	_, found, err := d.Sum(ctx, qry, &data, opts...)
	return data, found, err
}

func (d *Dao[T]) Sum(ctx context.Context, qry ddd_repository.FindPagingQuery, resData any, opts ...ddd_repository.Options) (any, bool, error) {
	if len(qry.GetValueCols()) == 0 {
		return nil, false, nil
	}

	var err error

	f1 := qry.GetFilter()
	f2 := qry.GetMustFilter()
	f3 := ""
	mustWhere, ok := qry.(ddd_repository.FindPagingQueryMustWhere)
	if ok {
		f3, err = mustWhere.GetMustWhere()
		if err != nil {
			return nil, false, err
		}
	}
	filter := getSqlAnds(f1, f2, f3)

	res, found, err := d.sum(ctx, qry.GetTenantId(), filter, qry.GetValueCols(), resData, opts...)
	return res, found, err
}

func (d *Dao[T]) sum(ctx context.Context, tenantId, rSql string, valueCols []*ddd_repository.ValueCol, resData any, opts ...ddd_repository.Options) (any, bool, error) {
	p := rsql_sql.NewProcess(tenantId)
	if err := rsql.ParseProcess(rSql, p); err != nil {
		return nil, false, err
	}
	sql := p.GetSQL()
	table := d.table(ctx)
	sumFields := make([]string, 0)
	for _, col := range valueCols {
		sumFields = append(sumFields, fmt.Sprintf("sum(%s) as %s", col.Field, col.Field))
	}
	sumSql := strings.Join(sumFields, ", ")

	res := table.Where(sql).Select(sumSql).Scan(resData)
	return resData, res.Error != nil, nil
}

func (d *Dao[T]) CountRows(ctx context.Context, tenantId string, filterData any, opts ...ddd_repository.Options) (int64, error) {
	var count int64
	filterMap, ok := filterData.(map[string]any)
	var sql string
	if ok {
		sql = d.mapAsSql(tenantId, filterMap)
	} else if where, ok := filterData.(string); ok {
		sql = where
	}
	table := d.table(ctx)
	if sql == "" {
		sql = fmt.Sprintf("tenant_id='%s'", tenantId)
	} else {
		sql = fmt.Sprintf("tenant_id='%s' and (%s)", tenantId, sql)
	}
	res := table.Where(sql).Select("count(*) as total").Scan(&count)
	return count, res.Error
}

func (d *Dao[T]) Count(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) (int64, error) {
	var count int64
	res := d.DoFilter(tenantId, rsql, func(sqlWhere string) (*ddd_repository.FindPagingResult[T], bool, error) {
		table := d.table(ctx)
		res := table.Where(sqlWhere).Select("count(*) as total").Scan(&count)
		return ddd_repository.NewFindPagingResultEmpty[T]().SetError(res.Error), res.Error == nil, nil
	})
	return count, res.Error
}

func (d *Dao[T]) table(ctx context.Context, opts ...ddd_repository.Options) *gorm.DB {
	tx := GetTx(ctx, d.db.Name())
	if tx == nil {
		tx = d.db
	}
	items := plugins.Items()
	hasPlugin := false
	for _, item := range items {
		if plugin, ok := item.(GetTablePlugin); ok {
			hasPlugin = true
			tx = plugin.GetTable(ctx, d, tx, d.tableName, opts...)
		}
	}
	if hasPlugin {
		return tx
	}
	return tx.Table(d.tableName)
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

func (d *Dao[T]) GetFilterMap(tenantId string, rSql string) map[string]any {
	return map[string]any{}
}

func (d *Dao[T]) mapAsSql(tenantId string, mapData any) string {
	filterMap, ok := mapData.(map[string]any)
	if ok {
		panic(errors.New("filter data type is not string"))
	}
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

func (d *Dao[T]) getSql(tenantId, rSql string) (string, error) {
	proc := rsql_sql.NewProcess(tenantId)
	err := rsql.ParseProcess(rSql, proc)
	return proc.GetSQL(), err
}

func (d *Dao[T]) DoFilter(tenantId, filter string, fun func(sqlWhere string) (*ddd_repository.FindPagingResult[T], bool, error)) *ddd_repository.FindPagingResult[T] {
	if tenantId == "" {
		return ddd_repository.NewFindPagingResultWithError[T](errors.New("tenantId can not be empty"))
	}
	var sqlWhere string
	if filter == "" {
		sqlWhere = fmt.Sprintf("tenant_id='%s'", tenantId)
	} else {
		process := rsql_sql.NewProcess(tenantId)
		if err := rsql.ParseProcess(filter, process); err != nil {
			return ddd_repository.NewFindPagingResultWithError[T](err)
		}
		sqlWhere = fmt.Sprintf("tenant_id='%s' and (%s)", tenantId, process.GetSQL())
	}
	data, _, err := fun(sqlWhere)
	if err != nil {
		err = nil
	}
	return data
}

func (d *Dao[T]) StartTx(ctx context.Context, fun ddd_repository.TxFunc, options ...*ddd_repository.SessionOptions) (err error) {
	return StartTx(ctx, d.db, d.dbKey, fun, options...)
}

func getSqlAnds(s ...string) string {
	res := ""
	for _, item := range s {
		res = getSqlAnd(res, item)
	}
	return res
}

func getSqlAnd(s1 string, s2 string) string {
	b1 := len(s1) > 0
	b2 := len(s2) > 0
	if b1 && b2 {
		return fmt.Sprintf("(%s) and (%s)", s1, s2)
	} else if b1 {
		return s1
	}
	return s2
}
