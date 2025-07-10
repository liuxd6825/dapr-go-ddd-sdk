package xmodel

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

type Base struct {
	Id          string      `json:"id" gorm:"type:varchar(255);primary_key" bson:"id"`
	CaseId      string      `json:"caseId" gorm:"type:varchar(255);case_id" bson:"case_id"`
	TenantId    string      `json:"tenantId" gorm:"type:varchar(255);tenant_id" bson:"tenant_id"`
	CreatedTime *times.Time `json:"createdTime" gorm:"created_time;<-:create" bson:"created_time"`
	CreatorId   string      `json:"creatorId" gorm:"creator_id;<-:create" bson:"creator_id"`
	CreatorName string      `json:"creatorName" gorm:"creator_name;<-:create" bson:"creator_name"`
	UpdatedTime *times.Time `json:"updatedTime" gorm:"updated_time" bson:"updated_time"`
	UpdaterId   string      `json:"updaterId" gorm:"updater_id" bson:"updater_id"`
	UpdaterName string      `json:"updaterName" gorm:"updater_name" bson:"updater_name"`
	DeletedTime *times.Time `json:"deletedTime" gorm:"deleted_time" bson:"deleted_time"`
	DeleterId   string      `json:"deleterId" gorm:"deleter_id" bson:"deleter_id"`
	DeleterName string      `json:"deleterName" gorm:"deleter_name" bson:"deleter_name"`
	IsDeleted   bool        `json:"isDeleted" gorm:"is_deleted" bson:"is_deleted"`
	Remark      string      `json:"remark" gorm:"remark" bson:"remark"`
}
