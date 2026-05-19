package dao

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_FindGraphByDrawId(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest2.NewContext()
		env.SetEnv(xtest2.NewEnvConfigNeo4j())
		dao := NewGraphDao("neo4j")
		graph := dao.FindGraphByDrawId(ctx, "1001", "XSj4nnZikcd02N3ntBVeSKugYX")
		t.Log(graph)
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}
