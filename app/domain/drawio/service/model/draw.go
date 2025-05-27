package model

import "time"

type Base struct {
	Id          string     `json:"id" gorm:"type:varchar(255);primary_key"`
	CaseId      string     `json:"caseId" gorm:"type:varchar(255);case_id"`
	TenantId    string     `json:"tenantId" gorm:"type:varchar(255);tenant_id"`
	CreatedTime *time.Time `json:"createdAt" gorm:"created_time"`
	CreatorId   string     `json:"creatorId" gorm:"creatorId"`
	CreatorName string     `json:"creatorName" gorm:"creator_name"`
	UpdatedTime *time.Time `json:"updateTime" gorm:"updated_time"`
	UpdaterId   string     `json:"updaterId" gorm:"updater_id"`
	UpdaterName string     `json:"updaterName" gorm:"updater_name"`
	DeletedTime *time.Time `json:"deletedTime" gorm:"deleted_time"`
	DeleterId   string     `json:"deleterId" gorm:"deleter_id"`
	DeleterName string     `json:"deleterName" gorm:"deleter_name"`
	IsDeleted   bool       `json:"isDeleted" gorm:"is_deleted"`
	Remark      string     `json:"remark" gorm:"remark"`
}

type Draw struct {
	Base
	Status   string `json:"status" gorm:"-"`
	Name     string `json:"name" gorm:"name"`
	FileName string `json:"fileName" json:"file_name"`
}
