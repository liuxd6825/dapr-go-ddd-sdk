package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase/xmodel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_Record_Insert(t *testing.T) {
	/*
		opts := &xtest.MongoOptions{DBName: stringutils.PStr("master")}
		xtest.InitEnv_MongoRemote(opts)
	*/
	opts := &xtest.MongoOptions{DBName: stringutils.PStr("test")}
	xtest.InitEnv_MongoLocal(opts)

	count := int64(10)
	accounts := getAccounts(int(count))
	recordDao := NewRecordDao(xtest.MongoDBKey)
	randomutils.RangeRand(2012, 2014)

	var list []*model.Record
	for i := 0; i < 2000; i++ {
		accIndex := randomutils.Int64Max(count)
		oppIndex := randomutils.Int64Max(count)

		acc := accounts[accIndex]
		opp := accounts[oppIndex]

		year := randomutils.RangeRand(2012, 2013)
		date := randomutils.NewYear(int(year))

		isPayout := randomutils.Boolean()
		amount := randomutils.PFloat64()
		var payout, income *float64
		if isPayout {
			payout = amount
		} else {
			income = amount
		}

		record := &model.Record{
			Base: xmodel.Base{
				Id:       randomutils.NewId(),
				TenantId: xtest.TenantId,
				CaseId:   "1001",
				Remark:   randomutils.String(10),
			},
			RowNum:      int64(i),
			TaskId:      "taskId",
			DocId:       "docId",
			FileId:      "fieldId",
			Iden:        "001",
			Date:        &date,
			Year:        date.Year(),
			Month:       int(date.Month()),
			Day:         date.Day(),
			Name:        acc.Name,
			Acct:        acc.Account,
			AcctType:    acc.AccountType,
			OppName:     opp.Name,
			OppAcct:     opp.Account,
			OppAcctType: opp.AccountType,
			Payout:      payout,
			Income:      income,
			Amount:      amount,
			BankName:    BankNames[randomutils.IntMax(len(BankNames))],
		}
		list = append(list, record)
	}

	ctx := xtest.NewContext()
	recordDao.CreateMany(ctx, list)
}

type Account struct {
	Name        string
	Account     string
	AccountType string
}

func getAccounts(count int) []Account {
	names := []string{"张宇", "刘建新"}
	names = append(names, randomutils.NameCN())
	names = append(names, randomutils.NameCN())
	names = append(names, randomutils.NameCN())

	var list []Account
	for i := 0; i < count; i++ {
		idx := randomutils.Int64Max(int64(len(names)))
		account := Account{
			Name:        names[idx],
			Account:     randomutils.StringNumber(6),
			AccountType: "",
		}
		list = append(list, account)
	}
	return list
}

var BankNames = []string{"吉林银行", "长春银行", "邮政银行", "中国银行", "工商银行"}
