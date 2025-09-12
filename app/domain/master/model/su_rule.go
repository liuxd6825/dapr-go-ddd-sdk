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
	IsEnable    bool `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	NumberValue int  `json:"number" gorm:"number" bson:"number" title:"整数倍数"`
}

// AmountCollarRule 对敲
type AmountCollarRule struct {
	IsEnable   bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	CollarDays int     `json:"collarDays" gorm:"collar_days" bson:"collar_days" title:"对敲天数"`
	TxAmount   float64 `json:"nearAmount" gorm:"near_amount" bson:"near_amount" title:"交易额度"`
	Percent    float64 `json:"nearPercent" gorm:"near_percent" bson:"near_percent" title:"兼容度"`
}

// AmountNearRule 近似金额特征
type AmountNearRule struct {
	IsEnable    bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	NearAmount  float64 `json:"nearAmount" gorm:"near_amount" bson:"near_amount" title:"临界额度"`
	NearPercent float64 `json:"nearPercent" gorm:"near_percent" bson:"near_percent" title:"临界百分比"`
}

type FreqAbnormalRule struct {
	IsEnable        bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	PeriodDays      int     `json:"periodDays" gorm:"period_days" bson:"period_days" title:"周期天数"`                   // 30 (周期天数)
	DayTolerance    int     `json:"dayTolerance" gorm:"day_tolerance" bson:"day_tolerance" title:"周期容差"`             // 2  (周期容差 ± 天数)
	AmountTolerance float64 `json:"amountTolerance" gorm:"amount_tolerance" bson:"amount_tolerance" title:"金额容差百分比"` // 5.0 (金额容差百分比 %)
	MinOccurrences  int     `json:"minOccurrences" gorm:"min_occurrences" bson:"min_occurrences" title:"最小发生次数"`     // 3  (最小发生次数)
}

type FreqHighRule struct {
	IsEnable      bool             `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	HighThreshold int              `json:"HighThreshold" gorm:"high_threshold" bson:"high_threshold" title:"交易次数阈值"  ` // 10 (交易次数阈值)
	HighPeriod    SuFreqHighPeriod `json:"HighPeriod" gorm:"high_period" bson:"high_period" title:"统计周期" `             // month, quarter, year"
}

// FreqSleepRule 休眠账号
type FreqSleepRule struct {
	IsEnable              bool `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	HibernationPeriodDays int  // 180 (休眠期天数)
	ActivationPeriodDays  int  // 30  (激活期天数)
	ActivationTxThreshold int  // 3   (激活期内交易次数阈值)
}

type PartHighRiskRule struct {
	IsEnable                 bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	NewCompanyDaysThreshold  int     `json:"newCompanyDaysThreshold" gorm:"new_company_days_threshold" bson:"new_company_days_threshold" title:"新公司定义天数"`      // 180 (新公司定义天数)
	NewCompanyLargeAmount    float64 `json:"newCompanyLargeAmount" gorm:"new_company_large_amount" bson:"new_company_large_amount" title:"对新公司的大额交易阈值"`        // 100000.00 (对新公司的大额交易阈值)
	AbnormalStatusTxAmount   float64 `json:"abnormalStatusTxAmount" gorm:"abnormal_status_tx_amount" bson:"abnormal_status_tx_amount" title:"与状态异常公司交易的金额阈值"`  // 50000.00 (与状态异常公司交易的金额阈值)
	LegalCasesTxAmount       float64 `json:"legalCasesTxAmount" gorm:"legal_cases_tx_amount" bson:"legal_cases_tx_amount" title:"与高司法风险公司交易的金额阈值"`             // 50000.00 (与高司法风险公司交易的金额阈值)
	HighLegalCasesThreshold  int     `json:"highLegalCasesThreshold" gorm:"high_legal_cases_threshold" bson:"high_legal_cases_threshold" title:"高司法风险定义：案件数量"` // 10 (高司法风险定义：案件数量)
	BusinessMismatchTxAmount float64 // 50000.00 (业务不匹配场景的金额阈值)
}

type PartAggregateRule struct {
	IsEnable bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	TopN     int     // 5 (取Top N个供应商)
	Percent  float64 // 30.0 (资金集中度百分比阈值)
}

type PartPrivateRule struct {
	IsEnable    bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	TxThreshold float64 // e.g., 50000.00 (对私交易金额阈值)
}

type PartRelationPartyRule struct {
	IsEnable bool `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
}

// TimeConcentratedPaymentsRule 集中支付规则
type TimeConcentratedPaymentsRule struct {
	IsEnable  bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	StartTime int     `json:"startTime"  title:"开始时间"`                                 // e.g., 23 (for 23:00)
	EndTime   int     `json:"endTime"  title:"结束时间"`                                   // e.g., 2  (for 02:00)
	TxCount   int     `json:"txCount"  title:"交易数量"`                                   // e.g., 10 (txs count)
	TxAmount  float64 `json:"txAmount" gorm:"tx_amount" bson:"tx_amount" title:"金额阈值"` // 可以复用 LargeValueThreshold，但独立出来更灵活
}

type TimeFastInOutRule struct {
	IsEnable       bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	FastInOutHours float64 `json:"fastInOutHours" gorm:"fast_in_out_hours" bson:"fast_in_out_hours" title:"快进快出时间"`
	Percent        float64 `json:"percent" gorm:"percent" bson:"percent" title:"金额相近度百分比"`  // 0.8
	TxAmount       float64 `json:"txAmount" gorm:"tx_amount" bson:"tx_amount" title:"金额阈值"` // 可以复用 LargeValueThreshold，但独立出来更灵活
}

type TimeNonWorkingHoursRule struct {
	IsEnable           bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	TxAmount           float64 `json:"txAmount" gorm:"tx_amount" bson:"tx_amount" title:"金额阈值"`
	IsNonWorkingSunday bool    `json:"isNonWorkingSunday" gorm:"is_non_working_sunday" bson:"is_non_working_sunday" title:"是否检查非工作日"`
	IsNonWorkingHours  bool    `json:"isNonWorkingHours" gorm:"is_non_working_hours" bson:"is_non_working_hours" title:"是否检查工作日非工时"`
	NonWorkingHoursMin int     `json:"nonWorkingHoursMin" gorm:"non_working_hours_min" bson:"non_working_hours_min" title:"非工作日最小时长"`
	NonWorkingHoursMax int     `json:"nonWorkingHoursMax" gorm:"non_working_hours_max" bson:"non_working_hours_max" title:"非工作日最长时间"`
}

type TimeSignificantDateRule struct {
	IsEnable        bool    `json:"isEnable" gorm:"is_enable" bson:"is_enable" title:"是否启用"`
	TxAmount        float64 // 金额阈值，可以复用 LargeValueThreshold，但独立出来更灵活
	DaysBefore      int     // e.g., 5 (期前天数)
	DaysAfter       int     // e.g., 2 (期后天数)
	CheckMonthEnd   bool    // [√] 月末
	CheckQuarterEnd bool    // [√] 季末
	CheckYearEnd    bool    // [√] 年末
}
