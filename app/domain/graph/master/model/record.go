package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"time"
)

type Record struct {
	xbase.BaseModel `bson:",inline"`

	RowNum    int64  `json:"rowNum" gorm:"row_num" bson:"row_num" index:""  title:"行号"`
	TaskId    string `json:"taskId" gorm:"task_id" bson:"task_id" index:""  title:"任务id"`
	DocId     string `json:"docId" gorm:"doc_id" bson:"doc_id" index:""  title:"文档id"`
	FileId    string `json:"fileId" gorm:"file_id" bson:"file_id" index:""  title:"文件id"`
	FileName  string `json:"fileName" gorm:"file_name" bson:"file_name" title:""`
	SheetName string `json:"sheetName" gorm:"sheet_name" bson:"sheet_name" title:""`

	Iden     string   `json:"iden"  gorm:"iden"  bson:"iden"  validate:"-" title:"我方标识"`                 // 标识
	Name     string   `json:"name"   gorm:"name" bson:"name" index:""  validate:"-" title:"我方名称"`        // 名称
	Acct     string   `json:"acct"   gorm:"acct"  bson:"acct" index:""  validate:"-" title:"我方账号"`       // 账号
	AcctType string   `json:"acctType"   gorm:"acct_type"  bson:"acct_type" validate:"-" title:"我方账号类型"` // 账号类型
	Category string   `json:"category"   gorm:"category"  bson:"category" validate:"-" title:"我方类别"`     // 类别Id 公司或个人
	BankName string   `json:"bankName"  gorm:"bank_name"  bson:"bank_name"  validate:"-" title:"我方开户银行"` // 开户银行
	Balance  *float64 `json:"balance"   gorm:"balance" bson:"balance" validate:"-" title:"我方余额账户"`       // 余额账户

	OppIden     string `json:"oppIden"   gorm:"opp_iden" bson:"opp_iden" validate:"-" title:"对方标识"`                 // 对方标识
	OppName     string `json:"oppName"   gorm:"opp_name"  bson:"opp_name" validate:"-" title:"对方名称"`                // 对方名称
	OppAcct     string `json:"oppAcct"   gorm:"opp_acct" bson:"opp_acct" validate:"-" title:"对方账号"`                 // 对方账号
	OppAcctType string `json:"oppAcctType"   gorm:"opp_acct_type" bson:"opp_acct_type" validate:"-" title:"对方账号类型"` // 对方账号类型
	OppCategory string `json:"oppCategory"   gorm:"opp_category" bson:"opp_category"  validate:"-" title:"对方类别"`    // 对方类别
	OppBankName string `json:"oppBankName"  gorm:"opp_bank_name" bson:"opp_bank_name"  validate:"-" title:"对方开户银行"` // 对方开户银行

	Serial string     `json:"serial"   gorm:"serial" bson:"serial" validate:"-" title:"流水号"`            // 流水号
	Payout *float64   `json:"payout"   gorm:"payout" bson:"payout" index:""  validate:"-" title:"支出金额"` // 借方发生额（支出）
	Income *float64   `json:"income"  gorm:"income" bson:"income" index:""   validate:"-" title:"收入金额"` // 贷方发生额（收入）
	Amount *float64   `json:"amount"   gorm:"amount" bson:"amount" index:""  validate:"-" title:"交易金额"` // 交易金额
	Date   *time.Time `json:"date"   gorm:"date" bson:"date" index:"" validate:"-" title:"交易时间"`        // 交易时间
	Year   int        `json:"year" gorm:"year" bson:"year" title:"交易年份"`
	Month  int        `json:"month" gorm:"month" bson:"month" title:"交易月份"`
	Day    int        `json:"day" gorm:"day" bson:"day" title:"交易天份"`

	Type    string `json:"type"   gorm:"type" bson:"type" validate:"-" title:"交易类型"`        // 交易类型
	Ccy     string `json:"ccy"  gorm:"ccy" bson:"ccy" validate:"-" title:"交易币种" `           // 交易币种
	Place   string `json:"place"   gorm:"place"  bson:"place" validate:"-" title:"地点"`      // 交易地点
	Summary string `json:"summary"   gorm:"summary" bson:"summary" validate:"-" title:"摘要"` // 摘要
	Notes   string `json:"notes"   gorm:"notes"  bson:"notes" validate:"-" title:"备注"`      // 备注

}
