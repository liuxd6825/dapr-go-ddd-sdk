package entities

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
)

type Base struct {
	Id          string      `gorm:"primary_key" json:"id" jsonschema:"required"`
	TenantId    string      `gorm:"type:varchar(64);not null;" json:"tenantId"  jsonschema:"required"`
	CaseId      string      `gorm:"type:varchar(64);" json:"caseId" jsonschema:""`
	CreatedTime *times.Time `gorm:"type:date" json:"createdTime" jsonschema:"date"`
	CreatorId   string      `gorm:"type:varchar(64);not null;" json:"creatorId" jsonschema:""`
	CreatorName string      `gorm:"type:varchar(64);not null;" json:"creatorName" jsonschema:""`
	UpdatedTime *times.Time `gorm:"type:date" json:"updatedTime"  jsonschema:""`
	UpdaterId   string      `gorm:"type:varchar(64);not null;" json:"updaterId" jsonschema:""`
	UpdaterName string      `gorm:"type:varchar(64);not null;" json:"updaterName" jsonschema:""`
	DeletedTime *times.Time `gorm:"type:date" json:"deletedTime" jsonschema:""`
	DeleterId   string      `gorm:"type:varchar(64);" json:"deleterId" jsonschema:""`
	DeleterName string      `gorm:"type:varchar(64);" json:"deleterName" jsonschema:""`
	IsDeleted   bool        `gorm:"type:boolean;default:0;" json:"isDeleted" jsonschema:""`
}
