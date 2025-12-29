package event

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
)

// DocumentCreateEvent
// @Description: 新建知识文档事件
type DocumentCreateEvent = events.Event[*DocumentCreateEventData]

const DocumentCreateEventEventType = "rag.document.create-event"

type DocumentCreateEventData struct {
	Id         string `json:"id" ` // 文档ID
	CaseId     string `json:"caseId"`
	FsKey      string `json:"fsKey"`      // 文件存储KEY
	FileId     string `json:"fileId" `    // 文件ID
	FileName   string `json:"fileName" `  // 内容
	FilePath   string `json:"filePath"`   // 文件名称
	SourceId   string `json:"sourceId"`   // 来源数据ID
	SourceType string `json:"sourceType"` // 来源类型
	SourceApp  string `json:"sourceApp"`  // 来源应用
}

func NewDocumentCreateEvent(ctx context.Context, appId string, data *DocumentCreateEventData) *DocumentCreateEvent {
	event := &DocumentCreateEvent{}
	event.SetData(ctx, appId, data, &events.EventOptions{EventType: DocumentCreateEventEventType})
	return event
}
