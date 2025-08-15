package model

// SuAmountType 金额可疑类型
type SuType string

const (
	SuType_AmountLarge       SuType = "大额交易"
	SuType_AmountRoundNumber SuType = "整数交易"
	SuType_AmountNear        SuType = "临界金额"
	SuType_AmountCollar      SuType = "对敲交易"
)

const (
	SuType_TimeNonWorkingHours     SuType = "非工作时间"
	SuType_TimeFastInFastOut       SuType = "快进快出"
	SuType_TimeConcentratedPayment SuType = "集中支付"
	SuType_TimeSignificantDate     SuType = "重大日期"
)

// SuType 频率可疑类型

const (
	SuType_Highuency               SuType = "高频交易"
	SuType_AbnormalRegularity      SuType = "异常规律"
	SuType_DormantAccountActivated SuType = "休眠账户激活"
)

// SuType 对手方可疑类型

const (
	SuType_RelatedParty      SuType = "关联方"
	SuType_PersonalAccount   SuType = "个人账户"
	SuType_HighRiskEntity    SuType = "高风险实体"
	SuType_BusinessMismatch  SuType = "业务不匹配"
	SuType_FundConcentration SuType = "资金集中"
)

// SuTaskStatus 调查任务状态
type SuTaskStatus string

const (
	SuTaskStatus_New         SuTaskStatus = "新建"
	SuTaskStatus_InProgress  SuTaskStatus = "进行中"
	SuTaskStatus_PendingInfo SuTaskStatus = "等待信息"
	SuTaskStatus_Completed   SuTaskStatus = "已完成"
	SuTaskStatus_Closed      SuTaskStatus = "已关闭"
)

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
