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

type DictService struct {
	dao *dao.DictDao
	xbase.Service
}

var (
	_dictOnce          sync.Once
	_dictDomainService *DictService
)

func NewDictService() *DictService {
	_dictOnce.Do(func() {
		_dictDomainService = &DictService{
			dao: dao.NewDictDao(config.DBKey),
		}
	})
	return _dictDomainService
}

func (t *DictService) Create(ctx context.Context, cmd *command.DictCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *DictService) Delete(ctx context.Context, cmd *command.DictDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *DictService) DeleteBatch(ctx context.Context, cmd *command.DictDeleteBatchCommand) error {
	return t.dao.DeleteByIds(ctx, cmd.Data).GetError()
}

func (t *DictService) Update(ctx context.Context, cmd *command.DictUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *DictService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Dict, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *DictService) FindPaging(ctx context.Context, caseId string, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Dict], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *DictService) FindByDictTypeAndCode(ctx context.Context, dictType, code string) ([]*model.Dict, error) {
	return t.dao.FindByRSQL(ctx, fmt.Sprintf("dict_type_id=='%s' and code=='%s'", dictType, code))
}

func (t *DictService) FindByDictTypeCodeAndParams(ctx context.Context, dictTypeCode, castTypeId, caseId string) ([]*model.Dict, error) {
	rSql := fmt.Sprintf("dict_type_code=='%s'", dictTypeCode)
	if castTypeId != "" {
		rSql = fmt.Sprintf("%s and cast_type_id=='%s'", rSql, castTypeId)
	}
	if caseId != "" {
		rSql = fmt.Sprintf("%s and case_id=='%s'", rSql, caseId)
	}
	return t.dao.FindByRSQL(ctx, rSql)
}
