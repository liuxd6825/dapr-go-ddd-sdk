package service

import (
	"testing"

	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_NewCode(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	ctx := xtest2.NewContext()
	service := NewCodeService()
	code := service.NewCode(ctx, "MC")
	t.Log(code)
}
