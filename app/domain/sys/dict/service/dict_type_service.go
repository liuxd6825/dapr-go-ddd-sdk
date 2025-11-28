package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type DictTypeService struct {
	dao *dao.DictTypeDao
	xbase.Service
}

var (
	_dictTypeOnce          sync.Once
	_dictTypeDomainService *DictTypeService
)

func NewDictTypeService() *DictTypeService {
	_dictTypeOnce.Do(func() {
		_dictTypeDomainService = &DictTypeService{
			dao: dao.NewDictTypeDao(config.DBKey),
		}
	})
	return _dictTypeDomainService
}

func (t *DictTypeService) Create(ctx context.Context, cmd *command.DictTypeCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *DictTypeService) Delete(ctx context.Context, cmd *command.DictTypeDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *DictTypeService) DeleteBatch(ctx context.Context, cmd *command.DictTypeDeleteBatchCommand) error {
	return t.dao.DeleteByIds(ctx, cmd.Data).GetError()
}

func (t *DictTypeService) Update(ctx context.Context, cmd *command.DictTypeUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *DictTypeService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.DictType, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *DictTypeService) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.DictType], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *DictTypeService) FindByCode(ctx context.Context, code string) ([]*model.DictType, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("code=='%s'", code))
}

func (t *DictTypeService) FindAll(ctx context.Context, qry store.FindPagingQuery) ([]*model.DictType, error) {
	res := t.dao.FindPaging(ctx, qry)
	return res.GetData(), res.GetError()
}
