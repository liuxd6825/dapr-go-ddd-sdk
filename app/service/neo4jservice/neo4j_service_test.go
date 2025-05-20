package neo4jservice

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/neo4jservice/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/neo4jservice/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_NodeInsert(t *testing.T) {
	ctx := xtest.NewContext()
	env.SetEnv(xtest.NewEnvConfig_Neo4j())
	/*
		nodeCfg := &dao.NewConfig{
			DBKey:              "neo4j",
			IsPubEvent:         dao.IsFalse(),
			GraphType:          idao.GraphType_Node,
			GraphLabels:        []string{"company_test"},
			IsCancelModified:   true,
			IsCancelSoftDelete: true,
			TableName:          "company",
			Env:                envInst,
		}*/

	nodeDao := dao.NewNodeDao([]string{"company_test"}, nil)

	node := &model.Node{
		Id:       randomutils.NewId(),
		CaseId:   "test",
		Name:     "张公公",
		TenantId: "test",
	}
	t.Log(node)
	gp.Try(func() error {
		nodeDao.MergeByName(ctx, node, node.Name)
		return nil
	}).Catch(func(err error) {
		t.Error(err)
	})

	nodeDao.FindNodeAndRelationsByName(ctx, "company_test")
}

func Test_FindNodeAndRelationsByName(t *testing.T) {
	ctx := xtest.NewContext()
	env.SetEnv(xtest.NewEnvConfig_Neo4j())
	nodeDao := dao.NewNodeDao([]string{"company_test"}, nil)
	nodeDao.GetConfig().Env = xtest.NewEnvConfig_Neo4j()

	gp.Try(func() error {
		nodeDao.FindNodeAndRelationsByName(ctx, "张公公", "company_test")
		return nil
	}).Catch(func(err error) {
		t.Error(err)
	})
}

func Test_MergeAndCreateRel(t *testing.T) {
	ctx := xtest.NewContext()
	env.SetEnv(xtest.NewEnvConfig_Neo4j())
	nodeDao := dao.NewNodeDao([]string{"company_test"}, nil)
	nodeDao.GetConfig().Env = xtest.NewEnvConfig_Neo4j()

	node := &model.Node{
		Id:       "zgg",
		Name:     "张公公",
		CaseId:   "test",
		TenantId: "test",
	}

	rel := &model.Relation{
		Id:      "randomutils.NewId()",
		RelType: "等于",
		StartId: "zww",
		EndId:   "zgg",
	}

	relNode := &model.Node{
		Id:       "zww",
		CaseId:   "test",
		TenantId: "test",
		Name:     "张维为",
	}

	gp.Try(func() error {
		nodeDao.MergeAndCreateRel(ctx, node, rel, relNode)
		return nil
	}).Catch(func(err error) {
		t.Error(err)
	})
}
