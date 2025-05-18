package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
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

type Base struct {
	UpdatedTime *time.Time `json:"updatedTime" gorm:"column:updated_time"`
	UpdaterName string     `json:"updaterName" gorm:"column:updater_name"`
	UpdaterId   string     `json:"updaterId" gorm:"column:updater_id"`
	CreatedTime *time.Time `json:"createdTime" gorm:"column:created_time"`
	CreatorName string     `json:"creatorName" gorm:"column:creator_name"`
	CreatorId   string     `json:"creatorId" gorm:"column:creator_id"`
}
type Node struct {
	Id       string `json:"id" gorm:"column:id"`
	CaseId   string `json:"caseId" gorm:"column:case_id;nodeLabel:true"`
	Name     string `json:"name" gorm:"column:name"`
	TenantId string `json:"tenantId" gorm:"column:tenant_id;nodeLabel:true"`
}

type Relation struct {
	Id       string `json:"id" gorm:"column:id"`
	CaseId   string `json:"caseId" gorm:"column:case_id"`
	RelType  string `json:"relType" gorm:"column:rel_type;relType:true"`
	StartId  string `json:"startId" gorm:"column:start_id;relStartId:true"`
	EndId    string `json:"endId" gorm:"column:end_id;relEndId:true"`
	TenantId string `json:"tenantId" gorm:"column:tenant_id"`
}

type Neo4jService struct {
	nodeDaoMap  *types.CMap[idao.Dao[*Node]]
	relDaoMap   *types.CMap[idao.Dao[*Relation]]
	dbSchemaMap *types.CMap[*dbschema.DBSchema]
}

func NewNeo4jService() *Neo4jService {
	ser := &Neo4jService{
		nodeDaoMap:  types.NewCMap[idao.Dao[*Node]](),
		relDaoMap:   types.NewCMap[idao.Dao[*Relation]](),
		dbSchemaMap: types.NewCMap[*dbschema.DBSchema](),
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
	tableName := record.Table
	nodeDao := s.getNodeDao(tableName)
	if nodeDao == nil {
		return
	}

	entity := newEntity(nodeDao, s.dbSchemaMap, tableName, record.After)
	ctx := s.newCtx(record)
	node := s.NewNode(entity)
	nodeDao.CreateUpdate(ctx, node)

	relDao := s.getRelDao(record)
	if relDao == nil {
		return
	}
	relation := s.NewRelation(entity, tableName)
	relDao.CreateUpdate(ctx, relation)

}

func (s *Neo4jService) Update(record *RawRecord) {
	tableName := record.Table
	nodeDao := s.getNodeDao(record.Table)
	if nodeDao == nil {
		return
	}

	entity := newEntity(nodeDao, s.dbSchemaMap, tableName, record.After)
	ctx := s.newCtx(record)
	node := s.NewNode(entity)
	nodeDao.Update(ctx, node)

	relDao := s.getRelDao(record)
	if relDao != nil {
		beforeEntity := newEntity(relDao, s.dbSchemaMap, tableName, record.Before)
		afterRel := s.NewRelation(entity, tableName)
		beforeRel := s.NewRelation(beforeEntity, tableName)
		if afterRel.RelType != beforeRel.RelType {
			relDao.Delete(ctx, beforeRel)
			relDao.Create(ctx, afterRel)
		}
	}

}

func (s *Neo4jService) Delete(record *RawRecord) {
	tableName := record.Table
	nodeDao := s.getNodeDao(record.Table)
	if nodeDao == nil {
		return
	}

	entity := newEntity(nodeDao, s.dbSchemaMap, tableName, record.After)
	node := s.NewNode(entity)
	ctx := s.newCtx(record)
	nodeDao.Delete(ctx, node)

	relDao := s.getRelDao(record)
	if relDao != nil {
		entity := newEntity(relDao, s.dbSchemaMap, tableName, record.After)
		rel := s.NewRelation(entity, tableName)
		relDao.Delete(ctx, rel)
	}
}

func (s *Neo4jService) getNodeDao(tableName string) idao.Dao[*Node] {
	dao, ok := s.nodeDaoMap.Get(tableName)
	if !ok {
		return nil
	}
	return dao
}

func (s *Neo4jService) getRelDao(record *RawRecord) idao.Dao[*Relation] {
	dao, ok := s.relDaoMap.Get(record.Table)
	if !ok {
		return nil
	}
	return dao
}

func (s *Neo4jService) addDao(sch *jsonschema.Schema) {
	meta := schema.GetMetaExtension(sch)
	tableName := meta.DBTable.Name
	if tableName == "company" {
		println(tableName)
	}
	if _, ok := s.nodeDaoMap.Get(tableName); ok {
		return
	}
	cfg, err := maputils.GetMap(meta.Attributes, "neo4j", nil)
	if err != nil {
		panic(err)
	}
	if cfg == nil {
		return
	}

	graphTypes, _ := maputils.GetStrings(cfg, "type", nil)
	if len(graphTypes) == 0 {
		return
	}

	dbSch := dbschema.NewDBSchemaWithJsonSchema(sch)
	if dbSch != nil {
		s.dbSchemaMap.Add(tableName, dbSch)
	}

	if isNode := idao.IsGraphType(graphTypes, idao.GraphType_Node); isNode {
		labels, _ := maputils.GetStrings(cfg, "labels", []string{tableName})
		nodeCfg := &dao.NewConfig{
			DBKey:              "neo4j",
			IsPubEvent:         dao.IsFalse(),
			GraphType:          idao.GraphType_Node,
			GraphLabels:        labels,
			IsCancelModified:   true,
			IsCancelSoftDelete: true,
		}
		nodeDao := dao.NewDao[*Node](nodeCfg)
		s.nodeDaoMap.Add(tableName, nodeDao)
	}

	if isRel := idao.IsGraphType(graphTypes, idao.GraphType_Rel); isRel {
		relCfg := &dao.NewConfig{
			DBKey:              "neo4j",
			IsPubEvent:         dao.IsFalse(),
			GraphType:          idao.GraphType_Rel,
			IsCancelModified:   true,
			IsCancelSoftDelete: true,
		}
		relDao := dao.NewDao[*Relation](relCfg)
		s.relDaoMap.Add(tableName, relDao)
	}
}

func newEntity[T any](dao idao.Dao[T], dbSchemaMap *types.CMap[*dbschema.DBSchema], tableName string, src map[string]any) map[string]any {
	target := map[string]any{}
	dbSch, ok := dbSchemaMap.Get(tableName)
	if !ok {
		return target
	}
	for _, field := range dbSch.Fields {
		v, ok := src[field.DBName]
		if ok {
			target[field.Name] = v
		}
	}
	return target
}

func (s *Neo4jService) NewNode(data map[string]any) *Node {
	tenantId, _ := maputils.GetString(data, "tenantId", "")
	name, _ := maputils.GetString(data, "name", "")
	caseId, _ := maputils.GetString(data, "caseId", "")
	return &Node{
		Id:       name,
		Name:     name,
		CaseId:   caseId,
		TenantId: tenantId,
	}
}

func (s *Neo4jService) NewRelation(data map[string]any, tableName string) *Relation {
	dbSchema, ok := s.dbSchemaMap.Get(tableName)
	if !ok {
		panic(errors.New("table " + tableName + " not exist"))
	}
	tenantId, _ := maputils.GetString(data, "tenantId", "")
	caseId, _ := maputils.GetString(data, "caseId", "")
	id, _ := maputils.GetString(data, "id", "")
	relStartId, _ := getRelStartId(dbSchema, data)
	relEndId, _ := getRelEndId(dbSchema, data)
	relType, _ := getRelType(dbSchema, data)
	return &Relation{
		Id:       id,
		CaseId:   caseId,
		TenantId: tenantId,
		RelType:  relType,
		StartId:  relStartId,
		EndId:    relEndId,
	}
}

func getRelStartId(dbSch *dbschema.DBSchema, data map[string]any) (string, error) {
	field := dbSch.GetRelStartIdField()
	if field == nil {
		return "", errors.New("start id not found")
	}
	val, err := maputils.GetString(data, field.Name, "")
	if err != nil {
		return "", err
	}
	return val, nil
}

func getRelEndId(dbSch *dbschema.DBSchema, data map[string]any) (string, error) {
	field := dbSch.GetRelTypeField()
	if field == nil {
		return "", errors.New("start id not found")
	}
	val, err := maputils.GetString(data, field.Name, "")
	if err != nil {
		return "", err
	}
	return val, nil
}

func getRelType(dbSch *dbschema.DBSchema, data map[string]any) (string, error) {
	field := dbSch.GetRelTypeField()
	if field == nil {
		return "", errors.New("start id not found")
	}
	val, err := maputils.GetString(data, field.Name, "")
	if err != nil {
		return "", err
	}
	return val, nil
}
