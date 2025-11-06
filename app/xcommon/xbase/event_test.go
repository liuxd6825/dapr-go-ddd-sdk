package xbase

import (
	"testing"

	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
	"github.com/stretchr/testify/assert"
)

type RecordImportMasterEvent struct {
	Command[*TestCommandData]
}
type TestCommandData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Test_PublishEvent(t *testing.T) {
	xtest2.InitEnv_MongoLocal(xtest2.NewMongoOptions().SetDBName("master"))
	ctx := xtest2.NewContext()
	cmd := &RecordImportMasterEvent{}
	cmd.SetData(&TestCommandData{
		Id:   "1",
		Name: "张三",
		Age:  18,
	})
	meta := map[string]any{
		"key": "value",
	}
	err := PublishEvent(ctx, "test-app", cmd, meta)
	assert.Nil(t, err)
}
