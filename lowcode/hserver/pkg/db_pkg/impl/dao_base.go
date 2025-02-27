package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
	"time"
)

type DaoBase struct {
	cfg         *db.DaoConfig
	dbKey       string                             // 配置中的数据库Key
	tableName   string                             // 表名
	appId       string                             // 应用ID
	aggField    string                             // 聚合根字段
	isPubEvent  bool                               // 是否发布事件
	eventPrefix string                             // 事件前缀
	dao         ddd_repository.Dao[map[string]any] // 数据访问
	env         common.IEnvConfig                  // 环境变量
}

const (
	CreatedTime = "createdTime"
	CreatorId   = "creatorId"
	CreatorName = "creatorName"
	UpdatedTime = "updatedTime"
	UpdaterId   = "updaterId"
	UpdaterName = "updaterName"
	DeletedTime = "deletedTime"
	DeleterId   = "deleterId"
	DeleterName = "deleterName"
	IsDeleted   = "isDeleted"
	TenantId    = "tenantId"
	Id          = "id"
)

func NewDaoBase(dao ddd_repository.Dao[map[string]any], cfg *db.DaoConfig) *DaoBase {
	if cfg == nil {
		panic("dao base config is nil")
	}
	aggField := cfg.AggField
	if aggField == "" {
		aggField = "id"
	}
	tableName := stringutils.AsFieldName(cfg.Schema.Name)
	return &DaoBase{
		dao:         dao,
		dbKey:       cfg.DbKey,
		tableName:   tableName,
		appId:       cfg.GetEnv().GetAppId(),
		isPubEvent:  cfg.GetIsPubEvent(),
		aggField:    aggField,
		eventPrefix: "eventPrefix",
		env:         cfg.Env,
	}
}

func (d *DaoBase) GetEnv() common.IEnvConfig {
	return d.env
}

func (d *DaoBase) GetDbKey() string {
	return d.dbKey
}

func (d *DaoBase) GetTableName() string {
	return d.tableName
}

func (d *DaoBase) SetAggField(val string) {
	d.aggField = val
}

func (d *DaoBase) GetAggField() string {
	return d.aggField
}

func (d *DaoBase) GetIsPubEvent() bool {
	return d.isPubEvent
}

func (d *DaoBase) GetEventPrefix() string {
	return d.eventPrefix
}

func (d *DaoBase) GetSchema() *jsonschema.Schema {
	return d.cfg.Schema
}

func (d *DaoBase) Create(ctx context.Context, entity map[string]any, opts ...*db.CallOptions) {
	if entity == nil {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	d.dao.SetTenantId(entity, tenantId)
	err := d.dao.Insert(ctx, entity, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
	d.PublishEvent(ctx, db.AccessTypeCreate, entity, opts...)
}

func (d *DaoBase) CreateMany(ctx context.Context, entity []map[string]any, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	for _, e := range entity {
		e[TenantId] = tenantId
	}
	err := d.dao.InsertMany(ctx, entity, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) Update(ctx context.Context, entity map[string]any, opts ...*db.CallOptions) {
	if entity == nil {
		panic(fmt.Errorf("Dao.Update() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	d.dao.SetTenantId(entity, tenantId)
	err := d.dao.Update(ctx, entity, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}

	d.PublishEvent(ctx, db.AccessTypeUpdate, entity, opts...)
}

func (d *DaoBase) SoftDeleteById(ctx context.Context, id string, opts ...*db.CallOptions) {
	entity := map[string]any{}
	entity[Id] = id
	entity[TenantId] = d.GetTenantId(ctx)
	err := d.dao.Update(ctx, entity, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}

	d.PublishEvent(ctx, db.AccessTypeUpdate, entity, opts...)
}

func (d *DaoBase) DeleteById(ctx context.Context, id string, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)

	err := d.dao.DeleteById(ctx, tenantId, id, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}

	if d.GetIsPubEvent() {
		entity := map[string]any{
			TenantId: tenantId,
			Id:       id,
		}
		d.PublishEvent(ctx, db.AccessTypeDelete, entity, opts...)
	}
}

func (d *DaoBase) DeleteByIds(ctx context.Context, ids []string, opts ...*db.CallOptions) {
	if d.GetIsPubEvent() {
		for _, id := range ids {
			d.DeleteById(ctx, id, opts...)
		}
		return
	}

	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByIds(ctx, tenantId, ids, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}

	if d.GetIsPubEvent() {
		list := []map[string]any{}
		for _, id := range ids {
			e := map[string]any{}
			e[Id] = id
			e[TenantId] = tenantId
			list = append(list, e)
		}

		d.PublishBatchEvent(ctx, db.AccessTypeBatchDelete, list, opts...)
	}
}

func (d *DaoBase) UpdateByMap(ctx context.Context, filterMap map[string]any, data map[string]any, opts ...*db.CallOptions) {
	if d.GetIsPubEvent() {
		res := d.FindListByMap(ctx, filterMap, opts...)
		if res.Error != nil {
			panic(res.Error.Error())
		}
		for _, entity := range res.Data {
			err := d.dao.UpdateMapById(ctx, d.dao.GetTenantId(entity), d.dao.GetId(entity), data)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	tenantId := d.GetTenantId(ctx)
	err := d.dao.UpdateMap(ctx, tenantId, filterMap, data, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) UpdateMany(ctx context.Context, entities []map[string]any, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	for _, entity := range entities {
		entity[TenantId] = tenantId
	}
	err := d.dao.UpdateManyById(ctx, entities, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) UpdateManyByFilter(ctx context.Context, filterRSQL string, data interface{}, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.UpdateManyByFilter(ctx, tenantId, filterRSQL, data, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteAll(ctx context.Context, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteAll(ctx, tenantId, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteByFilter(ctx context.Context, filterRSQL string, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByFilter(ctx, tenantId, filterRSQL, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*db.CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByMap(ctx, tenantId, filterMap, db.NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) FindById(ctx context.Context, id string, opts ...*db.CallOptions) map[string]any {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindById(ctx, tenantId, id, db.NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.Data
}

func (d *DaoBase) FindByIds(ctx context.Context, ids []string, opts ...*db.CallOptions) []map[string]any {
	tenantId := d.GetTenantId(ctx)
	data, _, err := d.dao.FindByIds(ctx, tenantId, ids, db.NewRepositoryOptions(opts)...).Result()
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) FindByRSQL(ctx context.Context, rsql string, opts ...*db.CallOptions) []map[string]any {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindByRSQL(ctx, tenantId, rsql, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res.Data
}

func (d *DaoBase) FindAll(ctx context.Context, opts ...*db.CallOptions) *ddd_repository.FindListResult[map[string]any] {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindAll(ctx, tenantId, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindListByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*db.CallOptions) *ddd_repository.FindListResult[map[string]any] {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindListByMap(ctx, tenantId, filterMap, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindPaging(ctx context.Context, findPaging *ddd_repository.FindPagingQueryRequest, opts ...*db.CallOptions) *ddd_repository.FindPagingResult[map[string]any] {
	findQuery := d.NewFindPagingQuery(ctx, findPaging)
	res := d.dao.FindPaging(ctx, findQuery, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*db.CallOptions) *ddd_repository.FindPagingResult[map[string]any] {
	if qry == nil {
		panic(errors.New("FindAutoComplete query is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	qry.SetTenantId(tenantId)
	res := d.dao.FindAutoComplete(ctx, qry, db.NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*db.CallOptions) *ddd_repository.FindPagingResult[map[string]any] {
	if qry == nil {
		panic(errors.New("FindDistinctQueryRequest query is nil"))
	}
	data := d.dao.FindDistinct(ctx, qry, db.NewRepositoryOptions(opts)...)
	if data.GetError() != nil {
		panic(data.GetError())
	}
	return data
}

func (d *DaoBase) SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*db.CallOptions) []map[string]any {
	data, _, err := d.dao.SumEntity(ctx, qry, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) SumMap(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*db.CallOptions) []map[string]any {
	data, _, err := d.dao.SumMap(ctx, qry, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) Sum(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, data any, opts ...*db.CallOptions) any {
	data, _, err := d.dao.Sum(ctx, qry, data, db.NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) SumByRSQL(ctx context.Context, rSql string, valueCols []*ddd_repository.ValueCol, opts ...*db.CallOptions) map[string]any {
	opt := db.NewCallOptions(opts...)
	tenantId := d.GetTenantId(ctx)
	return d.dao.SumByRSQL(ctx, tenantId, rSql, valueCols, opt)
}

func (d *DaoBase) NewAggregateAndEvent(ctx context.Context, operateType db.AccessType, entity map[string]any, opts ...*db.CallOptions) (*server.Aggregate, *common.Event, error) {
	opt := db.NewCallOptions(opts...)
	event, err := d.NewEvent(ctx, operateType, entity, opt)
	if err != nil {
		return nil, nil, err
	}
	agg, err := d.NewAggregate(entity, opt)
	if err != nil {
		return nil, nil, err
	}
	return agg, event, nil
}

func (d *DaoBase) NewEvent(ctx context.Context, operateType db.AccessType, entity map[string]any, opt *db.CallOptions) (*common.Event, error) {
	o := db.NewCallOptions(opt)
	eventId := idutils.NewId()
	tenantId := d.dao.GetTenantId(entity)
	aggId, err := d.GetAggregateId(entity, opt)
	if err != nil {
		return nil, err
	}
	eventType := d.GetEventType(operateType, opt)

	event := common.NewEvent()
	event.CommandId = o.GetCommandId(idutils.NewId())
	event.EventId = eventId
	event.EventType = eventType
	event.TenantId = tenantId
	event.CreatedTime = time.Now()
	event.AggregateId = aggId
	event.Data = entity
	event.EventVersion = o.GetEventVersion("v1.0")

	return event, nil
}

func (d *DaoBase) GetEventType(accessType db.AccessType, opts *db.CallOptions) string {
	eventType := d.tableName
	if opts != nil && opts.EventType != nil {
		eventType = *opts.EventType
	}
	return common.GetEventType(d.appId, eventType, string(accessType))
}

func (d *DaoBase) CountByMap(ctx context.Context, filterData any, opts ...*db.CallOptions) int64 {
	opt := db.NewCallOptions(opts...)

	tenantId := d.GetTenantId(ctx)
	count, err := d.dao.CountByMap(ctx, tenantId, filterData, opt)
	if err != nil {
		panic(err)
	}
	return count
}

func (d *DaoBase) CountByRSQL(ctx context.Context, rsql string, opts ...*db.CallOptions) int64 {
	opt := db.NewCallOptions(opts...)
	tenantId := d.GetTenantId(ctx)
	count, err := d.dao.CountByRSQL(ctx, tenantId, rsql, opt)
	if err != nil {
		panic(err)
	}
	return count
}

func (d *DaoBase) NewAggregate(entity map[string]any, opt *db.CallOptions) (*server.Aggregate, error) {
	tenantId := d.dao.GetTenantId(entity)
	aggregateId, err := d.GetAggregateId(entity, opt)
	if err != nil {
		return nil, err
	}
	agg := common.NewAggregate()
	agg.TenantId = tenantId
	agg.AggregateId = aggregateId
	agg.AggregateVersion = "v1.0"
	agg.AggregateType = d.tableName
	return agg, nil
}

func (d *DaoBase) GetAggregateId(entity map[string]any, opts *db.CallOptions) (string, error) {
	var aggId string
	if opts != nil && opts.AggId != nil {
		aggId = *opts.AggId
	} else if d.aggField != "" {
		if id, ok := entity[d.aggField].(string); ok {
			aggId = id
		} else {
			return "", errors.New(fmt.Sprintf("Aggregate field %s is not string", d.aggField))
		}
	} else {
		aggId = d.dao.GetId(entity)
	}
	return aggId, nil
}

func (d *DaoBase) GetTenantId(ctx context.Context) string {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user.GetTenantId()
	}
	panic("token is error")
}

func (d *DaoBase) GetAuthUser(ctx context.Context) appctx.AuthUser {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user
	}
	panic("token is error")
}

func (d *DaoBase) PublishEvent(ctx context.Context, opeType db.AccessType, entity map[string]any, opts ...*db.CallOptions) {
	if !d.isPubEvent {
		return
	}

	agg, event, err := d.NewAggregateAndEvent(ctx, opeType, entity, opts...)
	if err != nil {
		panic(err)
	}
	logs.Debug(ctx, "", logs.Fields{
		"eventId":     event.EventId,
		"eventType":   event.EventType,
		"commandId":   event.CommandId,
		"aggregateId": event.AggregateId,
		"tenantId":    d.GetTenantId(ctx),
	})
	switch opeType {
	case db.AccessTypeCreate:
		server.GetEventPkg().CreateEvent(ctx, agg, event)
	case db.AccessTypeUpdate:
		server.GetEventPkg().ApplyEvent(ctx, agg, event)
	case db.AccessTypeDelete:
		server.GetEventPkg().ApplyEvent(ctx, agg, event)
	}
}

func (d *DaoBase) PublishBatchEvent(ctx context.Context, opeType db.AccessType, list []map[string]any, opts ...*db.CallOptions) {
	if !d.isPubEvent {
		return
	}

	/*
		agg, event, err := d.NewAggregateAndEvent(opeType, list, opts...)
		if err != nil {
			panic(err)
		}
		logs.Debug(ctx, "", logs.Fields{
			"eventId":     event.EventId,
			"eventType":   event.EventType,
			"commandId":   event.CommandId,
			"aggregateId": event.AggregateId,
			"tenantId":    d.GetTenantId(ctx),
		})
		switch opeType {
		case db.AccessTypeCreate:
			server.GetEventPkg().CreateEvent(ctx, agg, event)
		case db.AccessTypeUpdate:
			server.GetEventPkg().ApplyEvent(ctx, agg, event)
		case db.AccessTypeDelete:
			server.GetEventPkg().ApplyEvent(ctx, agg, event)
		}

	*/
}

func (d *DaoBase) NewFindPagingQuery(ctx context.Context, findPagingMap any) ddd_repository.FindPagingQuery {
	var qry ddd_repository.FindPagingQuery
	if v, ok := findPagingMap.(ddd_repository.FindPagingQuery); ok {
		qry = v
	} else if mapData, ok := findPagingMap.(map[string]any); ok {
		builder := ddd_repository.NewFindPagingQueryBuilder()
		qry = builder.SetMapToQuery(mapData).Build()
	} else if query, ok := findPagingMap.(ddd_repository.FindPagingQuery); ok {
		qry = query
	} else {
		panic("FindPaging(findPagingMap:any) findPagingMap is map[string]any or ddd_repository.FindPagingQuery ")
	}
	tenantId := d.GetTenantId(ctx)
	qry.SetTenantId(tenantId)
	return qry
}
