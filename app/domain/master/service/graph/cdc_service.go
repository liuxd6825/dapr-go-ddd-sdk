package graph

import (
	"context"
	dao2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/dao"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/model"
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

type CdcService struct {
	nodeDaoMap *types.CMap[*dao2.NodeDao]
	relDaoMap  *types.CMap[*dao2.BusRelationDao]
}

func NewCdcService() *CdcService {
	ser := &CdcService{
		nodeDaoMap: types.NewCMap[*dao2.NodeDao](),
		relDaoMap:  types.NewCMap[*dao2.BusRelationDao](),
	}
	return ser
}

func (s *CdcService) Create(record *model2.Record) {
	ctx := s.newCtx(record)
	gp.Try(func() error {
		tableName := record.Table
		nodeDao := s.getNodeDao(tableName)
		if nodeDao == nil {
			return nil
		}

		afterData := record.AfterMap()
		node := model2.NewNode(tableName, afterData)
		if record.IsMaster() {
			nodeDao.CreateMain(ctx, node)
		} else if record.IsRelation() {
			// 是关系数据
			relDao := s.getRelDao(record)
			if relDao == nil {
				return nil
			}
			rel := model2.NewRelation(relDao.DBSchema, afterData)
			nodeDao.CreateRelNode(ctx, rel, node)
		}
		return nil
	}).Catch(func(e error) {
		logs.Errorf(ctx, nil, e.Error())
	})
}

func (s *CdcService) Update(record *model2.Record) {
	nodeDao := s.getNodeDao(record.Table)
	if nodeDao == nil {
		return
	}
	ctx := s.newCtx(record)
	afterMap := record.AfterMap()
	afterNode := model2.NewNode(record.Table, afterMap)
	if record.IsMaster() {
		// 是主数据表
		nodeDao.UpdateMain(ctx, afterNode)
	} else {
		// 更新关系与节点
		nodeDao.UpdateRelNode(ctx, record)
	}
}

func (s *CdcService) Delete(record *model2.Record) {
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

func (s *CdcService) Init() *CdcService {
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
	return s
}

func (s *CdcService) newCtx(record *model2.Record) context.Context {
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

func (s *CdcService) getNodeDao(tableName string) *dao2.NodeDao {
	get, ok := s.nodeDaoMap.Get(tableName)
	if !ok {
		return nil
	}
	return get
}

func (s *CdcService) getRelDao(record *model2.Record) *dao2.BusRelationDao {
	get, ok := s.relDaoMap.Get(record.Table)
	if !ok {
		return nil
	}
	return get
}

func (s *CdcService) AddDao(sch *jsonschema.Schema) {
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

	var nodeDao *dao2.NodeDao
	if isNode := idao.IsGraphType(graphTypes, idao.GraphType_Node); isNode {
		labels, _ := maputils.GetStrings(cfg, "labels", []string{tableName})
		labels = append(labels, "master")
		nodeDao = dao2.NewNodeDao(labels, dbSch)
		s.nodeDaoMap.Add(tableName, nodeDao)
	}

	if isRel := idao.IsGraphType(graphTypes, idao.GraphType_Rel); isRel {
		relDao := dao2.NewBusRelationDao(dbSch, nodeDao)
		s.relDaoMap.Add(tableName, relDao)
	}
}

func (s *CdcService) clearAll(ctx context.Context) {
	var nodeDao *dao2.NodeDao
	for _, d := range s.nodeDaoMap.Items() {
		nodeDao = d
		break
	}
	nodeDao.ClearAll(ctx)
}
