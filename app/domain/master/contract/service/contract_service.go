package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/query"
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

// ContractService 合同业务服务
type ContractService struct {
	dao         *dao.ContractDao
	codeSvc     *code_service.CodeService
	feign       any
	companyRelDao *dao.ContractCompanyDao
	contractDao *dao.ContractContractDao
	humanDao    *dao.ContractHumanDao
	productDao  *dao.ContractProductDao
}

var _contractService *ContractService
var _contractServiceOnce sync.Once

// NewContractService 单例构造函数
func NewContractService() *ContractService {
	_contractServiceOnce.Do(func() {
		_contractService = &ContractService{
			dao:           dao.NewContractDao(config.DBKey),
			codeSvc:       code_service.NewCodeService(),
			feign:         nil,
			companyRelDao: dao.NewContractCompanyDao(config.DBKey),
			contractDao:   dao.NewContractContractDao(config.DBKey),
			humanDao:      dao.NewContractHumanDao(config.DBKey),
			productDao:    dao.NewContractProductDao(config.DBKey),
		}
	})
	return _contractService
}

// Create 创建合同
func (s *ContractService) Create(ctx context.Context, cmd *command.ContractCreateCommand, opts ...idao.CallOptions) error {
	contract := cmd.Data
	if contract.Code == "" {
		contract.Code = s.codeSvc.NewContractCode(ctx, contract.CaseId)
	}
	if err := s.dao.Create(ctx, &contract, opts...).GetError(); err != nil {
		return err
	}
	s.publishFolderCreateEvent(ctx, &contract)
	return nil
}

// Update 更新合同
func (s *ContractService) Update(ctx context.Context, cmd *command.ContractUpdateCommand, opts ...idao.CallOptions) error {
	oldContract, err := s.dao.FindById(ctx, cmd.Data.Id)
	if err != nil {
		return err
	}
	if err := s.dao.Update(ctx, &cmd.Data, opts...).GetError(); err != nil {
		return err
	}
	if oldContract != nil && oldContract.Name != cmd.Data.Name {
		s.publishFolderUpdateAliasEvent(ctx, &cmd.Data)
	}
	return nil
}

// Submit 创建或更新合同
func (s *ContractService) Submit(ctx context.Context, cmd *command.ContractSubmitCommand, opts ...idao.CallOptions) error {
	old, err := s.dao.FindById(ctx, cmd.Data.Id, opts...)
	if err != nil {
		return err
	}
	if old != nil {
		return s.Update(ctx, &command.ContractUpdateCommand{CommandId: cmd.CommandId, Data: cmd.Data}, opts...)
	}
	return s.Create(ctx, &command.ContractCreateCommand{CommandId: cmd.CommandId, Data: cmd.Data}, opts...)
}

// DeleteById 按 ID 删除合同
func (s *ContractService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

// DeleteByIds 批量删除合同，并级联清理子实体
func (s *ContractService) DeleteByIds(ctx context.Context, cmd *command.ContractDeleteByIdsCommand, opts ...idao.CallOptions) error {
	if err := s.dao.DeleteByIds(ctx, cmd.Data.Ids, opts...).GetError(); err != nil {
		return err
	}

	rsqlStr := rsql.NewBuilder().In("contractId", cmd.Data.Ids).Build()
	logs.Info(ctx, logs.Fields{"rsql": rsqlStr})

	if err := s.companyRelDao.DeleteByContractIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.contractDao.DeleteByContractIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.humanDao.DeleteByContractIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.productDao.DeleteByContractIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}

	s.publishFolderDeleteEvent(ctx, cmd.Data.Ids, cmd.Data.CaseId)
	return nil
}

// FindById 按 ID 查询合同
func (s *ContractService) FindById(ctx context.Context, qry *query.ContractFindByIdQuery, opts ...idao.CallOptions) (*model.Contract, error) {
	return s.dao.FindById(ctx, qry.Id, opts...)
}

// FindPaging 分页查询合同
func (s *ContractService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Contract] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

// FindPagingByCaseId 按案件 ID 分页查询合同
func (s *ContractService) FindPagingByCaseId(ctx context.Context, qry *query.ContractFindByCaseIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Contract] {
	return s.dao.FindPagingByCaseId(ctx, qry, qry.CaseId, opts...)
}

// FindPagingByTagId 按标签 ID 查询合同
func (s *ContractService) FindPagingByTagId(ctx context.Context, qry *query.ContractFindByTagIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Contract] {
	paging := store.NewFindPagingQuery()
	paging.SetPageSize(store.PagingMaxPageSize)
	return s.dao.FindPagingByTagId(ctx, paging, qry.CaseId, qry.TagId, opts...)
}

// FindTotal 统计案件下的合同数量
func (s *ContractService) FindTotal(ctx context.Context, qry *query.ContractFindTotalQuery, opts ...idao.CallOptions) (*query.ContractFindTotalResult, error) {
	count, err := s.dao.CountByCaseId(ctx, qry.CaseId, opts...)
	if err != nil {
		return nil, err
	}
	return &query.ContractFindTotalResult{
		Count: count,
		Icon:  qry.Icon,
	}, nil
}

// publishFolderCreateEvent 发布创建目录事件
func (s *ContractService) publishFolderCreateEvent(ctx context.Context, contract *model.Contract) {
	tenantId := appctx.GetTenantId2(ctx)
	data := &folder_event.FolderCreateEventData{
		Id:                "D" + contract.Id,
		CaseId:            contract.CaseId,
		BusId:             "case",
		EntityId:          contract.CaseId,
		RootId:            fmt.Sprintf("%s_case_%s", tenantId, contract.CaseId),
		RootPath:          fmt.Sprintf("/%s/case/%s", tenantId, contract.CaseId),
		FolderPath:        fmt.Sprintf("/%s/case/%s/主数据附件/合同/%s", tenantId, contract.CaseId, contract.Code),
		ParentId:          fmt.Sprintf("%s_case_%s_master_contract", tenantId, contract.CaseId),
		Name:              contract.Code,
		Alias:             contract.Name,
		DisabledFrontEdit: true,
		Meta: []*folder_event.FolderCreateEventMeta{
			{
				SourceType: "合同",
				Source:     contract.Id,
				Name:       "SourceId",
				Value:      contract.Id,
			},
		},
	}
	evt := folder_event.NewFolderCreateEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderCreateEvent error: %s", err.Error())
	}
}

// publishFolderUpdateAliasEvent 发布更新目录别名事件
func (s *ContractService) publishFolderUpdateAliasEvent(ctx context.Context, contract *model.Contract) {
	data := &folder_event.FolderUpdateAliasEventData{
		Id:    "D" + contract.Id,
		Alias: contract.Name,
	}
	evt := folder_event.NewFolderUpdateAliasEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderUpdateAliasEvent error: %s", err.Error())
	}
}

// publishFolderDeleteEvent 发布删除目录事件
func (s *ContractService) publishFolderDeleteEvent(ctx context.Context, ids []string, caseId string) {
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