package impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/jsonschema/v6"
)

type DaoBase struct {
	cfg         *db.DaoConfig
	dbKey       string                             // 配置中的数据库Key
	tableName   string                             // 表名
	appId       string                             // 应用ID
	aggField    string                             // 聚合根字段
	isPubEvent  bool                               // 是否发布事件
	eventPrefix string                             // 事件前缀
	dao         ddd_repository.Dao[map[string]any] // 数据访问
	env         common.IEnvConfig                  // 环境变量
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

func NewDaoBase(dao ddd_repository.Dao[map[string]any], cfg *db.DaoConfig) *DaoBase {
	if cfg == nil {
		panic("dao base config is nil")
	}
	aggField := cfg.AggField
	if aggField == "" {
		aggField = "id"
	}
	tableName := stringutils.AsFieldName(cfg.Schema.Name)
	return &DaoBase{
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

func (d *DaoBase) GetEnv() common.IEnvConfig {
	return d.env
}

func (d *DaoBase) GetDbKey() string {
	return d.dbKey
}

func (d *DaoBase) GetTableName() string {
	return d.tableName
}

func (d *DaoBase) SetAggField(val string) {
	d.aggField = val
}

func (d *DaoBase) GetAggField() string {
	return d.aggField
}

func (d *DaoBase) GetIsPubEvent() bool {
	return d.isPubEvent
}

func (d *DaoBase) GetEventPrefix() string {
	return d.eventPrefix
}

func (d *DaoBase) GetSchema() *jsonschema.Schema {
	return d.cfg.Schema
}

func (d *DaoBase) GetAggregateId(entity map[string]any, opts *db.CallOptions) (string, error) {
	var aggId string
	if opts != nil && opts.AggId != nil {
		aggId = *opts.AggId
	} else if d.aggField != "" {
		if id, ok := entity[d.aggField].(string); ok {
			aggId = id
		} else {
			return "", errors.New(fmt.Sprintf("Aggregate field %s is not string", d.aggField))
		}
	} else {
		aggId = d.dao.GetId(entity)
	}
	return aggId, nil
}

func (d *DaoBase) GetTenantId(ctx context.Context) string {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user.GetTenantId()
	}
	panic("token is error")
}

func (d *DaoBase) GetAuthUser(ctx context.Context) appctx.AuthUser {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user
	}
	panic("token is error")
}

func (d *DaoBase) NewFindPagingQuery(ctx context.Context, findPagingMap any) ddd_repository.FindPagingQuery {
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
