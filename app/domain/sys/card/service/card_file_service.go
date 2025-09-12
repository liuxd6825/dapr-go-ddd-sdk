package service

import (
	"context"
	"fmt"
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

type CardFileService struct {
	dao *dao.CardFileDao
	xbase.Service
}

var (
	_cardFileOnce          sync.Once
	_cardFileDomainService *CardFileService
)

func NewCardFileService() *CardFileService {
	_cardFileOnce.Do(func() {
		_cardFileDomainService = &CardFileService{
			dao: dao.NewCardFileDao(config.DBKey),
		}
	})
	return _cardFileDomainService
}

func (t *CardFileService) Create(ctx context.Context, cmd *command.CardFileCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *CardFileService) Delete(ctx context.Context, cmd *command.CardFileDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *CardFileService) Update(ctx context.Context, cmd *command.CardFileUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *CardFileService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.CardFile, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *CardFileService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.CardFile], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *CardFileService) FindByAppId(ctx context.Context, appId string) ([]*model.CardFile, error) {
	rsql := fmt.Sprintf("app_id=='%s'", appId)
	return t.dao.FindByRSQL(ctx, rsql)
}

func (t *CardFileService) FindByFunId(ctx context.Context, funId string) ([]*model.CardFile, error) {
	rsql := fmt.Sprintf("fun_id=='%s'", funId)
	return t.dao.FindByRSQL(ctx, rsql)
}

func (t *CardFileService) FindCards(ctx context.Context, qry *query.FindCardFilesQuery) ([]*model.CardFile, error) {
	rsql := fmt.Sprintf("app_id=='%s'", qry.AppId)
	if len(qry.FunId) > 0 {
		rsql = rsql + fmt.Sprintf(" and fun_id=in=(%s)", qry.FunId)
	}
	return t.dao.FindByRSQL(ctx, rsql)
}
