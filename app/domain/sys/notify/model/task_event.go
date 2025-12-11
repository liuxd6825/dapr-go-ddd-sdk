package model

import "time"

// Message
// @Description: 信息通知
type Message struct {
	Id         string    `json:"id"`
	TenantId   string    `json:"tenantId"`
	UserId     string    `json:"userId"`
	CreateDate time.Time `json:"createDate"`
	Namespace  string    `json:"namespace"`
	Type       string    `json:"type"`
	Message    string    `json:"message"`
}

// Status
// @Description: 任务状态通知
type Status struct {
	Id         string    `json:"id"`
	TenantId   string    `json:"tenantId"`
	UserId     string    `json:"userId"`
	CreateDate time.Time `json:"createDate"`
	TaskType   string    `json:"taskType"`
	TaskId     string    `json:"taskId"`
	Status     string    `json:"status"` // e.g., "COMPLETED", "FAILED"
	Message    string    `json:"message"`
	ResultData string    `json:"resultData,omitempty"`
}
