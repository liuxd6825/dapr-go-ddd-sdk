package mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dao/mongo_dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Dao struct {
	db             *DB
	tableName      string
	dao            *ddd_mongodb.Dao[ddd.MapEntity]
	aggregateField string
}

type ModelOptions = mongo_dao.RepositoryOptions

func NewDao(db *DB, tableName string, opts ...*ModelOptions) *Dao {
	initTableName := tableName
	opt := mongo_dao.NewRepositoryOptions(opts...)
	var mongodb *ddd_mongodb.MongoDB
	var coll *mongo.Collection

	getCollCallback := func(ctx context.Context) (*ddd_mongodb.MongoDB, *mongo.Collection) {
		if mongodb == nil || coll == nil {
			mongodb = opt.MongoDB
			coll = opt.MongoDB.GetCollection(initTableName)
		}
		return mongodb, coll
	}

	if opt.GetCollCallback != nil {
		getCollCallback = opt.GetCollCallback
	}
	entBuilder := ddd.NewMapEntityBuilder[ddd.MapEntity]()
	daoOpts := ddd_mongodb.NewOptions[ddd.MapEntity]().SetAutoCreateCollection(true).SetAutoCreateIndex(true).SetEntityBuilder(entBuilder)

	dao := ddd_mongodb.NewDao[ddd.MapEntity](getCollCallback, daoOpts)
	return &Dao{db: db, dao: dao, tableName: tableName}
}

func (d *Dao) SetAggregateField(val string) *Dao {
	d.aggregateField = val
	return d
}

func (d *Dao) Table(ctx context.Context, schema *schema.Schema, opts ...*ddd_repository.RepositoryOptions) *Table {
	return NewTable(d.db, d.tableName, schema)
}

func (d *Dao) Save(ctx context.Context, setData *ddd.SetData[ddd.MapEntity], opts ...*OperateOptions) {
	err := d.dao.Save(ctx, setData, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) Create(ctx context.Context, entity ddd.MapEntity, opts ...*OperateOptions) {
	if entity == nil {
		panic(fmt.Errorf("Dao.Create() entity is nil"))
	}
	err := d.dao.Insert(ctx, entity, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
	return

	agg, event, err := d.newAggregateAndEvent(OperateType_Create, entity, opts...)
	if err != nil {
		panic(err)
	}
	server.GetEventPkg().CreateEvent(ctx, agg, event)
}

func (d *Dao) Update(ctx context.Context, entity ddd.MapEntity, opts ...*OperateOptions) {
	if entity == nil {
		panic(fmt.Errorf("Dao.Update() entity is nil"))
	}

	err := d.dao.Update(ctx, entity, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
	return

	agg, event, err := d.newAggregateAndEvent(OperateType_Update, entity, opts...)
	if err != nil {
		panic(err)
	}
	server.GetEventPkg().ApplyEvent(ctx, agg, event)
}

func (d *Dao) DeleteById(ctx context.Context, tenantId string, id string, opts ...*OperateOptions) {
	err := d.dao.DeleteById(ctx, tenantId, id, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
	return

	entity := ddd.MapEntity{
		"tenantId": tenantId,
		"id":       id,
	}
	agg, event, err := d.newAggregateAndEvent(OperateType_Delete, entity, opts...)
	if err != nil {
		panic(err)
	}
	server.GetEventPkg().ApplyEvent(ctx, agg, event)
}

func (d *Dao) CreateMany(ctx context.Context, entity []ddd.MapEntity, opts ...*OperateOptions) {
	err := d.dao.InsertMany(ctx, entity, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...*OperateOptions) {
	err := d.dao.DeleteByIds(ctx, tenantId, ids, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *Dao) UpdateByMap(ctx context.Context, tenantId string, filterMap map[string]any, data any, opts ...*OperateOptions) {
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

func (d *Dao) UpdateManyByFilter(ctx context.Context, tenantId, filter string, data interface{}, opts ...*OperateOptions) {
	err := d.dao.UpdateManyByFilter(ctx, tenantId, filter, data, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) DeleteAll(ctx context.Context, tenantId string, opts ...*OperateOptions) {
	err := d.dao.DeleteAll(ctx, tenantId, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...*OperateOptions) {
	err := d.dao.DeleteByFilter(ctx, tenantId, filter, newOptions(opts)...)
	if err != nil {
		panic(err)
	}
}

func (d *Dao) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...*OperateOptions) {
	err := d.dao.DeleteByMap(ctx, tenantId, filterMap, newOptions(opts)...).GetError()
	if err != nil {
		panic(err)
	}
}

func (d *Dao) FindById(ctx context.Context, tenantId string, id string, opts ...*OperateOptions) *common.Result[ddd.MapEntity] {
	res := d.dao.FindById(ctx, tenantId, id, newOptions(opts)...)
	if res.Err != nil {
		panic(res.Err)
	}
	return common.NewResult[ddd.MapEntity](res.Data, res.Err)
}

func (d *Dao) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...*OperateOptions) *common.Result[[]ddd.MapEntity] {
	data, _, err := d.dao.FindByIds(ctx, tenantId, ids, newOptions(opts)...).Result()
	if err != nil {
		panic(err)
	}
	return common.NewResult[[]ddd.MapEntity](data, err)
}

func (d *Dao) FindAll(ctx context.Context, tenantId string, opts ...*OperateOptions) *ddd_repository.FindListResult[ddd.MapEntity] {
	res := d.dao.FindAll(ctx, tenantId, newOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *Dao) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...*OperateOptions) *ddd_repository.FindListResult[ddd.MapEntity] {
	res := d.dao.FindListByMap(ctx, tenantId, filterMap, newOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *Dao) FindPaging(ctx context.Context, query *ddd_repository.FindPagingQueryRequest, opts ...*OperateOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	res := d.dao.FindPaging(ctx, query, newOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *Dao) FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*OperateOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	res := d.dao.FindAutoComplete(ctx, qry, newOptions(opts)...)
	if res.GetError() != nil {
		panic(res.GetError())
	}
	return res
}

func (d *Dao) FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*OperateOptions) *common.Result[*ddd_repository.FindPagingResult[ddd.MapEntity]] {
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
	o := opt
	if o == nil {
		o = &OperateOptions{}
	}
	eventId := entity.GetId()
	tenantId := entity.GetTenantId()
	aggId, err := d.getAggregateId(entity, opt)
	if err != nil {
		return nil, err
	}
	eventType := d.GetEventType(operateType, opt)

	event := common.NewEvent()
	event.EventId = eventId
	event.EventType = eventType
	event.TenantId = tenantId
	event.CreatedTime = time.Now()
	event.AggregateId = aggId
	event.Data = entity
	event.CommandId = eventId
	event.EventVersion = o.GetVersion("v1.0")

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
	var aggregateId string
	if opts != nil && opts.AggId != nil {
		aggregateId = *opts.AggId
	} else if d.aggregateField != "" {
		if id, ok := entity[d.aggregateField].(string); ok {
			aggregateId = id
		} else {
			return "", errors.New(fmt.Sprintf("Aggregate field %s is not string", d.aggregateField))
		}
	} else {
		aggregateId = entity.GetId()
	}
	return aggregateId, nil
}
