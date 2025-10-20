package model

import "time"

type Base struct {
	Id          string     `json:"id" gorm:"primaryKey;title:主键" bson:"id" title:"主键" validate:"required" `                       // 主键
	TenantId    string     `json:"tenantId" gorm:"index:idx_tenant_id;title:租户ID"  bson:"tenant_id"  title:"租户ID" `               // 租户ID
	CreatedTime *time.Time `json:"createdTime" gorm:"created_time;<-:create;title:创建时间" bson:"created_time"`                      // 创建时间
	CreatorId   string     `json:"creatorId" gorm:"creator_id;<-:create;title:创建人ID" bson:"creator_id,index:idx_creator_id"`      // 创建人ID
	CreatorName string     `json:"creatorName" gorm:"creator_name;<-:create;title:创建人" bson:"creator_name"`                       // 创建人名称
	UpdatedTime *time.Time `json:"updatedTime"  gorm:"updated_time;title:修改人ID" bson:"updated_time" title:"修改时间"`                 // 修改时间
	UpdaterId   string     `json:"updaterId"  gorm:"updater_id;title:修改人名称" bson:"updater_id,index:idx_updater_id" title:"修改人ID"` // 修改人ID
	UpdaterName string     `json:"updaterName"  gorm:"updater_name;title:修改人名称" bson:"updater_name" title:"修改人名称"`                // 修改人名称
	DeletedTime *time.Time `json:"deletedTime"  gorm:"deleted_time;title:删除时间" bson:"deleted_time" title:"删除时间"`                  // 删除时间
	DeleterId   string     `json:"deleterId"  gorm:"deleter_id;title:删除人ID" bson:"deleter_id,index:idx_deleter_id" title:"删除人ID"` // 删除人ID
	DeleterName string     `json:"deleterName"  gorm:"deleter_name;title:删除人名称" bson:"deleter_name" title:"删除人名称"`                // 删除人名称
	IsDeleted   bool       `json:"isDeleted"  gorm:"is_deleted;title:是否删除" bson:"is_deleted" title:"是否删除"`                        // 是否删除
	Remark      string     `json:"remark"  gorm:"remark;title:备注" bson:"remark" title:"备注"`
}
