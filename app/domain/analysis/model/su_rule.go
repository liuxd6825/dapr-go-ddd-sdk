package model

// SuTaskRule 审计规则
type SuTaskRule struct {
	AmountLarge  AmountLargeRule  `json:"amountLarge" gorm:"amount_large,json" bson:"amount_large" title:"大额"`
	AmountCollar AmountCollarRule `json:"amountCollar" gorm:"amount_collar,json" bson:"amount_collar" title:"对敲"`
	AmountNear   AmountNearRule   `json:"amountNear" gorm:"amount_near,json" bson:"amount_near" title:"近似金额特征"`
	AmountNumber AmountNumberRule `json:"amountNumber" gorm:"amount_number,json" bson:"amount_number" title:"整数金额"`

	FreqSleep    FreqSleepRule    `json:"freqSleep" gorm:"freq_sleep,json" bson:"freq_sleep" title:"休眠账号"`
	FreqAbnormal FreqAbnormalRule `json:"freqAbnormal" gorm:"freq_abnormal,json" bson:"freq_abnormal" title:"异常规律性支付"`
	FreqHigh     FreqHighRule     `json:"freqHigh" gorm:"freq_high,json" bson:"freq_high" title:"高频率交易"`

	PartAggregate PartAggregateRule     `json:"partAggregate" gorm:"part_aggregate,json" bson:"part_aggregate" title:"集中支付"`
	PartPrivate   PartPrivateRule       `json:"partPrivate" gorm:"part_private,json" bson:"part_private" title:"对私交易"`
	PartHighRisk  PartHighRiskRule      `json:"partHighRisk" gorm:"part_high_risk,json" bson:"part_high_risk" title:"高风险对手方"`
	PartRelation  PartRelationPartyRule `json:"partRelation" gorm:"part_relation,json" bson:"part_relation" title:"相关方交易"`

	TimeConcentratedPayments TimeConcentratedPaymentsRule `json:"timeConcentratedPayments" gorm:"time_concentrated_payments,json"  bson:"time_concentrated_payments"  title:"集中支付"`
	TimeFastInOut            TimeFastInOutRule            `json:"timeFastInOut"  gorm:"time_fast_in_out,json" bson:"time_fast_in_out"  title:"快进快出"`
	TimeNonWorkingHours      TimeNonWorkingHoursRule      `json:"timeNonWorkingHours" gorm:"time_non_working_hours,json" bson:"time_non_working_hours" title:"非工作时"`
	TimeSignificantDate      TimeSignificantDateRule      `json:"timeSignificantDate" gorm:"time_significant_date,json" bson:"time_significant_date" title:"重大日期"`
}

// AmountLargeRule 大额
type AmountLargeRule struct {
	IsEnable   bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	LargeValue float64 `json:"largeValue" gorm:"large_value" bson:"large_value" title:"大额金额"`
}

// AmountNumberRule 整数金额
type AmountNumberRule struct {
	IsEnable bool `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	Number   int  `json:"number" gorm:"number" bson:"number" title:"整数倍数"`
}

// AmountCollarRule 对敲
type AmountCollarRule struct {
	IsEnable bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	Days     int     `json:"days" gorm:"days" bson:"days" title:"对敲天数"`
	Amount   float64 `json:"amount" gorm:"amount" bson:"amount" title:"交易额度"`
	Percent  float64 `json:"percent" gorm:"percent" bson:"percent" title:"兼容度"`
}

// AmountNearRule 近似金额特征
type AmountNearRule struct {
	IsEnable    bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	NearAmount  float64 `json:"nearAmount" gorm:"near_amount" bson:"near_amount" title:"临界额度"`
	NearPercent float64 `json:"nearPercent" gorm:"near_percent" bson:"near_percent" title:"临界百分比"`
}

// FreqAbnormalRule
// @Description: 异常规律性支付
type FreqAbnormalRule struct {
	IsEnable        bool               `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	PeriodDays      int                `json:"periodDays" gorm:"period_days" bson:"period_days" title:"周期天数"`                   // 30 (周期天数)
	Period          FreqAbnormalPeriod `json:"period" gorm:"period" bson:"period" title:"周期单位"`                                 //
	DayTolerance    int                `json:"dayTolerance" gorm:"day_tolerance" bson:"day_tolerance" title:"周期容差"`             // 2  (周期容差 ± 天数)
	AmountTolerance float64            `json:"amountTolerance" gorm:"amount_tolerance" bson:"amount_tolerance" title:"金额容差百分比"` // 5.0 (金额容差百分比 %)
	MinOccurrences  int                `json:"minOccurrences" gorm:"min_occurrences" bson:"min_occurrences" title:"最小发生次数"`     // 3  (最小发生次数)
}

// FreqHighRule 高频交易
type FreqHighRule struct {
	IsEnable  bool             `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	Threshold int              `json:"threshold" gorm:"threshold" bson:"threshold" title:"交易次数阈值"  ` // 10 (交易次数阈值)
	Period    SuFreqHighPeriod `json:"period" gorm:"period" bson:"period" title:"统计周期" `             // month, quarter, year"
}

// FreqSleepRule 休眠账号
type FreqSleepRule struct {
	IsEnable          bool `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	SleepDays         int  `json:"sleepDays" gorm:"sleep_days" bson:"sleep_days" title:"休眠期天数"`                              // 180 (休眠期天数)
	ActivateDays      int  `json:"activateDays" gorm:"activate_days" bson:"activate_days" title:"激活期天数"`                     // 30  (激活期天数)
	ActivateThreshold int  `json:"activateThreshold" gorm:"activate_threshold" bson:"activate_threshold" title:"激活期内交易次数阈值"` // 3   (激活期内交易次数阈值)
}

// PartHighRiskRule
// @Description:高风险对手方
type PartHighRiskRule struct {
	IsEnable                 bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	NewCompanyDays           int     `json:"newCompanyDays" gorm:"new_company_days" bson:"new_company_days" title:"新公司定义天数"`                                      // 180 (新公司定义天数)
	NewCompanyLargeAmount    float64 `json:"newCompanyLargeAmount" gorm:"new_company_large_amount" bson:"new_company_large_amount" title:"对新公司的大额交易阈值"`           // 100000.00 (对新公司的大额交易阈值)
	AbnormalStatusTxAmount   float64 `json:"abnormalStatusTxAmount" gorm:"abnormal_status_tx_amount" bson:"abnormal_status_tx_amount" title:"与状态异常公司交易的金额阈值"`     // 50000.00 (与状态异常公司交易的金额阈值)
	LegalCasesTxAmount       float64 `json:"legalCasesTxAmount" gorm:"legal_cases_tx_amount" bson:"legal_cases_tx_amount" title:"与高司法风险公司交易的金额阈值"`                // 50000.00 (与高司法风险公司交易的金额阈值)
	HighLegalCasesThreshold  int     `json:"highLegalCasesThreshold" gorm:"high_legal_cases_threshold" bson:"high_legal_cases_threshold" title:"高司法风险定义：案件数量"`    // 10 (高司法风险定义：案件数量)
	BusinessMismatchTxAmount float64 `json:"businessMismatchTxAmount" gorm:"business_mismatch_tx_amount" bson:"business_mismatch_tx_amount" title:"业务不匹配场景的金额阈值"` // 50000.00 (业务不匹配场景的金额阈值)
}

type PartAggregateRule struct {
	IsEnable bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	TopN     int     `json:"topN" gorm:"top_n" bson:"top_n" title:"供应商"`               // 5 (取Top N个供应商)
	Percent  float64 `json:"percent" gorm:"percent" bson:"percent" title:"资金集中度百分比阈值"` // 30.0 (资金集中度百分比阈值)
}

// PartPrivateRule
// @Description:关联方交易
type PartPrivateRule struct {
	IsEnable    bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	TxThreshold float64 `json:"txThreshold" gorm:"tx_threshold" bson:"tx_threshold" title:"对私交易金额阈值"` // e.g., 50000.00 (对私交易金额阈值)
}

type PartRelationPartyRule struct {
	IsEnable bool `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
}

// TimeConcentratedPaymentsRule 集中支付规则
type TimeConcentratedPaymentsRule struct {
	IsEnable  bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	StartTime int     `json:"startTime"  gorm:"startTime" bson:"startTime" title:"开始时间"` // e.g., 23 (for 23:00)
	EndTime   int     `json:"endTime"   gorm:"endTime" bson:"endTime" title:"结束时间"`      // e.g., 2  (for 02:00)
	Count     int     `json:"count"  gorm:"count" bson:"count" title:"交易数量"`             // e.g., 10 (txs count)
	Amount    float64 `json:"amount" gorm:"amount" bson:"amount" title:"金额阈值"`           // 可以复用 LargeValueThreshold，但独立出来更灵活
}

type TimeFastInOutRule struct {
	IsEnable       bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	FastInOutHours float64 `json:"fastInOutHours" gorm:"fast_in_out_hours" bson:"fast_in_out_hours" title:"快进快出时间"`
	Percent        float64 `json:"percent" gorm:"percent" bson:"percent" title:"金额相近度百分比"` // 0.8
	Amount         float64 `json:"amount" gorm:"amount" bson:"amount" title:"金额阈值"`        // 可以复用 LargeValueThreshold，但独立出来更灵活
}

type TimeNonWorkingHoursRule struct {
	IsEnable           bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	Amount             float64 `json:"amount" gorm:"amount" bson:"amount" title:"金额阈值"`
	IsNonWorkingSunday bool    `json:"isNonWorkingSunday" gorm:"is_non_working_sunday" bson:"is_non_working_sunday" title:"是否检查非工作日"`
	IsNonWorkingHours  bool    `json:"isNonWorkingHours" gorm:"is_non_working_hours" bson:"is_non_working_hours" title:"是否检查工作日非工时"`
	HoursMin           int     `json:"hoursMin" gorm:"hours_min" bson:"hours_min" title:"最小时长"`
	HoursMax           int     `json:"hoursMax" gorm:"hours_max" bson:"hours_max" title:"最长时间"`
}

type TimeSignificantDateRule struct {
	IsEnable     bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	Amount       float64 `json:"amount" gorm:"amount" bson:"amount" title:"金额阈值"`                   // 金额阈值，可以复用 LargeValueThreshold，但独立出来更灵活
	DaysBefore   int     `json:"daysBefore" gorm:"days_before" bson:"days_before" title:"期前天数"`     // e.g., 5 (期前天数)
	DaysAfter    int     `json:"daysAfter" gorm:"days_after" bson:"days_after" title:"期后天数"`        // e.g., 2 (期后天数)
	CheckMonth   bool    `json:"checkMonth" gorm:"check_month" bson:"check_month" title:"月末"`       // [√] 月末
	CheckQuarter bool    `json:"checkQuarter" gorm:"check_quarter" bson:"check_quarter" title:"季末"` // [√] 季末
	CheckYear    bool    `json:"checkYear" gorm:"check_year" bson:"check_year" title:"年末"`          // [√] 年末
}
