package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/status/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/status/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/status/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/status/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type StatusService struct {
	dao *dao.StatusDao
}

func NewStatusService() *StatusService {
	statusDao := dao.NewStatusDao(config.DBKey)
	return &StatusService{
		dao: statusDao,
	}
}

func (s *StatusService) Create(ctx context.Context, cmd *command.StatusCreateCommand) error {
	return s.dao.Create(ctx, cmd.Data).GetError()
}

func (s *StatusService) Update(ctx context.Context, cmd *command.StatusUpdateCommand) error {
	return s.dao.Update(ctx, cmd.Data).GetError()
}

func (s *StatusService) Delete(ctx context.Context, cmd *xbase.DeleteByIdCommand) error {
	return s.dao.DeleteById(ctx, cmd.Data.Id).GetError()
}

func (s *StatusService) FindById(ctx context.Context, qry *store.FindByIdRequest) (*model.Status, error) {
	return s.dao.FindById(ctx, qry.Id)
}

func (s *StatusService) FindAll(ctx context.Context) ([]*model.Status, error) {
	result := s.dao.FindAll(ctx)
	return result.GetData(), result.GetError()
}

func (s *StatusService) FindPaging(ctx context.Context, qry *store.FindPagingQueryRequest) store.FindPagingResult[*model.Status] {
	return s.dao.FindPaging(ctx, qry)
}

func (s *StatusService) FindByUserId(ctx context.Context, qry *query.StatusFindByUserId) ([]*model.Status, error) {
	build := rsql.NewBuilder(func(b *rsql.Builder) rsql.Condition {
		return b.Eq("user_id", qry.UserId)
	})
	return s.dao.FindByRSQL(ctx, build.Build())
}
