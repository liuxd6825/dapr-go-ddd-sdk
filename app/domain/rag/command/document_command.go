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

type DocumentDeleteByDocCommand struct {
	CommandId string                  `json:"commandId"`
	Data      DocumentDeleteByDocData `json:"data"`
}

type DocumentScanCommand struct {
	CommandId string           `json:"commandId"`
	Data      DocumentScanData `json:"data"`
}

type DocumentDeleteData struct {
	CaseId string `json:"caseId"`
	Id     string `json:"id"`
}

type DocumentDeleteByDocData struct {
	DocId string `json:"docId"`
}

type DocumentCreateData struct {
	Id         string `json:"id" title:"主键" `
	CaseId     string `json:"caseId" title:"案件ID"`
	FsKey      string `json:"fsKey" title:"文件存储KEY"`
	FilePath   string `json:"filePath" title:"文件目录"`
	FileId     string `json:"fileId" title:"文件ID"`
	FileName   string `json:"fileName" title:"文件名称"`
	SourceId   string `json:"sourceId" title:"来源ID"`
	SourceType string `json:"sourceType" title:"来源类型"`
	SourceApp  string `json:"sourceApp" title:"来源应用"`
	SourceUrl  string `json:"sourceUrl" title:"来源链接"`
	SourceName string `json:"sourceName" title:"来源链接"`
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
