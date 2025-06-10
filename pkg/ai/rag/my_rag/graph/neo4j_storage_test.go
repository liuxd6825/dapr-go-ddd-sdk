package graph

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
		dao := NewNeo4jGraphStorage()
		names := []string{"张三"}
		graph := dao.FindNodes(ctx, "test", "1001", names, 5)
		t.Log(graph)
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}
