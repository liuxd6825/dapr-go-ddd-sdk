package model

import "time"

type Base struct {
	Id          string     `gorm:"id;primaryKey" json:"id" bson:"id"`
	TenantId    string     `gorm:"tenant_id" json:"tenantId,omitempty" bson:"tenant_id"` // 租户ID
	BusId       string     `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`          //业务ID  项目、某某附件
	EntityId    string     `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"` //实体ID  项目Id、某某Id
	Remark      string     `gorm:"remark" json:"remark,omitempty" bson:"remark"`
	CreatedTime *time.Time `gorm:"createdTime" json:"createdTime,omitempty" bson:"created_time"`
	CreatorId   string     `gorm:"creatorId" json:"creatorId,omitempty" bson:"creator_id"`
	CreatorName string     `gorm:"creatorName" json:"creatorName,omitempty" bson:"creator_name"`
	UpdatedTime *time.Time `gorm:"updatedTime" json:"updatedTime,omitempty" bson:"updated_time"`
	UpdaterId   string     `gorm:"updaterId" json:"updaterId,omitempty" bson:"updater_id"`
	UpdaterName string     `gorm:"updaterName" json:"updaterName,omitempty" bson:"updater_name"`
}
