package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type FolderCreateCommand struct {
	xbase.Command[model.Folder]
}

type FolderUpdateCommand struct {
	xbase.Command[model.Folder]
}

type FolderDeleteCommand struct {
	xbase.Command[model.Folder]
}

type FolderRenameCommand struct {
	xbase.Command[*RenameFolder]
}

type FolderMoveCommand struct {
	xbase.Command[model.MoveFolder]
}

type RenameFolder struct {
	Id         string `json:"id" gorm:"primaryKey;title:主键" bson:"id" title:"主键" validate:"required" ` // 主键
	BusId      string `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`
	EntityId   string `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"`
	Name       string `bson:"name" json:"name,omitempty" bson:"name"`
	OldName    string `bson:"old_name" json:"oldName,omitempty" bson:"old_name"`
	FolderPath string `bson:"folder_path" json:"folderPath,omitempty" bson:"folder_path"`
}
