package model

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/pkg/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

type RecordTemplate struct {
	Iden     *readexcel.Field `json:"iden" bson:"iden"`
	Name     *readexcel.Field `json:"name" bson:"name"`
	Acct     *readexcel.Field `json:"acct" bson:"acct"`
	AcctType *readexcel.Field `json:"acctType" bson:"acct_type"`
	//Category *readexcel.Field `json:"category" bson:"category"`
	BankName *readexcel.Field `json:"bankName" bson:"bank_name"`
	Balance  *readexcel.Field `json:"balance" bson:"balance"`

	OppIden     *readexcel.Field `json:"oppIden" bson:"opp_iden"`
	OppName     *readexcel.Field `json:"oppName" bson:"opp_name"`
	OppAcct     *readexcel.Field `json:"oppAcct" bson:"opp_acct"`
	OppAcctType *readexcel.Field `json:"oppAcctType" bson:"opp_acct_type"`
	//OppCategory *readexcel.Field `json:"oppCategory" bson:"opp_category"`
	OppBankName *readexcel.Field `json:"oppBankName" bson:"opp_bank_name"`

	Cash       *readexcel.Field     `json:"cash" bson:"cash"`
	Serial     *readexcel.Field     `json:"serial" bson:"serial"`
	Income     *readexcel.Field     `json:"income" bson:"income"`
	Payout     *readexcel.Field     `json:"payout" bson:"payout"`
	Amount     *readexcel.Field     `json:"amount" bson:"amount"`
	Date       *readexcel.Field     `json:"date" bson:"date"`
	Type       *readexcel.Field     `json:"type" bson:"type"`
	Ccy        *readexcel.Field     `json:"ccy" bson:"ccy"`
	Place      *readexcel.Field     `json:"place" bson:"place"`
	Summary    *readexcel.Field     `json:"summary" bson:"summary"`
	Notes      *readexcel.Field     `json:"notes" bson:"notes"`
	MapHeadRow int64                `json:"mapHeadRow" bson:"map_head_row"`
	MapHeads   []*readexcel.MapHead `json:"mapHeads" bson:"map_heads"`
	Consts     []*readexcel.Const   `json:"consts" bson:"consts"`

	fields map[string]*readexcel.Field
}

func NewRecordTemplate() *RecordTemplate {
	t := &RecordTemplate{
		Iden:     readexcel.NewField("iden", "标识", readexcel.DataType_String, true),
		Name:     readexcel.NewField("name", "名称", readexcel.DataType_String, false),
		Acct:     readexcel.NewField("acct", "账号", readexcel.DataType_String, false),
		AcctType: readexcel.NewField("acctType", "账号类型", readexcel.DataType_String, true),
		//Category: readexcel.NewField("category", "类别", readexcel.DataType_String, true),
		BankName: readexcel.NewField("bankName", "开户行", readexcel.DataType_String, false),
		Balance:  readexcel.NewField("balance", "余额", readexcel.DataType_Money, true),

		OppIden:     readexcel.NewField("oppIden", "对方标识", readexcel.DataType_String, true),
		OppName:     readexcel.NewField("oppName", "对方名称", readexcel.DataType_String, false),
		OppAcct:     readexcel.NewField("oppAcct", "对方账号", readexcel.DataType_String, false),
		OppAcctType: readexcel.NewField("oppAcctType", "对方账号类型", readexcel.DataType_String, true),
		//OppCategory: readexcel.NewField("oppCategory", "对方类别", readexcel.DataType_String, true),
		OppBankName: readexcel.NewField("oppBankName", "对方开户行", readexcel.DataType_String, false),

		Cash:    readexcel.NewField("cash", "现金标识", readexcel.DataType_Boolean, true),
		Serial:  readexcel.NewField("serial", "流水号", readexcel.DataType_String, true),
		Income:  readexcel.NewField("income", "收入金额(贷)", readexcel.DataType_Money, false),
		Payout:  readexcel.NewField("payout", "支出金额(借)", readexcel.DataType_Money, false),
		Amount:  readexcel.NewField("amount", "交易金额", readexcel.DataType_String, true),
		Date:    readexcel.NewField("date", "交易日期", readexcel.DataType_DateTime, false),
		Type:    readexcel.NewField("type", "交易类型", readexcel.DataType_String, true),
		Ccy:     readexcel.NewField("ccy", "币种", readexcel.DataType_String, true),
		Place:   readexcel.NewField("place", "交易地点", readexcel.DataType_String, true),
		Summary: readexcel.NewField("summary", "摘要", readexcel.DataType_String, true),
		Notes:   readexcel.NewField("notes", "备注", readexcel.DataType_String, true),
	}
	t.Name.AddReplace("（", "(").AddReplace("）", ")").AddReplace(" ", "")
	t.OppName.AddReplace("（", "(").AddReplace("）", ")").AddReplace(" ", "")
	fields := make(map[string]*readexcel.Field, 0)

	fields[t.Iden.Name] = t.Iden
	fields[t.Name.Name] = t.Name
	fields[t.Acct.Name] = t.Acct
	fields[t.AcctType.Name] = t.AcctType
	//fields[t.Category.Name] = t.Category
	fields[t.BankName.Name] = t.BankName
	fields[t.Balance.Name] = t.Balance

	fields[t.OppIden.Name] = t.OppIden
	fields[t.OppName.Name] = t.OppName
	fields[t.OppAcct.Name] = t.OppAcct
	fields[t.OppAcctType.Name] = t.OppAcctType
	//fields[t.OppCategory.Name] = t.OppCategory
	fields[t.OppBankName.Name] = t.OppBankName

	fields[t.Serial.Name] = t.Serial
	fields[t.Income.Name] = t.Income
	fields[t.Payout.Name] = t.Payout
	fields[t.Amount.Name] = t.Amount
	fields[t.Date.Name] = t.Date
	fields[t.Type.Name] = t.Type
	fields[t.Ccy.Name] = t.Ccy
	fields[t.Place.Name] = t.Place
	fields[t.Summary.Name] = t.Summary
	fields[t.Notes.Name] = t.Notes

	t.fields = fields
	return t
}

func (t *RecordTemplate) NewTemplate() (*readexcel.Template, error) {
	fields := []*readexcel.Field{
		t.Iden,
		t.Name,
		t.Acct,
		t.AcctType,
		//t.Category,
		t.BankName,
		t.Balance,

		t.OppIden,
		t.OppName,
		t.OppAcct,
		t.OppAcctType,
		//t.OppCategory,
		t.OppBankName,

		t.Serial,
		t.Income,
		t.Payout,
		t.Amount,
		t.Date,
		t.Type,
		t.Ccy,
		t.Place,
		t.Summary,
		t.Notes,
	}
	temp, err := readexcel.NewTemplate(fields, t.MapHeads, t.Consts)
	return temp, err
}

func (t *RecordTemplate) AddHead(head *readexcel.MapHead) {
	t.MapHeads = append(t.MapHeads, head)
}

func (t *RecordTemplate) AddConst(c *readexcel.Const) {
	t.Consts = append(t.Consts, c)
}

func (t *RecordTemplate) GetField(name string) (*readexcel.Field, bool) {
	field, ok := t.fields[name]
	return field, ok
}

// Validate
// @Description: 命令数据验证
func (t *RecordTemplate) Validate() error {
	err := errors.NewVerifyError()
	for _, f := range t.fields {
		if !f.AllowNull {
			if len(f.MapKeys) == 0 && len(f.Script) == 0 {
				err.AppendField(f.Name, "MapKeys 或 Script 不可为空 ")
			}
		}
	}

	if len(t.MapHeads) == 0 {
		err.AppendField("MapHeads", "不可为空")
	}

	for i, h := range t.MapHeads {
		if h.RowNum < 0 {
			err.AppendField(fmt.Sprintf("MapHeads[%v].RowNum", i), "不可为空")
		}
		if len(h.Columns) == 0 {
			err.AppendField(fmt.Sprintf("MapHeads[%v].Columns", i), "不可为空")
		}
	}

	if len(err.Errors) == 0 {
		return nil
	}
	return err
}
