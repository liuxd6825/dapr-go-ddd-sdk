package model

type Document struct {
	Base          `bson:",inline"`
	RootId        string `gorm:"root_id" json:"rootId,omitempty" bson:"root_id"`                      //根目录Id    tenant_id_bus_id_entity_id形式拼接
	RootPath      string `gorm:"root_path" json:"rootPath,omitempty" bson:"root_path"`                //物理根目录
	FolderId      string `gorm:"folder_id" json:"folderId,omitempty" bson:"folder_id"`                //目录Id
	FolderPath    string `gorm:"folder_path" json:"folderPath,omitempty" bson:"folder_path"`          //物理目录
	FileId        string `gorm:"file_id" json:"fileId,omitempty" bson:"file_id"`                      //主文件Id
	Name          string `gorm:"name" json:"name,omitempty" bson:"name"`                              //文档名称
	ObjectName    string `gorm:"object_name" json:"objectName,omitempty" bson:"object_name"`          //实际存储的文件名称
	ExtName       string `gorm:"ext_name" json:"extName,omitempty" bson:"ext_name"`                   //扩展名
	Size          int64  `gorm:"size,type:bigint" json:"size,omitempty" bson:"size"`                  //文件大小
	SizeTitle     string `gorm:"size_title" json:"sizeTitle,omitempty" bson:"size_title"`             //文件大小
	DownloadTotal int64  `gorm:"download_total" json:"downloadTotal,omitempty" bson:"download_total"` //下载次数
	DownloadUrl   string `gorm:"download_url" json:"downloadUrl,omitempty" bson:"download_url"`       //下载地址
	PreviewUrl    string `gorm:"preview_url" json:"previewUrl,omitempty" bson:"preview_url"`          //预览地址
	Thumbnail     string `gorm:"thumbnail;size:8000" json:"thumbnail,omitempty" bson:"thumbnail"`     //缩略图
	Md5           string `gorm:"md5" json:"md5,omitempty" bson:"md5"`                                 //md5串
	TagName       string `gorm:"tag_name" json:"tagName,omitempty" bson:"tag_name"`                   //标签名  逗号分隔
	TagId         string `gorm:"tag_id" json:"tagId,omitempty" bson:"tag_id"`                         //标签Id  逗号分隔
	TagColor      string `gorm:"tag_color" json:"tagColor,omitempty" bson:"tag_color"`                //标签颜色  逗号分隔
	FsKey         string `gorm:"fs_key" json:"fsKey,omitempty" bson:"fs_key"`
}

type RenameDocument struct {
	Base       `bson:",inline"`
	FileId     string `gorm:"file_id" json:"fileId,omitempty" bson:"file_id"`
	Name       string `gorm:"name" json:"name,omitempty" bson:"name"`
	ObjectName string `gorm:"object_name" json:"objectName,omitempty" bson:"object_name"`
	OldName    string `gorm:"old_name" json:"oldName,omitempty" bson:"old_name"`
	FolderPath string `gorm:"folder_path" json:"folderPath,omitempty" bson:"folder_path"`
}
