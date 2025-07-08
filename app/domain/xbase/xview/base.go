package xview

import "time"

// BaseView
// @Description: 视图基类
type Base struct {
	Id       string `json:"id,omitempty"  bson:"_id" gorm:"primaryKey" desc:"主键"`                        // 主键
	TenantId string `json:"tenantId,omitempty"  bson:"tenant_id" gorm:"index:idx_tenant_id" desc:"租户ID"` // 租户ID
	CaseId   string `json:"caseId,omitempty"  bson:"case_id" gorm:"index:idx_case_id" desc:"案件ID" `      // 案件ID

	CreatedTime *time.Time `json:"createdTime,omitempty"  bson:"created_time" gorm:"created_time" desc:"创建时间" `                // 创建时间
	CreatorId   string     `json:"creatorId,omitempty"  bson:"creator_id" gorm:"creator_id,index:idx_creator_id" desc:"创建人ID"` // 创建人ID
	CreatorName string     `json:"creatorName,omitempty"  bson:"creator_name" gorm:"creator_name" desc:"创建人名称"`                // 创建人名称

	DeletedTime *time.Time `json:"deletedTime,omitempty"  bson:"deleted_time" gorm:"deleted_time" desc:"删除时间"`                 // 删除时间
	DeleterId   string     `json:"deleterId,omitempty"  bson:"deleter_id" gorm:"deleter_id,index:idx_deleter_id" desc:"删除人ID"` // 删除人ID
	DeleterName string     `json:"deleterName,omitempty"  bson:"deleter_name" gorm:"deleter_name" desc:"删除人名称"`                // 删除人名称

	UpdatedTime *time.Time `json:"updatedTime,omitempty"  bson:"updated_time" gorm:"updated_time" desc:"修改时间"`                 // 修改时间
	UpdaterId   string     `json:"updaterId,omitempty"  bson:"updater_id" gorm:"updater_id,index:idx_updater_id" desc:"修改人ID"` // 修改人ID
	UpdaterName string     `json:"updaterName,omitempty"  bson:"updater_name" gorm:"updater_name" desc:"修改人名称"`                // 修改人名称

	IsDeleted bool   `json:"isDeleted,omitempty"  bson:"is_deleted" gorm:"is_deleted" desc:"是否删除"` // 是否删除
	Remarks   string `json:"remarks,omitempty"  bson:"remarks" gorm:"remarks" desc:"备注"`           // 备注
}
