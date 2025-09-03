package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"sync"
)

type FunService struct {
	dao *dao.FunDao
	xbase.Service
}

var (
	_funOnce          sync.Once
	_funDomainService *FunService
)

func NewFunService() *FunService {
	_funOnce.Do(func() {
		_funDomainService = &FunService{
			dao: dao.NewFunDao(config.DBKey),
		}
	})
	return _funDomainService
}

func (t *FunService) Create(ctx context.Context, cmd *command.FunCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *FunService) Delete(ctx context.Context, cmd *command.FunDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *FunService) Update(ctx context.Context, cmd *command.FunUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *FunService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Fun, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *FunService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Fun], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *FunService) FindAll(ctx context.Context) ([]*model.Fun, error) {
	return t.dao.FindByRSQL(ctx, "")
}
