package field

import (
	"time"
)

type RecordIeUpdateFields struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id,omitempty"  bson:"_id"   validate:"required"  desc:"租户标识"` // 行Id
	RowNum   int64  `json:"rowNum,omitempty" bson:"rowNum" desc:"租户标识"`

	CaseId string `json:"caseId,omitempty" bson:"case_id" desc:"案件id"`
	TaskId string `json:"taskId,omitempty" bson:"task_id"  desc:"任务id"`
	DocId  string `json:"docId,omitempty" bson:"doc_id"  desc:"文档id"`
	FileId string `json:"fileId,omitempty" bson:"file_id"  desc:"文件id"`

	Iden     string   `json:"iden,omitempty"  bson:"iden"   validate:"-" desc:"我方标识"`            // 标识
	Name     string   `json:"name,omitempty"   bson:"name"  validate:"-" desc:"我方名称"`            // 名称
	Acct     string   `json:"acct,omitempty"   bson:"acct"  validate:"-" desc:"我方账号"`            // 账号
	AcctType string   `json:"acctType,omitempty"   bson:"acct_type"  validate:"-" desc:"我方账号类型"` // 账号类型
	Category string   `json:"category,omitempty"   bson:"category"  validate:"-" desc:"我方类别"`    // 类别Id 公司或个人
	BankName string   `json:"bankName,omitempty"  bson:"bank_name"   validate:"-" desc:"我方开户银行"` // 开户银行
	Balance  *float64 `json:"balance,omitempty"   bson:"balance"  validate:"-" desc:"我方余额账户"`    // 余额账户

	OppIden     string `json:"oppIden,omitempty"   bson:"opp_iden"  validate:"-" desc:"对方标识"`            // 对方标识
	OppName     string `json:"oppName,omitempty"   bson:"opp_name"  validate:"-" desc:"对方名称"`            // 对方名称
	OppAcct     string `json:"oppAcct,omitempty"   bson:"opp_acct"  validate:"-" desc:"对方账号"`            // 对方账号
	OppAcctType string `json:"oppAcctType,omitempty"   bson:"opp_acct_type"  validate:"-" desc:"对方账号类型"` // 对方账号类型
	OppCategory string `json:"oppCategory,omitempty"   bson:"opp_category"  validate:"-" desc:"对方类别"`    // 对方类别
	OppBankName string `json:"oppBankName,omitempty"  bson:"opp_bank_name"   validate:"-" desc:"对方开户银行"` // 对方开户银行

	Serial  string     `json:"serial,omitempty"   bson:"serial"  validate:"-" desc:"流水号"`  // 流水号
	Payout  *float64   `json:"payout,omitempty"   bson:"payout"  validate:"-" desc:"支出金额"` // 借方发生额（支出）
	Income  *float64   `json:"income,omitempty"  bson:"income"   validate:"-" desc:"收入金额"` // 贷方发生额（收入）
	Amount  *float64   `json:"amount,omitempty"   bson:"amount"  validate:"-" desc:"交易金额"` // 交易金额
	Date    *time.Time `json:"date,omitempty"   bson:"date"  validate:"-" desc:"交易时间"`     // 交易时间
	Type    string     `json:"type,omitempty"   bson:"type"  validate:"-" desc:"交易类型"`     // 交易类型
	Ccy     string     `json:"ccy,omitempty"  bson:"ccy"  validate:"-" desc:"交易币种" `       // 交易币种
	Place   string     `json:"place,omitempty"   bson:"place"  validate:"-" desc:"地点"`     // 交易地点
	Summary string     `json:"summary,omitempty"   bson:"summary"  validate:"-" desc:"摘要"` // 摘要
	Notes   string     `json:"notes,omitempty"   bson:"notes"  validate:"-" desc:"备注"`     // 备注
}
