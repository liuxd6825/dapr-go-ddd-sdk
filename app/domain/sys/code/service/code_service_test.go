package service

import (
	"testing"

	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_NewCode(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	ctx := xtest2.NewContext()
	service := NewCodeService()

	humanCode := service.NewHumanCode(ctx, "1001")
	t.Log(humanCode)

	caseCode := service.NewCaseCode(ctx)
	t.Log(caseCode)
}
