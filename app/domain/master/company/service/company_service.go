package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/query"
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

// CompanyService 公司业务服务
type CompanyService struct {
	dao         *dao.CompanyDao
	codeSvc     *code_service.CodeService
	accountDao  *dao.CompanyAccountDao
	companyRelDao   *dao.CompanyCompanyDao
	contractDao *dao.CompanyContractDao
	humanDao    *dao.CompanyHumanDao
	productDao  *dao.CompanyProductDao
}

var _companyService *CompanyService
var _companyServiceOnce sync.Once

// NewCompanyService 单例构造函数
func NewCompanyService() *CompanyService {
	_companyServiceOnce.Do(func() {
		_companyService = &CompanyService{
			dao:           dao.NewCompanyDao(config.DBKey),
			codeSvc:       code_service.NewCodeService(),
			accountDao:    dao.NewCompanyAccountDao(config.DBKey),
			companyRelDao: dao.NewCompanyCompanyDao(config.DBKey),
			contractDao:   dao.NewCompanyContractDao(config.DBKey),
			humanDao:      dao.NewCompanyHumanDao(config.DBKey),
			productDao:    dao.NewCompanyProductDao(config.DBKey),
		}
	})
	return _companyService
}

// Create 创建公司
func (s *CompanyService) Create(ctx context.Context, cmd *command.CompanyCreateCommand, opts ...idao.CallOptions) error {
	company := cmd.Data
	if company.Code == "" {
		company.Code = s.codeSvc.NewCompanyCode(ctx, company.CaseId)
	}
	if err := s.dao.Create(ctx, &company, opts...).GetError(); err != nil {
		return err
	}
	s.publishFolderCreateEvent(ctx, &company)
	return nil
}

// Update 更新公司
func (s *CompanyService) Update(ctx context.Context, cmd *command.CompanyUpdateCommand, opts ...idao.CallOptions) error {
	oldCompany, err := s.dao.FindById(ctx, cmd.Data.Id)
	if err != nil {
		return err
	}
	if err := s.dao.Update(ctx, &cmd.Data, opts...).GetError(); err != nil {
		return err
	}
	if oldCompany != nil && oldCompany.Name != cmd.Data.Name {
		s.publishFolderUpdateAliasEvent(ctx, &cmd.Data)
	}
	return nil
}

// Submit 创建或更新公司
func (s *CompanyService) Submit(ctx context.Context, cmd *command.CompanySubmitCommand, opts ...idao.CallOptions) error {
	old, err := s.dao.FindById(ctx, cmd.Data.Id, opts...)
	if err != nil {
		return err
	}
	if old != nil {
		return s.Update(ctx, &command.CompanyUpdateCommand{CommandId: cmd.CommandId, Data: cmd.Data}, opts...)
	}
	return s.Create(ctx, &command.CompanyCreateCommand{CommandId: cmd.CommandId, Data: cmd.Data}, opts...)
}

// DeleteById 按 ID 删除公司
func (s *CompanyService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

// DeleteByIds 批量删除公司，并级联清理子实体
func (s *CompanyService) DeleteByIds(ctx context.Context, cmd *command.CompanyDeleteByIdsCommand, opts ...idao.CallOptions) error {
	if err := s.dao.DeleteByIds(ctx, cmd.Data.Ids, opts...).GetError(); err != nil {
		return err
	}

	rsqlStr := rsql.NewBuilder().In("companyId", cmd.Data.Ids).Build()
	logs.Info(ctx, logs.Fields{"rsql": rsqlStr})

	if err := s.accountDao.DeleteByRSQL(ctx, rsqlStr, opts...).GetError(); err != nil {
		return err
	}
	if err := s.companyRelDao.DeleteByRSQL(ctx, rsqlStr, opts...).GetError(); err != nil {
		return err
	}
	if err := s.contractDao.DeleteByRSQL(ctx, rsqlStr, opts...).GetError(); err != nil {
		return err
	}
	if err := s.humanDao.DeleteByRSQL(ctx, rsqlStr, opts...).GetError(); err != nil {
		return err
	}
	if err := s.productDao.DeleteByRSQL(ctx, rsqlStr, opts...).GetError(); err != nil {
		return err
	}

	s.publishFolderDeleteEvent(ctx, cmd.Data.Ids, cmd.Data.CaseId)
	return nil
}

// FindById 按 ID 查询公司
func (s *CompanyService) FindById(ctx context.Context, qry *query.CompanyFindByIdQuery, opts ...idao.CallOptions) (*model.Company, error) {
	v, err := s.dao.FindById(ctx, qry.Id, opts...)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// FindPaging 分页查询公司
func (s *CompanyService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Company] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

// FindPagingByCaseId 按案件 ID 分页查询公司
func (s *CompanyService) FindPagingByCaseId(ctx context.Context, qry *query.CompanyFindByCaseIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Company] {
	return s.dao.FindPagingByCaseId(ctx, qry, qry.CaseId, opts...)
}

// FindPagingByTagId 按标签 ID 查询公司
func (s *CompanyService) FindPagingByTagId(ctx context.Context, qry *query.CompanyFindByTagIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Company] {
	paging := store.NewFindPagingQuery()
	paging.SetPageSize(store.PagingMaxPageSize)
	return s.dao.FindPagingByTagId(ctx, paging, qry.CaseId, qry.TagId, opts...)
}

// FindTotal 统计案件下的公司数量
func (s *CompanyService) FindTotal(ctx context.Context, qry *query.CompanyFindTotalQuery, opts ...idao.CallOptions) (*query.CompanyFindTotalResult, error) {
	count, err := s.dao.CountByCaseId(ctx, qry.CaseId, opts...)
	if err != nil {
		return nil, err
	}
	return &query.CompanyFindTotalResult{
		Count: count,
		Icon:  qry.Icon,
	}, nil
}

// publishFolderCreateEvent 发布创建目录事件
func (s *CompanyService) publishFolderCreateEvent(ctx context.Context, company *model.Company) {
	tenantId := appctx.GetTenantId2(ctx)
	data := &folder_event.FolderCreateEventData{
		Id:                "D" + company.Id,
		CaseId:            company.CaseId,
		BusId:             "case",
		EntityId:          company.CaseId,
		RootId:            fmt.Sprintf("%s_case_%s", tenantId, company.CaseId),
		RootPath:          fmt.Sprintf("/%s/case/%s", tenantId, company.CaseId),
		FolderPath:        fmt.Sprintf("/%s/case/%s/主数据附件/公司/%s", tenantId, company.CaseId, company.Code),
		ParentId:          fmt.Sprintf("%s_case_%s_master_company", tenantId, company.CaseId),
		Name:              company.Code,
		Alias:             company.Name,
		DisabledFrontEdit: true,
		Meta: []*folder_event.FolderCreateEventMeta{
			{
				SourceType: "公司",
				Source:     company.Id,
				Name:       "SourceId",
				Value:      company.Id,
			},
		},
	}
	evt := folder_event.NewFolderCreateEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderCreateEvent error: %s", err.Error())
	}
}

// publishFolderUpdateAliasEvent 发布更新目录别名事件
func (s *CompanyService) publishFolderUpdateAliasEvent(ctx context.Context, company *model.Company) {
	data := &folder_event.FolderUpdateAliasEventData{
		Id:    "D" + company.Id,
		Alias: company.Name,
	}
	evt := folder_event.NewFolderUpdateAliasEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderUpdateAliasEvent error: %s", err.Error())
	}
}

// publishFolderDeleteEvent 发布删除目录事件
func (s *CompanyService) publishFolderDeleteEvent(ctx context.Context, ids []string, caseId string) {
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