package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_Insert(t *testing.T) {
	xtest.InitEnv_MongoRemoteMaster()
	ctx := xtest.NewContext()
	dao := NewSuTaskDao("db")

	suTask := model.NewSuTask()
	suTask.Id = "0001"
	suTask.Name = "测试"
	suTask.Code = "0001"
	suTask.OwnerName = "管理员"
	suTask.OwnerId = "admin"
	suTask.TargetId = "c001"
	suTask.TargetName = "北京XX科技公司"
	suTask.Status = model.SuTaskStatus_New

	if result := dao.Create(ctx, suTask); result.Error != nil {
		t.Error(result.Error.Error())
	}
}
