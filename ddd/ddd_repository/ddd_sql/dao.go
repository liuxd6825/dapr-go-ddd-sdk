package ddd_sql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql/rsql_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
	"strings"
)

type Dao[T any] struct {
	options    *Options[T]
	eb         ddd.EntityBuilder[T] // 实体构造器
	db         *gorm.DB
	dbKey      string
	entity     T
	tableName  string
	dbSchema   *dbschema.Schema
	gormSchema *gormschema.Schema
}

const (
	TenantId = "tenant_id"
	Id       = "id"
)

var initializePlugin bool = false

var fields = ddd.GetFields()

func NewDaoWithDbKey[T any](cfg *NewConfig) ddd_repository.Dao[T] {
	dbKey := cfg.DbKey
	item := restapp.GetDb(dbKey)
	if item == nil {
		panic(errors.New(fmt.Sprintf("db key %s not found", dbKey)))
	}
	db := cfg.Db
	if db == nil {
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
	}
	cfg.Db = db
	return NewDao[T](cfg)
}

type NewConfig struct {
	DbKey      string
	TableName  string
	Db         *gorm.DB
	DBSchema   *dbschema.Schema
	GormSchema *gormschema.Schema
}

func (c *NewConfig) Check() error {
	if c.Db == nil {
		return errors.New(fmt.Sprintf("ddd_sql.NewConfig.DB is nil"))
	}
	if c.DBSchema == nil {
		return errors.New(fmt.Sprintf("ddd_sql.NewConfig.DbSchema is nil"))
	}
	if c.TableName == "" || c.DBSchema.TableName == "" {
		return errors.New(fmt.Sprintf("ddd_sql.NewConfig.TableName and ddd_sql.NewConfig.DBSchema.TableName is empty "))
	}
	return nil
}

func NewDao[T any](cfg *NewConfig) *Dao[T] {
	if !initializePlugin {
		initializePlugin = true
		if err := cfg.Db.Use(NewGormFieldPlugin()); err != nil {
			panic(err)
		}
	}
	return newDao[T](cfg)
}

func newDao[T any](cfg *NewConfig) *Dao[T] {
	var err error

	if err = cfg.Check(); err != nil {
		panic(err)
	}

	eb := ddd.NewAnyEntityBuilder[T](cfg.DBSchema)
	entity := eb.NewEntity()

	gormSch := cfg.GormSchema
	if gormSch == nil {
		gormSch, err = NewGormSchema(cfg.DBSchema)
		if err != nil {
			panic(fmt.Sprintf("newDao() gorm schema error:%s", err.Error()))
		}
	}

	tableName := cfg.TableName
	if tableName == "" {
		tableName = cfg.DBSchema.TableName
	}

	return &Dao[T]{
		eb:         eb,
		db:         cfg.Db,
		dbKey:      cfg.DbKey,
		entity:     entity,
		tableName:  cfg.TableName,
		dbSchema:   cfg.DBSchema,
		gormSchema: gormSch,
	}
}

func (d *Dao[T]) getIds(entities []T) []string {
	var ids []string
	for _, entity := range entities {
		ids = append(ids, d.GetId(entity))
	}
	return ids
}

func (d *Dao[T]) GetSchema() *dbschema.Schema {
	return d.dbSchema
}

func (d *Dao[T]) ExecSql(sql string) error {
	return d.db.Exec(sql).Error
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

func (d *Dao[T]) Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) (res *ddd_repository.SetResult[T]) {
	res = ddd_repository.NewSetResult[T]()
	res.Error = gp.Try(func() error {
		db := d.table(ctx)
		d.eb.SetCreatedInfo(ctx, entity)
		db = db.Create(entity)
		res.SetRowsAffected(db.RowsAffected)
		return db.Error
	}).Error
	return res
}

func (d *Dao[T]) InsertMap(ctx context.Context, tenantId string, data map[string]any, opts ...ddd_repository.Options) (res *ddd_repository.SetResult[T]) {
	res = ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		d.eb.SetCreatedInfo(ctx, data)
		db := d.table(ctx).Model(data).Create(data)
		res.SetRowsAffected(db.RowsAffected)
		return db.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) InsertMany(ctx context.Context, tenantId string, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		for _, entity := range entities {
			d.eb.SetCreatedInfo(ctx, entity)
			d.SetTenantId(entity, tenantId)
		}
		db := d.table(ctx).CreateInBatches(entities, len(entities))
		res.SetRowsAffected(db.RowsAffected)
		return db.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res

}

func (d *Dao[T]) updateTable(ctx context.Context, opts ...ddd_repository.Options) *gorm.DB {
	table := d.table(ctx)
	opt := ddd_repository.NewOptions(opts...)
	// 指定更新字段
	updateFields := opt.GetUpdateFields()
	if len(updateFields) > 0 {
		for _, v := range updateFields {
			table = table.Select(v)
		}
	}

	table = table.Omit(fields.CreatedTime, fields.CreatorId, fields.CreatorName)

	// 指定取消更新的字段
	cancelFields := opt.GetUpdateCancel()
	if len(cancelFields) > 0 {
		for _, name := range cancelFields {
			table = table.Omit(name)
		}
	}
	return table
}

func (d *Dao[T]) Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResult[T]()
	_ = gp.Try(func() error {
		opt := ddd_repository.NewOptions(opts...)
		d.eb.SetUpdatedInfo(ctx, entity)
		tenantId := d.eb.GetTenantId(entity)
		table := d.updateTable(ctx, opt).Where("id=? and tenant_id=?", d.GetId(entity), tenantId)

		// 是否空值更新
		if !opt.GetNullUpdate() {
			if e, ok := any(entity).(map[string]any); ok {
				for k, v := range e {
					if v == nil {
						table = table.Omit(k)
					}
				}
			}
		}

		db := table.UpdateColumns(entity)
		res.SetRowsAffected(db.RowsAffected)
		return db.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})

	return res
}

func (d *Dao[T]) UpdateByRSQL(ctx context.Context, tenantId string, filterRSQL string, data T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	v := d.NewEntity()
	res := ddd_repository.NewSetResultEmpty[T]()
	_ = gp.Try(func() error {
		where, err := d.getSql(tenantId, filterRSQL)
		if err != nil {
			return err
		}
		d.eb.SetUpdatedInfo(ctx, data)
		db := d.updateTable(ctx, opts...).Where(where).Model(v).Updates(data)
		res.SetRowsAffected(db.RowsAffected)
		return db.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) UpdateMany(ctx context.Context, tenantId string, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResultEmpty[T]()
	gp.Try(func() error {
		for _, e := range entities {
			d.eb.SetUpdatedInfo(ctx, e)
		}
		db := d.updateTable(ctx, opts...).Model(d.NewEntity()).Save(entities)
		res.SetRowsAffected(db.RowsAffected)
		return db.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) UpdateManyMaskById(ctx context.Context, entities []T, mask []string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	var err error
	var res = ddd_repository.NewSetResult[T]()
	gp.Try(func() error {
		opt := ddd_repository.NewOptions(opts...)
		model := d.updateTable(ctx, opt).Model(d.entity)
		for _, e := range entities {
			id := d.GetId(e)
			d.eb.SetUpdatedInfo(ctx, e)
			db := model.Where("id = ?", id).UpdateColumn(strings.Join(mask, ","), e)
			if err != nil {
				return db.Error
			}
		}
		return nil
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) UpdateMap(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResult[T]()
	gp.Try(func() error {
		data[TenantId] = tenantId
		d.eb.SetUpdatedInfo(ctx, data)
		db := d.updateTable(ctx, opts...).Model(d.entity).Where("id", id).Updates(data)
		res.Error = db.Error
		res.RowsAffected = db.RowsAffected
		return res.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) UpdateMapAndGetCount(ctx context.Context, tenantId string, filter any, data any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResult[T]()
	gp.Try(func() error {
		var db *gorm.DB
		var count int64
		table := d.updateTable(ctx, opts...)
		res.Error = d.asFilter(filter, func(data map[string]any) error {
			d.eb.SetUpdatedInfo(ctx, data)
			db = table.Where(filter).Updates(data)
			return res.Error
		}, func(sql string) error {
			sql = fmt.Sprintf("%s='%s' and (%s)", TenantId, tenantId, sql)
			db = table.Where(sql).Count(&count)
			return res.Error
		})

		res.SetError(db.Error)
		res.SetRowsAffected(db.RowsAffected)
		return res.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})

	return res
}

func (d *Dao[T]) Delete(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResult[T]()
	gp.Try(func() error {
		tenantId := d.GetTenantId(entity)
		id := d.GetId(entity)
		d.eb.SetDeletedInfo(ctx, entity)
		db := d.table(ctx).Where("tenant_id=? and id=?", tenantId, id).Delete(entity)
		res.SetRowsAffected(db.RowsAffected)
		return res.Error
	}).Catch(func(err error) {
		res.SetError(err)
	})
	return res
}

func (d *Dao[T]) DeleteByRSQL(ctx context.Context, tenantId, rSQL string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	res := ddd_repository.NewSetResult[T]()
	where, err := d.getSql(tenantId, rSQL)
	if err != nil {
		return res.SetError(err)
	}
	table := d.table(ctx, opts...)
	db := table.Where("tenant_id=?", tenantId).Where(where).Delete(nil)
	return res.SetError(db.Error).SetRowsAffected(db.RowsAffected)
}

func (d *Dao[T]) DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	table := d.table(ctx, opts...)
	db := table.Where("tenant_id=? and id=?", tenantId, id).Delete(id)
	return ddd_repository.NewSetResult[T]().SetError(db.Error).SetRowsAffected(db.RowsAffected)
}

func (d *Dao[T]) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	var null = d.NewEntity()
	table := d.table(ctx, opts...)
	db := table.Where("tenant_id=? and id in ?", tenantId, ids).Delete(null)
	return ddd_repository.NewSetResult[T]().SetError(db.Error).SetRowsAffected(db.RowsAffected)
}

func (d *Dao[T]) DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	table := d.table(ctx, opts...)
	db := table.Where("tenant_id=?", tenantId).Delete("")
	return ddd_repository.NewSetResult[T]().SetError(db.Error).SetRowsAffected(db.RowsAffected)
}

func (d *Dao[T]) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]any, opts ...ddd_repository.Options) *ddd_repository.SetResult[T] {
	sql := d.mapAsSql(tenantId, filterMap)
	table := d.table(ctx, opts...)
	db := table.Where("tenant_id=?", tenantId).Delete(sql)
	return ddd_repository.NewSetResult[T]().SetError(db.Error).SetRowsAffected(db.RowsAffected)
}

func (d *Dao[T]) FindById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	var data T
	table := d.table(ctx, opts...)
	res := table.Where("id=? and tenant_id=?", id, tenantId).Find(&data)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindOneResult[T](data, isFound, res.Error)
}

func (d *Dao[T]) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	var data []T
	table := d.table(ctx, opts...)
	res := table.Where("id in ? and tenant_id=?", ids, tenantId).Find(&data)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindListResult[T](data, isFound, res.Error)
}

func (d *Dao[T]) FindOneAndUpdateById(ctx context.Context, tenantId string, id string, data map[string]any, opts ...ddd_repository.Options) (T, error) {
	d.eb.SetUpdatedInfo(ctx, data)
	d.UpdateMap(ctx, tenantId, id, data, opts...)
	res := d.FindById(ctx, tenantId, id, opts...)
	return res.Data, res.Error
}

func (d *Dao[T]) FindOneByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindOneResult[T] {
	var data T
	sql := d.mapAsSql(tenantId, filterMap)
	table := d.table(ctx, opts...)
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
	table := d.table(ctx, opts...)
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
	table := d.table(ctx, opts...)
	res := table.Where("tenant_id=?", tenantId).Find(&list, sql)
	isFound := false
	if res.RowsAffected > 0 {
		isFound = true
	}
	return ddd_repository.NewFindListResult[T](list, isFound, res.Error)
}

func (d *Dao[T]) FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	var list []T
	table := d.table(ctx, opts...)
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

func (d *Dao[T]) getAllFilter(qry ddd_repository.FindPagingQuery) string {
	filter := qry.GetFilter()
	mustFilter := qry.GetMustFilter()
	if filter != "" && mustFilter != "" {
		return fmt.Sprintf("(%s) and (%s)", filter, mustFilter)
	} else if filter != "" {
		return filter
	}
	return mustFilter
}

func (d *Dao[T]) findPaging(ctx context.Context, query ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	filter := d.getAllFilter(query)
	return d.DoFilter(query.GetTenantId(), filter, func(sqlWhere string) (findRes *ddd_repository.FindPagingResult[T], isFound bool, err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()

		qryDb := d.table(ctx, opts...)

		if len(query.GetFields()) > 0 {
			fs := strings.Split(query.GetFields(), ",")
			for i, field := range fs {
				fs[i] = stringutils.AsFieldName(field)
			}
			qryDb = qryDb.Select(strings.Join(fs, ","))
		}

		if len(sqlWhere) > 0 {
			qryDb = qryDb.Where(sqlWhere)
		}

		if query.GetPageSize() > 0 {
			qryDb = qryDb.Limit(int(query.GetPageSize()))
		}

		if query.GetPageNum() > 0 {
			qryDb = qryDb.Offset(int(query.GetPageSize() * query.GetPageNum()))
		}

		if len(query.GetSort()) > 0 {
			// 将sort格式  name:asc, sex:asc
			// 转为sql格式  name asc, sex asc格式
			sort := strings.ReplaceAll(query.GetSort(), ":", " ")
			qryDb = qryDb.Order(sort)
		}

		countDb := d.table(ctx)
		sumDb := d.table(ctx)
		isGroup := false
		colLen := len(query.GetGroupCols())
		keyLen := len(query.GetGroupKeys())
		// 是分组模式
		if colLen > 0 {
			isGroup = true
			// 分组where条件
			for i, value := range query.GetGroupKeys() {
				field := query.GetGroupCols()[i]
				qryDb = qryDb.Where(field.Field+"=?", value)
				countDb = countDb.Where(field.Field+"=?", value)
				sumDb = sumDb.Where(field.Field+"=?", value)
			}
			// 以groupKey位置的上个字段为分组字段
			if colLen-keyLen > 0 {
				field := query.GetGroupCols()[keyLen]
				qryDb = qryDb.Group(field.Field).Select(field.Field)
				countDb.Group(field.Field).Select(field.Field)

			}
			if isGroup && colLen-keyLen > 0 {
				d.setDbValueCols(qryDb, query.GetValueCols())
			}

		}

		list := d.eb.NewEntityList()
		if err = qryDb.Find(&list).Error; err != nil {
			return nil, false, err
		}

		// 是分组时，需要添加虚拟id
		if isGroup && colLen-keyLen > 0 {
			d.setListIds(list)
		}

		var totalRows int64 = -1
		if query.GetIsTotalRows() {
			if len(sqlWhere) > 0 {
				countDb = countDb.Where(sqlWhere)
			}
			if err := countDb.Count(&totalRows).Error; err != nil {
				return nil, false, err
			}
		}

		findData := ddd_repository.NewFindPagingResult[T](list, totalRows, query, err)
		// 不是分组模式，并且有汇总数据
		if !isGroup && len(query.GetValueCols()) > 0 {
			if len(sqlWhere) > 0 {
				sumDb = sumDb.Where(sqlWhere)
			}
			d.setDbValueCols(sumDb, query.GetValueCols())
			sumList := d.NewEntityList()
			sumDb.Find(&sumList)
			d.setListIds(sumList)

			findData.SetSum(true, sumList, sumDb.Error)
		}

		return findData, true, err
	})

}

func (d *Dao[T]) setListIds(list []T) {
	for _, item := range list {
		uid, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		d.eb.SetId(item, uid.String())
	}
}

// setDbValueCols 设置sum,count,avg,max,min,
func (d *Dao[T]) setDbValueCols(db *gorm.DB, valCols []*ddd_repository.ValueCol) {
	if len(valCols) > 0 {
		fields := db.Statement.Selects
		for _, valCol := range valCols {
			switch valCol.AggFunc {
			case ddd_repository.AggFuncSum:
				fields = append(fields, fmt.Sprintf("sum(%s) as %s", valCol.Field, valCol.Field))
				break
			case ddd_repository.AggFuncCount:
				fields = append(fields, fmt.Sprintf("count(%s) as %s", valCol.Field, valCol.Field))
				break
			case ddd_repository.AggFuncAvg:
				fields = append(fields, fmt.Sprintf("avg(%s) as %s", valCol.Field, valCol.Field))
				break
			case ddd_repository.AggFuncFirst:
				//qryDb.Select(fmt.Sprintf("fisrt(%s)", valCol.Field))
				break
			case ddd_repository.AggFuncLast:
				//qryDb.Select(fmt.Sprintf("last(%s)", valCol.Field))
				break
			case ddd_repository.AggFuncMax:
				fields = append(fields, fmt.Sprintf("max(%s) as %s", valCol.Field, valCol.Field))
				break
			case ddd_repository.AggFuncMin:
				fields = append(fields, fmt.Sprintf("min(%s) as %s", valCol.Field, valCol.Field))
				break
			case ddd_repository.AggFuncZero:
				break
			}
		}
		db.Statement.Selects = fields
	}
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
	_, found, err := d.SumByQuery(ctx, qry, &data, opts...)
	return data, found, err
}

func (d *Dao[T]) SumByQuery(ctx context.Context, qry ddd_repository.FindPagingQuery, resData any, opts ...ddd_repository.Options) (any, bool, error) {
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

func (d *Dao[T]) SumByRSQL(ctx context.Context, tenantId, rSql string, valueCols []*ddd_repository.ValueCol, opts ...ddd_repository.Options) map[string]any {
	data := make(map[string]any)
	_, _, err := d.sum(ctx, tenantId, rSql, valueCols, &data, opts...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *Dao[T]) sum(ctx context.Context, tenantId, rSql string, valueCols []*ddd_repository.ValueCol, resData any, opts ...ddd_repository.Options) (any, bool, error) {
	p := rsql_sql.NewProcess(tenantId)
	if err := rsql.ParseProcess(rSql, p); err != nil {
		return nil, false, err
	}
	sql := p.GetSQL()
	table := d.table(ctx, opts...)
	sumFields := make([]string, 0)
	for _, col := range valueCols {
		sumFields = append(sumFields, fmt.Sprintf("sum(%s) as %s", col.Field, col.Field))
	}
	sumSql := strings.Join(sumFields, ", ")

	res := table.Where(sql).Select(sumSql).Scan(resData)
	return resData, res.Error != nil, nil
}

func (d *Dao[T]) CountByMap(ctx context.Context, tenantId string, filterData any, opts ...ddd_repository.Options) (int64, error) {
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

func (d *Dao[T]) CountByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) (int64, error) {
	var count int64
	res := d.DoFilter(tenantId, rsql, func(sqlWhere string) (*ddd_repository.FindPagingResult[T], bool, error) {
		table := d.table(ctx, opts...)
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
	return tx.Table(d.tableName).MapSchema(d.gormSchema).Unscoped()
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
	if !ok {
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

func (d *Dao[T]) DoFilter(tenantId, filter string, fun func(sqlWhere string) (*ddd_repository.FindPagingResult[T], bool, error)) (findRes *ddd_repository.FindPagingResult[T]) {
	findRes = ddd_repository.NewFindPagingResultEmpty[T]()
	if tenantId == "" {
		return findRes.SetError(errors.New("tenantId can not be empty"))
	}
	var sqlWhere string
	if filter == "" {
		sqlWhere = fmt.Sprintf("tenant_id='%s'", tenantId)
	} else {
		process := rsql_sql.NewProcess(tenantId)
		if err := rsql.ParseProcess(filter, process); err != nil {
			return findRes.SetError(err)
		}
		sqlWhere = fmt.Sprintf("tenant_id='%s' and (%s)", tenantId, process.GetSQL())
	}
	data, _, err := fun(sqlWhere)
	if data != nil {
		findRes = data
	}
	return findRes.SetError(err)
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
