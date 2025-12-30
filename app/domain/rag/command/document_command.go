package command

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
)

type DocumentCreateCommand struct {
	CommandId string             `json:"commandId"`
	Data      DocumentCreateData `json:"data"`
}

type DocumentUpdateCommand struct {
	CommandId string       `json:"commandId"`
	Data      DocumentData `json:"data"`
}

type DocumentDeleteCommand struct {
	CommandId string             `json:"commandId"`
	Data      DocumentDeleteData `json:"data"`
}

type DocumentScanCommand struct {
	CommandId string           `json:"commandId"`
	Data      DocumentScanData `json:"data"`
}

type DocumentDeleteData struct {
	CaseId string `json:"caseId"`
	Id     string `json:"id"`
}

type DocumentCreateData struct {
	Id         string `json:"id"`
	CaseId     string `json:"caseId"`
	FsKey      string `json:"fsKey"`    // 文件fs key
	FilePath   string `json:"filePath"` // 文件目录
	FileId     string `json:"fileId"`   // 文件ID
	FileName   string `json:"fileName"` // 内容
	SourceId   string `json:"sourceId"`
	SourceType string `json:"sourceType"`
	SourceApp  string `json:"sourceApp"`
}

type DocumentData struct {
	Id       string `json:"id"`
	FsKey    string `json:"fsKey" `
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"` // 内容
	//FilePath   string     `json:"filePath"`            // 文件目录
	State      model.DocumentState `json:"state,omitempty"`       // 状态
	ChunkCount int64               `json:"chunkCount,omitempty" ` // 块数量
	DoneChunk  int64               `json:"doneChunk,omitempty" `  // 完成数理
	StartTime  *time.Time          `json:"startTime,omitempty" `  // 开始时间
	EndTime    *time.Time          `json:"endTime,omitempty" `    // 结束时间
	Message    string              `json:"message,omitempty" `    // 消息
}

type DocumentScanData struct {
	TenantId string `json:"tenantId"  gorm:"tenant_id" bson:"tenant_id"`
	//CaseId   string `json:"caseId" gorm:"case_id" bson:"case_id"`
}
