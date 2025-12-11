package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type CaseService struct {
	dao *dao.CaseDao
}

var (
	_caseOnce          sync.Once
	_caseDomainService *CaseService
)

func NewCaseService() *CaseService {
	_caseOnce.Do(func() {
		_caseDomainService = &CaseService{
			dao: dao.NewCaseDao(config.DBKey),
		}
	})
	return _caseDomainService
}

func (s *CaseService) GetConfig() *idao.DaoConfig {
	return s.dao.GetConfig()
}

func (t *CaseService) Create(ctx context.Context, cmd *command.CaseCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *CaseService) Delete(ctx context.Context, cmd *command.CaseDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *CaseService) DeleteBatch(ctx context.Context, cmd *command.CaseDeleteBatchCommand) error {
	return t.dao.DeleteByIds(ctx, cmd.Data.Ids).GetError()
}

func (t *CaseService) Update(ctx context.Context, cmd *command.CaseUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *CaseService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Case, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *CaseService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Case], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}
