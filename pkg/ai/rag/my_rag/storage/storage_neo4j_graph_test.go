package storage

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_FindNodes(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		dao := NewNeo4jGraphStorage("neo4j")

		names := []string{"张三"}
		query := GraphQueryParam{
			Keys: names,
		}
		graph := dao.FindNodes(ctx, query, Options{
			TenantId:  "test",
			CaseId:    "1001",
			NodeLabel: "all",
		})

		t.Log(graph)

		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_GraphEntity(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		dao := NewNeo4jGraphStorage("neo4j")
		entity, err := dao.GraphEntity(ctx, "孙悟空", Options{
			TenantId:  "test",
			CaseId:    "1001",
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
		ctx := xtest.NewContext()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		dao := NewNeo4jGraphStorage("neo4j")
		entity, err := dao.GraphRelationship(ctx, "孙悟空", "唐僧", Options{
			TenantId:  "test",
			CaseId:    "1001",
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
