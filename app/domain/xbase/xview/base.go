package xview

import "time"

// Base
// @Description: 视图基类
type Base struct {
	Id          string     `json:"id,omitempty" gorm:"primaryKey;title:主键"  bson:"id"  desc:"主键"`                                          // 主键
	TenantId    string     `json:"tenantId,omitempty" gorm:"index:idx_tenant_id;title:租户ID"  bson:"tenant_id"  desc:"租户ID"`                // 租户ID
	CaseId      string     `json:"caseId,omitempty" gorm:"index:idx_case_id;title:案件ID"  bson:"case_id" desc:"案件ID" `                      // 案件ID
	CreatedTime *time.Time `json:"createdTime" gorm:"created_time;<-:create;title:创建时间" bson:"created_time"`                               // 创建时间
	CreatorId   string     `json:"creatorId" gorm:"creatorId;<-:create;title:创建人ID" bson:"creator_id,index:idx_creator_id"`                // 创建人ID
	CreatorName string     `json:"creatorName" gorm:"creator_name;<-:create;title:创建人" bson:"creator_name"`                                // 创建人名称
	UpdatedTime *time.Time `json:"updatedTime,omitempty"  gorm:"updated_time;title:修改人ID" gorm:"updated_time" desc:"修改时间"`                 // 修改时间
	UpdaterId   string     `json:"updaterId,omitempty"  gorm:"updater_id;title:修改人名称" bson:"updater_id,index:idx_updater_id" desc:"修改人ID"` // 修改人ID
	UpdaterName string     `json:"updaterName,omitempty"  gorm:"updater_name;title:修改人名称" bson:"updater_name" desc:"修改人名称"`                // 修改人名称
	DeletedTime *time.Time `json:"deletedTime,omitempty"  gorm:"deleted_time;title:删除时间" bson:"deleted_time" desc:"删除时间"`                  // 删除时间
	DeleterId   string     `json:"deleterId,omitempty"  gorm:"deleter_id;title:删除人ID" bson:"deleter_id,index:idx_deleter_id" desc:"删除人ID"` // 删除人ID
	DeleterName string     `json:"deleterName,omitempty"  gorm:"deleter_name;title:删除人名称" bson:"deleter_name" desc:"删除人名称"`                // 删除人名称
	IsDeleted   bool       `json:"isDeleted,omitempty"  gorm:"is_deleted;title:是否删除" bson:"is_deleted" desc:"是否删除"`                        // 是否删除
	Remarks     string     `json:"remarks,omitempty"  gorm:"remarks;title:备注" bson:"remarks" desc:"备注"`                                    // 备注
}
