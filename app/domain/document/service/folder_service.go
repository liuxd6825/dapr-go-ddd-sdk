package service

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

const DBKey string = "$db"

type FolderService struct {
	dao        *dao.FolderDao
	docService *DocumentService
}

func NewFolderService() *FolderService {
	return &FolderService{
		dao:        dao.NewFolderDao(DBKey),
		docService: NewDocumentService(),
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
	ft.Alias = folder.Alias
	ft.FolderPath = folder.FolderPath
	ft.Color = folder.Color
	ft.CreatedTime = folder.CreatedTime
	ft.CreatorId = folder.CreatorId
	ft.CreatorName = folder.CreatorName
	ft.UpdatedTime = folder.UpdatedTime
	ft.UpdaterId = folder.UpdaterId
	ft.UpdaterName = folder.UpdaterName
	/*	ft.DeletedTime = folder.DeletedTime
		ft.DeleterId = folder.DeleterId
		ft.DeleterName = folder.DeleterName
		ft.IsDeleted = folder.IsDeleted*/
	ft.Children = make([]model.FolderTree, 0)
	return ft
}

//func (t *FolderService) Create(ctx context.Context, cmd *command.FolderCreateCommand) error {
//	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
//		return t.dao.Create(ctx, &cmd.Data).GetError()
//	})
//}

func (t *FolderService) CreateData(ctx context.Context, data *model.Folder) error {
	return t.dao.Create(ctx, data).GetError()
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

func (t *FolderService) FindChildrenFolderCount(ctx context.Context, folderId string) (int64, error) {
	return t.dao.CountByRSQL(ctx, fmt.Sprintf("parent_id=='%s'", folderId))
}

func (s *FolderService) GetChildrenCount(ctx context.Context, folderId string) (int64, error) {
	fCount, err := s.FindChildrenFolderCount(ctx, folderId)
	if err != nil {
		return 0, err
	}
	cCount, err := s.docService.FindDocumentCountByFolderId(ctx, folderId)
	if err != nil {
		return 0, err
	}
	return fCount + cCount, nil
}

func (t *FolderService) Folder2FolderView(folder *model.Folder) *model.FolderView {
	view := &model.FolderView{}
	view.Id = folder.Id
	view.TenantId = folder.TenantId
	view.CaseId = folder.CaseId
	view.CreatedTime = folder.CreatedTime
	view.CreatorId = folder.CreatorId
	view.CreatorName = folder.CreatorName
	view.UpdatedTime = folder.UpdatedTime
	view.UpdaterId = folder.UpdaterId
	view.UpdaterName = folder.UpdaterName
	/*	view.DeletedTime = folder.DeletedTime
		view.DeleterId = folder.DeleterId
		view.DeleterName = folder.DeleterName
		view.IsDeleted = folder.IsDeleted*/
	view.Remark = folder.Remark
	view.BusId = folder.BusId
	view.EntityId = folder.EntityId
	view.RootId = folder.RootId
	view.RootPath = folder.RootPath
	view.FolderPath = folder.FolderPath
	view.ParentId = folder.ParentId
	view.Name = folder.Name
	view.Color = folder.Color
	view.Alias = folder.Alias
	view.DisabledFrontEdit = folder.DisabledFrontEdit
	return view
}

func (t *FolderService) FolderView2Folder(folderView *model.FolderView) *model.Folder {
	folder := &model.Folder{}
	folder.Id = folderView.Id
	folder.TenantId = folderView.TenantId
	folder.CaseId = folderView.CaseId
	folder.CreatedTime = folderView.CreatedTime
	folder.CreatorId = folderView.CreatorId
	folder.CreatorName = folderView.CreatorName
	folder.UpdatedTime = folderView.UpdatedTime
	folder.UpdaterId = folderView.UpdaterId
	folder.UpdaterName = folderView.UpdaterName
	/*	folder.DeletedTime = folderView.DeletedTime
		folder.DeleterId = folderView.DeleterId
		folder.DeleterName = folderView.DeleterName
		folder.IsDeleted = folderView.IsDeleted*/
	folder.Remark = folderView.Remark
	folder.BusId = folderView.BusId
	folder.EntityId = folderView.EntityId
	folder.RootId = folderView.RootId
	folder.RootPath = folderView.RootPath
	folder.FolderPath = folderView.FolderPath
	folder.ParentId = folderView.ParentId
	folder.Name = folderView.Name
	folder.Color = folderView.Color
	folder.Alias = folderView.Alias
	folder.DisabledFrontEdit = folderView.DisabledFrontEdit
	return folder
}
