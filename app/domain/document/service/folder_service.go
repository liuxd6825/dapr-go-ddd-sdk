package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
)

const DBKey string = "$db"

type FolderService struct {
	*dao.FolderDao
}

func NewFolderService() *FolderService {
	return &FolderService{
		FolderDao: dao.NewFolderDao(DBKey),
	}
}

func (s *FolderService) HasChildren(ctx context.Context, parentId string) bool {
	return s.CountByRSQL(ctx, "parent_id=='"+parentId+"'") > 0
}
