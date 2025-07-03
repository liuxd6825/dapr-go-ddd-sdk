package model

import "time"

type Base struct {
	Id          string     `json:"id" gorm:"type:varchar(255);primary_key" bson:"id"`
	CaseId      string     `json:"caseId" gorm:"type:varchar(255);case_id" bson:"case_id"`
	TenantId    string     `json:"tenantId" gorm:"type:varchar(255);tenant_id" bson:"tenant_id"`
	CreatedTime *time.Time `json:"createdAt" gorm:"created_time" bson:"created_time"`
	CreatorId   string     `json:"creatorId" gorm:"creatorId" bson:"creator_id"`
	CreatorName string     `json:"creatorName" gorm:"creator_name" bson:"creator_name"`
	UpdatedTime *time.Time `json:"updateTime" gorm:"updated_time" bson:"update_time"`
	UpdaterId   string     `json:"updaterId" gorm:"updater_id" bson:"updater_id"`
	UpdaterName string     `json:"updaterName" gorm:"updater_name" bson:"updater_name"`
	DeletedTime *time.Time `json:"deletedTime" gorm:"deleted_time" bson:"deleted_time"`
	DeleterId   string     `json:"deleterId" gorm:"deleter_id" bson:"deleter_id"`
	DeleterName string     `json:"deleterName" gorm:"deleter_name" bson:"deleter_name"`
	IsDeleted   bool       `json:"isDeleted" gorm:"is_deleted" bson:"is_deleted"`
	Remark      string     `json:"remark" gorm:"remark" bson:"remark"`
}

type Draw struct {
	Base     `bson:",inline"`
	Status   string `json:"status" gorm:"status"`
	Name     string `json:"name" gorm:"name"`
	FileName string `json:"fileName" gorm:"file_name"`
}
