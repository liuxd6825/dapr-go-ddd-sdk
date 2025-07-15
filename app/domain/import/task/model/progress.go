package model

import "time"

type Progress struct {
	Base    `bson:",inline"`
	TaskId  string    `json:"taskId" bson:"task_id" index:""`
	Type    string    `json:"type" bson:"type" index:""`
	Time    time.Time `json:"time" bson:"time" index:""`
	Message string    `json:"message" bson:"message"`
}
