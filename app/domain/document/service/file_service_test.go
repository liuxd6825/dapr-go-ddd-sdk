package service

import (
	"testing"

	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
	"github.com/stretchr/testify/assert"
)

func Test_FindById(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	ctx := xtest2.NewContext()

	service := NewFileService()
	file, err := service.FindById(ctx, "FtzV5g8WKbJWXInZtbL9YxJ2t")
	assert.NoError(t, err)
	assert.NotNil(t, file)
	t.Log(file.Id)

	data, err := service.ReadByteByFile(ctx, file)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	t.Log(string(data))
}
