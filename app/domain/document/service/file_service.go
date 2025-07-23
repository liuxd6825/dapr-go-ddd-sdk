package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
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

func (s *FileService) ReadByteByFile(ctx context.Context, file *model.File) ([]byte, error) {
	if file == nil {
		return nil, errors.New("file参数不能为nil")
	}

	folder, err := s.folderService.FindById(ctx, file.FolderId)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, errors.New("没有找到目录id:%s", file.FolderId)
	}

	bytes, err := s.fsService.ReadFile(folder.FolderPath, file.ObjectName)
	return bytes, err
}

// ReadByteByFileId 根据文件ID从存储中读取文件的字节内容。
//
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期和取消操作。
//   - fileId: 要读取的文件的唯一标识符。
//
// 返回值:
//   - []byte: 文件内容的字节切片。
//   - error: 如果在查找文件或读取过程中发生错误，则返回相应的错误信息。
func (s *FileService) ReadByteByFileId(ctx context.Context, fileId string) ([]byte, error) {
	// 通过文件ID调用DAO层获取文件对象
	file, err := s.Dao.FindById(ctx, fileId)
	if err != nil {
		return nil, err
	}
	// 检查文件是否存在
	if file == nil {
		return nil, errors.New("文档中心不存在id为%s的文件", fileId)
	}

	// 调用ReadByteByFile方法实际读取文件内容
	return s.ReadByteByFile(ctx, file)
}
