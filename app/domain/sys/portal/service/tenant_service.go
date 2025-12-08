package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/enum"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver/pkg/orm_pkg/dao"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type TenantService struct {
	dao       *dao.TenantDao
	tuDao     *dao.TenantUserDao
	tenantOpt idao.CallOptions
	xbase.Service
}

var (
	_tenantOnce          sync.Once
	_tenantDomainService *TenantService
)

func NewTenantService() *TenantService {
	_tenantOnce.Do(func() {
		_tenantDomainService = &TenantService{
			tenantOpt: idao.NewCallOptions().SetTenantId(SystemTenantId),
			dao:       dao.NewTenantDao(config.DBKey),
			tuDao:     dao.NewTenantUserDao(config.DBKey),
		}
	})
	return _tenantDomainService
}

func (t *TenantService) CreateSysTenant(ctx context.Context) (*model.Tenant, error) {
	sysTenant, err := t.FindById(ctx, SystemTenantId)
	if err != nil {
		return nil, err
	}
	if sysTenant != nil {
		return nil, errors.New("租户只能初始化一次")
	}

	sysTenant, _ = model.NewTenant()
	sysTenant.Id = SystemTenantId
	sysTenant.Name = "系统租户"
	sysTenant.Status = enum.Using

	err = t.Create(ctx, sysTenant)
	if err != nil {
		return nil, err
	}
	return sysTenant, nil
}

func (t *TenantService) Create(ctx context.Context, data *model.Tenant) error {
	return t.dao.Create(ctx, data, t.tenantOpt).GetError()
}

func (t *TenantService) Delete(ctx context.Context, cmd *command.TenantDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id, t.tenantOpt).GetError()
	})
}

func (t *TenantService) Update(ctx context.Context, cmd *command.TenantUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions(t.tenantOpt, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask))).GetError()
	})
}

func (t *TenantService) FindById(ctx context.Context, id string) (*model.Tenant, error) {
	return t.dao.FindById(ctx, id, t.tenantOpt)
}

func (t *TenantService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Tenant], error) {
	res := t.dao.FindPaging(ctx, qry, t.tenantOpt)
	return res, res.GetError()
}

func (t *TenantService) FindCountByCode(ctx context.Context, code string) (int64, error) {
	return t.dao.CountByRSQL(ctx, fmt.Sprintf("code=='%s'", code), t.tenantOpt)
}

func (t *TenantService) FindByUserId(ctx context.Context, userId string) ([]*model.Tenant, error) {
	tus, err := t.tuDao.FindByRSQL(ctx, fmt.Sprintf("user_id=='%s'", userId), t.tenantOpt)
	if err != nil {
		return nil, err
	}
	if tus == nil {
		return []*model.Tenant{}, nil
	}

	tenantIds := make([]string, 0, len(tus))
	for _, tus := range tus {
		tenantIds = append(tenantIds, tus.TenId)
	}

	builder := dao2.NewRSQLBuilder()
	rSql := builder.And(builder.In("id", tenantIds), builder.Eq("status", "Using")).Build()

	return t.dao.FindByRSQL(ctx, rSql, t.tenantOpt)
}
