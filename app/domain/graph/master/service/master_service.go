package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
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
	"time"
)

// MasterService
// @Description: 主数据关系图
type MasterService struct {
	nodeDaoMap *types.CMap[*dao.MasterNodeDao]
	relDaoMap  *types.CMap[*dao.BusRelationDao]
	dbSchMap   *types.CMap[*dbschema.DBSchema]
}

func NewMasterService() *MasterService {
	ser := &MasterService{
		nodeDaoMap: types.NewCMap[*dao.MasterNodeDao](),
		relDaoMap:  types.NewCMap[*dao.BusRelationDao](),
		dbSchMap:   types.NewCMap[*dbschema.DBSchema](),
	}
	return ser
}

func (s *MasterService) DataChange(ctx context.Context, record *model.CDCRecord) error {
	logs.InfoMsg(ctx, "record", " opType=", record.OpType, " table=", record.Table)

	dbSch := s.GetDBSchema(record.Table)
	if dbSch == nil {
		return nil
	}
	record.DBSchema = dbSch

	// 根据操作类型处理数据
	switch record.OpType {
	case "c":
		s.Create(record)
	case "u":
		s.Update(record)
	case "d":
		s.Delete(record)
	}
	return nil
}

func (s *MasterService) GetDBSchema(tableName string) *dbschema.DBSchema {
	dbSch, ok := s.dbSchMap.Get(tableName)
	if !ok {
		return nil
	}
	return dbSch
}

func (s *MasterService) Create(record *model.CDCRecord) {
	ctx := s.newCtx(record)
	gp.Try(func() error {
		tableName := record.Table
		nodeDao := s.getNodeDao(tableName)
		if nodeDao == nil {
			return nil
		}

		afterData := record.AfterMap()
		node := model.NewMasterNode(afterData, record.DBSchema)
		if record.IsMaster() {
			nodeDao.CreateMain(ctx, node)
		} else if record.IsRelation() {
			// 是关系数据
			relDao := s.getRelDao(record)
			if relDao == nil {
				return nil
			}
			rel := model.NewMasterRelation(afterData, record.DBSchema)
			nodeDao.CreateRelNode(ctx, rel, node)
		}
		return nil
	}).Catch(func(e error) {
		logs.Errorf(ctx, nil, e.Error())
	})
}

func (s *MasterService) Update(record *model.CDCRecord) {
	nodeDao := s.getNodeDao(record.Table)
	if nodeDao == nil {
		return
	}
	ctx := s.newCtx(record)
	afterMap := record.AfterMap()
	afterNode := model.NewMasterNode(afterMap, record.DBSchema)
	if record.IsMaster() {
		// 是主数据表
		nodeDao.UpdateMain(ctx, afterNode)
	} else {
		// 更新关系与节点
		nodeDao.UpdateRelNode(ctx, record)
	}
}

func (s *MasterService) Delete(record *model.CDCRecord) {
	tableName := record.Table
	nodeDao := s.getNodeDao(tableName)
	if nodeDao == nil {
		return
	}
	ctx := s.newCtx(record)
	if record.IsMaster() {
		nodeDao.DeleteMain(ctx, record, nodeDao.GetSchema())
	} else {
		nodeDao.DeleteRelNode(ctx, record)
	}
}

func (s *MasterService) Init() *MasterService {
	srcFs, err := fs_pkg.NewFsPkg(env.GetEnv(), "src")
	if err != nil {
		panic("src fs not exist")
	}
	fileInfos := srcFs.ReadAllPath("/definition/db/master")
	start := time.Now()
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir {
			continue
		}
		fileName := fileInfo.Path + "/" + fileInfo.Name
		data := srcFs.ReadFile(fileName)
		jsonSch := schema.NewJsonSchemaWithBytes(fileName, data)
		dbSch := dbschema.NewDBSchemaWithJsonSchema(jsonSch)
		if dbSch != nil && dbSch.TableName != "" {
			s.dbSchMap.Add(dbSch.TableName, dbSch)
		}
		s.AddDao(jsonSch, dbSch)
	}
	elapsed := time.Since(start)
	fmt.Printf("CDC初始化耗时: %v\n", elapsed)
	return s
}

func (s *MasterService) newCtx(record *model.CDCRecord) context.Context {
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

func (s *MasterService) getNodeDao(tableName string) *dao.MasterNodeDao {
	get, ok := s.nodeDaoMap.Get(tableName)
	if !ok {
		return nil
	}
	return get
}

func (s *MasterService) getRelDao(record *model.CDCRecord) *dao.BusRelationDao {
	get, ok := s.relDaoMap.Get(record.Table)
	if !ok {
		return nil
	}
	return get
}

func (s *MasterService) AddDao(jsonSch *jsonschema.Schema, dbSchema *dbschema.DBSchema) {
	meta := schema.GetMetaExtension(jsonSch)
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

	var nodeDao *dao.MasterNodeDao
	if isNode := idao.IsGraphType(graphTypes, idao.GraphType_Node); isNode {
		labels, _ := maputils.GetStrings(cfg, "labels", []string{tableName})
		labels = append(labels, "master")
		nodeDao = dao.NewMasterNodeDao(labels, dbSchema)
		s.nodeDaoMap.Add(tableName, nodeDao)
	}

	if isRel := idao.IsGraphType(graphTypes, idao.GraphType_Rel); isRel {
		relDao := dao.NewBusRelationDao(dbSchema, nodeDao)
		s.relDaoMap.Add(tableName, relDao)
	}
}

func (s *MasterService) clearAll(ctx context.Context) {
	var nodeDao *dao.MasterNodeDao
	for _, d := range s.nodeDaoMap.Items() {
		nodeDao = d
		break
	}
	nodeDao.ClearAll(ctx)
}
