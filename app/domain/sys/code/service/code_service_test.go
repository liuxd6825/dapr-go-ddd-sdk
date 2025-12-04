package service

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/model"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_NewCode(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	ctx := xtest2.NewContext()
	service := NewCodeService()

	/*	humanCode := service.NewHumanCode(ctx, "1001")
		t.Log(humanCode)

		caseCode := service.NewCaseCode(ctx)
		t.Log(caseCode)*/

	cmd := command.NewCodeNewCommand()
	cmd.Data.Type = "Year"
	cmd.Data.Style = model.CodeStyle_Year
	cmd.Data.NoHead = true
	cmd.Data.Count = 3
	codes, err := service.newCode(ctx, cmd)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(codes)
}
