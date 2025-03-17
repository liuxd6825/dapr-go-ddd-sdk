package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

type DaoBase[T any] struct {
	cfg         *idao.DaoConfig
	dbKey       string                // 配置中的数据库Key
	tableName   string                // 表名
	appId       string                // 应用ID
	aggField    string                // 聚合根字段
	isPubEvent  bool                  // 是否发布事件
	eventPrefix string                // 事件前缀
	dao         ddd_repository.Dao[T] // 数据访问
	env         idao.IEnvConfig       // 环境变量
}

const (
	CreatedTime = "createdTime"
	CreatorId   = "creatorId"
	CreatorName = "creatorName"
	UpdatedTime = "updatedTime"
	UpdaterId   = "updaterId"
	UpdaterName = "updaterName"
	DeletedTime = "deletedTime"
	DeleterId   = "deleterId"
	DeleterName = "deleterName"
	IsDeleted   = "isDeleted"
	TenantId    = "tenantId"
	Id          = "id"
)

func NewDaoBase[T any](dao ddd_repository.Dao[T], cfg *idao.DaoConfig) *DaoBase[T] {
	if cfg == nil {
		panic("dao base config is nil")
	}
	aggField := cfg.AggField
	if aggField == "" {
		aggField = "id"
	}
	tableName := stringutils.AsFieldName(cfg.Schema.Name)
	return &DaoBase[T]{
		dao:         dao,
		dbKey:       cfg.DbKey,
		tableName:   tableName,
		appId:       cfg.GetEnv().GetAppId(),
		isPubEvent:  cfg.GetIsPubEvent(),
		aggField:    aggField,
		eventPrefix: "eventPrefix",
		env:         cfg.Env,
	}
}

func (d *DaoBase[T]) GetEnv() restapp.IEnvConfig {
	return d.env
}

func (d *DaoBase[T]) GetDbKey() string {
	return d.dbKey
}

func (d *DaoBase[T]) GetTableName() string {
	return d.tableName
}

func (d *DaoBase[T]) SetAggField(val string) {
	d.aggField = val
}

func (d *DaoBase[T]) GetAggField() string {
	return d.aggField
}

func (d *DaoBase[T]) GetIsPubEvent() bool {
	return d.isPubEvent
}

func (d *DaoBase[T]) GetEventPrefix() string {
	return d.eventPrefix
}

func (d *DaoBase[T]) GetSchema() *dbschema.Schema {
	return d.cfg.Schema
}

func (d *DaoBase[T]) GetAggregateId(entity T, opts *idao.CallOptions) (string, error) {
	aggId := d.dao.GetAggId(entity)
	return aggId, nil
}

func (d *DaoBase[T]) GetTenantId(ctx context.Context) string {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user.GetTenantId()
	}
	panic("token is error")
}

func (d *DaoBase[T]) GetAuthUser(ctx context.Context) appctx.AuthUser {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user
	}
	panic("token is error")
}

func (d *DaoBase[T]) NewFindPagingQuery(ctx context.Context, findPagingMap any) ddd_repository.FindPagingQuery {
	var qry ddd_repository.FindPagingQuery
	if v, ok := findPagingMap.(ddd_repository.FindPagingQuery); ok {
		qry = v
	} else if mapData, ok := findPagingMap.(map[string]any); ok {
		builder := ddd_repository.NewFindPagingQueryBuilder()
		qry = builder.SetMapToQuery(mapData).Build()
	} else if query, ok := findPagingMap.(ddd_repository.FindPagingQuery); ok {
		qry = query
	} else {
		panic("FindPaging(findPagingMap:any) findPagingMap is map[string]any or ddd_repository.FindPagingQuery ")
	}
	tenantId := d.GetTenantId(ctx)
	qry.SetTenantId(tenantId)
	return qry
}
