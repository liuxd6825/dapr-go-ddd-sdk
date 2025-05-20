package neo4jservice

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/neo4jservice/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/neo4jservice/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/jsonschema/v6"
)

type Neo4jService struct {
	nodeDaoMap *types.CMap[*dao.NodeDao]
	relDaoMap  *types.CMap[*dao.RelationDao]
}

func NewNeo4jService() *Neo4jService {
	ser := &Neo4jService{
		nodeDaoMap: types.NewCMap[*dao.NodeDao](),
		relDaoMap:  types.NewCMap[*dao.RelationDao](),
	}
	ser.Init()
	return ser
}

func (s *Neo4jService) Create(record *model.Record) {
	ctx := s.newCtx(record)
	gp.Try(func() error {
		tableName := record.Table
		nodeDao := s.getNodeDao(tableName)
		if nodeDao == nil {
			return nil
		}

		afterData := record.AfterMap()
		node := s.NewNode(tableName, afterData)
		// 是主数据
		if record.IsMaster() {
			nodeDao.MergeByName(ctx, node, node.Name)
		} else { //是关系数据
			relDao := s.getRelDao(record)
			if relDao == nil {
				return nil
			}
			rel := model.NewRelation(relDao.DBSchema, afterData)
			nodeDao.CreateByRelationData(ctx, rel, node)
		}
		return nil
	}).Catch(func(e error) {
		logs.Errorf(ctx, nil, e.Error())
	})
}

func (s *Neo4jService) Update(record *model.Record) {
	nodeDao := s.getNodeDao(record.Table)
	if nodeDao == nil {
		return
	}
	ctx := s.newCtx(record)
	afterMap := record.AfterMap()
	afterNode := s.NewNode(record.Table, record.AfterMap())
	dbSch := nodeDao.DBSchema
	if record.IsMaster() {
		// 是主数据表
		nodeDao.Update(ctx, afterNode)
	} else {
		// 是关系数据表
		relDao := s.getRelDao(record)
		if relDao == nil {
			return
		}
		afterRel := model.NewRelation(dbSch, afterMap)

		// 是否改名
		if record.IsRename() {
			beforeMap := record.BeforeMap()
			beforeRel := model.NewRelation(dbSch, beforeMap)
			relDao.Delete(ctx, beforeRel) //删除旧关系
		}
		// 创建或更新关系与节点
		relNode := model.NewNode(dbSch.TableName, afterMap)
		nodeDao.MergeAndCreateRel(ctx, afterNode, afterRel, relNode)
	}
}

func (s *Neo4jService) Delete(record *model.Record) {
	tableName := record.Table
	nodeDao := s.getNodeDao(tableName)
	if nodeDao == nil {
		return
	}
	ctx := s.newCtx(record)
	if record.IsMaster() {
		node := model.NewNode(tableName, record.AfterMap())
		nodeDao.Delete(ctx, node)
	} else {
		nodeDao.DeleteNodeAndRelationByRecord(ctx, record)
	}
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

func (s *Neo4jService) newCtx(record *model.Record) context.Context {
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

func (s *Neo4jService) getNodeDao(tableName string) *dao.NodeDao {
	get, ok := s.nodeDaoMap.Get(tableName)
	if !ok {
		return nil
	}
	return get
}

func (s *Neo4jService) getRelDao(record *model.Record) *dao.RelationDao {
	get, ok := s.relDaoMap.Get(record.Table)
	if !ok {
		return nil
	}
	return get
}

func (s *Neo4jService) addDao(sch *jsonschema.Schema) {
	meta := schema.GetMetaExtension(sch)
	tableName := meta.DBTable.Name

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
	relStartField := dbSch.GetRelStartIdField()
	relEndField := dbSch.GetRelEndIdField()
	relTypeField := dbSch.GetRelTypeField()
	nodeLabelFields := dbSch.GetNodeLabelFields()
	println("relStartField", relStartField, "relEndField", relEndField, "relTypeField", relTypeField, "nodeLabelFields", nodeLabelFields)

	if isNode := idao.IsGraphType(graphTypes, idao.GraphType_Node); isNode {
		labels, _ := maputils.GetStrings(cfg, "labels", []string{tableName})
		labels = append(labels, "master")
		nodeDao := dao.NewNodeDao(labels, dbSch)
		s.nodeDaoMap.Add(tableName, nodeDao)
	}

	if isRel := idao.IsGraphType(graphTypes, idao.GraphType_Rel); isRel {
		relDao := dao.NewRelationDao(dbSch)
		s.relDaoMap.Add(tableName, relDao)
	}
}

func newEntity(dbSch *dbschema.DBSchema, src map[string]any) map[string]any {
	target := map[string]any{}
	for _, field := range dbSch.Fields {
		v, ok := src[field.DBName]
		if ok {
			target[field.Name] = v
		}
	}
	return target
}

func (s *Neo4jService) NewNode(tableName string, data map[string]any) *model.Node {
	id, _ := maputils.GetString(data, "id", "")
	tenantId, _ := maputils.GetString(data, "tenantId", "")
	name, _ := maputils.GetString(data, "name", "")
	caseId, _ := maputils.GetString(data, "caseId", "")
	return &model.Node{
		Id:       id,
		Name:     name,
		CaseId:   caseId,
		TenantId: tenantId,
		Table:    tableName,
	}
}
