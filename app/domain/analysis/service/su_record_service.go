package service

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd/ddd_query"
)

type SuRecordService struct {
	dao *dao.SuRecordDao
}

func NewSuRecordService() *SuRecordService {
	return &SuRecordService{
		dao: dao.NewSuRecordDao(config.DBKey),
	}
}

func (s *SuRecordService) Create(ctx context.Context, v *model.SuRecord, opts ...idao.CallOptions) error {
	return s.dao.Create(ctx, v, opts...).GetError()
}

func (s *SuRecordService) CreateMany(ctx context.Context, v []*model.SuRecord, opts ...idao.CallOptions) error {
	return s.dao.CreateMany(ctx, v, opts...).GetError()
}

func (s *SuRecordService) Update(ctx context.Context, v *model.SuRecord, opts ...idao.CallOptions) error {
	return s.dao.Update(ctx, v, opts...).GetError()
}

func (s *SuRecordService) UpdateMany(ctx context.Context, v []*model.SuRecord, opts ...idao.CallOptions) error {
	return s.dao.UpdateMany(ctx, v, opts...).GetError()
}

func (s *SuRecordService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

func (s *SuRecordService) DeleteByIds(ctx context.Context, listId []string, opts ...idao.CallOptions) error {
	return s.dao.DeleteByIds(ctx, listId, opts...).GetError()
}

func (s *SuRecordService) FindById(ctx context.Context, qry *query.SuRecordFindByIdQuery, opts ...idao.CallOptions) (*model.SuRecord, error) {
	return s.dao.FindById(ctx, qry.Id, opts...)
}

func (s *SuRecordService) FindPaging(ctx context.Context, qry *ddd_query.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.SuRecord] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

func (s *SuRecordService) SuRecordFindPaging(ctx context.Context, qry *query.SuRecordFindByCaseIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.SuRecord] {
	mustFilter := fmt.Sprintf("case_id=='%s'", qry.CaseId)
	if qry.TaskId != "" {
		mustFilter = fmt.Sprintf("%s and task_id=='%s'", mustFilter, qry.TaskId)
	}
	qry.SetMustFilter(mustFilter)
	return s.dao.FindPaging(ctx, qry, opts...)
}
