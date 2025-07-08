package model

import "time"

type RecordIe struct {
	Id     string `json:"id"  gorm:"_id"   validate:"required"  desc:"id"` // 行Id
	RowNum int64  `json:"rowNum" gorm:"row_num" index:""  desc:"行号"`

	TenantId string `json:"tenantId" gorm:"tenant_id"  index:"" desc:"租户标识"`
	CaseId   string `json:"caseId" gorm:"case_id" index:""  desc:"案件id"`
	TaskId   string `json:"taskId" gorm:"task_id" index:""  desc:"任务id"`
	DocId    string `json:"docId" gorm:"doc_id" index:""  desc:"文档id"`
	FileId   string `json:"fileId" gorm:"file_id" index:""  desc:"文件id"`

	Iden     string   `json:"iden"  gorm:"iden"   validate:"-" desc:"我方标识"`            // 标识
	Name     string   `json:"name"   gorm:"name" index:""  validate:"-" desc:"我方名称"`   // 名称
	Acct     string   `json:"acct"   gorm:"acct" index:""  validate:"-" desc:"我方账号"`   // 账号
	AcctType string   `json:"acctType"   gorm:"acct_type"  validate:"-" desc:"我方账号类型"` // 账号类型
	Category string   `json:"category"   gorm:"category"  validate:"-" desc:"我方类别"`    // 类别Id 公司或个人
	BankName string   `json:"bankName"  gorm:"bank_name"   validate:"-" desc:"我方开户银行"` // 开户银行
	Balance  *float64 `json:"balance"   gorm:"balance"  validate:"-" desc:"我方余额账户"`    // 余额账户

	OppIden     string `json:"oppIden"   gorm:"opp_iden"  validate:"-" desc:"对方标识"`            // 对方标识
	OppName     string `json:"oppName"   gorm:"opp_name"  validate:"-" desc:"对方名称"`            // 对方名称
	OppAcct     string `json:"oppAcct"   gorm:"opp_acct"  validate:"-" desc:"对方账号"`            // 对方账号
	OppAcctType string `json:"oppAcctType"   gorm:"opp_acct_type"  validate:"-" desc:"对方账号类型"` // 对方账号类型
	OppCategory string `json:"oppCategory"   gorm:"opp_category"  validate:"-" desc:"对方类别"`    // 对方类别
	OppBankName string `json:"oppBankName"  gorm:"opp_bank_name"   validate:"-" desc:"对方开户银行"` // 对方开户银行

	Serial  string     `json:"serial"   gorm:"serial"  validate:"-" desc:"流水号"`            // 流水号
	Payout  *float64   `json:"payout"   gorm:"payout"  index:""  validate:"-" desc:"支出金额"` // 借方发生额（支出）
	Income  *float64   `json:"income"  gorm:"income"  index:""   validate:"-" desc:"收入金额"` // 贷方发生额（收入）
	Amount  *float64   `json:"amount"   gorm:"amount" index:""  validate:"-" desc:"交易金额"`  // 交易金额
	Date    *time.Time `json:"date"   gorm:"date" index:"" validate:"-" desc:"交易时间"`       // 交易时间
	Type    string     `json:"type"   gorm:"type"  validate:"-" desc:"交易类型"`               // 交易类型
	Ccy     string     `json:"ccy"  gorm:"ccy"  validate:"-" desc:"交易币种" `                 // 交易币种
	Place   string     `json:"place"   gorm:"place"  validate:"-" desc:"地点"`               // 交易地点
	Summary string     `json:"summary"   gorm:"summary"  validate:"-" desc:"摘要"`           // 摘要
	Notes   string     `json:"notes"   gorm:"notes"  validate:"-" desc:"备注"`               // 备注

	Errors map[string][]string `json:"errors" gorm:"errors"`
	Cells  map[string]Cells    `json:"cells" gorm:"cells"`
}

type Cells []Cell

type Cell struct {
	Key    string   `json:"key" gorm:"key"`
	Value  string   `json:"value" gorm:"value"`
	Col    int64    `json:"col" gorm:"col"`
	Row    int64    `json:"row" gorm:"row"`
	Errors []string `json:"errors" gorm:"errors"`
}

type FieldName string

type Keyword map[string]FieldName

const (
	FieldName_Iden        FieldName = "iden"     // 标识
	FieldName_Name        FieldName = "name"     // 名称
	FieldName_Account     FieldName = "acct"     // 账号
	FieldName_AccountType FieldName = "acctType" // 账号类型
	FieldName_Category    FieldName = "category" // 类别Id 公司或个人
	FieldName_BankName    FieldName = "bankName" // 开户银行
	FieldName_Balance     FieldName = "balance"  // 余额账户

	FieldName_OppIden        FieldName = "oppIden"     // 对方标识
	FieldName_OppName        FieldName = "oppName"     // 对方名称
	FieldName_OppAccount     FieldName = "oppAcct"     // 对方账号
	FieldName_OppAccountType FieldName = "oppAcctType" // 对方账号类型
	FieldName_OppCategory    FieldName = "oppCategory" // 对方类别
	FieldName_OppBankName    FieldName = "oppBankName" // 对方开户银行

	FieldName_Serial  FieldName = "serial"  // 流水号
	FieldName_Payout  FieldName = "payout"  // 借方发生额（支取）
	FieldName_Income  FieldName = "income"  // 贷方发生额（收入）
	FieldName_Amount  FieldName = "amount"  // 交易金额
	FieldName_Date    FieldName = "date"    // 交易时间
	FieldName_Type    FieldName = "type"    // 交易类型
	FieldName_Ccy     FieldName = "ccy"     // 交易币种
	FieldName_Place   FieldName = "place"   // 交易地点
	FieldName_Summary FieldName = "summary" // 摘要
	FieldName_Notes   FieldName = "notes"   // 备注
)

func (f FieldName) String() string {
	return string(f)
}

func NewKeywords() Keyword {
	keyword := Keyword{}
	keyword["账号"] = FieldName_Account
	keyword["账号名称"] = FieldName_Name
	keyword["币种"] = FieldName_Ccy
	keyword["交易日"] = FieldName_Date
	keyword["交易时间"] = FieldName_Date
	keyword["交易类型"] = FieldName_Type
	keyword["借方金额"] = FieldName_Payout
	keyword["贷方金额"] = FieldName_Income
	keyword["余额"] = FieldName_Balance
	keyword["摘要"] = FieldName_Summary
	keyword["流水号"] = FieldName_Serial
	keyword["用途"] = FieldName_Notes
	keyword["业务摘要"] = FieldName_Summary
	keyword["其它摘要"] = FieldName_Summary
	keyword["收(付)方分行名"] = FieldName_OppBankName
	keyword["收(付)方名称"] = FieldName_OppName
	keyword["收(付)方账号"] = FieldName_OppAccount
	return keyword
}

func (r *RecordIe) GetId() string {
	return r.Id
}

func (r *RecordIe) SetId(v string) {
	r.Id = v
}

func (r *RecordIe) GetTenantId() string {
	return r.TenantId
}

func (r *RecordIe) SetTenantId(v string) {
	r.TenantId = v
}

func (r *RecordIe) AddErrors(errs ...error) {
	// r.Errors = append(r.Errors, errs...)
}
