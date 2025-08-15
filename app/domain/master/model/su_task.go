package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"time"
)

// SuTask 可疑分析任务
type SuTask struct {
	xbase.BaseModel `bson:",inline"`
	SuTaskRule      `json:"rule" gorm:"rule" bson:"rule" title:"规则"`
	Name            string       `json:"taskName" gorm:"task_name" bson:"task_name" title:"任务名称"`
	StartTime       *time.Time   `json:"startTime" gorm:"start_time" bson:"start_time" title:"审计开始时间"`
	EndTime         *time.Time   `json:"endTime" gorm:"end_time" bson:"end_time" title:"审计结束时间"`
	Status          SuTaskStatus `json:"status" gorm:"status"  bson:"status" title:"状态"`
	OwnerName       string       `json:"ownerName" gorm:"owner_name"  bson:"owner_name" title:"操作人"`
	OwnerId         string       `json:"ownerId" gorm:"owner_id"  bson:"owner_id" title:"操作人ID"`
}

// SuTaskRule 审计规则
type SuTaskRule struct {
	SuReasonDetail `bson:",inline"`
	AmountRule     SuAmountRule `json:"amountRule" gorm:"amount_rule,json" bson:"amount_rule" title:"金额特征"`
	TimeRule       SuTimeRule   `json:"timeRule" gorm:"time_rule,json" bson:"time_rule" title:"时间特征"`
}

type SuAmountRule struct {
	LargeValue  float64 `json:"largeValue" gorm:"large_value" bson:"large_value" title:"大额金额"`
	IntValue    float64 `json:"intValue" gorm:"int_value" bson:"int_value" title:"整数倍数"`
	NearPercent float64 `json:"nearPercent" gorm:"near_percent" bson:"near_percent" title:"临界百分比"`
	NearMin     float64 `json:"nearMin" gorm:"near_min" bson:"near_min" title:"临界最小值"`
	NearMax     float64 `json:"nearMax" gorm:"near_max" bson:"near_max" title:"临界最大值"`
	CollarDays  int     `json:"collarDays" gorm:"collar_days" bson:"collar_days" title:"对敲天数"`
}

type SuTimeRule struct {
	FastInOutHours     float64 `json:"fastInOutHours" gorm:"fast_in_out_hours" bson:"fast_in_out_hours" title:"对敲天数"`
	IsNonWorkingSunday bool    `json:"isNonWorkingSunday" gorm:"is_non_working_sunday" bson:"is_non_working_sunday" title:"是否检查非工作日"`
	IsNonWorkingHours  bool    `json:"isNonWorkingHours" gorm:"is_non_working_hours" bson:"is_non_working_hours" title:"是否检查工作日非工时"`
	NonWorkingHoursMin int     `json:"nonWorkingHoursMin" gorm:"non_working_hours_min" bson:"non_working_hours_min" title:"非工作日最小时长"`
	NonWorkingHoursMax int     `json:"nonWorkingHoursMax" gorm:"non_working_hours_max" bson:"non_working_hours_max" title:"非工作日最长时间"`
}

// SuReasonDetail 审计规则
type SuReasonDetail struct {
	// 金额可疑
	IsAmountLarge  bool `json:"isAmountLarge"  gorm:"is_amount_large" bson:"is_amount_large"  title:"大额可疑" `
	IsAmountInt    bool `json:"isAmountInt"  gorm:"is_amount_int" bson:"is_amount_int" title:"整数可疑" `
	IsAmountCollar bool `json:"isAmountCollar" gorm:"is_amount_collar"  bson:"is_amount_collar" title:"对敲可疑" `
	IsAmountNear   bool `json:"isAmountNear"  gorm:"is_amount_near" bson:"is_amount_near"  title:"临界可疑" `
	// 时间可疑
	IsTimeNonWorkingHours bool `json:"isNonWorkingHours" gorm:"is_non_working_hours" bson:"is_non_working_hours" title:"非工作时间"`
	IsTimeFastInOut       bool `json:"isFastInFastOut" gorm:"is_fast_inout" bson:"is_fast_inout" title:"快速进快出"`
	IsTimePayment         bool `json:"isTimePayment" gorm:"is_time_payment" bson:"is_time_payment" title:"集中支付"`
	IsTimeSignificantDate bool `json:"isTimeSignificantDate" gorm:"is_time_significant_date" bson:"is_time_significant_date" title:"重大日期"`
	// 频率可疑
	IsFreqHigh     bool `json:"isFreqHigh" gorm:"is_freq_high" bson:"is_freq_high" title:"高频交易"`
	IsFreqAbnormal bool `json:"isFreqAbnormal" gorm:"is_freq_abnormal" bson:"is_freq_abnormal" title:"异常频率"`
	IsFreqSleep    bool `json:"isFreqSleep" gorm:"is_freq_sleep" bson:"is_freq_sleep" title:"休眠账户激活"`
	// 对手方可疑
	IsOppRelatedParty    bool `json:"IsOppRelatedParty" gorm:"is_opp_related_party" bson:"is_opp_related_party" title:"关联方"`
	IsOppPersonalAccount bool `json:"IsOppPersonalAccount" gorm:"is_opp_personal_account" bson:"is_opp_personal_account" title:"个人账户"`
	IsOppHighRiskEntity  bool `json:"IsOppHighRiskEntity" gorm:"is_opp_high_risk_entity" bson:"is_opp_high_risk_entity" title:"高风险实体"`
	IsOppBusiness        bool `json:"IsOppBusiness" gorm:"is_opp_business" bson:"is_opp_business" title:"业务不匹配"`
	IsOppConc            bool `json:"IsOppConc" gorm:"is_opp_conc" bson:"is_opp_conc" title:"资金集中"`
}

// SuTaskAccount 分析的账户
type SuTaskAccount struct {
	xbase.BaseModel `bson:",inline"`
	TaskId          string      `json:"taskId" gorm:"task_id" bson:"_task_id"`
	Account         string      `json:"account" gorm:"account"  bson:"account" title:"账户"`
	AccountType     AccountType `json:"accountType" gorm:"account_type"  bson:"account_type" title:"账户类型"`
	Name            string      `json:"name" gorm:"name"  bson:"name" title:"账户名称"`
	BankName        string      `json:"bankName" gorm:"bank_name" bson:"bank_name" title:"开户银行"`
}

type SuTran struct {
	xbase.BaseModel `bson:",inline"`
	SuReasonDetail  `bson:",inline"`
	TaskId          string      `json:"taskId" gorm:"task_id" bson:"task_id" title:"可疑任务ID" `
	TranId          string      `json:"tranId" gorm:"tran_id" bson:"tran_id" title:"交易id"`
	Name            string      `json:"name" gorm:"name"  bson:"name" title:"我方名称"`
	Acct            string      `json:"acct" gorm:"acct"  bson:"acct" title:"我方账号"`
	AcctType        AccountType `json:"acctType" gorm:"acct_type" bson:"acct_type"  title:"我方账号类型"`
	BankName        string      `json:"bankName" gorm:"bank_name" bson:"bank_name" title:"开户银行"`
	OppName         string      `json:"oppName" gorm:"opp_name" bson:"opp_name" title:"对方名称"`
	OppAcct         string      `json:"oppAcct" gorm:"opp_acct" bson:"opp_acct" title:"对方账号"`
	OppAcctType     AccountType `json:"oppAcctType" gorm:"opp_acct_type" bson:"opp_acct_type" title:"对方账号类型"`
	OppBankName     string      `json:"oppBankName" gorm:"opp_bank_name"  bson:"opp_bank_name"  title:"对方开户银行"`
	Cash            CashType    `json:"cash" gorm:"cash" bson:"cash" title:"现金"`
	Amount          *float64    `json:"amount" gorm:"amount"  bson:"amount"  title:"交易金额"`
	Date            *time.Time  `json:"date" gorm:"date"  bson:"date"  title:"交易时间"`
	Ccy             string      `json:"ccy" gorm:"ccy"  bson:"ccy"  title:"交易币种" `
	Place           string      `json:"place" gorm:"place"  bson:"place"  title:"地点"`
	Summary         string      `json:"summary" gorm:"summary" bson:"summary"  title:"摘要"`
	Notes           string      `json:"notes" gorm:"notes"  bson:"notes"  title:"备注"`
	Year            int         `json:"year" gorm:"year" bson:"year" title:"年"`
	Month           int         `json:"month" gorm:"month"  bson:"month" title:"月"`
	Day             int         `json:"day" gorm:"day" bson:"day" title:"日"`
	SuRisk          int         `json:"suRisk" gorm:"su_risk" bson:"su_risk"  title:"风险数" `
}

type SuTaskResult struct {
	TaskId  string
	Account string
	Items   map[string]*SuTaskItem
}

func NewSuTaskResult(taskId string, account string) *SuTaskResult {
	return &SuTaskResult{
		TaskId:  taskId,
		Account: account,
		Items:   make(map[string]*SuTaskItem),
	}
}

func (s *SuTaskItemReason) AddAmountCollar(tx *Tran) {
	s.AmountCollarItems = append(s.AmountCollarItems, tx)
}

func (s *SuTaskResult) AddItem(tx *Tran, suType SuType, reason string) *SuTaskItemReason {
	reasonItem := SuTaskItemReason{TranId: tx.Id, Account: tx.Acct, Reason: reason, Timestamp: tx.Date, AmountCollarItems: make([]*Tran, 0)}
	var item *SuTaskItem
	if val, ok := s.Items[tx.Id]; ok {
		item = val
		item.Reasons = append(item.Reasons, reasonItem)
	} else {
		item = &SuTaskItem{
			Tran:    *tx,
			TaskId:  s.TaskId,
			Reasons: []SuTaskItemReason{reasonItem},
		}
		s.Items[tx.Id] = item
	}
	switch suType {
	case SuType_AmountLarge:
		item.IsAmountLarge = true
	case SuType_AmountRoundNumber:
		item.IsAmountInt = true
	case SuType_AmountCollar:
		item.IsAmountCollar = true
	case SuType_AmountNear:
		item.IsAmountNear = true
	case SuType_TimeNonWorkingHours:
		item.IsTimeNonWorkingHours = true
	}
	return &reasonItem
}

// SuTaskItem 可疑交易项
type SuTaskItem struct {
	Tran           `bson:",inline"`
	SuReasonDetail `bson:",inline"`
	TaskId         string             `json:"taskId" gorm:"task_id" bson:"task_id" title:"任务ID"`
	Reasons        []SuTaskItemReason `title:"可疑原因"`
}

type SuTaskItemReason struct {
	TranId            string
	Account           string
	Reason            string
	Timestamp         time.Time
	AmountCollarItems []*Tran
}
