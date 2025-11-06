package dao

import (
	"testing"

	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_Record(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	ctx := xtest2.NewContext()
	dao := NewSuRecordDao("db")
	taskId := "0001"
	if count, err := dao.CountByTaskId(ctx, taskId); err != nil {
		t.Error(err)
	} else {
		t.Log("count=", count)
	}
	if sum, err := dao.SumByTaskId(ctx, taskId); err != nil {
		t.Error(err)
	} else {
		t.Log("sum=", sum)
	}
}
