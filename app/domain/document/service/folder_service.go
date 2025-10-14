package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

const DBKey string = "$db"

type FolderService struct {
	dao *dao.FolderDao
}

func NewFolderService() *FolderService {
	return &FolderService{
		dao: dao.NewFolderDao(DBKey),
	}
}

func (s *FolderService) HasChildren(ctx context.Context, parentId string) (bool, error) {
	count, err := s.dao.CountByRSQL(ctx, "parent_id=='"+parentId+"'")
	return count > 0, err
}

func (s *FolderService) GetConfig() *idao.DaoConfig {
	return s.dao.GetConfig()
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

func (t *FolderService) Create(ctx context.Context, cmd *command.FolderCreateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.Create(ctx, &cmd.Data).GetError()
	})
}

func (t *FolderService) Delete(ctx context.Context, cmd *command.FolderDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *FolderService) DeleteById(ctx context.Context, id string) *idao.Result {
	return t.dao.DeleteById(ctx, id)
}

func (t *FolderService) Update(ctx context.Context, data *model.Folder, opts ...idao.CallOptions) error {
	return t.dao.Update(ctx, data, opts...).GetError()
}

func (t *FolderService) UpdateMany(ctx context.Context, entities []*model.Folder, opts ...idao.CallOptions) error {
	return t.dao.UpdateMany(ctx, entities, opts...).GetError()
}

func (t *FolderService) FindById(ctx context.Context, id string) (*model.Folder, error) {
	return t.dao.FindById(ctx, id)
}

func (t *FolderService) FindPaging(ctx context.Context, qry store.FindPagingQuery) (idao.FindPagingResult[*model.Folder], error) {
	res := t.dao.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (t *FolderService) FindByRSQL(ctx context.Context, rsql string) ([]*model.Folder, error) {
	return t.dao.FindByRSQL(ctx, rsql)
}
