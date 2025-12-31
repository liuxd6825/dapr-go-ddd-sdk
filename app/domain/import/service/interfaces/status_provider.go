package interfaces

import "context"

type ImportStatusType string

const (
	ImportStatusType_Pending   ImportStatusType = "排队中"
	ImportStatusType_Importing ImportStatusType = "导入中"
	ImportStatusType_Succee    ImportStatusType = "完成"
	ImportStatusType_Failure   ImportStatusType = "失败"
)

func (s ImportStatusType) String() string {
	return string(s)
}

type ImportDoc interface {
	GetId() string         // 导入文档ID
	GetCaseId() string     // 案件ID
	GetSourceApp() string  // 来源系统名称
	GetSourceType() string // 来源数据类型
	GetSourceId() string   // 来源数据ID
	GetFilePath() string   // 文件目录
	GetFileId() string     // 文件id
	GetFileName() string   // 数据名称
	GetFsKey() string      // 存储KEY
}

type ImportStatusProvider interface {
	// UpdateStatus
	// @Description: 更新知识导入状态
	// @param fileID
	// @param status
	// @param statusMsg
	// @return error
	UpdateStatus(cxt context.Context, doc ImportDoc, status string, statusMsg string) error
}
