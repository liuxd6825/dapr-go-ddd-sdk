package action

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_applyRulesToAccount(t *testing.T) {
	xtest.InitEnv_MongoRemoteTest(xtest.NewMongoOptions().SetDBName("master"))
	ctx := xtest.NewContext()
	task := &model.SuTask{}
	task.Id = "001"
	task.IsAmountLarge = true
	task.AmountRule.LargeValue = 1000

	task.IsAmountNear = true
	task.AmountRule.NearPercent = 0.1
	task.AmountRule.NearMin = 1000
	task.AmountRule.NearMax = 1000

	task.IsAmountCollar = true
	task.AmountRule.CollarDays = 5

	task.IsAmountInt = true
	task.AmountRule.IntValue = 1000

	task.IsTimeNonWorkingHours = true
	task.TimeRule.NonWorkingHoursMin = 8
	task.TimeRule.NonWorkingHoursMax = 18

	account := "6235822099004087593"
	accounts := []*model.SuTaskAccount{}
	accounts = append(accounts, &model.SuTaskAccount{
		Account: account,
	})

	tranDao := dao.NewTranDao(config.DBKey)
	startTime := timeutils.NewDate(2019, 1, 1)
	endTime := timeutils.NewDate(2023, 1, 1)
	trans, err := tranDao.FindByAccountOppAccount(ctx, account, &startTime, &endTime)
	if err != nil {
		t.Fatal(err)
	}
	accTrans := &AccountTransactions{
		Account: account,
		Trans:   trans,
	}
	counterpartyMap := make(map[string]Counterparty)
	an := NewSuTaskAnalyse(task, accounts)

	result := an.applyRulesToAccount(ctx, accTrans, counterpartyMap)
	if result == nil {
		t.Fatal("result is nil")
	}
}
