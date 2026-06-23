package service

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_CreateManyTranDetail(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster(xtest2.NewMongoOptions().SetDBName("master"))
	ctx := xtest2.NewContext()
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
	xtest2.InitEnv_MongoRemoteMaster(xtest2.NewMongoOptions().SetDBName("master"))
	ctx := xtest2.NewContext()
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
