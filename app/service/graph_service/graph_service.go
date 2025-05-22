package graph_service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/graph_service/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/graph_service/model"
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

type GraphService struct {
	nodeDaoMap *types.CMap[*dao.NodeDao]
	relDaoMap  *types.CMap[*dao.BusRelationDao]
}

func NewGraphService() *GraphService {
	ser := &GraphService{
		nodeDaoMap: types.NewCMap[*dao.NodeDao](),
		relDaoMap:  types.NewCMap[*dao.BusRelationDao](),
	}
	return ser
}

func (s *GraphService) Create(record *model.Record) {
	ctx := s.newCtx(record)
	gp.Try(func() error {
		tableName := record.Table
		nodeDao := s.getNodeDao(tableName)
		if nodeDao == nil {
			return nil
		}

		afterData := record.AfterMap()
		node := model.NewNode(tableName, afterData)
		if record.IsMaster() {
			nodeDao.CreateMain(ctx, node)
		} else if record.IsRelation() {
			// 是关系数据
			relDao := s.getRelDao(record)
			if relDao == nil {
				return nil
			}
			rel := model.NewRelation(relDao.DBSchema, afterData)
			nodeDao.CreateRelNode(ctx, rel, node)
		}
		return nil
	}).Catch(func(e error) {
		logs.Errorf(ctx, nil, e.Error())
	})
}

func (s *GraphService) Update(record *model.Record) {
	nodeDao := s.getNodeDao(record.Table)
	if nodeDao == nil {
		return
	}
	ctx := s.newCtx(record)
	afterMap := record.AfterMap()
	afterNode := model.NewNode(record.Table, afterMap)
	if record.IsMaster() {
		// 是主数据表
		nodeDao.UpdateMain(ctx, afterNode)
	} else {
		// 更新关系与节点
		nodeDao.UpdateRelNode(ctx, record)
	}
}

func (s *GraphService) Delete(record *model.Record) {
	tableName := record.Table
	nodeDao := s.getNodeDao(tableName)
	if nodeDao == nil {
		return
	}
	ctx := s.newCtx(record)
	if record.IsMaster() {
		nodeDao.DeleteMain(ctx, record)
	} else {
		nodeDao.DeleteRelNode(ctx, record)
	}
}

func (s *GraphService) Init() {
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
		s.AddDao(sch)
	}
}

func (s *GraphService) newCtx(record *model.Record) context.Context {
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

func (s *GraphService) getNodeDao(tableName string) *dao.NodeDao {
	get, ok := s.nodeDaoMap.Get(tableName)
	if !ok {
		return nil
	}
	return get
}

func (s *GraphService) getRelDao(record *model.Record) *dao.BusRelationDao {
	get, ok := s.relDaoMap.Get(record.Table)
	if !ok {
		return nil
	}
	return get
}

func (s *GraphService) AddDao(sch *jsonschema.Schema) {
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

	var nodeDao *dao.NodeDao
	if isNode := idao.IsGraphType(graphTypes, idao.GraphType_Node); isNode {
		labels, _ := maputils.GetStrings(cfg, "labels", []string{tableName})
		labels = append(labels, "master")
		nodeDao = dao.NewNodeDao(labels, dbSch)
		s.nodeDaoMap.Add(tableName, nodeDao)
	}

	if isRel := idao.IsGraphType(graphTypes, idao.GraphType_Rel); isRel {
		relDao := dao.NewBusRelationDao(dbSch, nodeDao)
		s.relDaoMap.Add(tableName, relDao)
	}
}

func (s *GraphService) clearAll(ctx context.Context) {
	var nodeDao *dao.NodeDao
	for _, d := range s.nodeDaoMap.Items() {
		nodeDao = d
		break
	}
	nodeDao.ClearAll(ctx)
}
