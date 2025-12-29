package interfaces

import "context"

type ImportStatusType string

const (
	ImportStatusType_Pending   ImportStatusType = "排队中"
	ImportStatusType_Importing ImportStatusType = "导入中"
	ImportStatusType_Succee    ImportStatusType = "完成"
	ImportStatusType_Failure   ImportStatusType = "失败"
)

type ImportDoc interface {
	GetFileId() string
	GetSourceType() string
	GetSourceApp() string
	GetSourceId() string
	GetFilePath() string
	GetFileName() string
}

type ImportStatusProvider interface {
	// UpdateStatus
	// @Description: 更新知识导入状态
	// @param FileID
	// @param status
	// @param statusMsg
	// @return error
	UpdateStatus(cxt context.Context, doc ImportDoc, status ImportStatusType, statusMsg string) error
}
