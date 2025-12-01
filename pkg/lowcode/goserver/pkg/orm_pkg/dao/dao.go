package dao

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver/pkg/schema_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver/pkg/server_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/liuxd6825/jsonschema/v6"
)

type Dao struct {
	idao.Dao[map[string]any]
}

var daoMap *types.CMap[idao.Dao[map[string]any]]

type NewDaoConfig struct {
	DbKey      string             `json:"dbKey"`
	IsPubEvent *bool              `json:"isPubEvent"`
	Schema     *jsonschema.Schema `json:"schema"`
}

func NewDao(dbKey string, tableName string) *Dao {
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, map[string]any{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[map[string]any](newCfg)
	daoVal := &Dao{Dao: baseDao}
	return daoVal
}

func NewRSQLBuilder() *rsql.Builder {
	return rsql.NewBuilder()
}

func NewDao2(schFile string) idao.Dao[map[string]any] {
	sch := schema_pkg.LoadFile(schFile, "")
	aggField, aggType, tableName, _, dbKey := getInfos(sch)
	daoKey := GetKey(dbKey, tableName)
	if v, ok := daoMap.Get(daoKey); ok {
		return v
	}
	dbSch := dbschema.NewDBSchemaWithJsonSchema(sch)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		AggField:  aggField,
		AggType:   aggType,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	vDao := dao.NewDao[map[string]any](newCfg)
	daoMap.Set(daoKey, vDao)
	return vDao
}

func NewDaoWithCfg(cfg *NewDaoConfig) idao.Dao[map[string]any] {
	aggField, aggType, tableName, _, dbKey := getInfos(cfg.Schema)
	dbSch := dbschema.NewDBSchemaWithJsonSchema(cfg.Schema)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		AggField:  aggField,
		AggType:   aggType,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	return dao.NewDao[map[string]any](newCfg)
}

func GetDao(dbKey string, tableName string) idao.Dao[map[string]any] {
	daoKey := GetKey(dbKey, tableName)
	if v, ok := daoMap.Get(daoKey); ok {
		return v
	}
	panic("not found dao " + daoKey)
	return nil
}

func GetKey(dbKey string, tableName string) string {
	if dbKey == "" {
		dbKey = env.DefaultDBKey
	} else {
		dbKey = dao.GetDbKey(server_pkg.GetServer().Env(), dbKey)
	}
	return fmt.Sprintf("%s.%s", dbKey, tableName)
}

func getInfos(sch *jsonschema.Schema) (aggField string, aggType string, tableName string, isPubEvent bool, dbKey string) {
	meta := schema.GetMetaExtension(sch)
	tableName = sch.Name()
	isPubEvent = false
	if meta != nil && meta.DBTable != nil {
		tableName = meta.DBTable.Name
		dbKey = meta.DBTable.DBKey
	}
	if meta != nil && meta.DDD != nil {
		aggField = meta.DDD.AggField
		aggType = meta.DDD.AggType
		isPubEvent = meta.DDD.IsPubEvent
	}
	return
}

func getDbKey(envInst *env.Env, dbKey string) string {
	if dbKey == "" {
		return env.DefaultDBKey
	}
	return envInst.GetDBKeyValue(dbKey)
}
