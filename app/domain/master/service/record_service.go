package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"sync"
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

func (r *RecordService) FindPagingByCaseId(ctx context.Context, qry *query.RecordFindByCaseIdQuery, opts ...idao.CallOptions) (*idao.FindPagingResult[*model.Record], error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordService) FindByDocId(ctx context.Context, qry *query.RecordFindByDocIdQuery, opts ...idao.CallOptions) ([]*model.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordService) FindByTaskId(ctx context.Context, qry *query.RecordFindByTaskIdQuery, opts ...idao.CallOptions) ([]*model.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordService) FindByFileId(ctx context.Context, qry *query.RecordFindByFileIdQuery, opts ...idao.CallOptions) ([]*model.Record, error) {
	//TODO implement me
	panic("implement me")
}
