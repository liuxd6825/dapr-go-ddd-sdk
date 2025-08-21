package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

const account = "6235822099004087593"
const taskId = "01"

func init() {
	xtest.Init(xtest.TestType_MongoLocal)
}

func Test_ClearAnalyseResults(t *testing.T) {
	ctx := xtest.NewContext()
	service := NewSuTaskService()
	if err := service.ClearAnalyseResults(ctx, taskId); err != nil {
		t.Error(err)
		return
	}
}

func Test_analyse(t *testing.T) {
	startTime := timeutils.NewDate(2018, 1, 1)
	endTime := timeutils.NewDate(2024, 1, 1)

	ctx := xtest.NewContext()
	service := NewSuTaskService()
	task := &model.SuTask{
		Rules: model.SuTaskRule{
			AmountCollar: model.AmountCollarRule{
				IsEnable:   true,
				CollarDays: 5,
				TxAmount:   10,
				Percent:    90.0,
			},
			AmountLarge: model.AmountLargeRule{
				IsEnable:   true,
				LargeValue: 100000,
			},
		},
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	task.Id = taskId

	if err := service.ClearAnalyseResults(ctx, taskId); err != nil {
		t.Error(err)
		return
	}

	var taskAccounts []*model.SuTaskAccount
	taskAccounts = append(taskAccounts, &model.SuTaskAccount{
		TaskId:  task.Id,
		Account: account,
	})

	err := service.analyse(ctx, task, taskAccounts)
	if err != nil {
		t.Fatal(err)
	}
}
