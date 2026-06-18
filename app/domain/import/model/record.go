package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type CashType bool

const (
	CashType_None CashType = false // 非现金
	CashType_Cash CashType = true  // 现金
)

type AcctType string

const (
	AcctType_None     AcctType = ""
	AcctType_Personal AcctType = "个人"
	AcctType_Company  AcctType = "公司"
)

func NewAcctType(v string) AcctType {
	switch v {
	case "个人":
		return AcctType_Personal
	case "公司":
		return AcctType_Company
	default:
		return AcctType_None
	}
}

type RecordIe struct {
	xbase.BaseModel `bson:",inline"`

	RowNum  int64  `json:"rowNum" gorm:"row_num" bson:"row_num" index:"" title:"行号" parquet:"name=row_num, type=INT64, repetitiontype=REQUIRED"`
	TaskId  string `json:"taskId" gorm:"task_id" bson:"task_id" index:"" title:"任务id" parquet:"name=task_id, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	DocId   string `json:"docId" gorm:"doc_id" bson:"doc_id" index:"" title:"文档id" parquet:"name=doc_id, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	FileId  string `json:"fileId" gorm:"file_id" bson:"file_id" index:"" title:"文件id" parquet:"name=file_id, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	SheetId string `json:"sheetId" gorm:"sheet_id" bson:"sheet_id" index:"" title:"Sheet页Id" parquet:"name=sheet_id, type=BYTE_ARRAY, repetitiontype=REQUIRED"`

	Iden     string   `json:"iden" gorm:"iden" bson:"iden" validate:"-" title:"我方标识" parquet:"name=iden, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Name     string   `json:"name" gorm:"name" bson:"name" index:"" validate:"-" title:"我方名称" parquet:"name=name, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Acct     string   `json:"acct" gorm:"acct" bson:"acct" index:"" validate:"-" title:"我方账号" parquet:"name=acct, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	AcctType string   `json:"acctType" gorm:"acct_type" bson:"acct_type" validate:"-" title:"我方账号类型" parquet:"name=acct_type, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Category string   `json:"category" gorm:"category" bson:"category" validate:"-" title:"我方类别" parquet:"name=category, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	BankName string   `json:"bankName" gorm:"bank_name" bson:"bank_name" validate:"-" title:"我方开户银行" parquet:"name=bank_name, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Balance  *float64 `json:"balance" gorm:"balance" bson:"balance" validate:"-" title:"我方余额账户" parquet:"name=balance, type=DOUBLE, repetitiontype=OPTIONAL"`

	OppIden     string `json:"oppIden" gorm:"opp_iden" bson:"opp_iden" validate:"-" title:"对方标识" parquet:"name=opp_iden, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	OppName     string `json:"oppName" gorm:"opp_name" bson:"opp_name" validate:"-" title:"对方名称" parquet:"name=opp_name, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	OppAcct     string `json:"oppAcct" gorm:"opp_acct" bson:"opp_acct" validate:"-" title:"对方账号" parquet:"name=opp_acct, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	OppAcctType string `json:"oppAcctType" gorm:"opp_acct_type" bson:"opp_acct_type" validate:"-" title:"对方账号类型" parquet:"name=opp_acct_type, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	OppCategory string `json:"oppCategory" gorm:"opp_category" bson:"opp_category" validate:"-" title:"对方类别" parquet:"name=opp_category, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	OppBankName string `json:"oppBankName" gorm:"opp_bank_name" bson:"opp_bank_name" validate:"-" title:"对方开户银行" parquet:"name=opp_bank_name, type=BYTE_ARRAY, repetitiontype=REQUIRED"`

	Cash    string     `json:"cash" gorm:"cash" bson:"cash" index:"" title:"现金标识" parquet:"name=cash, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Serial  string     `json:"serial" gorm:"serial" bson:"serial" validate:"-" title:"流水号" parquet:"name=serial, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Payout  *float64   `json:"payout" gorm:"payout" bson:"payout" index:"" validate:"-" title:"支出金额" parquet:"name=payout, type=DOUBLE, repetitiontype=OPTIONAL"`
	Income  *float64   `json:"income" gorm:"income" bson:"income" index:"" validate:"-" title:"收入金额" parquet:"name=income, type=DOUBLE, repetitiontype=OPTIONAL"`
	Amount  *float64   `json:"amount" gorm:"amount" bson:"amount" index:"" validate:"-" title:"交易金额" parquet:"name=amount, type=DOUBLE, repetitiontype=OPTIONAL"`
	Date    *time.Time `json:"date" gorm:"date" bson:"date" index:"" validate:"-" title:"交易时间" parquet:"name=date, type=INT64, repetitiontype=OPTIONAL"`
	Type    string     `json:"type" gorm:"type" bson:"type" validate:"-" title:"交易类型" parquet:"name=type, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Ccy     string     `json:"ccy" gorm:"ccy" bson:"ccy" validate:"-" title:"交易币种" parquet:"name=ccy, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Place   string     `json:"place" gorm:"place" bson:"place" validate:"-" title:"交易地点" parquet:"name=place, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Summary string     `json:"summary" gorm:"summary" bson:"summary" validate:"-" title:"摘要" parquet:"name=summary, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	Notes   string     `json:"notes" gorm:"notes" bson:"notes" validate:"-" title:"备注" parquet:"name=notes, type=BYTE_ARRAY, repetitiontype=REQUIRED"`
	// 💡 修正后的 Errors Map 标签：
	Errors map[string][]string `json:"errors" gorm:"errors;json" bson:"errors" parquet:"name=errors, type=MAP, keytype=BYTE_ARRAY, key_converted_type=UTF8, valuetype=BYTE_ARRAY, value_converted_type=UTF8, value_repetitiontype=REPEATED"`
	// 💡 修正后的 Cells Map 标签：
	Cells map[string]RecordIeCells `json:"cells" gorm:"cells;json" bson:"cells" parquet:"name=cells, type=MAP,  keytype=BYTE_ARRAY, key_converted_type=UTF8, valuetype=STRUCT"`
}
type RecordIeCells []RecordIeCell

type RecordIeCell struct {
	Key    string   `json:"key" gorm:"key" bson:"key"`
	Value  string   `json:"value" gorm:"value" bson:"value"`
	Col    int64    `json:"col" gorm:"col" bson:"col"`
	Row    int64    `json:"row" gorm:"row" bson:"row"`
	Errors []string `json:"errors" gorm:"errors" bson:"errors"`
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

	FieldName_Cash    FieldName = "cash"    // 是否现金
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
	keyword["是否现金"] = FieldName_Cash
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
