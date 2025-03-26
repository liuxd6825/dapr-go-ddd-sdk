package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/qrysrv/vo"
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/xinfra/enums"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

// UserDao 用户领域Elasticsearch仓库接口
type UserDao interface {
	// FindByAccountAndPassword 根据账号和密码查找用户信息
	FindByAccountAndPassword(ctx context.Context, account, password string) *vo.UserVo

	// FindByIdAndPassword 根据ID和密码查找用户信息
	FindByIdAndPassword(ctx context.Context, id, password string) *vo.UserVo

	// FindUsers 根据用户ID列表和过滤条件查找用户列表
	FindUsers(ctx context.Context, usersId []string, filter string) []*vo.UserVo

	// FindByStatus 根据状态查找所有用户
	FindByStatus(ctx context.Context, status enums.UserStatus) []*vo.UserVo

	// FindByUserType 根据用户类型查找所有用户
	FindByUserType(ctx context.Context, userType enums.UserType) []*vo.UserVo

	// FindByAccount 根据账号查找用户信息
	FindByAccount(ctx context.Context, account string) *vo.UserVo

	// FindByAdminAndStatus 根据管理员状态和用户状态查找用户列表
	FindByAdminAndStatus(ctx context.Context, admin bool, status enums.UserStatus, filter string) []*vo.UserVo
}
type UserDaoImpl struct {
	dao idao.Dao[*vo.UserVo]
}

var DBKey string = "mysql"

func NewUserDao() UserDao {
	return &UserDaoImpl{
		dao: dao.NewDao[*vo.UserVo](&dao.NewConfig{DBKey: DBKey}),
	}
}

func (u *UserDaoImpl) FindByAccountAndPassword(ctx context.Context, account, password string) *vo.UserVo {
	bu := rsql.NewBuilder()
	filter := bu.And(bu.Eq("account", account), bu.Eq("password", password)).Build()
	list := u.dao.FindByRSQL(ctx, filter)
	if len(list) == 0 {
		return nil
	}
	if len(list) > 1 {
		return nil
	}
	return list[0]
}

func (u *UserDaoImpl) FindByIdAndPassword(ctx context.Context, id, password string) *vo.UserVo {
	bu := rsql.NewBuilder()
	newRSQL := bu.And(bu.Eq("id", id), bu.Eq("password", password)).Build()
	one := u.dao.FindOneByRSQL(ctx, newRSQL)
	return one
}

func (u *UserDaoImpl) FindUsers(ctx context.Context, usersId []string, filter string) []*vo.UserVo {
	bu := rsql.NewBuilder()
	newRSQL := bu.And(bu.In("id", usersId), bu.RSQL(filter)).Build()
	list := u.dao.FindByRSQL(ctx, newRSQL)
	return list
}

func (u *UserDaoImpl) FindByStatus(ctx context.Context, status enums.UserStatus) []*vo.UserVo {
	bu := rsql.NewBuilder()
	newRSQL := bu.And(bu.Eq("status", status)).Build()
	return u.dao.FindByRSQL(ctx, newRSQL)
}

func (u *UserDaoImpl) FindByUserType(ctx context.Context, userType enums.UserType) []*vo.UserVo {
	bu := rsql.NewBuilder()
	newRSQL := bu.And(bu.Eq("type", userType)).Build()
	return u.dao.FindByRSQL(ctx, newRSQL)
}

func (u *UserDaoImpl) FindByAccount(ctx context.Context, account string) *vo.UserVo {
	bu := rsql.NewBuilder()
	newRSQL := bu.And(bu.Eq("account", account)).Build()
	return u.dao.FindOneByRSQL(ctx, newRSQL)
}

func (u *UserDaoImpl) FindByAdminAndStatus(ctx context.Context, admin bool, status enums.UserStatus, filter string) []*vo.UserVo {
	bu := rsql.NewBuilder()
	newRSQL := bu.And(bu.Eq("admin", admin), bu.Eq("status", status), bu.RSQL(filter)).Build()
	return u.dao.FindByRSQL(ctx, newRSQL)
}
