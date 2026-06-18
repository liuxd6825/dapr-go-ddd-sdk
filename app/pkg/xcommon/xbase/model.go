package xbase

import (
	"time"
)

// BaseModel
// @Description: 实体基类
type BaseModel struct {
	Id          string     `json:"id" gorm:"primaryKey;title:主键" bson:"id" index:"unique" title:"主键" validate:"required" parquet:"name=id, type=BYTE_ARRAY,  repetitiontype=REQUIRED"`                 // 主键
	TenantId    string     `json:"tenantId" gorm:"index:idx_tenant_id;title:租户ID"  bson:"tenant_id"  title:"租户ID" parquet:"name=tenant_id, type=BYTE_ARRAY,  repetitiontype=REQUIRED"`                 // 租户ID
	CaseId      string     `json:"caseId" gorm:"index:idx_case_id;title:案件ID"  bson:"case_id" title:"案件ID" parquet:"name=case_id, type=BYTE_ARRAY,  repetitiontype=REQUIRED"`                          // 案件ID
	CreatedTime *time.Time `json:"createdTime" gorm:"created_time;<-:create;title:创建时间" bson:"created_time" parquet:"name=created_time, type=INT64,  repetitiontype=OPTIONAL"`                         // 创建时间
	CreatorId   string     `json:"creatorId" gorm:"creator_id;<-:create;title:创建人ID" bson:"creator_id,index:idx_creator_id" parquet:"name=creator_id, type=BYTE_ARRAY,  repetitiontype=REQUIRED"`      // 创建人ID
	CreatorName string     `json:"creatorName" gorm:"creator_name;<-:create;title:创建人" bson:"creator_name" parquet:"name=creator_name, type=BYTE_ARRAY,  repetitiontype=REQUIRED"`                     // 创建人名称
	UpdatedTime *time.Time `json:"updatedTime"  gorm:"updated_time;title:修改人ID" bson:"updated_time" title:"修改时间" parquet:"name=updated_time, type=INT64, repetitiontype=OPTIONAL"`                     // 修改时间
	UpdaterId   string     `json:"updaterId"  gorm:"updater_id;title:修改人名称" bson:"updater_id,index:idx_updater_id" title:"修改人ID" parquet:"name=updater_id, type=BYTE_ARRAY,  repetitiontype=REQUIRED"` // 修改人ID
	UpdaterName string     `json:"updaterName"  gorm:"updater_name;title:修改人名称" bson:"updater_name" title:"修改人名称" parquet:"name=updater_name, type=BYTE_ARRAY,  repetitiontype=REQUIRED"`              // 修改人名称
	Remark      string     `json:"remark"  gorm:"remark;title:备注" bson:"remark" title:"备注" parquet:"name=remark, type=BYTE_ARRAY,  repetitiontype=REQUIRED"`                                           // 备注
}
