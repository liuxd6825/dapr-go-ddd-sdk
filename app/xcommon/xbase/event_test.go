package xbase

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
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
	xtest.InitEnv_MongoLocal(xtest.NewMongoOptions().SetDBName("master"))
	ctx := xtest.NewContext()
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
