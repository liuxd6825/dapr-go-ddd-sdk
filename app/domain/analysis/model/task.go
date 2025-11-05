package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

// SuTask 可疑分析任务
type SuTask struct {
	xbase.BaseModel `bson:",inline"`
	Code            string       `json:"code" bson:"code" required:"true"`
	Name            string       `json:"name" gorm:"name" bson:"name" title:"任务名称"`
	Rules           SuTaskRule   `json:"rules" gorm:"rules;type:json" bson:"rules" title:"规则"` // 规则
	StartTime       *time.Time   `json:"startTime" gorm:"start_time" bson:"start_time" title:"审计开始时间"`
	EndTime         *time.Time   `json:"endTime" gorm:"end_time" bson:"end_time" title:"审计结束时间"`
	Status          SuTaskStatus `json:"status" gorm:"status"  bson:"status" title:"状态"`
	StatusName      string       `json:"statusName" gorm:"status_name"  bson:"status_name" title:"状态名称"`
	OwnerId         string       `json:"ownerId" gorm:"owner_id"  bson:"owner_id" title:"负责人ID"`
	OwnerName       string       `json:"ownerName" gorm:"owner_name"  bson:"owner_name" title:"负责人名称"`

	MasterId   string `json:"masterId" gorm:"master_id" bson:"master_id" title:"目标ID"`
	MasterName string `json:"masterName" gorm:"master_name" bson:"master_name" title:"目标名称"`
	MasterType string `json:"masterType" gorm:"master_type" bson:"master_type" title:"目标类型"`

	RecordCount int64   `json:"recordCount" gorm:"record_count" bson:"record_count" title:"流水总数"`
	TotalAmount float64 `json:"totalAmount" gorm:"total_amount" bson:"total_amount" title:"可疑总金额"`
	SuCount     int64   `json:"suCount" gorm:"su_count" bson:"su_count" title:"可疑交易数"`
	SuHighCount int64   `json:"suHighCount" gorm:"su_high_count" bson:"su_high_count" title:"高可疑交易数"`
	WorkflowId  string  `json:"workflowId" gorm:"workflow_id" bson:"workflow_id" title:"后台流程ID"`
}

type SuTaskBillView struct {
	Task     *SuTask          `json:"task"`
	Accounts []*SuTaskAccount `json:"accounts"`
}

func NewSuTask() *SuTask {
	return &SuTask{
		Status: SuTaskStatus_New,
	}
}

type SuTran struct {
	xbase.BaseModel `bson:",inline"`
	TaskId          string            `json:"taskId" gorm:"task_id" bson:"task_id" title:"可疑任务ID" `
	TranId          string            `json:"tranId" gorm:"tran_id" bson:"tran_id" title:"交易id"`
	Name            string            `json:"name" gorm:"name"  bson:"name" title:"我方名称"`
	Acct            string            `json:"acct" gorm:"acct"  bson:"acct" title:"我方账号"`
	AcctType        model.AccountType `json:"acctType" gorm:"acct_type" bson:"acct_type"  title:"我方账号类型"`
	BankName        string            `json:"bankName" gorm:"bank_name" bson:"bank_name" title:"开户银行"`
	OppName         string            `json:"oppName" gorm:"opp_name" bson:"opp_name" title:"对方名称"`
	OppAcct         string            `json:"oppAcct" gorm:"opp_acct" bson:"opp_acct" title:"对方账号"`
	OppAcctType     model.AccountType `json:"oppAcctType" gorm:"opp_acct_type" bson:"opp_acct_type" title:"对方账号类型"`
	OppBankName     string            `json:"oppBankName" gorm:"opp_bank_name"  bson:"opp_bank_name"  title:"对方开户银行"`
	Cash            model.CashType    `json:"cash" gorm:"cash" bson:"cash" title:"现金"`
	Amount          *float64          `json:"amount" gorm:"amount"  bson:"amount"  title:"交易金额"`
	Date            *time.Time        `json:"date" gorm:"date"  bson:"date"  title:"交易时间"`
	Ccy             string            `json:"ccy" gorm:"ccy"  bson:"ccy"  title:"交易币种" `
	Place           string            `json:"place" gorm:"place"  bson:"place"  title:"地点"`
	Summary         string            `json:"summary" gorm:"summary" bson:"summary"  title:"摘要"`
	Notes           string            `json:"notes" gorm:"notes"  bson:"notes"  title:"备注"`
	Year            int               `json:"year" gorm:"year" bson:"year" title:"年"`
	Month           int               `json:"month" gorm:"month"  bson:"month" title:"月"`
	Day             int               `json:"day" gorm:"day" bson:"day" title:"日"`
	SuRisk          int               `json:"suRisk" gorm:"su_risk" bson:"su_risk"  title:"风险数" `
}

// AccountRecords is the unit of work for our workers.
type AccountRecords struct {
	OwnerName string          `bson:"owner_name" json:"ownerName" bson:"owner_name" title:"账号拥有者"`
	Account   string          `bson:"account" json:"account" bson:"account" title:"账号"`
	Records   []*model.Record `bson:"records" json:"records;type:json" bson:"records" title:"交易流水"`
}
