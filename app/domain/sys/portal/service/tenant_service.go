package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"sync"
)

type TenantService struct {
	dao *dao.TenantDao
	xbase.Service
}

var (
	_tenantOnce          sync.Once
	_tenantDomainService *TenantService
)

func NewTenantService() *TenantService {
	_tenantOnce.Do(func() {
		_tenantDomainService = &TenantService{
			dao: dao.NewTenantDao(config.DBKey),
		}
	})
	return _tenantDomainService
}

func (t *TenantService) Create(ctx context.Context, cmd *command.TenantCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *TenantService) Delete(ctx context.Context, cmd *command.TenantDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *TenantService) Update(ctx context.Context, cmd *command.TenantUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask)).GetError()
	})
}

func (t *TenantService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Tenant, error) {
	return t.dao.FindById(ctx, qry.Id)
}

func (t *TenantService) FindPaging(ctx context.Context, caseId string, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Tenant], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}
