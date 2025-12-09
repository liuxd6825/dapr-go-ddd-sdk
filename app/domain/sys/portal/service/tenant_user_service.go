package service

import (
	"context"
	"fmt"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver/pkg/orm_pkg/dao"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type TenantUserService struct {
	dao *dao.TenantUserDao
	xbase.Service
	tenantOpt idao.CallOptions
}

var (
	_tenantUserOnce    sync.Once
	_tenantUserService *TenantUserService
)

func NewTenantUserService() *TenantUserService {
	_tenantUserOnce.Do(func() {
		_tenantUserService = &TenantUserService{
			tenantOpt: idao.NewCallOptions().SetTenantId(SystemTenantId),
			dao:       dao.NewTenantUserDao(config.DBKey),
		}
	})
	return _tenantUserService
}

func (t *TenantUserService) CreateSysTenantUser(ctx context.Context, tenantId, userId string) (*model.TenantUser, error) {
	tenantUser, _ := model.NewTenantUser()
	tenantUser.Id = fmt.Sprintf("%s_%s", tenantId, userId)
	tenantUser.TenId = tenantId
	tenantUser.UserId = userId
	tenantUser.IsAdmin = true

	err := t.dao.Create(ctx, tenantUser, t.tenantOpt).GetError()
	if err != nil {
		return nil, err
	}
	return tenantUser, nil
}

func (t *TenantUserService) CreateMany(ctx context.Context, data []*model.TenantUser) error {
	return t.dao.CreateMany(ctx, data, t.tenantOpt).GetError()
}

func (t *TenantUserService) Update(ctx context.Context, data *model.TenantUser, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, idao.NewCallOptions(t.tenantOpt, idao.NewCallOptions(opts...))).GetError()
}

func (t *TenantUserService) DeleteBatch(ctx context.Context, ids []string) error {
	return t.dao.DeleteByIds(ctx, ids, t.tenantOpt).GetError()
}

func (t *TenantUserService) DeleteByTenantIdsAndUserIds(ctx context.Context, tenantIds []string, userIds []string) error {
	builder := dao2.NewRSQLBuilder()
	rSql := builder.And(builder.In("ten_id", tenantIds), builder.In("user_id", userIds)).Build()
	return t.dao.DeleteByRSQL(ctx, rSql, t.tenantOpt).GetError()
}

func (t *TenantUserService) DeleteByTenantIds(ctx context.Context, tenantIds []string) error {
	builder := dao2.NewRSQLBuilder()
	rSql := builder.In("ten_id", tenantIds).Build()
	return t.dao.DeleteByRSQL(ctx, rSql, t.tenantOpt).GetError()
}

func (t *TenantUserService) DeleteByUserId(ctx context.Context, userId string) error {
	builder := dao2.NewRSQLBuilder()
	rSql := builder.Eq("user_id", userId).Build()
	return t.dao.DeleteByRSQL(ctx, rSql, t.tenantOpt).GetError()
}

func (t *TenantUserService) FindByTenantIdAndUserId(ctx context.Context, tenantId, userId string) (*model.TenantUser, error) {
	builder := dao2.NewRSQLBuilder()
	rSql := builder.And(builder.Eq("ten_id", tenantId), builder.Eq("user_id", userId)).Build()
	arr, err := t.dao.FindByRSQL(ctx, rSql, t.tenantOpt)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, nil
	}
	return arr[0], nil
}
