package command

import (
	"time"
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
	CaseId   string `json:"caseId"`
	FsKey    string `json:"fsKey"`    // 文件fs key
	FilePath string `json:"filePath"` // 文件目录
	FileId   string `json:"fileId"`   // 文件ID
	FileName string `json:"fileName"` // 内容
}

type DocumentData struct {
	Id       string `json:"id"`
	FsKey    string `json:"fsKey" `
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"` // 内容
	//FilePath   string     `json:"filePath"`                                                    // 文件目录
	State      int        `json:"state,omitempty"`                                             // 状态
	ChunkCount int        `json:"chunkCount,omitempty" gorm:"chunk_count"  bson:"chunk_count"` // 块数量
	DoneChunk  int        `json:"doneChunk,omitempty" gorm:"done_chunk"  bson:"done_chunk"`    // 完成数理
	StartTime  *time.Time `json:"startTime,omitempty" gorm:"start_time"  bson:"start_time"`    // 开始时间
	EndTime    *time.Time `json:"endTime,omitempty" gorm:"end_time"  bson:"end_time"`          // 结束时间
	Message    string     `json:"message,omitempty" gorm:"message"  bson:"message"`            // 消息
}

type DocumentScanData struct {
	TenantId string `json:"tenantId"  gorm:"tenant_id" bson:"tenant_id"`
	//CaseId   string `json:"caseId" gorm:"case_id" bson:"case_id"`
}
