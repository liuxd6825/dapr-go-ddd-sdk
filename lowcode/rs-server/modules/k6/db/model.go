package db

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dao/mongo_dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Model struct {
	db        *DB
	tableName string
	dao       *ddd_mongodb.Dao[ddd.MapEntity]
}

type ModelOptions = mongo_dao.RepositoryOptions

func NewModel(db *DB, tableName string, opts ...*ModelOptions) *Model {
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
	fmt.Print(opts)
	dao := ddd_mongodb.NewDao[ddd.MapEntity](getCollCallback, daoOpts)
	return &Model{db: db, dao: dao, tableName: tableName}
}

func (d *Model) Table(ctx context.Context, schema *schema.Schema, opts ...*ddd_repository.RepositoryOptions) *Table {
	return NewTable(d.db, d.tableName, schema)
}

func (d *Model) Save(ctx context.Context, setData *ddd.SetData[ddd.MapEntity], opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.Save(ctx, setData, newOptions(opts)...).GetError()
}

type EventOptions struct {
	AggId *string
}

func (d *Model) Create(ctx context.Context, entity ddd.MapEntity, eventOpt *EventOptions, opts ...*ddd_repository.RepositoryOptions) error {
	err := d.dao.Insert(ctx, entity, newOptions(opts)...).GetError()
	if err != nil {
		return err
	}

	event := d.newEvent(entity, eventOpt)
	agg := d.newAggregate(entity, eventOpt)
	res := server.GetEventPkg().CreateEvent(ctx, agg, event)
	return res.Error
}

func (d *Model) newEvent(entity ddd.MapEntity, eventOpt *EventOptions) *common.Event {
	eventType := fmt.Sprintf("%s.%s_create", "restapp.GetAppId()", d.tableName)
	eventId := entity.GetId()
	tenantId := entity.GetTenantId()
	aggId := d.getAggId(entity, eventOpt)

	event := common.NewEvent()
	event.EventId = eventId
	event.EventType = eventType
	event.TenantId = tenantId
	event.CreatedTime = time.Now()
	event.AggregateId = aggId
	event.Data = entity
	event.CommandId = eventId
	event.EventVersion = "v1.0"
	return event
}

func (d *Model) newAggregate(entity ddd.MapEntity, eventOpt *EventOptions) *server.Aggregate {
	tenantId := entity.GetTenantId()
	aggregateId := d.getAggId(entity, eventOpt)
	agg := server.NewAggregate()
	agg.TenantId = tenantId
	agg.AggregateId = aggregateId
	agg.AggregateVersion = "v1.0"
	agg.AggregateType = d.tableName
	return agg
}

func (d *Model) getAggId(entity ddd.MapEntity, eventOpt *EventOptions) string {
	aggregateId := entity.GetId()
	if eventOpt != nil {
		if eventOpt.AggId != nil {
			aggregateId = *eventOpt.AggId
		}
	}
	return aggregateId
}

func (d *Model) CreateByMap(ctx context.Context, tenantId string, data map[string]any, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.InsertMap(ctx, tenantId, data, newOptions(opts)...)
}

func (d *Model) CreateMany(ctx context.Context, entity []ddd.MapEntity, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.InsertMany(ctx, entity, newOptions(opts)...).GetError()
}

func (d *Model) Update(ctx context.Context, entity ddd.MapEntity, eventOpt *EventOptions, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.Update(ctx, entity, newOptions(opts)...).GetError()
}

func (d *Model) UpdateByMap(ctx context.Context, tenantId string, filterMap map[string]any, data any, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.UpdateMap(ctx, tenantId, filterMap, data, newOptions(opts)...)
}

func (d *Model) UpdateMany(ctx context.Context, entities []ddd.MapEntity, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.UpdateManyById(ctx, entities, newOptions(opts)...).GetError()
}

func (d *Model) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*ddd_repository.RepositoryOptions) (*ddd_repository.BulkWriteResult, error) {
	return d.dao.BulkWrite(ctx, models, newOptions(opts)...)
}

func (d *Model) UpdateManyByFilter(ctx context.Context, tenantId, filter string, data interface{}, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.UpdateManyByFilter(ctx, tenantId, filter, data, newOptions(opts)...).GetError()
}

func (d *Model) DeleteById(ctx context.Context, tenantId string, id string, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.DeleteById(ctx, tenantId, id, newOptions(opts)...).GetError()
}

func (d *Model) DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.DeleteByIds(ctx, tenantId, ids, newOptions(opts)...)
}

func (d *Model) DeleteAll(ctx context.Context, tenantId string, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.DeleteAll(ctx, tenantId, newOptions(opts)...).GetError()
}

func (d *Model) DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.DeleteByFilter(ctx, tenantId, filter, newOptions(opts)...)
}

func (d *Model) DeleteByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...*ddd_repository.RepositoryOptions) error {
	return d.dao.DeleteByMap(ctx, tenantId, filterMap, newOptions(opts)...).GetError()
}

func (d *Model) FindById(ctx context.Context, tenantId string, id string, opts ...*ddd_repository.RepositoryOptions) *common.Result[ddd.MapEntity] {
	res := d.dao.FindById(ctx, tenantId, id, newOptions(opts)...)
	return common.NewResult[ddd.MapEntity](res.Data, res.Err)
}

func (d *Model) FindByIds(ctx context.Context, tenantId string, ids []string, opts ...*ddd_repository.RepositoryOptions) *common.Result[[]ddd.MapEntity] {
	data, _, err := d.dao.FindByIds(ctx, tenantId, ids, newOptions(opts)...).Result()
	return common.NewResult[[]ddd.MapEntity](data, err)
}

func (d *Model) FindAll(ctx context.Context, tenantId string, opts ...*ddd_repository.RepositoryOptions) *ddd_repository.FindListResult[ddd.MapEntity] {
	return d.dao.FindAll(ctx, tenantId, newOptions(opts)...)
}

func (d *Model) FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...*ddd_repository.RepositoryOptions) *ddd_repository.FindListResult[ddd.MapEntity] {
	return d.dao.FindListByMap(ctx, tenantId, filterMap, newOptions(opts)...)
}

func (d *Model) FindPaging(ctx context.Context, query *ddd_repository.FindPagingQueryRequest, opts ...*ddd_repository.RepositoryOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	return d.dao.FindPaging(ctx, query, newOptions(opts)...)
}

func (d *Model) FindAutoComplete(ctx context.Context, qry *ddd_repository.FindAutoCompleteQueryRequest, opts ...*ddd_repository.RepositoryOptions) *ddd_repository.FindPagingResult[ddd.MapEntity] {
	return d.dao.FindAutoComplete(ctx, qry, newOptions(opts)...)
}

func (d *Model) FindDistinct(ctx context.Context, qry *ddd_repository.FindDistinctQueryRequest, opts ...*ddd_repository.RepositoryOptions) *common.Result[*ddd_repository.FindPagingResult[ddd.MapEntity]] {
	data := d.dao.FindDistinct(ctx, qry, newOptions(opts)...)
	return common.NewResult[*ddd_repository.FindPagingResult[ddd.MapEntity]](data, nil)
}

func (d *Model) AggregateByPipeline(ctx context.Context, pipeline mongo.Pipeline, data interface{}) error {
	return d.dao.AggregateByPipeline(ctx, pipeline, data)
}

func (d *Model) SumEntity(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*ddd_repository.RepositoryOptions) *common.Result[[]ddd.MapEntity] {
	data, _, err := d.dao.SumEntity(ctx, qry, newOptions(opts)...)
	return common.NewResult[[]ddd.MapEntity](data, err)
}

func (d *Model) SumMap(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, opts ...*ddd_repository.RepositoryOptions) *common.Result[[]map[string]any] {
	data, _, err := d.dao.SumMap(ctx, qry, newOptions(opts)...)
	return common.NewResult[[]map[string]any](data, err)
}

func (d *Model) Sum(ctx context.Context, qry *ddd_repository.FindPagingQueryRequest, data any, opts ...*ddd_repository.RepositoryOptions) *common.Result[any] {
	data, _, err := d.dao.Sum(ctx, qry, data, newOptions(opts)...)
	return common.NewResult[any](data, err)
}

func (d *Model) GetFilterMap(tenantId string, rsqlstr string) *common.Result[ddd.MapEntity] {
	data, err := d.dao.GetFilterMap(tenantId, rsqlstr)
	return common.NewResult[ddd.MapEntity](data, err)
}

func newOptions(opts []*ddd_repository.RepositoryOptions) []ddd_repository.Options {
	var res []ddd_repository.Options
	for _, o := range opts {
		res = append(res, o)
	}
	return res
}
