package doris

import (
	"strings"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
)

// RecordTableColumns Doris master_record 表的列顺序（与 DESC 输出严格一致）
// Stream Load 时通过 LoadOptions.Columns 显式映射，Doris 按位置匹配 Parquet 列
// （避免列名差异 / 类型差异导致的 0 行入库）。
const (
	RecordTableColumns = "id, tenant_id, case_id, creator_name, remark, doc_id, master_type, task_id, cash, amount, month, graph_id, master_name, name, opp_acct_type, place, updated_time, file_id, row_num, opp_acct, year, ccy, created_time, opp_bank_name, serial, day, summary, updater_id, bank_name, balance, payout, date, notes, sheet_id, master_id, acct, acct_type, income, io, creator_id, updater_name, opp_name"
)

// Record Doris master_record1 表对应的 Parquet 写入结构。
//
// 字段命名、类型与 Doris 表 master_record1 严格对齐（DESC 输出为准）。
//
// 字段类型映射（基于 github.com/parquet-go/parquet-go v0.30+ 的反射）：
//   - string              → BYTE_ARRAY (UTF8)        Doris varchar
//   - bool                → BOOLEAN                  Doris boolean
//   - int32               → INT32                    Doris int
//   - *float64            → DOUBLE (optional)        Doris double
//   - *time.Time + timestamp → INT64 + TIMESTAMP logical type  Doris datetime
//   - pointer 字段默认视为 optional（nullable）
//
// parquet-go 输出完全符合 Apache Parquet 规范（Doris 原生兼容，无需任何补丁）。
type Record struct {
	// === Basic IDs ===
	Id       string `parquet:"id"`
	TenantId string `parquet:"tenant_id"`
	CaseId   string `parquet:"case_id"`

	// === Creator / Updater ===
	CreatorId   string `parquet:"creator_id"`
	CreatorName string `parquet:"creator_name"`
	UpdaterId   string `parquet:"updater_id"`
	UpdaterName string `parquet:"updater_name"`

	// === Basic text fields ===
	Remark   string `parquet:"remark"`
	DocId    string `parquet:"doc_id"`
	TaskId   string `parquet:"task_id"`
	Name     string `parquet:"name"`
	Acct     string `parquet:"acct"`
	AcctType string `parquet:"acct_type"`
	BankName string `parquet:"bank_name"`

	// === Numeric (row_num is int32 to match Doris int; double fields are optional) ===
	RowNum  int32    `parquet:"row_num"`
	Balance *float64 `parquet:"balance,optional"`
	Payout  *float64 `parquet:"payout,optional"`
	Income  *float64 `parquet:"income,optional"`
	Amount  *float64 `parquet:"amount,optional"`

	// === Counterparty ===
	OppAcctType string `parquet:"opp_acct_type"`
	OppAcct     string `parquet:"opp_acct"`
	OppName     string `parquet:"opp_name"`
	OppBankName string `parquet:"opp_bank_name"`
	Cash        bool   `parquet:"cash"`
	Serial      string `parquet:"serial"`

	// === Master extra fields (Doris 表里有但原 Record 没有) ===
	MasterType string `parquet:"master_type"`
	MasterName string `parquet:"master_name"`
	MasterId   string `parquet:"master_id"`
	GraphId    string `parquet:"graph_id"`
	Month      int32  `parquet:"month"`
	Year       int32  `parquet:"year"`
	Day        int32  `parquet:"day"`
	Io         int32  `parquet:"io"`

	// === Time fields (optional, *time.Time; timestamp tag 让 parquet-go 写为 INT64 + TIMESTAMP) ===
	Date        *time.Time `parquet:"date,optional,timestamp"`
	CreatedTime *time.Time `parquet:"created_time,optional,timestamp"`
	UpdatedTime *time.Time `parquet:"updated_time,optional,timestamp"`

	// === Misc ===
	FileId  string `parquet:"file_id"`
	SheetId string `parquet:"sheet_id"`
	Place   string `parquet:"place"`
	Summary string `parquet:"summary"`
	Notes   string `parquet:"notes"`
	Ccy     string `parquet:"ccy"`
}

// NewRecord 从业务对象构造 Doris Record。
// 注意：当前仅设置业务可推导的列（Id/RowNum/OppAcctType 等），
// 其他列由调用方在写 Record 前自行填充。
func NewRecord(re *model.RecordIe) *Record {
	if re == nil {
		return nil
	}
	r := &Record{
		// BaseModel
		Id:          re.Id,
		TenantId:    re.TenantId,
		CaseId:      re.CaseId,
		CreatorId:   re.CreatorId,
		CreatorName: re.CreatorName,
		UpdaterId:   re.UpdaterId,
		UpdaterName: re.UpdaterName,
		Remark:      re.Remark,
		CreatedTime: re.CreatedTime,
		UpdatedTime: re.UpdatedTime,

		// RecordIe 字段
		RowNum:      int32(re.RowNum),
		TaskId:      re.TaskId,
		DocId:       re.DocId,
		FileId:      re.FileId,
		SheetId:     re.SheetId,
		Name:        re.Name,
		Acct:        re.Acct,
		AcctType:    re.AcctType,
		BankName:    re.BankName,
		Balance:     re.Balance,
		OppAcct:     re.OppAcct,
		OppAcctType: re.OppAcctType,
		OppName:     re.OppName,
		OppBankName: re.OppBankName,
		Serial:      re.Serial,
		Payout:      re.Payout,
		Income:      re.Income,
		Amount:      re.Amount,
		Date:        re.Date,
		Year:        int32(re.Date.Year()),
		Month:       int32(re.Date.Month()),
		Day:         int32(re.Date.Day()),
		Ccy:         re.Ccy,
		Place:       re.Place,
		Summary:     re.Summary,
		Notes:       re.Notes,

		// 字符串转 bool
		Cash: convertCashStringToBool(re.Cash),
		// Master*/Month/Year/Day/Io 保持零值
		MasterName: "",
		MasterId:   re.Id,
	}
	return r
}

// convertCashStringToBool 将 RecordIe.Cash 字符串转换为 Record.Cash 布尔。
// 当前实现：大小写不敏感的 "true" / "1" 视为 true，其他视为 false。
func convertCashStringToBool(s string) bool {
	s = strings.TrimSpace(s)
	return strings.EqualFold(s, "true") || s == "1"
}
