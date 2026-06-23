package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// Contract
// @Description: 合同基本信息
type Contract struct {
	xbase.BaseModel `bson:",inline"`

	Code           string     `json:"code" gorm:"code" bson:"code" title:"编码"`
	ContractCode   string     `json:"contractCode" gorm:"contract_code" bson:"contract_code" index:"" title:"合同编号"`
	Name           string     `json:"name" gorm:"name" bson:"name" index:"" title:"合同名称"`
	ContractType   string     `json:"contractType" gorm:"contract_type" bson:"contract_type" title:"合同类型"`
	SourceType     string     `json:"sourceType" gorm:"source_type" bson:"source_type" title:"数据来源"`
	ContractStatus string     `json:"contractStatus" gorm:"contract_status" bson:"contract_status" title:"合同状态"`
	Amount         float64    `json:"amount" gorm:"amount" bson:"amount" title:"合同金额"`
	PlanStartDate  *time.Time `json:"planStartDate" gorm:"plan_start_date" bson:"plan_start_date" title:"计划开始日期"`
	PlanEndDate    *time.Time `json:"planEndDate" gorm:"plan_end_date" bson:"plan_end_date" title:"计划完成日期"`
	RealStartDate  *time.Time `json:"realStartDate" gorm:"real_start_date" bson:"real_start_date" title:"实际开始日期"`
	RealEndDate    *time.Time `json:"realEndDate" gorm:"real_end_date" bson:"real_end_date" title:"实际完成日期"`
	TagId          string     `json:"tagId" gorm:"tag_id" bson:"tag_id" title:"tagId"`
	TagColor       string     `json:"tagColor" gorm:"tag_color" bson:"tag_color" title:"tagColor"`
	TagName        string     `json:"tagName" gorm:"tag_name" bson:"tag_name" title:"标签"`
	Content        string     `json:"content" gorm:"content" bson:"content" title:"合同内容"`
}

func NewContract() *Contract {
	return &Contract{}
}
