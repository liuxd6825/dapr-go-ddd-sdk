package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/pkg/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/command"
)

func CreateTestExcelData(ctx context.Context, taskId string, fileName string, sheet string, batchSize int64, temp *model.RecordTemplate) error {
	cmd := &command.RecordCreate4ExcelCommand{}
	cmd.CommandId = "001"
	cmd.Data = field.RecordCreate4ExcelCommandFields{
		CaseId:    "001",
		DocId:     "docId",
		FileId:    taskId,
		TaskId:    taskId,
		SheetName: sheet,
		FileName:  fileName,
		BatchSize: batchSize,
		Template:  temp,
	}

	appService := NewRecordService()
	if _, err := appService.Create4Excel(ctx, cmd, func(batch readexcel.Batching) error {
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func NewRecordTemplate() *model.RecordTemplate {
	tmp := model.NewRecordTemplate()
	tmp.MapHeads = []*readexcel.MapHead{
		{
			RowNum: 4,
			Columns: []*readexcel.MapColumn{
				{Key: "交易时间", Label: "A"}, {Key: "我方账号", Label: "B"}, {Key: "我方户名", Label: "C"},
				{Key: "我方开户行", Label: "D"}, {Key: "对方账号", Label: "E"}, {Key: "对方户名", Label: "F"},
				{Key: "对方开户行", Label: "G"}, {Key: "汇出金额", Label: "H"}, {Key: "汇入金额", Label: "I"},
				{Key: "余额", Label: "J"}, {Key: "摘要", Label: "K"}, {Key: "用途", Label: "L"},
				{Key: "备注", Label: "M"},
			},
		},
	}
	tmp.Date.SetMapKeys("交易时间")
	tmp.Name.SetMapKeys("我方户名").SetScript("我方户名")
	tmp.Acct.SetMapKeys("我方账号")
	tmp.BankName.SetMapKeys("我方开户行")
	tmp.Payout.SetMapKeys("汇出金额").SetScript("Math.abs(汇出金额)")
	tmp.Income.SetMapKeys("汇入金额").SetScript("Math.abs(汇入金额)")
	tmp.Balance.SetMapKeys("余额").SetScript("Math.abs(余额)")
	tmp.Amount.SetMapKeys("汇出金额", "汇入金额").SetScript(`parseFloat(汇出金额)!=0?Math.abs(汇出金额):Math.abs(汇入金额)`)
	tmp.OppName.SetMapKeys("对方户名")
	tmp.OppAcct.SetMapKeys("对方账号")
	tmp.OppBankName.SetMapKeys("对方开户行")
	tmp.Summary.SetMapKeys("摘要")
	tmp.Notes.SetMapKeys("备注", "用途")
	return tmp
}

func NewRecordTemplateByDSCF() *model.RecordTemplate {
	tmp := model.NewRecordTemplate()
	tmp.MapHeads = []*readexcel.MapHead{
		{
			RowNum: 1,
			Columns: []*readexcel.MapColumn{
				{Key: "交易日期", Label: "C"}, {Key: "交易时间", Label: "D"}, {Key: "我方账号", Label: "G"},
				{Key: "我方户名", Label: "F"}, {Key: "对方账号", Label: "X"}, {Key: "对方户名", Label: "Y"},
				{Key: "对方开户行", Label: "AB"}, {Key: "汇出金额", Label: "P"}, {Key: "汇入金额", Label: "O"},
				{Key: "余额", Label: "R"}, {Key: "进出", Label: "Q"}, {Key: "交易摘要", Label: "S"},
				{Key: "文字摘要", Label: "T"},
			},
		},
	}
	tmp.Date.SetMapKeys("交易日期", "交易时间").SetScript(" 交易日期 +' '+ 交易时间 ")
	tmp.Name.SetMapKeys("我方户名").SetScript("我方户名")
	tmp.Acct.SetMapKeys("我方账号")
	tmp.BankName.SetScript("'招商银行'")
	tmp.Payout.SetMapKeys("汇出金额").SetScript("汇出金额==''?0:abs(汇出金额)")
	tmp.Income.SetMapKeys("汇入金额").SetScript("汇入金额==''?0:abs(汇入金额)")
	tmp.Balance.SetMapKeys("余额").SetScript("Math.abs(余额)")
	tmp.Amount.SetMapKeys("进出").SetScript(`Math.abs(进出)`)
	tmp.OppName.SetMapKeys("对方户名")
	tmp.OppAcct.SetMapKeys("对方账号")
	tmp.OppBankName.SetMapKeys("对方开户行")
	tmp.Summary.SetMapKeys("交易摘要")
	tmp.Notes.SetMapKeys("文字摘要")

	return tmp
}

func NewRecordTemplate10w() *model.RecordTemplate {
	tmp := model.NewRecordTemplate()
	tmp.MapHeads = []*readexcel.MapHead{
		{
			RowNum: 1,
			Columns: []*readexcel.MapColumn{
				{Key: "ID", Label: "A"},
				{Key: "账户", Label: "B"},
				{Key: "记录代码", Label: "C"},
				{Key: "过账日期", Label: "D"},
				{Key: "交易日期", Label: "E"},
				{Key: "交易时间", Label: "F"},
				{Key: "流水号", Label: "G"},
				{Key: "交易类型", Label: "H"},
				{Key: "交易金额", Label: "I"},
				{Key: "余额", Label: "J"},
				{Key: "对方账号", Label: "K"},
				{Key: "对方户名", Label: "L"},
				{Key: "对方开户行", Label: "M"},
				{Key: "对方开户户名", Label: "N"},
				{Key: "摘要", Label: "O"},
				{Key: "交易来源", Label: "P"},
				{Key: "交易渠道", Label: "Q"},
				{Key: "类别", Label: "R"},
				{Key: "标记人", Label: "S"},
				{Key: "操作时间", Label: "T"},
				{Key: "创建时间/标记时间", Label: "U"},
			},
		},
	}
	tmp.Date.SetMapKeys("交易日期", "交易时间").SetScript("toDateTime(交易日期, 交易时间)")
	tmp.Name.SetScript("'北京玖富普惠信息技术有限公司'")
	tmp.Acct.SetMapKeys("账户")
	tmp.BankName.SetScript("'招商银行'")
	tmp.Payout.SetMapKeys("交易金额").SetScript("payout(交易金额)")
	tmp.Income.SetMapKeys("交易金额").SetScript("income(交易金额)")
	tmp.Balance.SetMapKeys("余额").SetScript("abs(余额)")
	tmp.Amount.SetMapKeys("交易金额").SetScript(`amount(交易金额)`)
	tmp.OppName.SetMapKeys("对方户名")
	tmp.OppAcct.SetMapKeys("对方账号")
	tmp.OppBankName.SetMapKeys("对方开户户名")
	tmp.Summary.SetMapKeys("摘要")
	//tmp.Notes.SetMapKeys("文字摘要")

	return tmp
}

func NewRecordTemplateByDSCF_LIJIA7190() *model.RecordTemplate {
	tmp := model.NewRecordTemplate()
	tmp.MapHeads = []*readexcel.MapHead{
		{
			RowNum: 1,
			Columns: []*readexcel.MapColumn{
				{Key: "交易时间", Label: "G"},
				{Key: "我方账号", Label: "B"}, {Key: "我方户名", Label: "A"},
				{Key: "对方账号", Label: "K"}, {Key: "对方户名", Label: "N"}, {Key: "对方开户行", Label: "P"},
				{Key: "收付标识", Label: "J"}, {Key: "交易金额", Label: "H"}, {Key: "交易余额", Label: "I"},
				{Key: "交易摘要", Label: "Q"}, {Key: "备注", Label: "AH"},
			},
		},
	}
	tmp.Date.SetMapKeys("交易日期", "交易时间").SetScript("交易时间")
	tmp.Name.SetMapKeys("我方户名").SetScript("我方户名")
	tmp.Acct.SetMapKeys("我方账号")
	tmp.BankName.SetScript("'招商银行'")
	tmp.Payout.SetMapKeys("交易金额", "收付标识").SetScript("收付标识=='出'?abs(交易金额):0")
	tmp.Income.SetMapKeys("交易金额", "收付标识").SetScript("收付标识=='进'?abs(交易金额):0")
	tmp.Balance.SetMapKeys("交易余额").SetScript("abs(交易余额)")
	tmp.Amount.SetMapKeys("交易金额").SetScript("abs(交易金额)")
	tmp.OppName.SetMapKeys("对方户名")
	tmp.OppAcct.SetMapKeys("对方账号")
	tmp.OppBankName.SetMapKeys("对方开户行")
	tmp.Summary.SetMapKeys("交易摘要")
	tmp.Notes.SetMapKeys("备注")

	return tmp
}

func NewRecordTemplateByDSCF_LRM_9615() *model.RecordTemplate {
	tmp := model.NewRecordTemplate()
	tmp.MapHeads = []*readexcel.MapHead{
		{
			RowNum: 1,
			Columns: []*readexcel.MapColumn{
				{Key: "交易日期", Label: "C"}, {Key: "交易时间", Label: "D"},
				{Key: "我方账号", Label: "G"}, {Key: "我方户名", Label: "F"},
				{Key: "对方账号", Label: "Z"}, {Key: "对方户名", Label: "AA"},
				{Key: "对方开户行", Label: "AB"},
				{Key: "入账", Label: "O"}, {Key: "出账", Label: "P"},
				{Key: "交易金额", Label: "N"}, {Key: "交易余额", Label: "T"},
				{Key: "交易摘要", Label: "U"}, {Key: "备注", Label: "V"},
			},
		},
	}
	tmp.Date.SetMapKeys("交易日期", "交易时间").SetScript("交易日期+' '+replace(交易时间, '.000', '')")
	tmp.Name.SetMapKeys("我方户名").SetScript("我方户名")
	tmp.Acct.SetMapKeys("我方账号")
	tmp.BankName.SetScript("'招商银行'")
	tmp.Payout.SetMapKeys("出账").SetScript("abs(出账)")
	tmp.Income.SetMapKeys("入账").SetScript("abs(入账)")
	tmp.Balance.SetMapKeys("交易余额").SetScript("abs(交易余额)")
	tmp.Amount.SetMapKeys("交易金额").SetScript("abs(交易金额)")
	tmp.OppName.SetMapKeys("对方户名")
	tmp.OppAcct.SetMapKeys("对方账号")
	tmp.OppBankName.SetMapKeys("对方开户行")
	tmp.Summary.SetMapKeys("交易摘要")
	tmp.Notes.SetMapKeys("备注")

	return tmp
}

func NewRecordTemplate70w() *model.RecordTemplate {
	tmp := model.NewRecordTemplate()
	tmp.MapHeads = []*readexcel.MapHead{
		{
			RowNum: 1,
			Columns: []*readexcel.MapColumn{
				{Key: "商户名称", Label: "A"}, // 交易时间
				{Key: "商户编号", Label: "B"}, // 我方账号
				{Key: "交易账号", Label: "C"},
				{Key: "交易户名", Label: "D"},
				{Key: "交易账户开户银行", Label: "E"},
				{Key: "交易日期", Label: "F"},
				{Key: "交易类型", Label: "G"},
				{Key: "交易金额", Label: "H"},
			},
		},
	}
	tmp.Date.SetMapKeys("交易日期")
	tmp.Name.SetMapKeys("商户名称")
	tmp.Acct.SetMapKeys("商户编号")
	tmp.BankName.SetScript("'招商银行'")
	tmp.Payout.SetMapKeys("交易类型", "交易金额").SetScript("交易类型=='代付'?abs(交易金额):0")
	tmp.Income.SetMapKeys("交易类型", "交易金额").SetScript("交易类型=='代收'?abs(交易金额):0")
	//tmp.Balance.SetMapKeys("交易余额").SetScript("Math.abs(交易余额)")
	tmp.Amount.SetMapKeys("交易金额")
	tmp.OppName.SetMapKeys("交易户名")
	tmp.OppAcct.SetMapKeys("交易账号")
	tmp.OppBankName.SetMapKeys("交易账户开户银行")
	//tmp.Summary.SetMapKeys("交易摘要")
	//tmp.Notes.SetMapKeys("备注")
	return tmp
}
