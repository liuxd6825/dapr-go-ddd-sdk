package dao

import (
	"testing"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_Insert(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	ctx := xtest2.NewContext()
	dao := NewSuTaskDao("db")

	suTask := model2.NewSuTask()
	suTask.Id = "0001"
	suTask.Name = "测试"
	suTask.Code = "0001"
	suTask.OwnerName = "管理员"
	suTask.OwnerId = "admin"
	suTask.TargetId = "c001"
	suTask.TargetName = "北京XX科技公司"
	suTask.Status = model2.SuTaskStatus_New

	if result := dao.Create(ctx, suTask); result.Error != nil {
		t.Error(result.Error.Error())
	}
}
