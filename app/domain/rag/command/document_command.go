package command

import (
	"time"
)

type DocumentCreateCommand struct {
	CommandId string       `json:"commandId"`
	Data      DocumentData `json:"data"`
}

type DocumentUpdateCommand struct {
	CommandId string       `json:"commandId"`
	Data      DocumentData `json:"data"`
}

type DocumentScanCommand struct {
	CommandId string       `json:"commandId"`
	Data      DocumentData `json:"data"`
}

type DocumentData struct {
	Id         string     `json:"id"`
	FsKey      string     `json:"fsKey"  gorm:"fs_key" bson:"fs_key"`
	FileId     string     `json:"fileId" gorm:"file_id" bson:"file_id"`
	FileName   string     `json:"fileName,omitempty" gorm:"file_name;size:200" bson:"file_name"` // 内容
	State      int        `json:"state,omitempty" gorm:"state"  bson:"state"`                    // 状态
	ChunkCount int        `json:"chunkCount,omitempty" gorm:"chunk_count"  bson:"chunk_count"`   // 块数量
	DoneChunk  int        `json:"doneChunk,omitempty" gorm:"done_chunk"  bson:"done_chunk"`      // 完成数理
	StartTime  *time.Time `json:"startTime,omitempty" gorm:"start_time"  bson:"start_time"`      // 开始时间
	EndTime    *time.Time `json:"endTime,omitempty" gorm:"end_time"  bson:"end_time"`            // 结束时间
	Message    string     `json:"message,omitempty" gorm:"message"  bson:"message"`              // 消息
}

type DocumentScanData struct {
	TenantId string `json:"tenantId"  gorm:"tenant_id" bson:"tenant_id"`
	CaseId   string `json:"caseId" gorm:"case_id" bson:"case_id"`
}
