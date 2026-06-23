package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// HumanSuspectAmount 嫌疑人金额
type HumanSuspectAmount struct {
	xbase.BaseModel `bson:",inline"`

	HumanId        string  `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	UnlawfulAmount float64 `json:"unlawfulAmount" gorm:"unlawful_amount" bson:"unlawful_amount" title:"违法所得金额"`
	ReportAmount   float64 `json:"reportAmount" gorm:"report_amount" bson:"report_amount" title:"报案金额"`
	ReturnAmount   float64 `json:"returnAmount" gorm:"return_amount" bson:"return_amount" title:"退赃金额"`
	Remark         string  `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanSuspectAmount() *HumanSuspectAmount {
	return &HumanSuspectAmount{}
}