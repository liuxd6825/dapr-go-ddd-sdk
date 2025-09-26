package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"

type SuRecord struct {
	model.Record            `bson:",inline"`
	RecordId                string           `json:"recordId" gorm:"record_id" bson:"record_id" title:"交易ID" `
	TaskId                  string           `json:"taskId" gorm:"task_id" bson:"task_id" title:"可疑任务ID" `
	AmountLarge             bool             `json:"amountLarge" gorm:"amount_large" bson:"amount_large"  `
	AmountNumber            bool             `json:"amountNumber" gorm:"amount_number" bson:"amount_number"  `
	AmountCollar            bool             `json:"amountCollar" gorm:"amount_collar" bson:"amount_collar"  `
	AmountNear              bool             `json:"amountNear" gorm:"amount_near" bson:"amount_near"  `
	TimeNonWorking          bool             `json:"timeNonWorking" gorm:"time_non_working" bson:"time_non_working" title:"非工作时间"`
	TimeFastInOut           bool             `json:"timeFastInOut" gorm:"time_fast_inout" bson:"time_fast_inout" title:"快进快出"`
	TimeConcentratedPayment bool             `json:"timeConcentratedPayment" gorm:"time_concentrated_payment"  bson:"time_concentrated_payment"  title:"集中支付"`
	TimeSignificantDate     bool             `json:"timeSignificantDate" gorm:"time_significant_date"  bson:"time_significant_date"  title:"重大日期"`
	FreqHigh                bool             `json:"freqHigh" gorm:"freq_high" bson:"freq_high"  title:"高频交易"`
	FreqAbnormal            bool             `json:"freqAbnormal" gorm:"freq_abnormal" bson:"freq_abnormal"  title:"异常规律"`
	FreqSleep               bool             `json:"freqSleep" gorm:"freq_sleep" bson:"freq_sleep"  title:"休眠账户激活"`
	PartyRelated            bool             `json:"partyRelated" gorm:"party_related" bson:"party_related"  title:"关联方"`
	PartyPrivate            bool             `json:"partyPersonal" gorm:"party_personal" bson:"party_personal"  title:"个人账户"`
	PartyHighRisk           bool             `json:"partyHighRisk" gorm:"party_high_risk" bson:"party_high_risk"  title:"高风险实体"`
	PartyBusiness           bool             `json:"partyBusiness" gorm:"party_business" bson:"party_business"  title:"业务不匹配"`
	PartyAggregate          bool             `json:"partyAggregate" gorm:"party_aggregate" bson:"party_aggregate"  title:"资金集中"`
	Risk                    int              `json:"risk" gorm:"risk" bson:"risk"  title:"风险数" `
	AmountTag               bool             `json:"amountTag" gorm:"amount_tag" bson:"amount_tag"  title:"金额标签" `
	TimeTag                 bool             `json:"timeTag" gorm:"time_tag" bson:"time_tag"  title:"时间标签"`
	FreqTag                 bool             `json:"freqTag" gorm:"freq_tag" bson:"freq_tag"  title:"频率标签"`
	PartyTag                bool             `json:"partyTag" gorm:"party_tag" bson:"party_tag"  title:"对手方标签"`
	Reasons                 []SuRecordReason `json:"reasons" gorm:"reasons;type:json" bson:"reasons" title:"原因"`
}

type SuRecordReason struct {
	Reason string `json:"reason" bson:"reason"`
	Type   SuType `json:"type" bson:"type"`
}
