package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"time"
)

type Tran struct {
	xbase.BaseModel `bson:",inline"`
	//	TranId          string      `json:"tranId" gorm:"tran_id" bson:"tran_id" index:"" title:"交易id"`
	Name        string      `json:"name"   gorm:"name"  bson:"name"  index:"" validate:"-" title:"我方名称"`                             // 名称
	Acct        string      `json:"acct"   gorm:"acct"  bson:"acct"  index:""  validate:"-" title:"我方账号"`                            // 账号
	AcctType    AccountType `json:"acctType" gorm:"acct_type"    bson:"acct_type"  validate:"-" title:"我方账号类型"`                      // 账号类型
	BankName    string      `json:"bankName"  gorm:"bank_name"  bson:"bank_name"   index:""   validate:"-" title:"开户银行"`             // 开户银行
	OppName     string      `json:"oppName"   gorm:"opp_name"  bson:"opp_name"  index:""   validate:"-" title:"对方名称"`                // 对方名称
	OppAcct     string      `json:"oppAcct"   gorm:"opp_acct"  bson:"opp_acct"  index:""   validate:"-" title:"对方账号"`                // 对方账号
	OppAcctType AccountType `json:"oppAcctType"  gorm:"opp_acct_type"   bson:"opp_acct_type"  validate:"-" title:"对方账号类型"`           // 对方账号类型
	OppBankName string      `json:"oppBankName"  gorm:"opp_bank_name"  bson:"opp_bank_name"  index:""   validate:"-" title:"对方开户银行"` // 对方开户银行
	Cash        CashType    `json:"cash" gorm:"cash" bson:"cash" index:"" title:"现金"`
	Amount      float64     `json:"amount"   gorm:"amount"  bson:"amount"  index:""  validate:"-" title:"交易金额"`  // 交易金额
	Date        time.Time   `json:"date"   gorm:"date"  bson:"date"  index:""  validate:"-" title:"交易时间"`        // 交易时间
	Ccy         string      `json:"ccy"  gorm:"ccy"  bson:"ccy"  index:""  validate:"-" title:"交易币种" `           // 交易币种
	Place       string      `json:"place"   gorm:"place"  bson:"place"  index:""  validate:"-" title:"地点"`       // 交易地点
	Summary     string      `json:"summary"  gorm:"summary"   bson:"summary" index:""   validate:"-" title:"摘要"` // 摘要
	Notes       string      `json:"notes"   gorm:"notes"  bson:"notes"  index:""  validate:"-" title:"备注"`       // 备注
	Year        int         `json:"year"  gorm:"year" bson:"year" index:"" title:"年"`
	Month       int         `json:"month" gorm:"month"  bson:"month" index:"" title:"月"`
	Day         int         `json:"day"  gorm:"day" bson:"day" index:"" title:"日"`
}

func NewTran() *Tran {
	return &Tran{}
}

func NewTranFromRecord(record *Record) *Tran {
	acct := record.Acct
	name := record.Name
	acctType := record.AcctType
	bankName := record.BankName

	oppAcct := record.OppAcct
	oppName := record.OppName
	oppAcctType := record.OppAcctType
	oppBankName := record.OppBankName

	if record.Payout != nil {
		acct = record.OppAcct
		name = record.OppName
		acctType = record.OppAcctType
		bankName = record.OppBankName

		oppAcct = record.Acct
		oppName = record.Name
		oppAcctType = record.AcctType
		oppBankName = record.BankName
	}

	res := &Tran{
		BaseModel:   record.BaseModel,
		Name:        name,
		Acct:        acct,
		AcctType:    acctType,
		BankName:    bankName,
		OppName:     oppName,
		OppAcct:     oppAcct,
		OppAcctType: oppAcctType,
		OppBankName: oppBankName,
		Cash:        record.Cash,
		Amount:      *record.Amount,
		Date:        *record.Date,
		Ccy:         record.Ccy,
		Place:       record.Place,
		Summary:     record.Summary,
		Notes:       record.Notes,
		Year:        record.Date.Year(),
		Month:       int(record.Date.Month()),
		Day:         record.Date.Day(),
	}
	res.Id = record.TranId
	return res
}
