package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// HumanReportedAmount 报案人金额
type HumanReportedAmount struct {
	xbase.BaseModel `bson:",inline"`

	HumanId              string  `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	ReportInvestAmount   float64 `json:"reportInvestAmount" gorm:"report_invest_amount" bson:"report_invest_amount" title:"报案投资金额"`
	ReportRefundAmount   float64 `json:"reportRefundAmount" gorm:"report_refund_amount" bson:"report_refund_amount" title:"报案返款金额"`
	ReportLossAmount     float64 `json:"reportLossAmount" gorm:"report_loss_amount" bson:"report_loss_amount" title:"报案损失金额"`
	ContractInvestAmount float64 `json:"contractInvestAmount" gorm:"contract_invest_amount" bson:"contract_invest_amount" title:"合同投资金额"`
	ContractRefundAmount float64 `json:"contractRefundAmount" gorm:"contract_refund_amount" bson:"contract_refund_amount" title:"合同返款金额"`
	ContractLossAmount   float64 `json:"contractLossAmount" gorm:"contract_loss_amount" bson:"contract_loss_amount" title:"合同损失金额"`
	DiffInvestAmount     float64 `json:"diffInvestAmount" gorm:"diff_invest_amount" bson:"diff_invest_amount" title:"差额投资金额"`
	DiffRefundAmount     float64 `json:"diffRefundAmount" gorm:"diff_refund_amount" bson:"diff_refund_amount" title:"差额返款金额"`
	DiffLossAmount       float64 `json:"diffLossAmount" gorm:"diff_loss_amount" bson:"diff_loss_amount" title:"差额损失金额"`
	Remark               string  `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanReportedAmount() *HumanReportedAmount {
	return &HumanReportedAmount{}
}