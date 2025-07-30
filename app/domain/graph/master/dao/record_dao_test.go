package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecordDao_Create(t *testing.T) {
	xtest.InitEnv_Neo4j(xtest.Neo4jRemoveOption)
	dao := NewRecordDao()
	ctx := xtest.NewContext()
	record := newRecord("1")
	err := dao.Create(ctx, record)
	assert.Nil(t, err)
}

func TestRecordDao_Delete(t *testing.T) {
	xtest.InitEnv_Neo4j(xtest.Neo4jRemoveOption)
	dao := NewRecordDao()
	ctx := xtest.NewContext()
	record := newRecord("1")
	err := dao.Delete(ctx, record)
	assert.Nil(t, err)
}

func TestRecordDao_Update(t *testing.T) {
	xtest.InitEnv_Neo4j(xtest.Neo4jRemoveOption)
	dao := NewRecordDao()
	ctx := xtest.NewContext()
	amount := 120.0
	record := newRecord("1")
	record.Amount = &amount
	err := dao.Update(ctx, record)
	assert.Nil(t, err)
}

func newRecord(id string) *model.Record {
	amount := 100.0
	date := randomutils.Date()
	record := &model.Record{
		Date:        &date,
		Name:        "张二",
		Acct:        "123456",
		AcctType:    "1",
		Amount:      &amount,
		Balance:     &amount,
		BankName:    "中国工商银行",
		Category:    "1",
		Ccy:         "CNY",
		OppName:     "张三",
		OppAcct:     "654321",
		OppAcctType: "1",
		OppBankName: "中国工商银行",
		OppCategory: "1",
		OppIden:     "1",
		Notes:       "notes",
	}
	record.Id = id
	record.TenantId = xtest.TenantId
	record.CaseId = xtest.CaseId
	return record
}
