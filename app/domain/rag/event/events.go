package event

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
)

// CreateDocumentEvent
// @Description: 新建知识文档事件
type CreateDocumentEvent = events.Event[*model.Document]

type CreateDocumentEventData struct {
	SourceId     string `json:"sourceId"`
	SourceType   string `json:"sourceType"`
	SourceApp    string `json:"sourceApp"`
	FsKey        string `json:"fsKey"`
	FileId       string `json:"fileId"`
	FileName     string `json:"fileName"`
	FullFileName string `json:"fullFileName" `
}
