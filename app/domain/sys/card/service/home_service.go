package service

import (
	"context"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type HomeService struct {
	dao *dao.HomeDao
	xbase.Service
}

var (
	_homeOnce          sync.Once
	_homeDomainService *HomeService
)

func NewHomeService() *HomeService {
	_homeOnce.Do(func() {
		_homeDomainService = &HomeService{
			dao: dao.NewHomeDao(config.DBKey),
		}
	})
	return _homeDomainService
}

func (t *HomeService) Create(ctx context.Context, cmd *command.HomeCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *HomeService) Delete(ctx context.Context, cmd *command.HomeDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *HomeService) Update(ctx context.Context, cmd *command.HomeUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *HomeService) CreateByModel(ctx context.Context, entity *model.Home) error {
	return t.dao.Create(ctx, entity).GetError()
}

func (t *HomeService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Home, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *HomeService) FindPaging(ctx context.Context, caseId string, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Home], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *HomeService) FindByCode(ctx context.Context, qry *query.FindByCodeQuery) (*model.Home, error) {
	builder := rsql.NewBuilder()
	conditions := []rsql.Condition{}
	conditions = append(conditions, rsql.Eq("code", qry.Code))
	if qry.CaseId != "" {
		conditions = append(conditions, rsql.Eq("case_id", qry.CaseId))
	}
	if qry.UserId != "" {
		conditions = append(conditions, rsql.Eq("user_id", qry.UserId))
	}
	res, err := t.dao.FindByRSQL(ctx, builder.And(conditions...).Build())
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res[0], nil
}

func (t *HomeService) FindDefaultByCode(ctx context.Context, code string) (*model.Home, error) {
	builder := rsql.NewBuilder().And(
		rsql.Eq("code", "default"),
		rsql.Eq("is_default", true),
	)
	res, err := t.dao.FindByRSQL(ctx, builder.Build())
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res[0], nil
}
