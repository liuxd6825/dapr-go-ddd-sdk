package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/view"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type RecordQueryService interface {
	Create(ctx context.Context, v *view.RecordView, opts ...idao.CallOptions) error
	CreateMany(ctx context.Context, v []*view.RecordView, opts ...idao.CallOptions) error
	Update(ctx context.Context, v *view.RecordView, opts ...idao.CallOptions) error
	UpdateMany(ctx context.Context, v []*view.RecordView, opts ...idao.CallOptions) error
	DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error
	DeleteByTaskId(ctx context.Context, taskId string, opts ...idao.CallOptions) error
	DeleteAll(ctx context.Context, opts ...idao.CallOptions) error

	FindById(ctx context.Context, qry *query.RecordFindByIdQuery, opts ...idao.CallOptions) (*view.RecordView, error)
	FindByIds(ctx context.Context, qry *query.RecordFindByIdsQuery, opts ...idao.CallOptions) ([]*view.RecordView, error)
	FindAll(ctx context.Context, qry *query.RecordFindAllQuery, opts ...idao.CallOptions) ([]*view.RecordView, error)
	FindPaging(ctx context.Context, qry *ddd_query.FindPagingQuery, opts ...idao.CallOptions) (*idao.FindPagingResult[*view.RecordView], error)
	FindPagingByCaseId(ctx context.Context, qry *query.RecordFindByCaseIdQuery, opts ...idao.CallOptions) (*idao.FindPagingResult[*view.RecordView], error)
	FindByDocId(ctx context.Context, qry *query.RecordFindByDocIdQuery, opts ...idao.CallOptions) ([]*view.RecordView, error)
	FindByTaskId(ctx context.Context, qry *query.RecordFindByTaskIdQuery, opts ...idao.CallOptions) ([]*view.RecordView, error)
	FindByFileId(ctx context.Context, qry *query.RecordFindByFileIdQuery, opts ...idao.CallOptions) ([]*view.RecordView, error)
}

type RecordQueryServiceImpl struct {
	dao idao.Dao[*view.RecordView]
}

func (r *RecordQueryServiceImpl) Create(ctx context.Context, v *view.RecordView, opts ...idao.CallOptions) error {
	return r.dao.Create(ctx, v, opts...).GetError()
}

func (r *RecordQueryServiceImpl) CreateMany(ctx context.Context, v []*view.RecordView, opts ...idao.CallOptions) error {
	return r.dao.CreateMany(ctx, v, opts...).GetError()
}

func (r *RecordQueryServiceImpl) Update(ctx context.Context, v *view.RecordView, opts ...idao.CallOptions) error {
	return r.dao.Update(ctx, v, opts...).GetError()
}

func (r *RecordQueryServiceImpl) UpdateMany(ctx context.Context, v []*view.RecordView, opts ...idao.CallOptions) error {
	return r.dao.UpdateMany(ctx, v, opts...).GetError()
}

func (r *RecordQueryServiceImpl) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return r.dao.DeleteById(ctx, id, opts...).GetError()
}

func (r *RecordQueryServiceImpl) DeleteByTaskId(ctx context.Context, taskId string, opts ...idao.CallOptions) error {
	return r.dao.DeleteById(ctx, taskId, opts...).GetError()
}

func (r *RecordQueryServiceImpl) DeleteAll(ctx context.Context, opts ...idao.CallOptions) error {
	return r.dao.DeleteAll(ctx, opts...).GetError()
}

func (r *RecordQueryServiceImpl) FindById(ctx context.Context, qry *query.RecordFindByIdQuery, opts ...idao.CallOptions) (*view.RecordView, error) {
	return r.dao.FindById(ctx, qry.Id, opts...)
}

func (r *RecordQueryServiceImpl) FindByIds(ctx context.Context, qry *query.RecordFindByIdsQuery, opts ...idao.CallOptions) ([]*view.RecordView, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordQueryServiceImpl) FindAll(ctx context.Context, qry *query.RecordFindAllQuery, opts ...idao.CallOptions) ([]*view.RecordView, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordQueryServiceImpl) FindPaging(ctx context.Context, qry *ddd_query.FindPagingQuery, opts ...idao.CallOptions) (*idao.FindPagingResult[*view.RecordView], error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordQueryServiceImpl) FindPagingByCaseId(ctx context.Context, qry *query.RecordFindByCaseIdQuery, opts ...idao.CallOptions) (*idao.FindPagingResult[*view.RecordView], error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordQueryServiceImpl) FindByDocId(ctx context.Context, qry *query.RecordFindByDocIdQuery, opts ...idao.CallOptions) ([]*view.RecordView, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordQueryServiceImpl) FindByTaskId(ctx context.Context, qry *query.RecordFindByTaskIdQuery, opts ...idao.CallOptions) ([]*view.RecordView, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordQueryServiceImpl) FindByFileId(ctx context.Context, qry *query.RecordFindByFileIdQuery, opts ...idao.CallOptions) ([]*view.RecordView, error) {
	//TODO implement me
	panic("implement me")
}

func NewRecordQueryService(dao idao.Dao[*view.RecordView]) RecordQueryService {
	return &RecordQueryServiceImpl{
		dao: dao,
	}
}
