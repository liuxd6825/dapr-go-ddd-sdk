package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"os"
)

type FileService struct {
	*dao.FileDao
	folderService *FolderService
	fsService     *FsService
}

func NewFileService() *FileService {
	return &FileService{
		FileDao:       dao.NewFileDao(DBKey),
		folderService: NewFolderService(),
		fsService:     NewFsService(),
	}
}

func (s *FileService) GetFile(document *model.Document) *model.File {
	file := &model.File{}
	file.Id = document.FileId
	file.TenantId = document.TenantId
	file.BusId = document.BusId
	file.EntityId = document.EntityId
	file.CaseId = document.CaseId
	file.RootId = document.RootId
	file.RootPath = document.RootPath
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

func (s *FsService) ReadFile(fullFileName string) ([]byte, error) {
	readFile, err := s.docFs.Open(fullFileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = readFile.Close()
	}()

	fileInfo, err := readFile.Stat()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, fileInfo.Size())
	_, err = readFile.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (s *FileService) ReadFile(ctx context.Context, fileId string) ([]byte, error) {
	file, err := s.Dao.FindById(ctx, fileId)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("没有找到文件id:%s", fileId)
	}

	folder, err := s.folderService.FindById(ctx, file.FolderId)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, errors.New("没有找到目录id:%s", file.FolderId)
	}

	fullName := folder.FolderPath + "\"" + file.ObjectName
	bytes, err := s.fsService.ReadFile(fullName)
	return bytes, err
}
