package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
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

func (s *FolderService) HasChildren(ctx context.Context, parentId string) (bool, error) {
	count, err := s.CountByRSQL(ctx, "parent_id=='"+parentId+"'")
	return count > 0, err
}

func (s *FolderService) TranslateTreeData(data []*model.Folder) []model.FolderTree {
	if len(data) == 0 {
		return []model.FolderTree{}
	}
	arr := make([]model.FolderTree, 0)
	for _, row := range data {
		if row.ParentId == "" {
			root := s.FolderToFolderTree(row)
			s.TranslateTreeChildData(data, &root)
			arr = append(arr, root)
		}
	}
	return arr
}

func (s *FolderService) TranslateTreeChildData(data []*model.Folder, parent *model.FolderTree) {
	for _, row := range data {
		if row.ParentId == parent.Id {
			node := s.FolderToFolderTree(row)
			s.TranslateTreeChildData(data, &node)
			parent.Children = append(parent.Children, node)
		}
	}
}

func (s *FolderService) FolderToFolderTree(folder *model.Folder) model.FolderTree {
	ft := model.FolderTree{}
	ft.Id = folder.Id
	ft.TenantId = folder.TenantId
	ft.BusId = folder.BusId
	ft.EntityId = folder.EntityId
	ft.CaseId = folder.CaseId
	ft.RootId = folder.RootId
	ft.ParentId = folder.ParentId
	ft.Name = folder.Name
	ft.FolderPath = folder.FolderPath
	ft.Color = folder.Color
	ft.CreatedTime = folder.CreatedTime
	ft.CreatorId = folder.CreatorId
	ft.CreatorName = folder.CreatorName
	ft.UpdatedTime = folder.UpdatedTime
	ft.UpdaterId = folder.UpdaterId
	ft.UpdaterName = folder.UpdaterName
	ft.DeletedTime = folder.DeletedTime
	ft.DeleterId = folder.DeleterId
	ft.DeleterName = folder.DeleterName
	ft.IsDeleted = folder.IsDeleted
	ft.Children = make([]model.FolderTree, 0)
	return ft
}
