package event

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
)

type FolderCreateEvent = events.Event[*FolderCreateEventData]

type FolderCreateEventData struct {
	BusId      string `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`
	EntityId   string `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"`
	RootId     string `gorm:"root_id" json:"rootId,omitempty" bson:"root_id"`             //根目录Id    tenant_id_bus_id_entity_id形式拼接
	RootPath   string `gorm:"root_path" json:"rootPath,omitempty" bson:"root_path"`       //物理根目录
	FolderPath string `gorm:"folder_path" json:"folderPath,omitempty" bson:"folder_path"` //物理目录
	ParentId   string `gorm:"parent_id" json:"parentId,omitempty" bson:"parent_id"`       //父目录Id
	Name       string `gorm:"name" json:"name,omitempty" bson:"name"`                     //目录名称
	Color      string `gorm:"color" json:"color,omitempty" bson:"color"`                  //目录颜色
	Alias      string `gorm:"alias" json:"alias,omitempty" bson:"alias"`                  // 目录别名
}

const FolderCreateEventType = "document.folder.create"

// NewFolderCreateEvent
// @Description: 新建目录
// @param ctx
// @param appId
// @param data
// @return *RecordImportMasterEvent
func NewFolderCreateEvent(ctx context.Context, appId string, data *FolderCreateEventData) *FolderCreateEvent {
	event := &FolderCreateEvent{}
	event.SetData(ctx, appId, data, &events.EventOptions{EventType: FolderCreateEventType})
	return event
}
