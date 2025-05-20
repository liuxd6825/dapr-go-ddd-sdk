package model

import "time"

type Base struct {
	UpdatedTime *time.Time `json:"updatedTime" gorm:"column:updated_time"`
	UpdaterName string     `json:"updaterName" gorm:"column:updater_name"`
	UpdaterId   string     `json:"updaterId" gorm:"column:updater_id"`
	CreatedTime *time.Time `json:"createdTime" gorm:"column:created_time"`
	CreatorName string     `json:"creatorName" gorm:"column:creator_name"`
	CreatorId   string     `json:"creatorId" gorm:"column:creator_id"`
}
