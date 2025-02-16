package db

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"time"
)

type DaoBase struct {
	cfg         *DaoConfig
	dbKey       string                            // 配置中的数据库Key
	tableName   string                            // 表名
	appId       string                            // 应用ID
	aggField    string                            // 聚合根字段
	isPubEvent  bool                              // 是否发布事件
	eventPrefix string                            // 事件前缀
	dao         ddd_repository.Dao[ddd.MapEntity] // 数据访问
	env         common.IEnvConfig                 // 环境变量
}

func NewDaoBase(dao ddd_repository.Dao[ddd.MapEntity], cfg *DaoConfig) *DaoBase {
	if cfg == nil {
		panic("dao base config is nil")
	}
	aggField := cfg.AggField
	if aggField == "" {
		aggField = "id"
	}
	return &DaoBase{
		dao:         dao,
		dbKey:       cfg.DbKey,
		tableName:   cfg.TableName,
		appId:       restapp.GetAppId(),
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

func (d *DaoBase) GetSchema() *schema.Schema {
	return d.cfg.Schema
}

func (d *DaoBase) Create(ctx context.Context, entity ddd.MapEntity, opts ...*CallOptions) {
	if entity == nil {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	entity.SetTenantId(tenantId)
	err := d.dao.Insert(ctx, entity, NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
	d.PublishEvent(ctx, AccessTypeCreate, entity, opts...)
}

func (d *DaoBase) Update(ctx context.Context, entity ddd.MapEntity, opts ...*CallOptions) {
	if entity == nil {
		panic(fmt.Errorf("Dao.Update() entity is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	entity.SetTenantId(tenantId)
	err := d.dao.Update(ctx, entity, NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}

	d.PublishEvent(ctx, AccessTypeUpdate, entity, opts...)
}

func (d *DaoBase) DeleteById(ctx context.Context, id string, opts ...*CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteById(ctx, tenantId, id, NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}

	if d.GetIsPubEvent() {
		entity := ddd.MapEntity{
			"tenantId": tenantId,
			"id":       id,
		}
		d.PublishEvent(ctx, AccessTypeDelete, entity, opts...)
	}
}

func (d *DaoBase) CreateMany(ctx context.Context, entity []ddd.MapEntity, opts ...*CallOptions) {
	if d.GetIsPubEvent() {
		for _, entity := range entity {
			d.Create(ctx, entity, opts...)
		}
		return
	}

	err := d.dao.InsertMany(ctx, entity, NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteByIds(ctx context.Context, ids []string, opts ...*CallOptions) {
	if d.GetIsPubEvent() {
		for _, id := range ids {
			d.DeleteById(ctx, id, opts...)
		}
		return
	}

	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByIds(ctx, tenantId, ids, NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) UpdateByMap(ctx context.Context, filterMap map[string]any, data map[string]any, opts ...*CallOptions) {
	if d.GetIsPubEvent() {
		res := d.FindListByMap(ctx, filterMap, opts...)
		if res.Error != nil {
			panic(res.Error.Error())
		}
		for _, entity := range res.Data {
			err := d.dao.UpdateMapById(ctx, entity.GetTenantId(), entity.GetId(), data)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	tenantId := d.GetTenantId(ctx)
	err := d.dao.UpdateMap(ctx, tenantId, filterMap, data, NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) UpdateMany(ctx context.Context, entities []ddd.MapEntity, opts ...*CallOptions) {
	err := d.dao.UpdateManyById(ctx, entities, NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) UpdateManyByFilter(ctx context.Context, filter string, data interface{}, opts ...*CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.UpdateManyByFilter(ctx, tenantId, filter, data, NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteAll(ctx context.Context, opts ...*CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteAll(ctx, tenantId, NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteByFilter(ctx context.Context, filter string, opts ...*CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByFilter(ctx, tenantId, filter, NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) DeleteByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*CallOptions) {
	tenantId := d.GetTenantId(ctx)
	err := d.dao.DeleteByMap(ctx, tenantId, filterMap, NewRepositoryOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *DaoBase) FindById(ctx context.Context, id string, opts ...*CallOptions) ddd.MapEntity {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindById(ctx, tenantId, id, NewRepositoryOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.Data
}

func (d *DaoBase) FindByIds(ctx context.Context, ids []string, opts ...*CallOptions) []ddd.MapEntity {
	tenantId := d.GetTenantId(ctx)
	data, _, err := d.dao.FindByIds(ctx, tenantId, ids, NewRepositoryOptions(opts)...).Result()
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) FindAll(ctx context.Context, opts ...*CallOptions) *ddd_repository.FindListResult[ddd.MapEntity] {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindAll(ctx, tenantId, NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindListByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*CallOptions) *ddd_repository.FindListResult[ddd.MapEntity] {
	tenantId := d.GetTenantId(ctx)
	res := d.dao.FindListByMap(ctx, tenantId, filterMap, NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindPaging(ctx context.Context, findPaging *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	findQuery := d.NewFindPagingQuery(ctx, findPaging)
	res := d.dao.FindPaging(ctx, findQuery, NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	if qry == nil {
		panic(errors.New("FindAutoComplete query is nil"))
	}
	tenantId := d.GetTenantId(ctx)
	qry.SetTenantId(tenantId)
	res := d.dao.FindAutoComplete(ctx, qry, NewRepositoryOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *DaoBase) FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*CallOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	if qry == nil {
		panic(errors.New("FindDistinctQueryRequest query is nil"))
	}
	data := d.dao.FindDistinct(ctx, qry, NewRepositoryOptions(opts)...)
	if data.GetError() != nil {
		panic(data.GetError())
	}
	return data
}

func (d *DaoBase) SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) []ddd.MapEntity {
	data, _, err := d.dao.SumEntity(ctx, qry, NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) SumMap(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*CallOptions) []map[string]any {
	data, _, err := d.dao.SumMap(ctx, qry, NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) Sum(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, data any, opts ...*CallOptions) any {
	data, _, err := d.dao.Sum(ctx, qry, data, NewRepositoryOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return data
}

func (d *DaoBase) NewAggregateAndEvent(operateType AccessType, entity ddd.MapEntity, opts ...*CallOptions) (*server.Aggregate, *common.Event, error) {
	opt := NewCallOptions(opts...)
	event, err := d.NewEvent(operateType, entity, opt)
	if err != nil {
		return nil, nil, err
	}
	agg, err := d.NewAggregate(entity, opt)
	if err != nil {
		return nil, nil, err
	}
	return agg, event, nil
}

func (d *DaoBase) NewEvent(operateType AccessType, entity ddd.MapEntity, opt *CallOptions) (*common.Event, error) {
	o := NewCallOptions(opt)
	eventId := idutils.NewId()
	tenantId := entity.GetTenantId()
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

func (d *DaoBase) GetEventType(accessType AccessType, opts *CallOptions) string {
	eventType := d.tableName
	if opts != nil && opts.EventType != nil {
		eventType = *opts.EventType
	}
	return common.GetEventType(d.appId, eventType, string(accessType))
}

func (d *DaoBase) NewAggregate(entity ddd.MapEntity, opt *CallOptions) (*server.Aggregate, error) {
	tenantId := entity.GetTenantId()
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

func (d *DaoBase) GetAggregateId(entity ddd.MapEntity, opts *CallOptions) (string, error) {
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
		aggId = entity.GetId()
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

func (d *DaoBase) PublishEvent(ctx context.Context, opeType AccessType, entity ddd.MapEntity, opts ...*CallOptions) {
	if !d.isPubEvent {
		return
	}

	agg, event, err := d.NewAggregateAndEvent(opeType, entity, opts...)
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
	case AccessTypeCreate:
		server.GetEventPkg().CreateEvent(ctx, agg, event)
	case AccessTypeUpdate:
		server.GetEventPkg().ApplyEvent(ctx, agg, event)
	case AccessTypeDelete:
		server.GetEventPkg().ApplyEvent(ctx, agg, event)
	}
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
