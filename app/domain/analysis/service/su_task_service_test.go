package service

import (
	"testing"
	"time"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

const account = "6235822099004087593"
const taskId = "01"

func init() {
	xtest2.Init(xtest2.TestType_MongoLocal)
}

func Test_ClearAnalyseResults(t *testing.T) {
	ctx := xtest2.NewContext()
	service := NewSuTaskService()
	if err := service.ClearAnalyseResults(ctx, taskId); err != nil {
		t.Error(err)
		return
	}
}

func Test_analyse(t *testing.T) {
	startTime := timeutils.NewDate(2018, 1, 1).Add(9 * time.Hour)
	endTime := timeutils.NewDate(2024, 1, 1).Add(9 * time.Hour)

	ctx := xtest2.NewContext()
	service := NewSuTaskService()
	task := &model2.SuTask{
		Rules: model2.SuTaskRule{
			AmountCollar: model2.AmountCollarRule{
				IsEnable: true,
				Days:     5,
				Amount:   10,
				Percent:  90.0,
			},
			AmountLarge: model2.AmountLargeRule{
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

	var taskAccounts []*model2.SuTaskAccount
	taskAccounts = append(taskAccounts, &model2.SuTaskAccount{
		TaskId:  task.Id,
		Account: account,
	})

	err := service.analyse(ctx, task, taskAccounts)
	if err != nil {
		t.Fatal(err)
	}
}
