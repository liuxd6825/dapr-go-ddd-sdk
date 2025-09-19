package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"time"
)

type SuTaskCreateCommand struct {
	xbase.Command[SuTaskCreateData]
}
type SuTaskUpdateCommand struct {
	xbase.Command[SuTaskCreateData]
}

// SuTaskCreateData 可疑分析任务
type SuTaskCreateData struct {
	xbase.BaseModel `bson:",inline"`
	Code            string           `json:"code" gorm:"code" bson:"code" title:"编号" validate:"required" `
	Name            string           `json:"taskName" gorm:"task_name" bson:"task_name" title:"任务名称"  `
	Rules           model.SuTaskRule `json:"rules" gorm:"rules;type:json" bson:"rules" title:"规则"` // 规则
	StartTime       *time.Time       `json:"startTime" gorm:"start_time" bson:"start_time" title:"审计开始时间"  `
	EndTime         *time.Time       `json:"endTime" gorm:"end_time" bson:"end_time" title:"审计结束时间"`
	OwnerId         string           `json:"ownerId" gorm:"owner_id"  bson:"owner_id" title:"负责人ID" `
	OwnerName       string           `json:"ownerName" gorm:"owner_name"  bson:"owner_name" title:"负责人名称"  `
	TargetId        string           `json:"targetId" gorm:"target_id" bson:"target_id" title:"目标ID" `
	TargetName      string           `json:"targetName" gorm:"target_name" bson:"target_name" title:"目标名称"`
	TargetType      string           `json:"targetType" gorm:"target_type" bson:"target_type" title:"目标类型"`
}

type SuTaskStartCommand struct {
	xbase.Command[SuTaskStartData]
}

// SuTaskStartData 可疑分析任务
type SuTaskStartData struct {
	Id string `json:"id" bson:"id"  validate:"required" title:"任务ID"`
}
