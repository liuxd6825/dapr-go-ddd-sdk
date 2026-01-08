package model

// SuAmountType 金额可疑类型
type SuType string

const (
	SuType_AmountLarge  SuType = "大额交易"
	SuType_AmountNumber SuType = "整数交易"
	SuType_AmountNear   SuType = "临界金额"
	SuType_AmountCollar SuType = "对敲交易"
)

const (
	SuType_TimeNonWorking          SuType = "非工作时间"
	SuType_TimeFastInOut           SuType = "快进快出"
	SuType_TimeConcentratedPayment SuType = "集中支付"
	SuType_TimeSignificantDate     SuType = "重大日期"
)

// SuType 频率可疑类型
const (
	SuType_FreqHigh     SuType = "高频交易"
	SuType_FreqAbnormal SuType = "异常规律"
	SuType_FreqSleep    SuType = "休眠账户激活"
)

// SuType 对手方可疑类型

const (
	SuType_PartyRelated   SuType = "关联方"
	SuType_PartyPrivate   SuType = "个人账户"
	SuType_PartyHighRisk  SuType = "高风险实体"
	SuType_PartyBusiness  SuType = "业务不匹配"
	SuType_PartyAggregate SuType = "资金集中"
)

// SuTaskStatus 调查任务状态
type SuTaskStatus int

const (
	SuTaskStatus_New         SuTaskStatus = iota // 任务新建
	SuTaskStatus_TaskQueuing                     // 任务排队
	SuTaskStatus_InProgress                      // 系统计算
	SuTaskStatus_Inspect                         // 流水调查
	SuTaskStatus_Completed                       // 任务完成
	SuTaskStatus_Closed                          // 任务关闭
)

func (s SuTaskStatus) String() string {
	switch s {
	case SuTaskStatus_New:
		return "任务新建"
	case SuTaskStatus_TaskQueuing:
		return "排队中"
	case SuTaskStatus_InProgress:
		return "计算中"
	case SuTaskStatus_Inspect:
		return "流水调查"
	case SuTaskStatus_Completed:
		return "任务完成"
	case SuTaskStatus_Closed:
		return "任务关闭"
	default:
		return "未知状态"
	}
}

// SuConclusionType 审计结论
type SuConclusionType string

const (
	SuConclusionType_ExplainedNoAbnormality    SuConclusionType = "已解释/无异常"
	SuConclusionType_InternalControlDeficiency SuConclusionType = "内部控制缺陷"
	SuConclusionType_SuspectedFraudOrViolation SuConclusionType = "疑似违规/舞弊"
)

// SuEvidenceType 证据类型
type SuEvidenceType string

const (
	SuEvidenceType_AccountingVoucher    SuEvidenceType = "会计凭证"
	SuEvidenceType_Contract             SuEvidenceType = "合同"
	SuEvidenceType_Invoice              SuEvidenceType = "发票"
	SuEvidenceType_WarehouseReceipt     SuEvidenceType = "入库单"
	SuEvidenceType_ConfirmationLetter   SuEvidenceType = "询证函回函"
	SuEvidenceType_InterviewRecord      SuEvidenceType = "访谈记录"
	SuEvidenceType_PublicInfoScreenshot SuEvidenceType = "公示信息截图"
)

// SuVerificationStatus 证据核对状态
type SuVerificationStatus string

const (
	SuVerificationStatus_VerifiedConsistent SuVerificationStatus = "核对一致"
	SuVerificationStatus_Inconsistent       SuVerificationStatus = "不一致"
	SuVerificationStatus_Missing            SuVerificationStatus = "缺失"
	SuVerificationStatus_ForgedSuspected    SuVerificationStatus = "疑似伪造"
)

// SuActivityType 验证活动类型
type SuActivityType string

const (
	SuActivityType_Confirmation SuActivityType = "发函询证"
	SuActivityType_Interview    SuActivityType = "约谈"
)

type FreqAbnormalPeriod string

const (
	FreqAbnormalPeriod_Day   FreqAbnormalPeriod = "day"
	FreqAbnormalPeriod_Month FreqAbnormalPeriod = "month"
	FreqAbnormalPeriod_Year  FreqAbnormalPeriod = "year"
)

type SuFreqHighPeriod string

const (
	SuFreqHighPeriod_Month   SuFreqHighPeriod = "month"
	SuFreqHighPeriod_Quarter SuFreqHighPeriod = "quarter"
	SuFreqHighPeriod_Year    SuFreqHighPeriod = "year"
)
