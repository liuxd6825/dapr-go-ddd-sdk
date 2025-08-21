package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type File struct {
	xbase.BaseModel `bson:",inline"`
	BusId           string `gorm:"bus_id" json:"busId,omitempty" bson:"bus_id"`
	EntityId        string `gorm:"entity_id" json:"entityId,omitempty" bson:"entity_id"`
	RootId          string `gorm:"root_id" json:"rootId,omitempty" bson:"root_id"`                      //根目录Id    tenant_id_bus_id_entity_id形式拼接
	RootPath        string `gorm:"root_path" json:"rootPath,omitempty" bson:"root_path"`                //物理根目录
	FolderId        string `gorm:"folder_id" json:"folderId,omitempty" bson:"folder_id"`                //目录Id
	DocumentId      string `gorm:"document_id" json:"documentId,omitempty" bson:"document_id"`          //文档Id
	Name            string `gorm:"name" json:"name,omitempty" bson:"name"`                              //文档名称
	ObjectName      string `gorm:"object_name" json:"objectName,omitempty" bson:"object_name"`          //实际存储的文件名称
	ExtName         string `gorm:"ext_name" json:"extName,omitempty" bson:"ext_name"`                   //扩展名
	Size            int64  `gorm:"size,type:bigint" json:"size,omitempty" bson:"size"`                  //文件大小
	SizeTitle       string `gorm:"size_title" json:"sizeTitle,omitempty" bson:"size_title"`             //文件大小
	DownloadTotal   int64  `gorm:"download_total" json:"downloadTotal,omitempty" bson:"download_total"` //下载次数
	DownloadUrl     string `gorm:"download_url" json:"downloadUrl,omitempty" bson:"download_url"`       //下载地址
	PreviewUrl      string `gorm:"preview_url" json:"previewUrl,omitempty" bson:"preview_url"`          //预览地址
	Thumbnail       string `gorm:"thumbnail;size:8000" json:"thumbnail,omitempty" bson:"thumbnail"`     //缩略图
	Md5             string `gorm:"md5" json:"md5,omitempty" bson:"md5"`                                 //md5串
	IsMain          bool   `gorm:"is_main" json:"isMain,omitempty" bson:"is_main"`                      //是否主文件
	VerInfo         string `gorm:"ver_info" json:"verInfo,omitempty" bson:"ver_info"`                   //版本信息
	FsKey           string `gorm:"fs_key" json:"fsKey,omitempty" bson:"fs_key"`
}
