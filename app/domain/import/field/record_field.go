package field

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/enum"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
)

type RecordCreateFields struct {
	TaskId    string        `json:"taskId,omitempty"  validate:"required" title:"任务ID"`
	CaseId    string        `json:"caseId" validate:"required"  title:"案件Id"`      // 案件Id
	DocId     string        `json:"docId" validate:"required"  title:"文档Id"`       // 文档Id
	FileId    string        `json:"fileId" validate:"required"  title:"文件Id"`      // 文档Id
	FileName  string        `json:"fileName" validate:"required"  title:"文件名称"`    // 文档Id
	SheetId   string        `json:"sheetId" validate:"required"  title:"SheetId"`  // sheet页ID
	SheetName string        `json:"sheetName" validate:"required"  title:"Sheet页"` // sheet页
	Record    *RecordFields `json:"records" validate:"required"  title:"流水明细"`     // 流水明细
	Remark    string        `json:"remark" title:"备注"`
}

type RecordPreviewCommandFields struct {
	CaseId    string                `json:"caseId,omitempty"  validate:"required" title:"案件ID"`
	DocId     string                `json:"docId,omitempty"  validate:"required"  title:"文档ID"`
	FileId    string                `json:"fileId,omitempty"  validate:"required"   title:"文件ID"`
	TaskId    string                `json:"taskId,omitempty"  validate:"required" title:"任务ID"`
	FileName  string                `json:"fileName,omitempty"  validate:"required" title:"文件名称"`
	SheetName string                `json:"sheetName,omitempty" validate:"required"  title:"Sheet页"`
	Template  *model.RecordTemplate `json:"template,omitempty"  validate:"required"  title:"模板"`
}

type RecordCreate4ExcelCommandFields struct {
	TaskId    string                `json:"taskId,omitempty"  validate:"required" title:"任务ID"`
	CaseId    string                `json:"caseId,omitempty"  validate:"required" title:"案件ID"`
	DocId     string                `json:"docId,omitempty"  validate:"required"  title:"文档ID"`
	FileId    string                `json:"fileId,omitempty"  validate:"required"   title:"文件ID"`
	FileName  string                `json:"fileName,omitempty"  validate:"required" title:"文件名称"`
	SheetId   string                `json:"sheetId" validate:"required"  title:"SheetId"` // sheet页ID
	SheetName string                `json:"sheetName,omitempty" validate:"required"  title:"Sheet页"`
	BatchSize int64                 `json:"batchSize,omitempty"  validate:"required"  title:"批大小"`
	IsView    bool                  `json:"isView,omitempty"  validate:"-" title:"是预览"`
	Template  *model.RecordTemplate `json:"template,omitempty"  validate:"required"  title:"模板"`
}

type RecordImport2MasterFields struct {
	TaskId    string `json:"taskId" validate:"required" `
	CaseId    string `json:"caseId"  validate:"required" `
	DocId     string `json:"docId"  validate:"required" `
	FileId    string `json:"fileId"  validate:"required" `
	FileName  string `json:"fileName" validate:"required" `
	SheetId   string `json:"sheetId" validate:"required" `
	SheetName string `json:"sheetName" validate:"required" `
	PageSize  int64  `json:"pageSize" validate:"required" `
}

type RecordIeCreateFields struct {
	Id    string `json:"id"`
	Field string `json:"field"`
	Value any    `json:"value"`
}

type RecordIeUpdateFieldFields struct {
	Id     string         `json:"id"`
	Values map[string]any `json:"values"`
}

type RecordIeUpdateFilterFields struct {
	Filter string         `json:"filter"`
	TaskId string         `json:"taskId"`
	Values map[string]any `json:"values"`
}

type RecordIeDeleteFields struct {
	TaskId string `json:"taskId"`
	Id     string `json:"id,omitempty"  title:"租户标识"` // 行Id
}

type RecordIeUpdateFilterFieldsValues struct {
	Iden        *string    `json:"iden,omitempty"  bson:"iden"   validate:"-" title:"我方标识"`                   // 标识
	Name        *string    `json:"name,omitempty"   bson:"name"  validate:"-" title:"我方名称"`                   // 名称
	Acct        *string    `json:"acct,omitempty"   bson:"acct"  validate:"-" title:"我方账号"`                   // 账号
	AcctType    *string    `json:"acctType,omitempty"   bson:"acct_type"  validate:"-" title:"我方账号类型"`        // 账号类型
	Category    *string    `json:"category,omitempty"   bson:"category"  validate:"-" title:"我方类别"`           // 类别Id 公司或个人
	BankName    *string    `json:"bankName,omitempty"  bson:"bank_name"   validate:"-" title:"我方开户银行"`        // 开户银行
	Balance     *float64   `json:"balance,omitempty"   bson:"balance"  validate:"-" title:"我方余额账户"`           // 余额账户
	OppIden     *string    `json:"oppIden,omitempty"   bson:"opp_iden"  validate:"-" title:"对方标识"`            // 对方标识
	OppName     *string    `json:"oppName,omitempty"   bson:"opp_name"  validate:"-" title:"对方名称"`            // 对方名称
	OppAcct     *string    `json:"oppAcct,omitempty"   bson:"opp_acct"  validate:"-" title:"对方账号"`            // 对方账号
	OppAcctType *string    `json:"oppAcctType,omitempty"   bson:"opp_acct_type"  validate:"-" title:"对方账号类型"` // 对方账号类型
	OppCategory *string    `json:"oppCategory,omitempty"   bson:"opp_category"  validate:"-" title:"对方类别"`    // 对方类别
	OppBankName *string    `json:"oppBankName,omitempty"  bson:"opp_bank_name"   validate:"-" title:"对方开户银行"` // 对方开户银行
	Serial      *string    `json:"serial,omitempty"   bson:"serial"  validate:"-" title:"流水号"`                // 流水号
	Payout      *float64   `json:"payout,omitempty"   bson:"payout"  validate:"-" title:"支出金额"`               // 借方发生额（支出）
	Income      *float64   `json:"income,omitempty"  bson:"income"   validate:"-" title:"收入金额"`               // 贷方发生额（收入）
	Amount      *float64   `json:"amount,omitempty"   bson:"amount"  validate:"-" title:"交易金额"`               // 交易金额
	Date        *time.Time `json:"date,omitempty"   bson:"date"  validate:"-" title:"交易时间"`                   // 交易时间
	Type        *string    `json:"type,omitempty"   bson:"type"  validate:"-" title:"交易类型"`                   // 交易类型
	Ccy         *string    `json:"ccy,omitempty"  bson:"ccy"  validate:"-" title:"交易币种" `                     // 交易币种
	Place       *string    `json:"place,omitempty"   bson:"place"  validate:"-" title:"地点"`                   // 交易地点
	Summary     *string    `json:"summary,omitempty"   bson:"summary"  validate:"-" title:"摘要"`               // 摘要
	Notes       *string    `json:"notes,omitempty"   bson:"notes"  validate:"-" title:"备注"`                   // 备注
}

type RecordIeUpdateFields struct {
	Id     string `json:"id,omitempty"  bson:"_id"   validate:"required"  title:"租户标识"` // 行Id
	RowNum int64  `json:"rowNum,omitempty" bson:"rowNum" title:"租户标识"`

	CaseId string `json:"caseId,omitempty" bson:"case_id" title:"案件id"`
	TaskId string `json:"taskId,omitempty" bson:"task_id"  title:"任务id"`
	DocId  string `json:"docId,omitempty" bson:"doc_id"  title:"文档id"`
	FileId string `json:"fileId,omitempty" bson:"file_id"  title:"文件id"`

	Iden     string   `json:"iden,omitempty"  bson:"iden"   validate:"-" title:"我方标识"`            // 标识
	Name     string   `json:"name,omitempty"   bson:"name"  validate:"-" title:"我方名称"`            // 名称
	Acct     string   `json:"acct,omitempty"   bson:"acct"  validate:"-" title:"我方账号"`            // 账号
	AcctType string   `json:"acctType,omitempty"   bson:"acct_type"  validate:"-" title:"我方账号类型"` // 账号类型
	Category string   `json:"category,omitempty"   bson:"category"  validate:"-" title:"我方类别"`    // 类别Id 公司或个人
	BankName string   `json:"bankName,omitempty"  bson:"bank_name"   validate:"-" title:"我方开户银行"` // 开户银行
	Balance  *float64 `json:"balance,omitempty"   bson:"balance"  validate:"-" title:"我方余额账户"`    // 余额账户

	OppIden     string `json:"oppIden,omitempty"   bson:"opp_iden"  validate:"-" title:"对方标识"`            // 对方标识
	OppName     string `json:"oppName,omitempty"   bson:"opp_name"  validate:"-" title:"对方名称"`            // 对方名称
	OppAcct     string `json:"oppAcct,omitempty"   bson:"opp_acct"  validate:"-" title:"对方账号"`            // 对方账号
	OppAcctType string `json:"oppAcctType,omitempty"   bson:"opp_acct_type"  validate:"-" title:"对方账号类型"` // 对方账号类型
	OppCategory string `json:"oppCategory,omitempty"   bson:"opp_category"  validate:"-" title:"对方类别"`    // 对方类别
	OppBankName string `json:"oppBankName,omitempty"  bson:"opp_bank_name"   validate:"-" title:"对方开户银行"` // 对方开户银行

	Serial  string     `json:"serial,omitempty"   bson:"serial"  validate:"-" title:"流水号"`  // 流水号
	Payout  *float64   `json:"payout,omitempty"   bson:"payout"  validate:"-" title:"支出金额"` // 借方发生额（支出）
	Income  *float64   `json:"income,omitempty"  bson:"income"   validate:"-" title:"收入金额"` // 贷方发生额（收入）
	Amount  *float64   `json:"amount,omitempty"   bson:"amount"  validate:"-" title:"交易金额"` // 交易金额
	Date    *time.Time `json:"date,omitempty"   bson:"date"  validate:"-" title:"交易时间"`     // 交易时间
	Type    string     `json:"type,omitempty"   bson:"type"  validate:"-" title:"交易类型"`     // 交易类型
	Ccy     string     `json:"ccy,omitempty"  bson:"ccy"  validate:"-" title:"交易币种" `       // 交易币种
	Place   string     `json:"place,omitempty"   bson:"place"  validate:"-" title:"地点"`     // 交易地点
	Summary string     `json:"summary,omitempty"   bson:"summary"  validate:"-" title:"摘要"` // 摘要
	Notes   string     `json:"notes,omitempty"   bson:"notes"  validate:"-" title:"备注"`     // 备注
}

// RecordMetaFields
// 资金记录元数据 实体类型
type RecordMetaFields struct {
	Id      string          `json:"id" bson:"_id"`
	SrcType enum.SourceType `json:"srcType" bson:"src_type,omitempty"`
	Meta    *RowMeta        `json:"meta" bson:"meta,omitempty"`

	Name     *FieldMeta `json:"name" bson:"name" validate:"-" title:"名称"`                     // 名称
	Acct     *FieldMeta `json:"acct"  bson:"acct" validate:"-" title:"账号"`                    // 账号
	AcctType *FieldMeta `json:"accType" bson:"accType" validate:"-" title:"账号类型"`             // 账号类型
	Category *FieldMeta `json:"category" bson:"category" validate:"-" title:"类别"`             // 类别
	Balance  *FieldMeta `json:"balance" bson:"balance" validate:"-" title:"余额账户"`             // 余额账户
	BankName *FieldMeta `json:"bankName" bson:"bank_name,omitempty" validate:"-" title:"开户行"` // 开户行

	OppAcct     *FieldMeta `json:"oppAcct"  bson:"oppAcct" validate:"-" title:"对方账号"`         // 对方账号
	OppAccType  *FieldMeta `json:"oppAccType" bson:"oppAccType" validate:"-" title:"对方账号类型"`  // 对方账号类型
	OppCategory *FieldMeta `json:"oppCategory" bson:"oppCategory" validate:"-" title:"对方类别"`  // 对方类别
	OppIden     *FieldMeta `json:"oppIden" bson:"oppIden"  validate:"-" title:"对方标识"`         // 对方标识
	OppName     *FieldMeta `json:"oppName" bson:"oppName" validate:"-" title:"对方名称"`          // 对方名称
	OppBankName *FieldMeta `json:"oppBankName" bson:"oppBankName" validate:"-" title:"对方开户行"` // 对方开户行

	Amount *FieldMeta `json:"amount" bson:"amount" validate:"-" title:"交易金额"` // 交易金额
	Date   *FieldMeta `json:"date" bson:"date" validate:"-" title:"交易时间"`     // 交易时间
	Notes  *FieldMeta `json:"notes" bson:"notes" validate:"-" title:"交易备注"`   // 交易备注
	Place  *FieldMeta `json:"place" bson:"place" validate:"-" title:"交易地点"`   // 交易地点
	Type   *FieldMeta `json:"type" bson:"type" validate:"-" title:"交易类型"`     // 交易类型
	Ccy    *FieldMeta `json:"ccy" bson:"ccy" validate:"-" title:"交易币种"`       // 交易币种

	Serial  *FieldMeta `json:"serial"   bson:"serial"  validate:"-" title:"流水号"`   // 流水号
	Payout  *FieldMeta `json:"payout"   bson:"payout"  validate:"-" title:"借方发生额"` // 借方发生额（支取）
	Income  *FieldMeta `json:"income"  bson:"income"   validate:"-" title:"贷方发生额"` // 贷方发生额（收入）
	Summary *FieldMeta `json:"summary"   bson:"summary"  validate:"-" title:"摘要"`  // 摘要

}

type RowMeta struct {
	SrcType enum.SourceType `json:"srcType" bson:"src_type"`
	Db      *RowDbMeta      `json:"db" bson:"db"`
	Excel   *RowExcelMate   `json:"excel" bson:"excel"`
	Scan    *RowScanMeta    `json:"scan" bson:"scan"`
}

type RowDbMeta struct {
	DbId   string         `json:"dbId" bson:"db_id"`
	Table  string         `json:"table" bson:"table"`
	Values map[string]any `json:"values" bson:"values"`
}

type RowExcelMate struct {
	DocId  string `json:"docId" bson:"doc_id"`
	FileId string `json:"fileId" bson:"file_id"`
	Sheet  string `json:"sheet" bson:"sheet"`
	RowNum int64  `json:"rowNum" bson:"row_num"`
}

type RowScanMeta struct {
	Region *Region `json:"region"`
}

type FieldMeta struct {
	Scan  *FieldScanValue `json:"scan,omitempty" bson:"scan,omitempty"`
	Excel *FieldExcelMeta `json:"excel,omitempty" bson:"excel,omitempty"`
	Db    *FieldDBMeta    `json:"db" bson:"db,omitempty"`
}

type FieldDBMeta map[string]any

type FieldExcelMeta struct {
	Cells []*Cell `json:"cells" bson:"cells"`
}

type Cell struct {
	Row int `json:"row" bson:"row"`
	Col int `json:"col" bson:"col"`
}

type FieldScanValue struct {
	FileId    string `json:"fileId" bson:"file_id"  validate:"required" title:"文件Id"`
	Key       string `json:"key" bson:"key"  validate:"-" title:"关键字"`
	PageId    string `json:"pageId"  bson:"page_id" validate:"required" title:"流水页ID"`
	Region    Region `json:"region" bson:"region"  validate:"-" title:"文字位置"`
	RowId     string `json:"rowId" bson:"row_id" validate:"required" title:"行Id"`
	ScanValue string `json:"scanValue" bson:"scan_value" validate:"-" title:"扫描的原始值"`
	Type      string `json:"type" bson:"type" validate:"-" title:"数据类型"`
	Value     string `json:"value" bson:"value" validate:"-" title:"校对值"`
}

type Region struct {
	X1 int64 `json:"x1" bson:"x1" validate:"-" `
	X2 int64 `json:"x2" bson:"x2" validate:"-" `
	Y1 int64 `json:"y1" bson:"y1" validate:"-" `
	Y2 int64 `json:"y2" bson:"y2" validate:"-" `
}

// RecordFields
// 资金记录 实体类型
type RecordFields struct {
	Id       string   `json:"id"  bson:"id" validate:"required" title:"行Id"` // 行Id
	RowNum   int64    `json:"rowNum" bson:"rowNum" index:""  title:"行号"`
	Iden     string   `json:"iden" bson:"iden"  validate:"-" title:"标识"`           // 标识
	Name     string   `json:"name"  bson:"name" validate:"-" title:"名称"`           // 名称
	Acct     string   `json:"acct" bson:"acct" validate:"-" title:"账号"`            // 账号
	AcctType string   `json:"accType" bson:"accType"  validate:"-" title:"账号类型"`   // 账号类型
	Category string   `json:"category"  bson:"category" validate:"-" title:"类别"`   // 类别Id 公司或个人
	Balance  *float64 `json:"balance" bson:"balance" validate:"-" title:"余额账户"`    // 余额账户
	BankName string   `json:"bankName" bson:"bankName"  validate:"-" title:"开户银行"` // 开户银行

	OppIden     string `json:"oppIden" bson:"oppIden" validate:"-" title:"对方标识"`            // 对方标识
	OppName     string `json:"oppName" bson:"oppName" validate:"-" title:"对方名称"`            // 对方名称
	OppAcct     string `json:"oppAcct"  bson:"oppAcct" validate:"-" title:"对方账号"`           // 对方账号
	OppAcctType string `json:"oppAccType" bson:"oppAccType"  validate:"-" title:"对方账号类型"`   // 对方账号类型
	OppCategory string `json:"oppCategory" bson:"oppCategory"  validate:"-" title:"对方类别"`   // 对方类别
	OppBankName string `json:"oppBankName" bson:"oppBankName"  validate:"-" title:"对方开户银行"` // 对方开户银行

	Serial  string            `json:"serial"   bson:"serial"  validate:"-" title:"流水号"`   // 流水号
	Payout  *float64          `json:"payout"   bson:"payout"  validate:"-" title:"借方发生额"` // 借方发生额（支取）
	Income  *float64          `json:"Income"  bson:"income"   validate:"-" title:"贷方发生额"` // 贷方发生额（收入）
	Amount  *float64          `json:"amount"   bson:"amount"  validate:"-" title:"交易金额"`  // 交易金额
	Date    *time.Time        `json:"date"   bson:"date"  validate:"-" title:"交易时间"`      // 交易时间
	Type    string            `json:"type"   bson:"type"  validate:"-" title:"交易类型"`      // 交易类型
	Ccy     string            `json:"ccy"  bson:"ccy"  validate:"-" title:"交易币种" `        // 交易币种
	Place   string            `json:"place"   bson:"place"  validate:"-" title:"地点"`      // 交易地点
	Summary string            `json:"summary"   bson:"summary"  validate:"-" title:"摘要"`  // 摘要
	Notes   string            `json:"notes"   bson:"notes"  validate:"-" title:"备注"`      // 备注
	Meta    *RecordMetaFields `json:"meta"  gorm:"meta;json" validate:"-" title:"元数据"`    // 元数据
}
