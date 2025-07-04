package model

import "time"

type Base struct {
	Id          string     `gorm:"id;primaryKey;title:主键" json:"id" bson:"id"`
	TenantId    string     `gorm:"tenant_id;title:租户ID" json:"tenantId,omitempty" bson:"tenant_id"` // 租户ID
	CaseId      string     `gorm:"case_id;title:案件ID" json:"caseId,omitempty" bson:"case_id"`
	Remark      string     `gorm:"remark;title:备注" json:"remark,omitempty" bson:"remark"`
	CreatedTime *time.Time `gorm:"created_time" json:"createdTime,omitempty" bson:"created_time"`
	CreatorId   string     `gorm:"creator_id" json:"creatorId,omitempty" bson:"creator_id"`
	CreatorName string     `gorm:"creator_name" json:"creatorName,omitempty" bson:"creator_name"`
	UpdatedTime *time.Time `gorm:"updated_time" json:"updatedTime,omitempty" bson:"updated_time"`
	UpdaterId   string     `gorm:"updater_id" json:"updaterId,omitempty" bson:"updater_id"`
	UpdaterName string     `gorm:"updater_name" json:"updaterName,omitempty" bson:"updater_name"`
}
