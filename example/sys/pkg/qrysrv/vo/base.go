package vo

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

// BaseVo 对应Java的BaseVo类，适配GORM
type BaseVo struct {
	Id           string      `gorm:"primaryKey;type:varchar(36)" json:"id" elastic:"type:keyword"` // GORM主键，ES keyword类型
	TenantId     string      `gorm:"type:varchar(36)" json:"tenant_id"`
	Version      int64       `gorm:"type:bigint" json:"version"`                                                       // 长整型
	CreatedBy    string      `gorm:"type:varchar(255)" json:"createdBy" elastic:"type:keyword"`                        // ES keyword类型
	Creator      string      `gorm:"type:varchar(255)" json:"creator" elastic:"type:keyword"`                          // ES keyword类型
	CreatedDate  *times.Time `gorm:"type:datetime" json:"createdDate" elastic:"type:date,format=yyyy-MM-dd HH:mm:ss"`  // 自定义日期格式
	ModifiedBy   string      `gorm:"type:varchar(255)" json:"modifiedBy"`                                              // 无特殊格式要求
	Modifier     string      `gorm:"type:varchar(255)" json:"modifier"`                                                // 无特殊格式要求
	ModifiedDate *times.Time `gorm:"type:datetime" json:"modifiedDate" elastic:"type:date,format=yyyy-MM-dd HH:mm:ss"` // 自定义日期格式
	DeletedDate  *times.Time `gorm:"type:datetime" json:"deletedDate" elastic:"type:date,format=yyyy-MM-dd HH:mm:ss"`  // 带时区声明
	DeletedAt    *times.Time `gorm:"index" json:"-"`                                                                   // GORM软删除字段，不暴露给JSON
}
