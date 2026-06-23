package factory

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
)

type RecordFactory struct {
}

func NewRecordFactory() *RecordFactory {
	return &RecordFactory{}
}

func (f *RecordFactory) NewByRecordImportMasterEvent(ctx context.Context, e *event.RecordImportMasterEvent) ([]*model.Record, error) {
	/*	list := make([]*model2.Record, len(e.Data.Items))
		for i, item := range e.Data.Items {
			v := model2.NewRecord()
			date := item.Date
			year, month, day := date.Date()
			v.Id = item.Id
			v.RowNum = item.RowNum

			v.MasterType = e.Data.MasterType
			v.MasterId = e.Data.MasterId
			v.TenantId = appctx.GetTenantId2(ctx)
			v.CaseId = e.Data.CaseId
			v.DocId = e.Data.DocId
			v.TaskId = e.Data.TaskId
			v.FileId = e.Data.FileId
			v.SheetId = e.Data.SheetId

			v.DocId = e.Data.DocId
			v.FileId = e.Data.FileId
			v.SheetId = e.Data.SheetId

			v.MasterId = e.Data.MasterId
			v.MasterType = e.Data.MasterType

			v.RowNum = item.RowNum
			v.TaskId = e.Data.TaskId
			v.Name = item.Name
			v.Acct = item.Acct
			v.AcctType = item.AcctType
			v.BankName = item.BankName
			v.Balance = item.Balance
			v.OppName = item.OppName
			v.OppAcct = item.OppAcct
			v.OppAcctType = item.OppAcctType
			v.OppBankName = item.OppBankName
			v.Serial = item.Serial
			v.Payout = getFloat(item.Payout)
			v.Income = getFloat(item.Income)
			v.Cash = 0 //model.CashType(item.Cash)
			if item.Cash {
				v.Cash = 1
			}
			v.Io = getIoType(item.Income, item.Payout)
			v.Amount = getFloat(item.Amount)

			v.Ccy = item.Ccy
			v.Place = item.Place
			v.Summary = item.Summary
			v.Notes = item.Notes

			v.Date = *date
			v.Year = year
			v.Month = int(month)
			v.Day = day

			list[i] = v
		}
		return list, nil*/
	return nil, nil
}

func getFloat(val *float64) float64 {
	if val == nil {
		return 0
	}
	return *val
}

func getIoType(income, payout *float64) model.IOType {
	if income != nil && *income != 0 {
		return model.IOType_In
	} else if payout != nil && *payout != 0 {
		return model.IOType_Out
	}
	return model.IOType_In
}
