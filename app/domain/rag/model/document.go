package model

import (
	"time"
)

type Document struct {
	Base       `bson:",inline"` // 基本
	FsKey      string           `json:"fsKey" gorm:"fs_key"  bson:"fs_key" title:"文件存储KEY" `                      // 文件存储KEY
	FileId     string           `json:"fileId"  gorm:"file_id" bson:"file_id" title:"文件ID"`                       // 文件ID
	FileName   string           `json:"fileName,omitempty" gorm:"file_name;size:200" bson:"file_name" title:"内容"` // 内容
	FilePath   string           `json:"filePath" gorm:"file_path"  bson:"file_path" title:"文件名称"`                 // 文件名称
	SourceId   string           `json:"sourceId,omitempty" gorm:"source_id"  bson:"source_id" title:"来源数据ID"`     // 来源数据ID
	SourceType string           `json:"sourceType,omitempty" gorm:"source_type"  bson:"source_type" title:"来源类型"` // 来源类型
	SourceApp  string           `json:"sourceApp,omitempty" gorm:"source_app"  bson:"source_app" title:"来源应用"`    // 来源应用
	SourceUrl  string           `json:"sourceUrl,omitempty" gorm:"source_url"  bson:"source_url" title:"来源的URL"`  // 来源的URL
	SourceName string           `json:"sourceName" gorm:"source_name" bson:"source_name" title:"来源名称"`
	State      DocumentState    `json:"state,omitempty" gorm:"state"  bson:"state" title:"状态"`                   // 状态
	ChunkCount int64            `json:"chunkCount,omitempty" gorm:"chunk_count"  bson:"chunk_count" title:"块数量"` // 块数量
	DoneChunk  int64            `json:"doneChunk,omitempty" gorm:"done_chunk"  bson:"done_chunk" title:"完成数理"`   // 完成数理
	StartTime  *time.Time       `json:"startTime,omitempty" gorm:"start_time"  bson:"start_time" title:"开始时间"`   // 开始时间
	EndTime    *time.Time       `json:"endTime,omitempty" gorm:"end_time"  bson:"end_time" title:"结束时间"`         // 结束时间
	Message    string           `json:"message,omitempty" gorm:"message"  bson:"message" title:"消息"`             // 消息

}

type DocumentState string

const (
	DocumentState_Pending            DocumentState = "排队中"
	DocumentState_Importing          DocumentState = "导入中"
	DocumentState_Succee             DocumentState = "已导入"
	DocumentState_Failure            DocumentState = "失败"
	DocumentState_FormatNotSupported DocumentState = "格式不支持"
)

func (s DocumentState) String() string {
	return string(s)
}

func (d *Document) GetId() string {
	return d.Id
}

func (d *Document) GetCaseId() string {
	return d.CaseId
}

func (d *Document) GetFileId() string {
	return d.FileId
}

func (d *Document) GetSourceType() string {
	return d.SourceType
}

func (d *Document) GetSourceApp() string {
	return d.SourceApp
}

func (d *Document) GetSourceId() string {
	return d.SourceId
}

func (d *Document) GetFilePath() string {
	return d.FilePath
}

func (d *Document) GetFileName() string {
	return d.FileName
}

func (d *Document) GetFsKey() string {
	return d.FsKey
}
