package mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dao/mongo_dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Dao struct {
	db        *DB
	tableName string
	dao       *ddd_mongodb.Dao[ddd.MapEntity]
	aggField  string // 聚合根ID字段
	daoType   DaoType
}

type ModelOptions = mongo_dao.RepositoryOptions

type DaoType string

const (
	DaoType_None  DaoType = ""
	DaoType_SQL   DaoType = "sql"
	DaoType_Event DaoType = "event"
)

type DaoOptions struct {
	DB        *DB
	TableName string
	DaoType   DaoType
	AggField  string
	// mongo
	MongoDB         *ddd_mongodb.MongoDB
	GetCollCallback mongo_dao.GetCollectionCallback
	RepositoryType  *mongo_dao.RepositoryType
}

func NewDao(opts *DaoOptions) *Dao {
	tableName := opts.TableName

	opt := mongo_dao.NewRepositoryOptions(&mongo_dao.RepositoryOptions{
		MongoDB:         opts.MongoDB,
		GetCollCallback: opts.GetCollCallback,
		RepositoryType:  opts.RepositoryType,
	})

	var mongodb *ddd_mongodb.MongoDB
	var coll *mongo.Collection

	getCollCallback := func(ctx context.Context) (*ddd_mongodb.MongoDB, *mongo.Collection) {
		if mongodb == nil || coll == nil {
			mongodb = opt.MongoDB
			coll = opt.MongoDB.GetCollection(tableName)
		}
		return mongodb, coll
	}

	if opt.GetCollCallback != nil {
		getCollCallback = opt.GetCollCallback
	}
	entBuilder := ddd.NewMapEntityBuilder[ddd.MapEntity]()
	daoOpts := ddd_mongodb.NewOptions[ddd.MapEntity]().SetAutoCreateCollection(true).SetAutoCreateIndex(true).SetEntityBuilder(entBuilder)

	daoType := DaoType_SQL
	if opts.DaoType == DaoType_Event {
		daoType = DaoType_Event
	}
	aggField := opts.AggField
	if aggField == "" {
		aggField = "id"
	}

	dao := ddd_mongodb.NewDao[ddd.MapEntity](getCollCallback, daoOpts)
	return &Dao{
		db:        opts.DB,
		dao:       dao,
		tableName: tableName,
		daoType:   daoType,
		aggField:  aggField,
	}
}

func (d *Dao) SetAggField(val string) *Dao {
	d.aggField = val
	return d
}

func (d *Dao) GetAggField(val string) string {
	return d.aggField
}

func (d *Dao) Table(ctx context.Context, schema *schema.Schema, opts ...*ddd_repository.RepositoryOptions) *Table {
	return NewTable(d.db, d.tableName, schema)
}

/*
func (d *Dao) Save(ctx context.Context, setData *ddd.SetData[ddd.MapEntity], opts ...*OperateOptions) {
	err := d.dao.Save(ctx, setData, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}*/

func (d *Dao) Create(ctx context.Context, entity ddd.MapEntity, opts ...*OperateOptions) {
	if entity == nil {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	tenantId := d.getTenantId(ctx)
	entity.SetTenantId(tenantId)
	err := d.dao.Insert(ctx, entity, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
	if d.daoType == DaoType_Event {
		d.pub(ctx, OperateType_Create, entity, opts...)
	}
}

func (d *Dao) pub(ctx context.Context, opeType OperateType, entity ddd.MapEntity, opts ...*OperateOptions) {
	agg, event, err := d.newAggregateAndEvent(opeType, entity, opts...)
	if err != nil {
		panic(err)
	}
	logs.Info(ctx, "", logs.Fields{"eventId": event.EventId, "eventType": event.EventType, "commandId": event.CommandId, "aggregateId": event.AggregateId, "tenantId": d.getTenantId(ctx)})
	switch opeType {
	case OperateType_Create:
		server.GetEventPkg().CreateEvent(ctx, agg, event)
	case OperateType_Update:
		server.GetEventPkg().ApplyEvent(ctx, agg, event)
	case OperateType_Delete:
		server.GetEventPkg().ApplyEvent(ctx, agg, event)
	}

}

func (d *Dao) Update(ctx context.Context, entity ddd.MapEntity, opts ...*OperateOptions) {
	if entity == nil {
		panic(fmt.Errorf("Dao.Update() entity is nil"))
	}
	tenantId := d.getTenantId(ctx)
	entity.SetTenantId(tenantId)
	err := d.dao.Update(ctx, entity, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
	if d.daoType == DaoType_Event {
		d.pub(ctx, OperateType_Update, entity, opts...)
	}
}

func (d *Dao) DeleteById(ctx context.Context, id string, opts ...*OperateOptions) {
	tenantId := d.getTenantId(ctx)
	err := d.dao.DeleteById(ctx, tenantId, id, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}

	if d.daoType == DaoType_Event {
		entity := ddd.MapEntity{
			"tenantId": tenantId,
			"id":       id,
		}
		d.pub(ctx, OperateType_Delete, entity, opts...)
	}
}

func (d *Dao) CreateMany(ctx context.Context, entity []ddd.MapEntity, opts ...*OperateOptions) {
	if d.daoType == DaoType_Event {
		for _, entity := range entity {
			d.Create(ctx, entity, opts...)
		}
		return
	}

	err := d.dao.InsertMany(ctx, entity, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) DeleteByIds(ctx context.Context, ids []string, opts ...*OperateOptions) {
	if d.daoType == DaoType_Event {
		for _, id := range ids {
			d.DeleteById(ctx, id, opts...)
		}
		return
	}

	tenantId := d.getTenantId(ctx)
	err := d.dao.DeleteByIds(ctx, tenantId, ids, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *Dao) UpdateByMap(ctx context.Context, filterMap map[string]any, data map[string]any, opts ...*OperateOptions) {
	if d.daoType == DaoType_Event {
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

	tenantId := d.getTenantId(ctx)
	err := d.dao.UpdateMap(ctx, tenantId, filterMap, data, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *Dao) UpdateMany(ctx context.Context, entities []ddd.MapEntity, opts ...*OperateOptions) {
	err := d.dao.UpdateManyById(ctx, entities, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*OperateOptions) *ddd_repository.BulkWriteResult {
	res, err := d.dao.BulkWrite(ctx, models, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return res
}

func (d *Dao) UpdateManyByFilter(ctx context.Context, filter string, data interface{}, opts ...*OperateOptions) {
	tenantId := d.getTenantId(ctx)
	err := d.dao.UpdateManyByFilter(ctx, tenantId, filter, data, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) DeleteAll(ctx context.Context, opts ...*OperateOptions) {
	tenantId := d.getTenantId(ctx)
	err := d.dao.DeleteAll(ctx, tenantId, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) DeleteByFilter(ctx context.Context, filter string, opts ...*OperateOptions) {
	tenantId := d.getTenantId(ctx)
	err := d.dao.DeleteByFilter(ctx, tenantId, filter, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *Dao) DeleteByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*OperateOptions) {
	tenantId := d.getTenantId(ctx)
	err := d.dao.DeleteByMap(ctx, tenantId, filterMap, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) FindById(ctx context.Context, id string, opts ...*OperateOptions) *common.Result[ddd.MapEntity] {
	tenantId := d.getTenantId(ctx)
	res := d.dao.FindById(ctx, tenantId, id, newOptions(opts)...)
	if res.Error != nil {
		panic(res.Error)
	}
	return common.NewResult[ddd.MapEntity](res.Data, res.Error)
}

func (d *Dao) FindByIds(ctx context.Context, ids []string, opts ...*OperateOptions) *common.Result[[]ddd.MapEntity] {
	tenantId := d.getTenantId(ctx)
	data, _, err := d.dao.FindByIds(ctx, tenantId, ids, newOptions(opts)...).Result()
	if err != nil {
		panic(err)
	}
	return common.NewResult[[]ddd.MapEntity](data, err)
}

func (d *Dao) FindAll(ctx context.Context, opts ...*OperateOptions) *ddd_repository.FindListResult[ddd.MapEntity] {
	tenantId := d.getTenantId(ctx)
	res := d.dao.FindAll(ctx, tenantId, newOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *Dao) FindListByMap(ctx context.Context, filterMap map[string]interface{}, opts ...*OperateOptions) *ddd_repository.FindListResult[ddd.MapEntity] {
	tenantId := d.getTenantId(ctx)
	res := d.dao.FindListByMap(ctx, tenantId, filterMap, newOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *Dao) FindPaging(ctx context.Context, findPagingMap any, opts ...*OperateOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	var findQuery ddd_repository.FindPagingQuery
	if mapData, ok := findPagingMap.(map[string]any); ok {
		builder := ddd_repository.NewFindPagingQueryBuilder()
		findQuery = builder.SetMapToQuery(mapData).Build()
	} else if query, ok := findPagingMap.(ddd_repository.FindPagingQuery); ok {
		findQuery = query
	} else {
		panic("FindPaging(findPagingMap:any) findPagingMap is map[string]any or ddd_repository.FindPagingQuery ")
	}
	tenantId := d.getTenantId(ctx)
	findQuery.SetTenantId(tenantId)
	res := d.dao.FindPaging(ctx, findQuery, newOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *Dao) FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*OperateOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	if qry == nil {
		panic(errors.New("FindAutoComplete query is nil"))
	}
	tenantId := d.getTenantId(ctx)
	qry.SetTenantId(tenantId)
	res := d.dao.FindAutoComplete(ctx, qry, newOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *Dao) FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*OperateOptions) *common.Result[*ddd_repository.FindPagingResult[ddd.MapEntity]] {
	if qry == nil {
		panic(errors.New("FindDistinctQueryRequest query is nil"))
	}
	data := d.dao.FindDistinct(ctx, qry, newOptions(opts)...)
	if data.GetError() != nil {
		panic(data.GetError())
	}
	return common.NewResult[*ddd_repository.FindPagingResult[ddd.MapEntity]](data, nil)
}

func (d *Dao) AggregateByPipeline(ctx context.Context, pipeline mongo.Pipeline, data interface{}) {
	err := d.dao.AggregateByPipeline(ctx, pipeline, data)
	if err != nil {
		panic(err)
	}
}

func (d *Dao) getTenantId(ctx context.Context) string {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user.GetTenantId()
	}
	panic("token is error")
}

func (d *Dao) getAuthUser(ctx context.Context) appctx.AuthUser {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user
	}
	panic("token is error")
}

func (d *Dao) SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*OperateOptions) *common.Result[[]ddd.MapEntity] {
	data, _, err := d.dao.SumEntity(ctx, qry, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return common.NewResult[[]ddd.MapEntity](data, err)
}

func (d *Dao) SumMap(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*OperateOptions) *common.Result[[]map[string]any] {
	data, _, err := d.dao.SumMap(ctx, qry, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return common.NewResult[[]map[string]any](data, err)
}

func (d *Dao) Sum(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, data any, opts ...*OperateOptions) *common.Result[any] {
	data, _, err := d.dao.Sum(ctx, qry, data, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
	return common.NewResult[any](data, err)
}

func (d *Dao) GetFilterMap(tenantId string, rsql string) *common.Result[ddd.MapEntity] {
	data, err := d.dao.GetFilterMap(tenantId, rsql)
	if err != nil {
		panic(err)
	}
	return common.NewResult[ddd.MapEntity](data, err)
}

func (d *Dao) newAggregateAndEvent(operateType OperateType, entity ddd.MapEntity, opts ...*OperateOptions) (*server.Aggregate, *common.Event, error) {
	opt := NewOperateOptions(opts...)
	event, err := d.newEvent(operateType, entity, opt)
	if err != nil {
		return nil, nil, err
	}
	agg, err := d.newAggregate(entity, opt)
	if err != nil {
		return nil, nil, err
	}
	return agg, event, nil
}

func (d *Dao) newEvent(operateType OperateType, entity ddd.MapEntity, opt *OperateOptions) (*common.Event, error) {
	o := NewOperateOptions(opt)
	eventId := idutils.NewId()
	tenantId := entity.GetTenantId()
	aggId, err := d.getAggregateId(entity, opt)
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

type OperateType string

const (
	OperateType_Create OperateType = "create"
	OperateType_Update OperateType = "update"
	OperateType_Delete OperateType = "delete"
)

func (d *Dao) GetEventType(operateType OperateType, opt *OperateOptions) string {
	eventType := d.tableName
	if opt == nil && opt.EventType != nil {
		eventType = *opt.EventType
	}
	return common.GetEventType(d.db.cfg.GetAppId(), eventType, string(operateType))
}

func (d *Dao) newAggregate(entity ddd.MapEntity, opt *OperateOptions) (*server.Aggregate, error) {
	tenantId := entity.GetTenantId()
	aggregateId, err := d.getAggregateId(entity, opt)
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

func (d *Dao) getAggregateId(entity ddd.MapEntity, opts *OperateOptions) (string, error) {
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
