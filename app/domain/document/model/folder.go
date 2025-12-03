package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// Folder
// @Description:  目录
type Folder struct {
	xbase.BaseModel `bson:",inline"`
	BusId           string `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`
	EntityId        string `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"`
	RootId          string `gorm:"root_id" json:"rootId,omitempty" bson:"root_id"`             //根目录Id    tenant_id_bus_id_entity_id形式拼接
	RootPath        string `gorm:"root_path" json:"rootPath,omitempty" bson:"root_path"`       //物理根目录
	FolderPath      string `gorm:"folder_path" json:"folderPath,omitempty" bson:"folder_path"` //物理目录
	ParentId        string `gorm:"parent_id" json:"parentId,omitempty" bson:"parent_id"`       //父目录Id
	Name            string `gorm:"name" json:"name,omitempty" bson:"name"`                     //目录名称
	Color           string `gorm:"color" json:"color,omitempty" bson:"color"`                  //目录颜色

	Alias string `gorm:"alias" json:"alias,omitempty" bson:"alias"` // 目录别名

	DisabledFrontEdit bool `json:"disabledFrontEdit" bson:"disabled_front_edit" gorm:"disabled_front_edit"`
}

func NewFolder() (*Folder, error) {
	return &Folder{}, nil
}

// FolderMeta
// @Description: 目录元数据
type FolderMeta struct {
	xbase.BaseModel `bson:",inline"`
	FolderId        string `gorm:"folder_id" json:"folderId,omitempty" bson:"folder_id" title:"来源"`
	SourceType      string `gorm:"source_type" json:"sourceType,omitempty" bson:"source_type" title:"来源类型"`
	Source          string `gorm:"source" json:"source,omitempty" bson:"source" title:"来源"`
	Name            string `gorm:"name" json:"name,omitempty" bson:"name" title:"名称"`
	Value           string `gorm:"value" json:"value,omitempty" bson:"value" title:"值"`
}

func NewFolderMeta() *FolderMeta {
	return &FolderMeta{}
}

type FolderView struct {
	xbase.BaseModel `bson:",inline"`
	BusId           string `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`
	EntityId        string `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"`
	RootId          string `gorm:"root_id" json:"rootId,omitempty" bson:"root_id"`             //根目录Id    tenant_id_bus_id_entity_id形式拼接
	RootPath        string `gorm:"root_path" json:"rootPath,omitempty" bson:"root_path"`       //物理根目录
	FolderPath      string `gorm:"folder_path" json:"folderPath,omitempty" bson:"folder_path"` //物理目录
	ParentId        string `gorm:"parent_id" json:"parentId,omitempty" bson:"parent_id"`       //父目录Id
	Name            string `gorm:"name" json:"name,omitempty" bson:"name"`                     //目录名称
	Color           string `gorm:"color" json:"color,omitempty" bson:"color"`                  //目录颜色
	Alias           string `gorm:"alias" json:"alias,omitempty" bson:"alias"`                  // 目录别名

	Meta []*FolderMeta `json:"meta"`

	DisabledFrontEdit bool `json:"disabledFrontEdit" bson:"disabled_front_edit" gorm:"disabled_front_edit"`
}

type FolderTree struct {
	xbase.BaseModel `bson:",inline"`
	BusId           string       `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`
	EntityId        string       `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"`
	RootId          string       `gorm:"root_id" json:"rootId,omitempty" bson:"root_id"`             //根目录Id    tenant_id_bus_id_entity_id形式拼接
	RootPath        string       `gorm:"root_path" json:"rootPath,omitempty" bson:"root_path"`       //物理根目录
	FolderPath      string       `gorm:"folder_path" json:"folderPath,omitempty" bson:"folder_path"` //物理目录
	ParentId        string       `gorm:"parent_id" json:"parentId,omitempty" bson:"parent_id"`       //父目录Id
	Name            string       `gorm:"name" json:"name,omitempty" bson:"name"`                     //目录名称
	Color           string       `gorm:"color" json:"color,omitempty" bson:"color"`                  //目录颜色
	Children        []FolderTree `gorm:"children" json:"children,omitempty" bson:"children"`

	Alias string `gorm:"alias" json:"alias,omitempty" bson:"alias"` // 目录别名

	//Meta []*FolderMeta `json:"meta"`
}

type RenameFolder struct {
	Id         string `json:"id" gorm:"primaryKey;title:主键" bson:"id" title:"主键" validate:"required" ` // 主键
	BusId      string `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`
	EntityId   string `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"`
	Name       string `bson:"name" json:"name,omitempty" bson:"name"`
	OldName    string `bson:"old_name" json:"oldName,omitempty" bson:"old_name"`
	FolderPath string `bson:"folder_path" json:"folderPath,omitempty" bson:"folder_path"`

	Alias string `gorm:"alias" json:"alias,omitempty" bson:"alias"` // 目录别名
}

type MoveFolder struct {
	xbase.BaseModel `bson:",inline"`
	BusId           string `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`
	EntityId        string `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"`
	Name            string `bson:"name" json:"name,omitempty" bson:"name"`
	ParentId        string `bson:"parent_id" json:"parentId,omitempty" bson:"parent_id"`
	SourcePath      string `gorm:"source_path" json:"sourcePath,omitempty" bson:"source_path"`
	TargetPath      string `gorm:"target_path" json:"targetPath,omitempty" bson:"target_path"`

	Alias string `gorm:"alias" json:"alias,omitempty" bson:"alias"` // 目录别名
}
