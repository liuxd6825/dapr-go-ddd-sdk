package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/jsonschema/v6"
	"time"
)

type RawRecord struct {
	DB           string         `json:"db"`     // 数据库
	Table        string         `json:"table"`  // 数据表
	Before       map[string]any `json:"before"` // 之前数据
	After        map[string]any `json:"after"`  // 之后数据
	OpType       string         `json:"opType"` // 操作状态 "r" for read/backfill, "c" for create, "u" for update, "d" for delete
	CdcTimestamp time.Time      `json:"cdcTimestamp"`
}

type Neo4jService struct {
	nodeDaoMap *types.CMap[idao.Dao[map[string]any]]
	relDaoMap  *types.CMap[idao.Dao[map[string]any]]
}

func NewNeo4jService() *Neo4jService {
	ser := &Neo4jService{
		nodeDaoMap: types.NewCMap[idao.Dao[map[string]any]](),
		relDaoMap:  types.NewCMap[idao.Dao[map[string]any]](),
	}
	ser.Init()
	return ser
}

func (s *Neo4jService) Init() {
	srcFs, err := fs_pkg.NewFsPkg(env.GetEnv(), "src")
	if err != nil {
		panic("src fs not exist")
	}
	fileInfos := srcFs.ReadAllPath("/definition/db/master")
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir {
			continue
		}
		fileName := fileInfo.Path + "/" + fileInfo.Name
		data := srcFs.ReadFile(fileName)
		sch := schema.NewJsonSchemaWithBytes(fileName, data)
		props := sch.GetAllProperties()
		for _, prop := range props {
			meta := schema.GetMetaExtension(prop)
			if meta != nil && meta.DBField != nil {
				meta.DBField.Updatable = true
				meta.DBField.Creatable = true
			}
		}
		s.addDao(sch)
	}
}
func (s *Neo4jService) newCtx(record *RawRecord) context.Context {
	parent := context.Background()
	tenantId, _ := maputils.GetString(record.After, "tenant_id", "")
	userName, _ := maputils.GetString(record.After, "updater_name", "")
	userId, _ := maputils.GetString(record.After, "updater_id", "")

	authToken := appctx.NewAuthToken()
	authToken.User.Id = userId
	authToken.User.Name = userName

	ctx := appctx.NewContextWidthAuthToken(parent, tenantId, authToken, nil)
	return ctx
}

func (s *Neo4jService) Create(record *RawRecord) {

	nodeDao := s.getNodeDao(record)
	if nodeDao == nil {
		return
	}
	entity := newEntity(nodeDao, record.After)
	ctx := s.newCtx(record)
	nodeDao.Create(ctx, entity)

	relDao := s.getRelDao(record)
	if relDao == nil {
		return
	}
	relDao.Create(ctx, entity)

}

func (s *Neo4jService) Update(record *RawRecord) {
	nodeDao := s.getNodeDao(record)
	if nodeDao == nil {
		return
	}
	afterEntity := newEntity(nodeDao, record.After)
	ctx := s.newCtx(record)
	nodeDao.Update(ctx, afterEntity)

	relDao := s.getRelDao(record)
	if relDao != nil {
		beforeEntity := newEntity(nodeDao, record.Before)
		newType, _ := maputils.GetString(afterEntity, "relationType", "")
		oldType, _ := maputils.GetString(beforeEntity, "relationType", "")
		if newType != oldType {
			relDao.Delete(ctx, beforeEntity)
			relDao.Create(ctx, afterEntity)
		}
	}

}

func (s *Neo4jService) Delete(record *RawRecord) {
	nodeDao := s.getNodeDao(record)
	if nodeDao == nil {
		return
	}

	afterEntity := newEntity(nodeDao, record.After)
	ctx := s.newCtx(record)
	nodeDao.Delete(ctx, afterEntity)

	relDao := s.getRelDao(record)
	if relDao != nil {
		relDao.Delete(ctx, afterEntity)
	}
}

func (s *Neo4jService) getNodeDao(record *RawRecord) idao.Dao[map[string]any] {
	dao, ok := s.nodeDaoMap.Get(record.Table)
	if !ok {
		return nil
	}
	return dao.(idao.Dao[map[string]any])
}

func newEntity[T any](dao idao.Dao[T], src map[string]any) map[string]any {
	target := map[string]any{}
	for _, field := range dao.GetSchema().Fields {
		v, ok := src[field.DBName]
		if ok {
			target[field.Name] = v
		}
	}
	return target
}

func (s *Neo4jService) getRelDao(record *RawRecord) idao.Dao[map[string]any] {
	dao, ok := s.relDaoMap.Get(record.Table)
	if !ok {
		return nil
	}
	return dao.(idao.Dao[map[string]any])
}

func (s *Neo4jService) addDao(sch *jsonschema.Schema) {
	meta := schema.GetMetaExtension(sch)
	tableName := meta.DBTable.Name
	if _, ok := s.nodeDaoMap.Get(tableName); ok {
		return
	}

	neo4jMap, err := maputils.GetMap(meta.DBTable.Properties, "neo4j", nil)
	if err != nil {
		panic(err)
	}
	if neo4jMap == nil {
		return
	}

	dataTypeVal, _ := maputils.GetString(neo4jMap, "type", "")
	if dataTypeVal == "" {
		return
	}
	relType, err := idao.GetRefType(dataTypeVal)
	if err != nil {
		panic(fmt.Sprintf("neo4j data type %s not exist node or rel", dataTypeVal))
	}

	nodeCfg := &idao.DaoConfig{
		DbKey:      "neo4j",
		IsPubEvent: false,
		Env:        env.GetEnv(),
		RefType:    idao.RelType_Node,
		DBSchema:   dbschema.NewDBSchemaWithJsonSchema(sch),
	}

	nodeDao := neo4j.NewDao[map[string]any](nodeCfg)
	s.nodeDaoMap.Add(tableName, nodeDao)

	if relType == idao.RelType_Rel {
		relCfg := &idao.DaoConfig{
			DbKey:      "neo4j",
			IsPubEvent: false,
			Env:        env.GetEnv(),
			RefType:    idao.RelType_Rel,
			DBSchema:   dbschema.NewDBSchemaWithJsonSchema(sch),
		}

		relDao := neo4j.NewDao[map[string]any](relCfg)
		s.relDaoMap.Add(tableName, relDao)
	}
}
