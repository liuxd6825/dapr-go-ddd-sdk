package model

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"math"
	"strconv"
	"strings"
	"time"
)

// Record
// @Description:  导入的原始流水
type Record struct {
	xbase.BaseModel `bson:",inline"`
	TranId          string `json:"tranId" gorm:"tran_id" bson:"tran_id" index:"" title:"交易id"`
	DocId           string `json:"docId" gorm:"doc_id" bson:"doc_id" index:"" title:"文档id"`
	FileId          string `json:"fileId"  gorm:"file_id" bson:"file_id" index:"" title:"文件id"`
	SheetId         string `json:"sheetId" gorm:"sheet_id" bson:"sheet_id" index:"" title:"工作表ID"`
	GraphId         string `json:"graphId" gorm:"graph_id" bson:"graph_id"`
	MasterId        string `json:"masterId" gorm:"master_id" bson:"master_id" title:"主数据ID"`
	MasterType      string `json:"masterType" gorm:"master_type" bson:"master_type" title:"主数据类型"`

	RowNum      int64       `json:"rowNum"  gorm:"row_num" bson:"row_num" index:"" title:"行号"`
	TaskId      string      `json:"taskId" gorm:"task_id"  bson:"task_id" index:""  title:"任务id"`
	Name        string      `json:"name"   gorm:"name"  bson:"name"  index:"" validate:"-" title:"我方名称"`                             // 名称
	Acct        string      `json:"acct"   gorm:"acct"  bson:"acct"  index:""  validate:"-" title:"我方账号"`                            // 账号
	AcctType    AccountType `json:"acctType" gorm:"acct_type"    bson:"acct_type"  validate:"-" title:"我方账号类型"`                      // 账号类型
	BankName    string      `json:"bankName"  gorm:"bank_name"  bson:"bank_name"   index:""   validate:"-" title:"开户银行"`             // 开户银行
	Balance     *float64    `json:"balance"  gorm:"balance"   bson:"balance"  index:""   validate:"-" title:"账户余额"`                  // 账户余额
	OppName     string      `json:"oppName"   gorm:"opp_name"  bson:"opp_name"  index:""   validate:"-" title:"对方名称"`                // 对方名称
	OppAcct     string      `json:"oppAcct"   gorm:"opp_acct"  bson:"opp_acct"  index:""   validate:"-" title:"对方账号"`                // 对方账号
	OppAcctType AccountType `json:"oppAcctType"  gorm:"opp_acct_type"   bson:"opp_acct_type"  validate:"-" title:"对方账号类型"`           // 对方账号类型
	OppBankName string      `json:"oppBankName"  gorm:"opp_bank_name"  bson:"opp_bank_name"  index:""   validate:"-" title:"对方开户银行"` // 对方开户银行
	Serial      string      `json:"serial"  gorm:"serial"   bson:"serial"  index:""   validate:"-" title:"流水号"`                      // 流水号
	Payout      float64     `json:"payout"  gorm:"payout"   bson:"payout"  index:""   validate:"-" title:"支出金额"`                     // 借方发生额（支出）
	Income      float64     `json:"income"  gorm:"income"  bson:"income"   index:""  validate:"-" title:"收入金额"`
	Cash        CashType    `json:"cash" gorm:"cash" bson:"cash" index:"" title:"现金标识"`
	Io          IOType      `json:"io" gorm:"io" bson:"io" index:"" title:"收付标志"`
	Amount      float64     `json:"amount"   gorm:"amount"  bson:"amount"  index:""  validate:"-" title:"交易金额"`  // 交易金额
	Date        time.Time   `json:"date"   gorm:"date"  bson:"date"  index:""  validate:"-" title:"交易时间"`        // 交易时间
	Ccy         string      `json:"ccy"  gorm:"ccy"  bson:"ccy"  index:""  validate:"-" title:"交易币种" `           // 交易币种
	Place       string      `json:"place"   gorm:"place"  bson:"place"  index:""  validate:"-" title:"地点"`       // 交易地点
	Summary     string      `json:"summary"  gorm:"summary"   bson:"summary" index:""   validate:"-" title:"摘要"` // 摘要
	Notes       string      `json:"notes"   gorm:"notes"  bson:"notes"  index:""  validate:"-" title:"备注"`       // 备注
}

func NewRecord() *Record {
	return &Record{}
}

func (r *Record) GetPayout() float64 {
	return math.Abs(r.Payout)
}

func (r *Record) GetIO() IOType {
	if r.Payout != 0 {
		return IOType_Out
	}
	return IOType_In
}

func NewTranId(record *Record) string {
	acct := strings.ToLower(record.Acct)
	oppAcct := strings.ToLower(record.OppAcct)
	date := fmt.Sprintf("%04d%02d%02d-%02d%02d%02d", record.Date.Year(), record.Date.Month(), record.Date.Day(), record.Date.Hour(), record.Date.Minute(), record.Date.Second())
	amount := strconv.FormatFloat(record.Amount, 'f', -1, 64)
	amount = strings.Replace(amount, ".", "_", -1)
	if record.Payout != 0 {
		x := oppAcct
		oppAcct = acct
		acct = x
	}
	return fmt.Sprintf("%s-%s-%s-%s", acct, oppAcct, amount, date)
}

func NewTranId2(tran *Tran) string {
	acct := strings.ToLower(tran.Acct)
	oppAcct := strings.ToLower(tran.OppAcct)
	date := fmt.Sprintf("%04d%02d%02d-%02d%02d%02d", tran.Date.Year(), tran.Date.Month(), tran.Date.Day(), tran.Date.Hour(), tran.Date.Minute(), tran.Date.Second())
	amount := strconv.FormatFloat(tran.Amount, 'f', -1, 64)
	amount = strings.Replace(amount, ".", "_", -1)
	return fmt.Sprintf("%s-%s-%s-%s", acct, oppAcct, amount, date)
}
