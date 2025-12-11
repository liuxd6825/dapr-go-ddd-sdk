package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type Status struct {
	xbase.BaseModel
	WfId     string `json:"wfId" bson:"wf_id" gorm:"wf_id" validate:"required" title:"工作流Id"`
	ActId    string `json:"actId" bson:"act_id" gorm:"act_id"`
	Name     string `json:"name" bson:"name" gorm:"name" title:"名称"`
	Status   string `json:"status" bson:"status" gorm:"status" title:"状态"`
	UserId   string `json:"userId" bson:"user_id" gorm:"user_id" title:"用户ID"`
	UserName string `json:"userName" bson:"user_name" gorm:"user_name" title:"用户"`
	OpenType string `json:"openType" bson:"open_type" gorm:"open_type" title:"打开方式"`
	OpenUrl  string `json:"openUrl" bson:"open_url" gorm:"open_url" title:"打开地址"`
}
