package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// HumanExt 人员-扩展信息
type HumanExt struct {
	xbase.BaseModel `bson:",inline"`

	HumanId string `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	ExtType string `json:"extType" gorm:"ext_type" bson:"ext_type" title:"扩展类型"`
	Content string `json:"content" gorm:"content" bson:"content" title:"内容"`
	Remark  string `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanExt() *HumanExt {
	return &HumanExt{}
}