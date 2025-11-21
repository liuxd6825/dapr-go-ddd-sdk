package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

// BankService
// @Description: 币种
type BankService struct {
	dao *dao.BankDao
	xbase.Service
}

var (
	_BankOnce    sync.Once
	_BankService *BankService
)

func NewBankService() *BankService {
	_BankOnce.Do(func() {
		_BankService = &BankService{
			dao: dao.NewBankDao(config.DBKey),
		}
	})
	return _BankService
}

func (s *BankService) GetConfig() *idao.DaoConfig {
	return s.dao.GetConfig()
}

func (s *BankService) Create(ctx context.Context, cmd *command.BankCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (s *BankService) CreateMany(ctx context.Context, cmd *command.BankCreateManyCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.CreateMany(ctx, cmd.Data).GetError()
	})
}

func (s *BankService) Delete(ctx context.Context, cmd *command.BankDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (s *BankService) DeleteAll(ctx context.Context, cmd *command.BankDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.DeleteAll(ctx).GetError()
	})
}

func (s *BankService) DeleteBatch(ctx context.Context, cmd *command.BankDeleteBatchCommand) error {
	return s.dao.DeleteByIds(ctx, cmd.Data).GetError()
}

func (s *BankService) Update(ctx context.Context, cmd *command.BankUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return s.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (s *BankService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Bank, error) {
	return s.dao.FindById(ctx, qry.Id)
}

func (s *BankService) FindByAll(ctx context.Context) ([]*model.Bank, error) {
	data := s.dao.FindAll(ctx)
	return data.Data, data.Error
}

func (s *BankService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Bank], error) {
	res := s.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}
