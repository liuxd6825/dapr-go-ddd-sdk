package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// ProductRecord 产品-记录
// ⚠ 严格按源保留：源 webscript 中 product-record/type.d.ts 错误使用 contractId
// 这里保留 source bug，DAO 级联删除使用 productId
type ProductRecord struct {
	xbase.BaseModel `bson:",inline"`

	ContractId string     `json:"contractId" gorm:"contract_id" bson:"contract_id" index:"" title:"合同ID"` // ⚠ source bug
	Title      string     `json:"title" gorm:"title" bson:"title" title:"标题"`
	Content    string     `json:"content" gorm:"content" bson:"content" title:"内容"`
	Status     bool       `json:"status" gorm:"status" bson:"status" title:"状态"`
	RecordTime *time.Time `json:"recordTime" gorm:"record_time" bson:"record_time" title:"记录时间"`
	Recorder   string     `json:"recorder" gorm:"recorder" bson:"recorder" title:"记录人"`
	Remark     string     `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewProductRecord() *ProductRecord {
	return &ProductRecord{}
}