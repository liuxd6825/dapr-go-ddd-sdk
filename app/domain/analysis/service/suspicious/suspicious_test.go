package suspicious

import (
	"context"
	"testing"

	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service/suspicious/action"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

const account = "6235822099004087593"

func init() {
	xtest2.Init(xtest2.TestType_MongoLocal)
}

func Test_FastInOutHours(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.TimeFastInOut.IsEnable = true
		task.Rules.TimeFastInOut.FastInOutHours = 24 * 5
		task.Rules.TimeFastInOut.Percent = 80.0
		task.Rules.TimeFastInOut.Amount = 1000
		return nil
	})
}

func Test_TimeSignificantDate(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.TimeSignificantDate.IsEnable = true
		task.Rules.TimeSignificantDate.CheckYear = true
		task.Rules.TimeSignificantDate.CheckQuarter = true
		task.Rules.TimeSignificantDate.CheckMonth = true
		task.Rules.TimeSignificantDate.Amount = 1000
		task.Rules.TimeSignificantDate.DaysBefore = 3
		task.Rules.TimeSignificantDate.DaysAfter = 2
		return nil
	})
}

func Test_TimeConcentratedPayments(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.TimeConcentratedPayments.IsEnable = true
		task.Rules.TimeConcentratedPayments.Amount = 1000
		task.Rules.TimeConcentratedPayments.StartTime = 18
		task.Rules.TimeConcentratedPayments.EndTime = 8
		task.Rules.TimeConcentratedPayments.Count = 10
		return nil
	})
}

func Test_TimeNonWorkingHours(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.TimeNonWorkingHours.IsEnable = true
		task.Rules.TimeNonWorkingHours.Amount = 1000
		task.Rules.TimeNonWorkingHours.HoursMax = 18
		task.Rules.TimeNonWorkingHours.HoursMin = 8
		task.Rules.TimeNonWorkingHours.IsNonWorkingHours = true
		task.Rules.TimeNonWorkingHours.IsNonWorkingSunday = true
		return nil
	})
}

func Test_AmountNear(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.AmountNear.IsEnable = true
		task.Rules.AmountNear.NearAmount = 100000
		task.Rules.AmountNear.NearPercent = 85.5
		return nil
	})
}

func Test_AmountLarge(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.AmountLarge.IsEnable = true
		task.Rules.AmountLarge.LargeValue = 100000
		return nil
	})
}

func Test_AmountCollar(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.AmountCollar.IsEnable = true
		task.Rules.AmountCollar.Days = 5
		task.Rules.AmountCollar.Amount = 100
		task.Rules.AmountCollar.Percent = 90.0
		return nil
	})
}

func Test_AmountNumber(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.AmountNumber.IsEnable = true
		task.Rules.AmountNumber.Number = 10000
		return nil
	})
}

func Test_FreqAbnormal(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.FreqAbnormal.IsEnable = true
		task.Rules.FreqAbnormal.DayTolerance = 2
		task.Rules.FreqAbnormal.MinOccurrences = 3
		task.Rules.FreqAbnormal.AmountTolerance = 5.0
		task.Rules.FreqAbnormal.PeriodDays = 30
		return nil
	})
}

func Test_FreqHigh(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.FreqHigh.IsEnable = true
		task.Rules.FreqHigh.Threshold = 2
		task.Rules.FreqHigh.Period = model2.SuFreqHighPeriod_Month
		return nil
	})
}

func Test_FreqSleep(t *testing.T) {
	_, _ = doAnalyse(t, account, func(task *model2.SuTask) error {
		task.Rules.FreqSleep.IsEnable = true
		task.Rules.FreqSleep.ActivateThreshold = 3
		task.Rules.FreqSleep.SleepDays = 180
		task.Rules.FreqSleep.ActivateDays = 30
		return nil
	})
}

func doAnalyse(t *testing.T, account string, do func(task *model2.SuTask) error) (*action.AnalyseResult, error) {
	ctx := xtest2.NewContext()
	task := &model2.SuTask{}
	task.Id = "001"
	err := do(task)
	if err != nil {
		t.Fatal(err)
		return nil, err
	}

	taskRecords := newTaskAccounts(t, ctx, account)
	accRecords := newAccountRecords(t, ctx, account)

	analyes := NewAnalyse(task, taskRecords, newDataRepo())
	result := analyes.applyRulesToAccount(ctx, accRecords)
	if result != nil {
		t.Log("suRecords count: ", len(result.SuRecords))
		t.Log("batches count: ", len(result.Batches))
	}
	return result, nil
}

func newTaskAccounts(t *testing.T, ctx context.Context, accounts ...string) []*model2.SuTaskAccount {
	var res []*model2.SuTaskAccount
	for _, account := range accounts {
		res = append(res, &model2.SuTaskAccount{
			Account:     account,
			TaskId:      "01",
			AccountType: model.AccountType_Company,
			Name:        "",
			BankName:    "",
		})
	}
	return res
}

func newAccountRecords(t *testing.T, ctx context.Context, account string) *model2.AccountRecords {
	recordDao := dao.NewRecordDao(config.DBKey)
	startTime := timeutils.NewDate(2018, 1, 1)
	endTime := timeutils.NewDate(2024, 1, 1)
	records, err := recordDao.FindByAccountOppAccount(ctx, account, startTime, endTime)
	if err != nil {
		t.Fatal(err)
	}
	accRecord := &model2.AccountRecords{
		Account: account,
		Records: records,
	}
	return accRecord
}

type dataRepo struct {
}

func newDataRepo() *dataRepo {
	return &dataRepo{}
}

func (d *dataRepo) GetCoreSuppliers() (map[string]action.Counterparty, error) {
	return map[string]action.Counterparty{}, nil
}

func (d *dataRepo) LoadCounterparties() (map[string]action.Counterparty, error) {
	return map[string]action.Counterparty{}, nil
}

func (d *dataRepo) GetAllTransactionsGroupedByAccount() (map[string][]*model.Tran, error) {
	return map[string][]*model.Tran{}, nil
}

func (d *dataRepo) LoadRelatedParties() ([]action.RelatedParty, error) {
	return []action.RelatedParty{}, nil
}

func (d *dataRepo) GetCounterparty(name string) (part *action.Counterparty, found bool) {
	return nil, false
}
