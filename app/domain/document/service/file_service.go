package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
)

type FileService struct {
	*dao.FileDao
}

func NewFileService() *FileService {
	return &FileService{
		FileDao: dao.NewFileDao(DBKey),
	}
}

func (s *FileService) GetFile(document *model.Document) *model.File {
	file := &model.File{}
	file.Id = document.FileId
	file.TenantId = document.TenantId
	file.BusId = document.BusId
	file.EntityId = document.EntityId
	file.RootId = document.RootId
	file.FolderId = document.FolderId
	file.Name = document.Name
	file.ObjectName = document.ObjectName
	file.DocumentId = document.Id
	file.IsMain = true
	file.Size = document.Size
	file.SizeTitle = document.SizeTitle
	file.ExtName = document.ExtName
	file.DownloadUrl = document.DownloadUrl
	file.PreviewUrl = document.PreviewUrl
	file.DownloadTotal = document.DownloadTotal
	file.Thumbnail = document.Thumbnail
	file.Md5 = document.Md5
	file.FsKey = document.FsKey

	return file
}
