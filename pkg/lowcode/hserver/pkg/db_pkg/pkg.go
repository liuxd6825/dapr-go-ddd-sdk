package db_pkg

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/schema_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/liuxd6825/jsonschema/v6"
)

type Pkg struct {
	server    element.Server
	daoMap    *types.CMap[idao.Dao[map[string]any]]
	schemaPkg *schema_pkg.SchemaPkg
}

type NewDaoConfig struct {
	DbKey      string             `json:"dbKey"`
	IsPubEvent *bool              `json:"isPubEvent"`
	Schema     *jsonschema.Schema `json:"schema"`
}

func New(server element.Server) *Pkg {
	return &Pkg{
		server:    server,
		daoMap:    types.NewCMap[idao.Dao[map[string]any]](),
		schemaPkg: schema_pkg.New(server),
	}
}

func (p *Pkg) NewRSQLBuilder() *rsql.Builder {
	return rsql.NewBuilder()
}

func (p *Pkg) NewDao(schFile string) idao.Dao[map[string]any] {
	sch := p.schemaPkg.LoadFile(schFile, "")
	aggField, aggType, tableName, _, dbKey := p.getInfos(sch)
	daoKey := p.GetKey(dbKey, tableName)
	if v, ok := p.daoMap.Get(daoKey); ok {
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
	p.daoMap.Set(daoKey, vDao)
	return vDao
}

func (p *Pkg) NewDaoWithCfg(cfg *NewDaoConfig) idao.Dao[map[string]any] {
	aggField, aggType, tableName, _, dbKey := p.getInfos(cfg.Schema)
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

func (p *Pkg) GetDao(dbKey string, tableName string) idao.Dao[map[string]any] {
	daoKey := p.GetKey(dbKey, tableName)
	if v, ok := p.daoMap.Get(daoKey); ok {
		return v
	}
	panic("not found dao " + daoKey)
	return nil
}

func (p *Pkg) GetKey(dbKey string, tableName string) string {
	if dbKey == "" {
		dbKey = env.DefaultDBKey
	} else {
		dbKey = dao.GetDbKey(p.server.GetEnvCfg(), dbKey)
	}
	return fmt.Sprintf("%s.%s", dbKey, tableName)
}

func (p *Pkg) getInfos(sch *jsonschema.Schema) (aggField string, aggType string, tableName string, isPubEvent bool, dbKey string) {
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
