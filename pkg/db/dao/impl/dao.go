package impl

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

type DaoBase[T any] struct {
	cfg         *idao.DaoConfig
	dbKey       string          // 配置中的数据库Key
	tableName   string          // 表名
	appId       string          // 应用ID
	aggField    string          // 聚合根字段
	aggType     string          // 聚合根类型名称
	eventPrefix string          // 事件前缀
	store       store.IStore[T] // 数据访问
	env         *env.Env        // 环境变量
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

func NewDaoBase[T any](store store.IStore[T], cfg *idao.DaoConfig) *DaoBase[T] {
	if cfg == nil {
		panic("dao base config is nil")
	}
	aggField := cfg.AggField
	if aggField == "" {
		aggField = "id"
	}
	tableName := stringutils.AsFieldName(cfg.DBSchema.Name)
	return &DaoBase[T]{
		store:       store,
		cfg:         cfg,
		dbKey:       cfg.DBKey,
		tableName:   tableName,
		appId:       cfg.GetEnv().App.AppId,
		aggField:    aggField,
		aggType:     cfg.AggType,
		eventPrefix: "eventPrefix",
	}
}

func (d *DaoBase[T]) GetEnv() *env.Env {
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

func (d *DaoBase[T]) GetEventPrefix() string {
	return d.eventPrefix
}

func (d *DaoBase[T]) GetSchema() *store.DBSchema {
	return d.cfg.DBSchema
}

func (d *DaoBase[T]) GetStore() store.IStore[T] {
	return d.store
}

func (d *DaoBase[T]) GetAggId(entity T, opts idao.CallOptions) (string, error) {
	aggId := d.store.GetAggId(entity)
	return aggId, nil
}

func (d *DaoBase[T]) GetTenantId(ctx context.Context, opts ...idao.CallOptions) string {
	opt := idao.NewCallOptions(opts...)
	return opt.GetTenantId2(ctx)
}

func (d *DaoBase[T]) GetAuthUser(ctx context.Context) appctx.AuthUser {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user
	}
	panic("token is error")
}

func (d *DaoBase[T]) NewFindPagingQuery(ctx context.Context, findPagingMap any) store.FindPagingQuery {
	var qry store.FindPagingQuery
	if v, ok := findPagingMap.(store.FindPagingQuery); ok {
		qry = v
	} else if mapData, ok := findPagingMap.(map[string]any); ok {
		builder := store.NewFindPagingQueryBuilder()
		qry = builder.SetMapToQuery(mapData).Build()
	} else if query, ok := findPagingMap.(store.FindPagingQuery); ok {
		qry = query
	} else {
		panic("FindPaging(findPagingMap:any) findPagingMap is map[string]any or ddd_repository.FindPagingQuery ")
	}
	return qry
}

func (d *DaoBase[T]) GetDbType() string {
	return d.store.GetDbType()
}
