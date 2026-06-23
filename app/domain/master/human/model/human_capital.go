package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// HumanCapital 人员-资产关联
type HumanCapital struct {
	xbase.BaseModel `bson:",inline"`

	HumanId string `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	Name    string `json:"name" gorm:"name" bson:"name" title:"资产名称"`
	Content string `json:"content" gorm:"content" bson:"content" title:"内容"`
	Cost    string `json:"cost" gorm:"cost" bson:"cost" title:"价值"`
	Remark  string `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanCapital() *HumanCapital {
	return &HumanCapital{}
}