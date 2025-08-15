package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_CreateManyTranDetail(t *testing.T) {
	xtest.InitEnv_MongoRemoteMaster(xtest.NewMongoOptions().SetDBName("master"))
	ctx := xtest.NewContext()
	tranDetailService := NewTranDetailService()
	recordService := NewRecordService()
	records, err := recordService.FindAll(ctx, &query.RecordFindAllQuery{TenantId: "test"})
	if err != nil {
		t.Error(err)
		return
	}

	details := make([]*model.Tran, len(records))
	for i, record := range records {
		details[i] = model.NewTranFromRecord(record)
	}

	err = tranDetailService.CreateMany(ctx, details)
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_FindThresholdQuery(t *testing.T) {
	xtest.InitEnv_MongoRemoteMaster(xtest.NewMongoOptions().SetDBName("master"))
	ctx := xtest.NewContext()
	tranDetailService := NewTranDetailService()
	qry := &query.TranDetailFindThresholdQuery{}
	qry.Name = "赵蕾"
	qry.StartDate, _ = timeutils.AsTime("2023-08-10")
	qry.EndDate, _ = timeutils.AsTime("2023-08-11")
	qry.Amount = 100.0
	items, err := tranDetailService.FindThresholdQuery(ctx, qry)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("items.count:", len(items))
}
