package model

type Folder struct {
	Base       `bson:",inline"`
	RootId     string `gorm:"root_id" json:"rootId,omitempty" bson:"root_id"`             //根目录Id    tenant_id_bus_id_entity_id形式拼接
	RootPath   string `gorm:"root_path" json:"rootPath,omitempty" bson:"root_path"`       //物理根目录
	FolderPath string `gorm:"folder_path" json:"folderPath,omitempty" bson:"folder_path"` //物理目录
	ParentId   string `gorm:"parent_id" json:"parentId,omitempty" bson:"parent_id"`       //父目录Id
	Name       string `gorm:"name" json:"name,omitempty" bson:"name"`                     //目录名称
	Color      string `gorm:"color" json:"color,omitempty" bson:"color"`                  //目录颜色
}

type RenameFolder struct {
	Base       `bson:",inline"`
	Name       string `bson:"name" json:"name,omitempty" bson:"name"`
	OldName    string `bson:"old_name" json:"oldName,omitempty" bson:"old_name"`
	FolderPath string `bson:"folder_path" json:"folderPath,omitempty" bson:"folder_path"`
}

type MoveFolder struct {
	Base       `bson:",inline"`
	Name       string `bson:"name" json:"name,omitempty" bson:"name"`
	ParentId   string `bson:"parent_id" json:"parentId,omitempty" bson:"parent_id"`
	SourcePath string `gorm:"source_path" json:"sourcePath,omitempty" bson:"source_path"`
	TargetPath string `gorm:"target_path" json:"targetPath,omitempty" bson:"target_path"`
}
