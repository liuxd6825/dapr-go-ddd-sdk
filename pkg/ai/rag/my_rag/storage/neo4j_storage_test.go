package storage

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/randomutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/sirupsen/logrus"
)

var (
	tenantId = "test"
	caseId   = "1001"
	docId    = "aa"
)

func Test_FindNodes(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest2.NewContext()
		logger := logrus.New()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		dao := NewNeo4jGraphStorage("neo4j", logger)

		names := []string{"张三"}
		query := GraphQueryParam{
			Keys: names,
		}
		nodes, edges, err := dao.FindNodes(ctx, query, Options{
			TenantId:  "test",
			CaseId:    "1001",
			NodeLabel: "all",
		})
		if err != nil {
			return err
		}

		t.Log(nodes)
		t.Log(edges)

		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_DeleteDoc(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest2.NewContext()
		logger := logrus.New()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		dao := NewNeo4jGraphStorage("neo4j", logger)
		graph := dao.DeleteDoc(ctx, tenantId, caseId, "aa")
		t.Log(graph)
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_GraphEntity(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest2.NewContext()
		logger := logrus.New()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		dao := NewNeo4jGraphStorage("neo4j", logger)
		entity, err := dao.GraphEntity(ctx, "孙悟空", Options{
			TenantId:  tenantId,
			CaseId:    caseId,
			NodeLabel: "doc",
		})
		if err == nil {
			t.Log(entity)
		}
		return err
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_GraphRelationship(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest2.NewContext()
		logger := logrus.New()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		dao := NewNeo4jGraphStorage("neo4j", logger)
		entity, err := dao.GraphRelationship(ctx, "孙悟空", "唐僧", Options{
			TenantId:  tenantId,
			CaseId:    caseId,
			NodeLabel: "doc",
		})
		if err == nil {
			t.Log(entity)
		}
		return err
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_graphSaveDocEntities(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest2.NewContext()
		logger := logrus.New()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		store := NewNeo4jGraphStorage("neo4j", logger)
		var entities []*GraphEntity
		var rels []*GraphRelationship

		for i := 0; i < 10; i++ {
			name := randomutils.NameCN()
			if i == 0 {
				name = "路薇"
			}
			e := &GraphEntity{
				TenantId:     tenantId,
				CaseId:       caseId,
				DocId:        docId,
				Name:         name,
				Type:         getType(),
				Descriptions: randomutils.String(20),
				SourceIDs:    "",
				CreatedAt:    randomutils.Time(),
			}
			e.Id = e.Name
			entities = append(entities, e)
		}
		for i := 0; i < 10; i++ {
			j := i + 1
			if j == 10 {
				j = 0
			}
			source := entities[i].Name
			target := entities[j].Name
			rel := &GraphRelationship{
				Source:       source,
				Target:       target,
				CaseId:       caseId,
				DocId:        docId,
				Descriptions: randomutils.String(20),
			}
			rels = append(rels, rel)
		}
		err := store.GraphSaveDoc(ctx, tenantId, caseId, "aa", entities, rels)
		return err
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func getType() string {
	id := randomutils.IntMax(3)
	switch id {
	case 1:
		return "人员"
	case 2:
		return "公司"
	case 3:
		return "合同"
	}
	return "其他"
}
