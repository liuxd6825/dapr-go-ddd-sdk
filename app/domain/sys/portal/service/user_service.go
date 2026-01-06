package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/enum"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver/pkg/orm_pkg/dao"
	client "github.com/ory/kratos-client-go"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/oryservice"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

const SystemTenantId = "sys"

type UserService struct {
	dao *dao.UserDao
	xbase.Service
	tenantOpt  idao.CallOptions
	oryService *service2.OryService
	tuDao      *dao.TenantUserDao
}

var (
	_userOnce          sync.Once
	_userDomainService *UserService
)

func NewUserService() *UserService {
	_userOnce.Do(func() {
		_userDomainService = &UserService{
			tenantOpt:  idao.NewCallOptions().SetTenantId(SystemTenantId),
			dao:        dao.NewUserDao(config.DBKey),
			oryService: service2.NewOryService(),
			tuDao:      dao.NewTenantUserDao(config.DBKey),
		}
	})
	return _userDomainService
}

func (t *UserService) GetConfig() *idao.DaoConfig {
	return t.dao.GetConfig()
}

func (t *UserService) CreateSysUser(ctx context.Context) (*model.User, error) {

	user, _ := model.NewUser()
	user.Id = "super_admin"
	user.TenantId = SystemTenantId
	user.Account = "super_admin"
	user.Email = "super_admin@163.com"
	user.Name = "系统管理员"
	user.Status = enum.Using

	user, err := t.CreateInitUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (t *UserService) CreateDemoUser(ctx context.Context) (*model.User, error) {

	user, _ := model.NewUser()
	user.Id = "user_demo"
	user.TenantId = SystemTenantId
	user.Account = "demo"
	user.Email = "demo@163.com"
	user.Name = "演示用户"
	user.Status = enum.Using

	user, err := t.CreateInitUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (t *UserService) CreateInitUser(ctx context.Context, user *model.User) (*model.User, error) {

	ident, err := t.oryService.GetIdentityByCode(ctx, user.Account)
	if err != nil {
		return nil, err
	}

	password := "123@@abc"
	if ident == nil {
		body := client.CreateIdentityBody{
			Credentials: &client.IdentityWithCredentials{
				Password: &client.IdentityWithCredentialsPassword{
					Config: &client.IdentityWithCredentialsPasswordConfig{
						Password: &password,
					},
				},
			},
			Traits: map[string]interface{}{
				"account": user.Account,
				"email":   user.Email,
			},
		}
		ident, err = t.oryService.CreateIdentityExecute(ctx, body)
		if err != nil {
			if err.Error() == "409 Conflict" {

			} else {
				return nil, err
			}
		}
	}

	user.OryIdentityId = ident.Id
	user.Password = password

	tempUser, err := t.FindById(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	if tempUser == nil {
		err = t.CreateData(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (t *UserService) Create(ctx context.Context, cmd *command.UserCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data, t.tenantOpt).GetError()
	})
}

func (t *UserService) CreateData(ctx context.Context, data *model.User) error {
	return t.dao.Create(ctx, data, t.tenantOpt).GetError()
}

func (t *UserService) Delete(ctx context.Context, cmd *command.UserDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id, t.tenantOpt).GetError()
	})
}

func (t *UserService) DeleteById(ctx context.Context, id string) error {
	return t.dao.DeleteById(ctx, id, t.tenantOpt).GetError()
}

func (t *UserService) Update(ctx context.Context, cmd *command.UserUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Update(ctx, &cmd.Data, idao.NewCallOptions(t.tenantOpt, idao.NewCallOptions().SetUpdateFields(cmd.UpdateMask))).GetError()
	})
}

func (t *UserService) UpdateData(ctx context.Context, data *model.User, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, idao.NewCallOptions(t.tenantOpt, idao.NewCallOptions(opts...))).GetError()
}

func (t *UserService) FindById(ctx context.Context, id string) (*model.User, error) {
	return t.dao.FindById(ctx, id, t.tenantOpt)
}

func (t *UserService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.User], error) {
	res := t.dao.FindPaging(ctx, qry, t.tenantOpt)
	return res, res.GetError()
}

func (t *UserService) FindPagingByTenantId(ctx context.Context, tenantId string, qry store.FindPagingQuery) (idao.FindPagingResult[*model.UserView], error) {
	tus, err := t.tuDao.FindByRSQL(ctx, fmt.Sprintf("ten_id=='%s'", tenantId), t.tenantOpt)
	if err != nil {
		return nil, err
	}
	if tus == nil {
		return store.NewFindPagingResult([]*model.UserView{}, 0, qry, nil), nil
	}

	userIds := make([]string, 0, len(tus))
	for _, tu := range tus {
		userIds = append(userIds, tu.UserId)
	}

	builder := dao2.NewRSQLBuilder()
	rSql := builder.And(builder.In("id", userIds), builder.Eq("status", "Using")).Build()

	qry.SetMustFilter(rSql)

	userRes := t.dao.FindPaging(ctx, qry, t.tenantOpt)
	if userRes.GetError() != nil {
		return nil, userRes.GetError()
	}

	uvs := make([]*model.UserView, 0)

	for _, user := range userRes.GetData() {
		tu := t.findTenantUserByUserId(tus, user.Id)
		if tu == nil {
			continue
		}
		uv := model.NewUserView()
		uv.Id = user.Id
		uv.TenantId = user.TenantId
		uv.Account = user.Account
		uv.Name = user.Name
		uv.Phone = user.Phone
		uv.Email = user.Email
		uv.Gender = user.Gender
		uv.Address = user.Address
		uv.Work = user.Work
		uv.Remark = user.Remark
		uv.Status = user.Status
		uv.OryIdentityId = user.OryIdentityId
		uv.HeadPicture = user.HeadPicture
		uv.TanId = tu.TenId
		uv.TenantUserId = tu.Id
		uv.IsAdmin = tu.IsAdmin
		uv.CreatorId = user.CreatorId
		uv.CreatorName = user.CreatorName
		uv.CreatedTime = user.CreatedTime
		uv.UpdaterId = user.UpdaterId
		uv.UpdaterName = user.UpdaterName
		uv.UpdatedTime = user.UpdatedTime

		uvs = append(uvs, uv)
	}

	res := store.NewFindPagingResult(uvs, int64(len(uvs)), qry, nil)
	return res, nil
}

func (t *UserService) FindPagingByNotInTenantUser(ctx context.Context, tenantId string, qry store.FindPagingQuery) (idao.FindPagingResult[*model.User], error) {
	tus, err := t.tuDao.FindByRSQL(ctx, fmt.Sprintf("ten_id=='%s'", tenantId), t.tenantOpt)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, 0, len(tus))
	for _, tu := range tus {
		userIds = append(userIds, tu.UserId)
	}

	builder := dao2.NewRSQLBuilder()
	rSql := builder.And(builder.Out("id", userIds), builder.Eq("status", "Using")).Build()

	qry.SetMustFilter(rSql)

	userRes := t.dao.FindPaging(ctx, qry, t.tenantOpt)
	return userRes, userRes.GetError()
}

func (t *UserService) findTenantUserByUserId(tus []*model.TenantUser, userId string) *model.TenantUser {
	for _, tu := range tus {
		if tu.UserId == userId {
			return tu
		}
	}
	return nil
}

func (t *UserService) FindUsingByAccount(ctx context.Context, account string) (*model.User, error) {
	arr, err := t.dao.FindByRSQL(ctx, fmt.Sprintf("account=='%s' and status=='Using'", account), t.tenantOpt)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, errors.New("NotFound")
	}
	return arr[0], nil
}
