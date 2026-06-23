package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// ContractRecord 合同-记录关联
type ContractRecord struct {
	xbase.BaseModel `bson:",inline"`

	ContractId string     `json:"contractId" gorm:"contract_id" bson:"contract_id" index:"" title:"合同ID"`
	Title      string     `json:"title" gorm:"title" bson:"title" title:"标题"`
	Content    string     `json:"content" gorm:"content" bson:"content" title:"内容"`
	Status     string     `json:"status" gorm:"status" bson:"status" title:"分析状态"`
	RecordTime *time.Time `json:"recordTime" gorm:"record_time" bson:"record_time" title:"记录时间"`
	Recorder   string     `json:"recorder" gorm:"recorder" bson:"recorder" title:"记录人"`
}

func NewContractRecord() *ContractRecord {
	return &ContractRecord{}
}