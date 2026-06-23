package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/query"
	folder_event "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/event"
	code_service "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

const folderAppId = "duxm-master-cmd-service"

// ProductService 产品业务服务
type ProductService struct {
	dao              *dao.ProductDao
	codeSvc          *code_service.CodeService
	feign            any
	productProductDao *dao.ProductProductDao
	contractDao      *dao.ProductContractDao
	companyDao       *dao.ProductCompanyDao
	humanDao         *dao.ProductHumanDao
	recordDao        *dao.ProductRecordDao
}

var _productService *ProductService
var _productServiceOnce sync.Once

// NewProductService 单例构造函数
func NewProductService() *ProductService {
	_productServiceOnce.Do(func() {
		_productService = &ProductService{
			dao:               dao.NewProductDao(config.DBKey),
			codeSvc:           code_service.NewCodeService(),
			feign:             nil,
			productProductDao: dao.NewProductProductDao(config.DBKey),
			contractDao:       dao.NewProductContractDao(config.DBKey),
			companyDao:        dao.NewProductCompanyDao(config.DBKey),
			humanDao:          dao.NewProductHumanDao(config.DBKey),
			recordDao:         dao.NewProductRecordDao(config.DBKey),
		}
	})
	return _productService
}

// Create 创建产品
func (s *ProductService) Create(ctx context.Context, cmd *command.ProductCreateCommand, opts ...idao.CallOptions) error {
	product := cmd.Data
	if product.Code == "" {
		product.Code = s.codeSvc.NewProductCode(ctx, product.CaseId)
	}
	if err := s.dao.Create(ctx, &product, opts...).GetError(); err != nil {
		return err
	}
	s.publishFolderCreateEvent(ctx, &product)
	return nil
}

// Update 更新产品
func (s *ProductService) Update(ctx context.Context, cmd *command.ProductUpdateCommand, opts ...idao.CallOptions) error {
	oldProduct, err := s.dao.FindById(ctx, cmd.Data.Id)
	if err != nil {
		return err
	}
	if err := s.dao.Update(ctx, &cmd.Data, opts...).GetError(); err != nil {
		return err
	}
	if oldProduct != nil && oldProduct.Name != cmd.Data.Name {
		s.publishFolderUpdateAliasEvent(ctx, &cmd.Data)
	}
	return nil
}

// Submit 创建或更新产品
func (s *ProductService) Submit(ctx context.Context, cmd *command.ProductSubmitCommand, opts ...idao.CallOptions) error {
	old, err := s.dao.FindById(ctx, cmd.Data.Id, opts...)
	if err != nil {
		return err
	}
	if old != nil {
		return s.Update(ctx, &command.ProductUpdateCommand{CommandId: cmd.CommandId, Data: cmd.Data}, opts...)
	}
	return s.Create(ctx, &command.ProductCreateCommand{CommandId: cmd.CommandId, Data: cmd.Data}, opts...)
}

// DeleteById 按 ID 删除产品
func (s *ProductService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

// DeleteByIds 批量删除产品，并级联清理 5 个子实体
func (s *ProductService) DeleteByIds(ctx context.Context, cmd *command.ProductDeleteByIdsCommand, opts ...idao.CallOptions) error {
	if err := s.dao.DeleteByIds(ctx, cmd.Data.Ids, opts...).GetError(); err != nil {
		return err
	}

	rsqlStr := rsql.NewBuilder().In("productId", cmd.Data.Ids).Build()
	logs.Info(ctx, logs.Fields{"rsql": rsqlStr})

	// 5 张子表级联删除
	if err := s.productProductDao.DeleteByProductIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.contractDao.DeleteByProductIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.companyDao.DeleteByProductIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.humanDao.DeleteByProductIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.recordDao.DeleteByProductIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}

	s.publishFolderDeleteEvent(ctx, cmd.Data.Ids, cmd.Data.CaseId)
	return nil
}

// FindById 按 ID 查询产品
func (s *ProductService) FindById(ctx context.Context, qry *query.ProductFindByIdQuery, opts ...idao.CallOptions) (*model.Product, error) {
	return s.dao.FindById(ctx, qry.Id, opts...)
}

// FindPaging 分页查询产品
func (s *ProductService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Product] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

// FindPagingByCaseId 按案件 ID 分页查询产品
func (s *ProductService) FindPagingByCaseId(ctx context.Context, qry *query.ProductFindByCaseIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Product] {
	return s.dao.FindPagingByCaseId(ctx, qry, qry.CaseId, opts...)
}

// FindPagingByTagId 按标签 ID 查询产品
func (s *ProductService) FindPagingByTagId(ctx context.Context, qry *query.ProductFindByTagIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Product] {
	paging := store.NewFindPagingQuery()
	paging.SetPageSize(store.PagingMaxPageSize)
	return s.dao.FindPagingByTagId(ctx, paging, qry.CaseId, qry.TagId, opts...)
}

// FindTotal 统计案件下的产品数量
func (s *ProductService) FindTotal(ctx context.Context, qry *query.ProductFindTotalQuery, opts ...idao.CallOptions) (*query.ProductFindTotalResult, error) {
	count, err := s.dao.CountByCaseId(ctx, qry.CaseId, opts...)
	if err != nil {
		return nil, err
	}
	return &query.ProductFindTotalResult{
		Count: count,
		Icon:  qry.Icon,
	}, nil
}

// publishFolderCreateEvent 发布创建目录事件
func (s *ProductService) publishFolderCreateEvent(ctx context.Context, product *model.Product) {
	tenantId := appctx.GetTenantId2(ctx)
	data := &folder_event.FolderCreateEventData{
		Id:                "D" + product.Id,
		CaseId:            product.CaseId,
		BusId:             "case",
		EntityId:          product.CaseId,
		RootId:            fmt.Sprintf("%s_case_%s", tenantId, product.CaseId),
		RootPath:          fmt.Sprintf("/%s/case/%s", tenantId, product.CaseId),
		FolderPath:        fmt.Sprintf("/%s/case/%s/主数据附件/产品/%s", tenantId, product.CaseId, product.Code),
		ParentId:          fmt.Sprintf("%s_case_%s_master_product", tenantId, product.CaseId),
		Name:              product.Code,
		Alias:             product.Name,
		DisabledFrontEdit: true,
		Meta: []*folder_event.FolderCreateEventMeta{
			{
				SourceType: "产品",
				Source:     product.Id,
				Name:       "SourceId",
				Value:      product.Id,
			},
		},
	}
	evt := folder_event.NewFolderCreateEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderCreateEvent error: %s", err.Error())
	}
}

// publishFolderUpdateAliasEvent 发布更新目录别名事件
func (s *ProductService) publishFolderUpdateAliasEvent(ctx context.Context, product *model.Product) {
	data := &folder_event.FolderUpdateAliasEventData{
		Id:    "D" + product.Id,
		Alias: product.Name,
	}
	evt := folder_event.NewFolderUpdateAliasEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderUpdateAliasEvent error: %s", err.Error())
	}
}

// publishFolderDeleteEvent 发布删除目录事件
func (s *ProductService) publishFolderDeleteEvent(ctx context.Context, ids []string, caseId string) {
	tenantId := appctx.GetTenantId2(ctx)
	folderIds := make([]string, 0, len(ids))
	for _, id := range ids {
		folderIds = append(folderIds, "D"+id)
	}
	data := &folder_event.FolderDeleteEventData{
		Ids:      folderIds,
		TenantId: tenantId,
		BusId:    "case",
		EntityId: caseId,
	}
	evt := folder_event.NewFolderDeleteEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderDeleteEvent error: %s", err.Error())
	}
}