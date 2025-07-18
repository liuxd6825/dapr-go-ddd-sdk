package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindById(t *testing.T) {
	xtest.InitEnv_MongoRemoteMaster()
	ctx := xtest.NewContext()

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
