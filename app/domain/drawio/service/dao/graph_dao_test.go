package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_FindGraphByDrawId(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		dao := NewGraphDao()
		nodes, rels := dao.FindGraphByDrawId(ctx, "XSj4nnZikcd02N3ntBVeSKugYX")
		t.Log(nodes, rels)
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}
