package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_NodeInsert(t *testing.T) {
	ctx := xtest.NewContext()
	envInst := xtest.NewEnvConfig_Neo4j()
	nodeCfg := &dao.NewConfig{
		DBKey:              "neo4j",
		IsPubEvent:         dao.IsFalse(),
		GraphType:          idao.GraphType_Node,
		GraphLabels:        []string{"company_test"},
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
		TableName:          "company",
		Env:                envInst,
	}
	nodeDao := dao.NewDao[*Node](nodeCfg)
	node := &Node{
		Id:       randomutils.NewId(),
		CaseId:   "test",
		Name:     randomutils.NameCN(),
		TenantId: "test",
	}
	t.Log(node)
	nodeDao.CreateUpdate(ctx, node)
}
