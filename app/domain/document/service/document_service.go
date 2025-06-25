package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
)

type DocumentService struct {
	*dao.DocumentDao
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		DocumentDao: dao.NewDocumentDao(DBKey),
	}
}

func (s *DocumentService) HasDocumentByFolder(ctx context.Context, folderId string) bool {
	return s.CountByRSQL(ctx, "folder_id=='"+folderId+"'") > 0
}
