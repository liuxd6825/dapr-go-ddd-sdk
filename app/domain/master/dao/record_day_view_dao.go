package dao

import (
	"context"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/view"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type IRecordDayViewDao interface {
	Create(ctx context.Context, view *view.RecordDayView, opts ...idao.CallOptions) error
	CreateMany(ctx context.Context, views []*view.RecordDayView, opts ...idao.CallOptions) error

	Update(ctx context.Context, view *view.RecordDayView, opts ...idao.CallOptions) error
	UpdateMany(ctx context.Context, views []*view.RecordDayView, opts ...idao.CallOptions) error

	Save(ctx context.Context, view *view.RecordDayView, opts ...idao.CallOptions) error
	SaveMany(ctx context.Context, views []*view.RecordDayView, opts ...idao.CallOptions) error

	IncAmountMany(ctx context.Context, views []*view.RecordDayView, opts ...idao.CallOptions) error

	Delete(ctx context.Context, view *view.RecordDayView, opts ...idao.CallOptions) error
	DeleteMany(ctx context.Context, tenantId string, views []*view.RecordDayView, opts ...idao.CallOptions) error
	DeleteById(ctx context.Context, tenantId string, id string, opts ...idao.CallOptions) error
	DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...idao.CallOptions) error
	DeleteByFilter(ctx context.Context, tenantId, filter string, opts ...idao.CallOptions) error
	DeleteAll(ctx context.Context, tenantId string, opts ...idao.CallOptions) error
	DeleteByGraphId(ctx context.Context, tenantId string, graphId string, opts ...idao.CallOptions) error

	FindById(ctx context.Context, tenantId string, id string, opts ...idao.CallOptions) (*view.RecordDayView, bool, error)
	FindByIds(ctx context.Context, tenantId string, ids []string, opts ...idao.CallOptions) ([]*view.RecordDayView, bool, error)
	FindAll(ctx context.Context, tenantId string, opts ...idao.CallOptions) ([]*view.RecordDayView, bool, error)
	FindPaging(ctx context.Context, query idao.FindPagingQuery, opts ...idao.CallOptions) (*idao.FindPagingResult[*view.RecordDayView], bool, error)
	FindByGraphId(ctx context.Context, tenantId string, graphId string, opts ...idao.CallOptions) ([]*view.RecordDayView, bool, error)
	FindByCaseId(ctx context.Context, tenantId string, graphId string, opts ...idao.CallOptions) ([]*view.RecordDayView, bool, error)

	FindRecordSumChart(ctx context.Context, tenantId string, summaryType string, filter string, opts ...idao.CallOptions) ([]*view.RecordSumChartView, bool, error)
	FindRecordSumTable(ctx context.Context, tenantId string, filter string, groupFilter string, sort string, pageSize int64, pageNum int64, opts ...idao.CallOptions) (*view.RecordSumTableQueryView, bool, error)
}

type RecordDayViewDao struct {
	idao.Dao[*view.RecordDayView]
}

func NewRecordDayViewDao(dbKey string) *RecordDayViewDao {
	tableName := "master_record_day_view"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &view.RecordDayView{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*view.RecordDayView](newCfg)
	daoVal := &RecordDayViewDao{Dao: baseDao}
	return daoVal
}

func (r *RecordDayViewDao) Create(ctx context.Context, view *view.RecordDayView, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.Create(ctx, view, opts...)
}

func (r *RecordDayViewDao) CreateMany(ctx context.Context, views []*view.RecordDayView, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.CreateMany(ctx, views, opts...)
}

func (r *RecordDayViewDao) Update(ctx context.Context, view *view.RecordDayView, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.Update(ctx, view, opts...)
}

func (r *RecordDayViewDao) UpdateMany(ctx context.Context, views []*view.RecordDayView, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.UpdateMany(ctx, views, opts...)
}

func (r *RecordDayViewDao) Delete(ctx context.Context, view *view.RecordDayView, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.Delete(ctx, view, opts...)
}

func (r *RecordDayViewDao) DeleteMany(ctx context.Context, views []*view.RecordDayView, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.DeleteMany(ctx, views, opts...)
}

func (r *RecordDayViewDao) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.DeleteById(ctx, id, opts...)
}

func (r *RecordDayViewDao) DeleteByIds(ctx context.Context, ids []string, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.DeleteByIds(ctx, ids, opts...)
}

func (r *RecordDayViewDao) DeleteByRSQL(ctx context.Context, rsql string, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.DeleteByRSQL(ctx, rsql, opts...)
}

func (r *RecordDayViewDao) DeleteAll(ctx context.Context, opts ...idao.CallOptions) *idao.Result {
	return r.Dao.DeleteAll(ctx, opts...)
}

func (r *RecordDayViewDao) DeleteByGraphId(ctx context.Context, tenantId string, graphId string, opts ...idao.CallOptions) error {
	build := rsql.NewBuilder().Eq("graph_id", graphId).Build()
	return r.Dao.DeleteByRSQL(ctx, build, opts...).GetError()
}

func (r *RecordDayViewDao) FindByGraphId(ctx context.Context, tenantId string, graphId string, opts ...idao.CallOptions) ([]*view.RecordDayView,
	error) {
	build := rsql.NewBuilder().Eq("graph_id", graphId).Build()
	return r.Dao.FindByRSQL(ctx, build, opts...)
}

func (r *RecordDayViewDao) FindByCaseId(ctx context.Context, caseId string, opts ...idao.CallOptions) ([]*view.RecordDayView, error) {
	build := rsql.NewBuilder().Eq("case_id", caseId).Build()
	return r.Dao.FindByRSQL(ctx, build, opts...)
}

func (r *RecordDayViewDao) getMongoDao() store_mongodb.IMongoDao[*view.RecordDayView] {
	dao, ok := r.Dao.GetStore().(store_mongodb.IMongoDao[*view.RecordDayView])
	if !ok {
		panic("RecordDayViewDao is not store_mongodb.IMongoDao[*view.RecordDayView]")
	}
	return dao
}

// IncAmountMany
// @Description: 累计日汇总
// @receiver r
// @param ctx
// @param views
// @param opts
// @return error
func (r *RecordDayViewDao) IncAmountMany(ctx context.Context, views []*view.RecordDayView, opts ...idao.CallOptions) error {
	mongoDao := r.getMongoDao()
	opts = append(opts, idao.NewCallOptions().SetUpsert(true))
	var list []mongo.WriteModel
	for _, v := range views {
		filter := bson.D{{"id", v.Id}}
		data := bson.M{
			"$set": bson.M{
				"id":          v.Id,
				"tenant_id":   v.TenantId,
				"case_id":     v.CaseId,
				"graph_id":    v.GraphId,
				"master_id":   v.MasterId,
				"master_type": v.MasterType,

				"name": v.Name,
				"acct": v.Acct,

				"opp_name": v.OppName,
				"opp_acct": v.OppAcct,

				"date":  v.Date,
				"year":  v.Year,
				"month": v.Month,
				"day":   v.Day,
				"ccy":   v.Ccy,
			},
			"$inc": bson.M{
				"payout": v.Payout,
				"income": v.Income,
				"amount": v.Amount,
				"count":  1,
			},
		}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(data).SetUpsert(true)
		list = append(list, model)
	}
	if len(list) > 0 {
		_, err := mongoDao.BulkWrite(ctx, list, opts...)
		return err
	}
	return nil
}

type SummaryType string

const (
	SummaryType_Year  SummaryType = "year"
	SummaryType_Month SummaryType = "month"
	SummaryType_Day   SummaryType = "day"
)

func (s SummaryType) String() string {
	return string(s)
}

func (r *RecordDayViewDao) FindRecordSumChart(ctx context.Context, tenantId string, caseId string, summaryType string, filter string, opts ...idao.CallOptions) ([]*view.RecordSumChartView, bool, error) {
	mongoDao := r.getMongoDao()
	filterMap := mongoDao.GetFilterMap(tenantId, filter)

	groupMap := make(map[string]interface{})
	groupMap1 := make(map[string]interface{})
	projectMap := make(map[string]interface{})
	idMap := make(map[string]interface{})
	idMap1 := make(map[string]interface{})
	sort := make([]bson.E, 0)

	switch summaryType {
	case SummaryType_Year.String():
		idMap["year"] = "$year"
		idMap1["year"] = "$_id.year"
		projectMap["year"] = "$_id.year"
		sort = append(sort, bson.E{Key: "year", Value: 1})
	case SummaryType_Month.String():
		idMap["year"] = "$year"
		idMap["month"] = "$month"
		idMap1["year"] = "$_id.year"
		idMap1["month"] = "$_id.month"
		projectMap["year"] = "$_id.year"
		projectMap["month"] = "$_id.month"
		sort = append(sort, bson.E{Key: "year", Value: 1})
		sort = append(sort, bson.E{Key: "month", Value: 1})
	case SummaryType_Day.String():
		idMap["year"] = "$year"
		idMap["month"] = "$month"
		idMap["day"] = "$day"
		idMap1["year"] = "$_id.year"
		idMap1["month"] = "$_id.month"
		idMap1["day"] = "$_id.day"
		projectMap["year"] = "$_id.year"
		projectMap["month"] = "$_id.month"
		projectMap["day"] = "$_id.day"
		sort = append(sort, bson.E{Key: "year", Value: 1})
		sort = append(sort, bson.E{Key: "month", Value: 1})
		sort = append(sort, bson.E{Key: "day", Value: 1})
	default:
		return nil, false, errors.New("")
	}

	idMap["opp_name"] = "$opp_name"

	groupMap["_id"] = idMap
	groupMap["summary_type"] = map[string]interface{}{"$max": summaryType}
	groupMap["ccy"] = map[string]interface{}{"$max": "$ccy"}
	groupMap["amount"] = map[string]interface{}{"$sum": "$amount"}
	groupMap["count"] = map[string]interface{}{"$sum": 1}

	groupMap1["_id"] = idMap1
	groupMap1["summary_type"] = map[string]interface{}{"$max": "$summary_type"}
	groupMap1["ccy"] = map[string]interface{}{"$max": "$ccy"}
	groupMap1["amount"] = map[string]interface{}{"$sum": "$amount"}
	groupMap1["count"] = map[string]interface{}{"$sum": "$count"}
	groupMap1["opp_num"] = map[string]interface{}{"$sum": 1}

	projectMap["_id"] = "$_id"
	projectMap["summary_type"] = "$summary_type"
	projectMap["ccy"] = "$ccy"
	projectMap["amount"] = "$amount"
	projectMap["count"] = "$count"
	projectMap["opp_num"] = "$opp_num"

	data := make([]*view.RecordSumChartView, 0)

	err := mongoDao.AggregateByPipeline(ctx, mongo.Pipeline{
		bson.D{{"$match", filterMap}},
		bson.D{{"$group", groupMap}},
		bson.D{{"$group", groupMap1}},
		bson.D{{"$project", projectMap}},
		bson.D{{"$sort", sort}},
	}, &data)
	if err != nil {
		return nil, false, err
	}

	return data, true, nil
}

func (r *RecordDayViewDao) FindRecordSumTable(ctx context.Context, tenantId string, caseId string, filter string, groupFilter string, sort string, pageSzie int64, pageNum int64, opts ...idao.CallOptions) (*view.RecordSumTableQueryView, bool, error) {
	mongoDao := r.getMongoDao()

	filterMap := mongoDao.GetFilterMap(tenantId, filter)
	groupFilterMap := mongoDao.GetFilterMap(tenantId, groupFilter)

	groupMap := make(map[string]interface{})
	groupMap["_id"] = map[string]interface{}{
		"tenant_id": "$tenant_id",
		"name":      "$name",
		"opp_name":  "$opp_name",
		"ccy":       "$ccy",
		"opp_acct":  "$opp_acct",
	}
	groupMap["begin_date"] = map[string]interface{}{"$min": "$date"}
	groupMap["end_date"] = map[string]interface{}{"$max": "$date"}
	groupMap["count"] = map[string]interface{}{"$sum": "$count"}
	groupMap["amount"] = map[string]interface{}{"$sum": "$amount"}

	groupMap1 := make(map[string]interface{})
	groupMap1["_id"] = map[string]interface{}{
		"tenant_id": "$_id.tenant_id",
		"name":      "$_id.name",
		"opp_name":  "$_id.opp_name",
		"ccy":       "$_id.ccy",
	}
	groupMap1["begin_date"] = map[string]interface{}{"$min": "$begin_date"}
	groupMap1["end_date"] = map[string]interface{}{"$max": "$end_date"}
	groupMap1["count"] = map[string]interface{}{"$sum": "$count"}
	groupMap1["amount"] = map[string]interface{}{"$sum": "$amount"}
	groupMap1["opp_acct_count"] = map[string]interface{}{"$sum": 1}

	data := make([]*view.RecordSumTableQueryView, 0)

	pipeline := mongo.Pipeline{
		bson.D{{"$match", filterMap}},
		bson.D{{"$group", groupMap}},
		bson.D{{"$group", groupMap1}},
	}
	projectMap := make(map[string]interface{})
	projectMap["_id"] = map[string]interface{}{"$concat": []interface{}{"$_id.tenant_id", "-", "$_id.name", "-", "$_id.opp_name", "-", "$_id.ccy"}}
	projectMap["tenant_id"] = "$_id.tenant_id"
	projectMap["name"] = "$_id.name"
	projectMap["opp_name"] = "$_id.opp_name"
	projectMap["ccy"] = "$_id.ccy"
	projectMap["begin_date"] = "$begin_date"
	projectMap["end_date"] = "$end_date"
	projectMap["count"] = "$count"
	projectMap["amount"] = "$amount"
	projectMap["opp_acct_count"] = "$opp_acct_count"
	pipeline = append(pipeline, bson.D{{"$project", projectMap}})
	if len(groupFilter) > 0 {
		pipeline = append(pipeline, bson.D{{"$match", groupFilterMap}})
	}
	if len(sort) > 0 {
		sortSlice := strings.Split(sort, ":")
		field := sortSlice[0]
		order := 1
		if strings.ToLower(sortSlice[1]) == "desc" {
			order = -1
		}
		pipeline = append(pipeline, bson.D{{"$sort", []bson.E{bson.E{Key: stringutils.SnakeString(field), Value: order}}}})
	} else {
		pipeline = append(pipeline, bson.D{{"$sort", []bson.E{bson.E{Key: "amount", Value: -1}}}})
	}

	pipeline = append(pipeline, bson.D{{"$group", map[string]interface{}{
		"_id": idutils.NewId(),
		"data": map[string]interface{}{
			"$push": map[string]interface{}{
				"_id":            "$_id",
				"name":           "$name",
				"opp_name":       "$opp_name",
				"ccy":            "$ccy",
				"begin_date":     "$begin_date",
				"end_date":       "$end_date",
				"count":          "$count",
				"amount":         "$amount",
				"opp_acct_count": "$opp_acct_count",
			},
		},
		"total_rows": map[string]interface{}{"$sum": 1},
	}}})
	pipeline = append(pipeline, bson.D{{"$project", map[string]interface{}{
		"_id":        "$_id",
		"data":       map[string]interface{}{"$slice": []interface{}{"$data", pageSzie * pageNum, pageSzie}},
		"total_rows": "$total_rows",
	}}})

	err := mongoDao.AggregateByPipeline(ctx, pipeline, &data)
	if err != nil {
		return nil, false, err
	}

	if len(data) > 0 {
		return data[0], true, nil
	}
	return &view.RecordSumTableQueryView{
		Data:      make([]*view.RecordSumTableView, 0),
		TotalRows: 0,
	}, false, nil
}
