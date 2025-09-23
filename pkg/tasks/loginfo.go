package tasks

import "time"

type LogInfo struct {
	Id     string    `json:"id" gorm:"id" bson:"id"`
	TaskId string    `json:"task_id" gorm:"task_id" bson:"task_id"`
	AppId  string    `json:"app_id"  gorm:"app_id" bson:"app_id"`
	Level  string    `json:"level"  gorm:"level" bson:"level"`
	Time   time.Time `json:"time" gorm:"time" bson:"time"`
	Log    string    `json:"log" gorm:"log" bson:"log"`
}
