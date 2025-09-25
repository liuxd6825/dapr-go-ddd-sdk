package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_Record(t *testing.T) {
	xtest.InitEnv_MongoRemoteMaster()
	ctx := xtest.NewContext()
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
