package service

import (
	"context"
	"errors"
	"fmt"
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

const SystemTenantId = "sys"

type UserService struct {
	dao *dao.UserDao
	xbase.Service
	tenantOpt idao.CallOptions
}

var (
	_userOnce          sync.Once
	_userDomainService *UserService
)

func NewUserService() *UserService {
	_userOnce.Do(func() {
		_userDomainService = &UserService{
			tenantOpt: idao.NewCallOptions().SetTenantId(SystemTenantId),
			dao:       dao.NewUserDao(config.DBKey),
		}
	})
	return _userDomainService
}

func (t *UserService) GetConfig() *idao.DaoConfig {
	return t.dao.GetConfig()
}

func (t *UserService) Create(ctx context.Context, cmd *command.UserCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data, t.tenantOpt).GetError()
	})
}

func (t *UserService) Delete(ctx context.Context, cmd *command.UserDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id, t.tenantOpt).GetError()
	})
}

func (t *UserService) Update(ctx context.Context, cmd *command.UserUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions(t.tenantOpt, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask))).GetError()
	})
}

func (t *UserService) UpdateData(ctx context.Context, data *model.User, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, idao.NewCallOptions(t.tenantOpt, idao.NewCallOptions(opts...))).GetError()
}

func (t *UserService) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.User, error) {
	return t.dao.FindById(ctx, qry.Id, t.tenantOpt)
}

func (t *UserService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.User], error) {
	res := t.dao.FindPaging(ctx, qry, t.tenantOpt)
	return res, res.GetError()
}

func (t *UserService) FindByAccount(ctx context.Context, account string) (*model.User, error) {
	arr, err := t.dao.FindByRSQL(ctx, fmt.Sprintf("account=='%s'", account), t.tenantOpt)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, errors.New(fmt.Sprintf("没有找到Account=[%s]的用户数据", account))
	}
	return arr[0], nil
}
