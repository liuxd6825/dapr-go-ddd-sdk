package model

import (
	"time"
)

type Document struct {
	Base       `bson:",inline"` // 基本
	FsKey      string           `json:"fsKey" gorm:"fs_key"   bson:"fs_key" `                          // 文件存储KEY
	FileId     string           `json:"fileId"  gorm:"file_id" bson:"file_id"`                         // 文件ID
	FileName   string           `json:"fileName,omitempty" gorm:"file_name;size:200" bson:"file_name"` // 内容
	FilePath   string           `json:"filePath" gorm:"file_path"  bson:"file_path"`                   // 文件名称
	SourceId   string           `json:"sourceId,omitempty" gorm:"source_id"  bson:"source_id"`         // 来源数据ID
	SourceType string           `json:"sourceType,omitempty" gorm:"source_type"  bson:"source_type"`   // 来源类型
	SourceApp  string           `json:"sourceApp,omitempty" gorm:"source_app"  bson:"source_app"`      // 来源应用
	State      int              `json:"state,omitempty" gorm:"state"  bson:"state"`                    // 状态
	ChunkCount int              `json:"chunkCount,omitempty" gorm:"chunk_count"  bson:"chunk_count"`   // 块数量
	DoneChunk  int              `json:"doneChunk,omitempty" gorm:"done_chunk"  bson:"done_chunk"`      // 完成数理
	StartTime  *time.Time       `json:"startTime,omitempty" gorm:"start_time"  bson:"start_time"`      // 开始时间
	EndTime    *time.Time       `json:"endTime,omitempty" gorm:"end_time"  bson:"end_time"`            // 结束时间
	Message    string           `json:"message,omitempty" gorm:"message"  bson:"message"`              // 消息

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
