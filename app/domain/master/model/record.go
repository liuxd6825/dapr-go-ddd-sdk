package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"time"
)

type Record struct {
	xbase.BaseModel `bson:",inline"`
	DocId           string     `json:"docId" bson:"doc_id" index:"" title:"文档id"`
	FileId          string     `json:"fileId" bson:"file_id" index:"" title:"文件id"`
	SheetId         string     `json:"sheetId" bson:"sheet_id" index:"" title:"工作表ID"`
	RowNum          int64      `json:"rowNum" bson:"row_num" index:"" title:"行号"`
	TaskId          string     `json:"taskId" bson:"task_id" index:""  title:"任务id"`
	Iden            string     `json:"iden,omitempty"  bson:"iden"   index:""  validate:"-" title:"我方标识"`                        // 标识
	Name            string     `json:"name,omitempty"   bson:"name"  index:"" validate:"-" title:"我方名称"`                         // 名称
	Acct            string     `json:"acct,omitempty"   bson:"acct"  index:""  validate:"-" title:"我方账号"`                        // 账号
	AcctType        string     `json:"acctType,omitempty"   bson:"acct_type"  validate:"-" title:"我方账号类型"`                       // 账号类型
	Category        string     `json:"category,omitempty"   bson:"category"  validate:"-" title:"我方类别"`                          // 类别Id 公司或个人
	BankName        string     `json:"bankName,omitempty"  bson:"bank_name"   index:""   validate:"-" title:"我方开户银行"`            // 开户银行
	Balance         *float64   `json:"balance,omitempty"   bson:"balance"  index:""   validate:"-" title:"我方余额账户"`               // 余额账户
	OppIden         string     `json:"oppIden,omitempty"   bson:"opp_iden"  index:""   validate:"-" title:"对方标识"`                // 对方标识
	OppName         string     `json:"oppName,omitempty"   bson:"opp_name"  index:""   validate:"-" title:"对方名称"`                // 对方名称
	OppAcct         string     `json:"oppAcct,omitempty"   bson:"opp_acct"  index:""   validate:"-" title:"对方账号"`                // 对方账号
	OppAcctType     string     `json:"oppAcctType,omitempty"   bson:"opp_acct_type"  validate:"-" title:"对方账号类型"`                // 对方账号类型
	OppCategory     string     `json:"oppCategory,omitempty"   bson:"opp_category"  validate:"-" title:"对方类别"`                   // 对方类别
	OppBankName     string     `json:"oppBankName,omitempty"  bson:"opp_bank_name"  index:""   validate:"-" title:"对方开户银行"`      // 对方开户银行
	Serial          string     `json:"serial,omitempty"   bson:"serial"  index:""   validate:"-" title:"流水号"`                    // 流水号
	Payout          *float64   `json:"payout,omitempty"   bson:"payout"  index:""   validate:"-" title:"支出金额"`                   // 借方发生额（支出）
	Income          *float64   `json:"income,omitempty"  bson:"income"   index:""  validate:"-" title:"收入金额"`                    // 贷方发生额（收入）
	Amount          *float64   `json:"amount,omitempty"   bson:"amount"  index:""  validate:"-" title:"交易金额"`                    // 交易金额
	Date            *time.Time `json:"date,omitempty"   bson:"date" time_format:"datetime"  index:""  validate:"-" title:"交易时间"` // 交易时间
	Type            string     `json:"type,omitempty"   bson:"type"  index:""  validate:"-" title:"交易类型"`                        // 交易类型
	Ccy             string     `json:"ccy,omitempty"  bson:"ccy"  index:""  validate:"-" title:"交易币种" `                          // 交易币种
	Place           string     `json:"place,omitempty"   bson:"place"  index:""  validate:"-" title:"地点"`                        // 交易地点
	Summary         string     `json:"summary,omitempty"   bson:"summary" index:""   validate:"-" title:"摘要"`                    // 摘要
	Notes           string     `json:"notes,omitempty"   bson:"notes"  index:""  validate:"-" title:"备注"`                        // 备注
	DataSource      string     `json:"dataSource,omitempty"   bson:"data_source"  index:""   validate:"-" title:"数据源"`
	DataSourceId    string     `json:"dataSourceId,omitempty"   bson:"data_source_id"  index:""   validate:"-" title:"数据源id"`
	SuLarge         bool       `json:"suLarge" bson:"su_large"  index:""   validate:"-" title:"金额可疑" `
	SuInt           bool       `json:"suInt" bson:"su_int"  index:""   validate:"-" title:"整数可疑" `
	SuCollar        bool       `json:"suCollar" bson:"su_collar"  index:""   validate:"-" title:"对敲可疑" `
	SuNear          bool       `json:"suNear" bson:"su_near"  index:""   validate:"-" title:"临界可疑" `
	RiskCount       int        `json:"riskCount" bson:"risk_count"  index:""   validate:"-" title:"风险数" `
}

func NewRecord() *Record {
	return &Record{}
}

func (r *Record) GetId() string {
	return r.Id
}

func (r *Record) SetId(v string) {
	r.Id = v
}

func (r *Record) GetTenantId() string {
	return r.Id
}

func (r *Record) SetTenantId(v string) {
	r.TenantId = v
}
