package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd/ddd_query"
)

type RecordService struct {
	dao *dao.RecordDao
}

var _recordService *RecordService
var _recordServiceOnce sync.Once

func NewRecordService() *RecordService {
	_recordServiceOnce.Do(func() {
		_recordService = newRecordService()
	})
	return _recordService
}

func newRecordService() *RecordService {
	return &RecordService{
		dao: dao.NewRecordDao(config.DBKey),
	}
}

func (r *RecordService) Create(ctx context.Context, v *model.Record, opts ...idao.CallOptions) error {
	return r.dao.Create(ctx, v, opts...).GetError()
}

func (r *RecordService) CreateMany(ctx context.Context, v []*model.Record, opts ...idao.CallOptions) error {
	return r.dao.CreateMany(ctx, v, opts...).GetError()
}

func (r *RecordService) Update(ctx context.Context, v *model.Record, opts ...idao.CallOptions) error {
	return r.dao.Update(ctx, v, opts...).GetError()
}

func (r *RecordService) UpdateMany(ctx context.Context, v []*model.Record, opts ...idao.CallOptions) error {
	return r.dao.UpdateMany(ctx, v, opts...).GetError()
}

func (r *RecordService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return r.dao.DeleteById(ctx, id, opts...).GetError()
}

func (r *RecordService) DeleteByTaskId(ctx context.Context, taskId string, opts ...idao.CallOptions) error {
	return r.dao.DeleteById(ctx, taskId, opts...).GetError()
}

func (r *RecordService) DeleteAll(ctx context.Context, opts ...idao.CallOptions) error {
	return r.dao.DeleteAll(ctx, opts...).GetError()
}

func (r *RecordService) FindById(ctx context.Context, qry *query.RecordFindByIdQuery, opts ...idao.CallOptions) (*model.Record, error) {
	return r.dao.FindById(ctx, qry.Id, opts...)
}

func (r *RecordService) FindByIds(ctx context.Context, qry *query.RecordFindByIdsQuery, opts ...idao.CallOptions) ([]*model.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordService) FindAll(ctx context.Context, qry *query.RecordFindAllQuery, opts ...idao.CallOptions) ([]*model.Record, error) {
	result := r.dao.FindAll(ctx, opts...)
	return result.Data, result.GetError()
}

func (r *RecordService) FindPaging(ctx context.Context, qry *ddd_query.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Record] {
	return r.dao.FindPaging(ctx, qry, opts...)
}

func (r *RecordService) FindPagingByCaseId(ctx context.Context, qry *query.RecordFindByCaseIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Record] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s'", qry.CaseId))
	return r.dao.FindPaging(ctx, qry, opts...)
}

func (r *RecordService) DistinctAccount(ctx context.Context, qry *query.DistinctAccountByNameQuery, opts ...idao.CallOptions) ([]*query.DistinctAccountResult, error) {
	myFilter := fmt.Sprintf(`name="%s"`, qry.Name)
	oppFilter := fmt.Sprintf(`opp_name="%s"`, qry.Name)
	return r.dao.DistinctAccount(ctx, qry.CaseId, qry.MasterType, qry.MasterId, myFilter, oppFilter)
}

func (r *RecordService) DistinctCompany(ctx context.Context, qry *query.DistinctNameQuery, opts ...idao.CallOptions) ([]*query.DistinctNameResult, error) {
	return r.dao.DistinctCompany(ctx, qry.CaseId, qry.MasterType, qry.MasterId, qry.MyFilter, qry.OppFilter)
}

func (r *RecordService) DistinctHuman(ctx context.Context, qry *query.DistinctNameQuery, opts ...idao.CallOptions) ([]*query.DistinctNameResult, error) {
	return r.dao.DistinctHuman(ctx, qry.CaseId, qry.MasterType, qry.MasterId, qry.MyFilter, qry.OppFilter)
}

func (r *RecordService) DistinctName(ctx context.Context, qry *query.DistinctNameQuery, opts ...idao.CallOptions) ([]*query.DistinctNameResult, error) {
	return r.dao.DistinctName(ctx, qry.CaseId, qry.MasterType, qry.MasterId, qry.MyFilter, qry.OppFilter)
}

func (r *RecordService) FindByDocId(ctx context.Context, qry *query.RecordFindByDocIdQuery, opts ...idao.CallOptions) ([]*model.Record, error) {
	rsqlStr := rsql.NewBuilder().And(rsql.Eq("case_id", qry.CaseId), rsql.Eq("doc_id", qry.DocId)).Build()
	return r.dao.FindByRSQL(ctx, rsqlStr, opts...)
}

func (r *RecordService) FindByTaskId(ctx context.Context, qry *query.RecordFindByTaskIdQuery, opts ...idao.CallOptions) ([]*model.Record, error) {
	rsqlStr := rsql.NewBuilder().And(rsql.Eq("case_id", qry.CaseId), rsql.Eq("task_id", qry.TaskId)).Build()
	return r.dao.FindByRSQL(ctx, rsqlStr, opts...)
}

func (r *RecordService) FindByFileId(ctx context.Context, qry *query.RecordFindByFileIdQuery, opts ...idao.CallOptions) ([]*model.Record, error) {
	rsqlStr := rsql.NewBuilder().And(rsql.Eq("case_id", qry.CaseId), rsql.Eq("file_id", qry.FileId)).Build()
	return r.dao.FindByRSQL(ctx, rsqlStr, opts...)
}
