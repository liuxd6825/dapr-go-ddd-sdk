package service

import (
	"context"
	"fmt"
	"sync"

	folder_event "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/query"
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

// HumanService 人员业务服务
type HumanService struct {
	dao               *dao.HumanDao
	codeSvc           *code_service.CodeService
	feign             any
	credentialDao     *dao.HumanCredentialDao
	accountDao        *dao.HumanAccountDao
	addressDao        *dao.HumanAddressDao
	capitalDao        *dao.HumanCapitalDao
	companyDao        *dao.HumanCompanyDao
	contractDao       *dao.HumanContractDao
	extDao            *dao.HumanExtDao
	humanRelDao       *dao.HumanHumanDao
	linkDao           *dao.HumanLinkDao
	productDao        *dao.HumanProductDao
	reportedAmountDao *dao.HumanReportedAmountDao
	suspectAmountDao  *dao.HumanSuspectAmountDao
}

var _humanService *HumanService
var _humanServiceOnce sync.Once

// NewHumanService 单例构造函数
func NewHumanService() *HumanService {
	_humanServiceOnce.Do(func() {
		_humanService = &HumanService{
			dao:               dao.NewHumanDao(config.DBKey),
			codeSvc:           code_service.NewCodeService(),
			feign:             nil,
			credentialDao:     dao.NewHumanCredentialDao(config.DBKey),
			accountDao:        dao.NewHumanAccountDao(config.DBKey),
			addressDao:        dao.NewHumanAddressDao(config.DBKey),
			capitalDao:        dao.NewHumanCapitalDao(config.DBKey),
			companyDao:        dao.NewHumanCompanyDao(config.DBKey),
			contractDao:       dao.NewHumanContractDao(config.DBKey),
			extDao:            dao.NewHumanExtDao(config.DBKey),
			humanRelDao:       dao.NewHumanHumanDao(config.DBKey),
			linkDao:           dao.NewHumanLinkDao(config.DBKey),
			productDao:        dao.NewHumanProductDao(config.DBKey),
			reportedAmountDao: dao.NewHumanReportedAmountDao(config.DBKey),
			suspectAmountDao:  dao.NewHumanSuspectAmountDao(config.DBKey),
		}
	})
	return _humanService
}

// Create 创建人员
func (s *HumanService) Create(ctx context.Context, cmd *command.HumanCreateCommand, opts ...idao.CallOptions) error {
	human := cmd.Data
	if human.Code == "" {
		human.Code = s.codeSvc.NewHumanCode(ctx, human.CaseId)
	}
	if err := s.dao.Create(ctx, &human, opts...).GetError(); err != nil {
		return err
	}
	s.publishFolderCreateEvent(ctx, &human)
	return nil
}

// Update 更新人员
func (s *HumanService) Update(ctx context.Context, cmd *command.HumanUpdateCommand, opts ...idao.CallOptions) error {
	oldHuman, err := s.dao.FindById(ctx, cmd.Data.Id)
	if err != nil {
		return err
	}
	if err := s.dao.Update(ctx, &cmd.Data, opts...).GetError(); err != nil {
		return err
	}
	if oldHuman != nil && oldHuman.Name != cmd.Data.Name {
		s.publishFolderUpdateAliasEvent(ctx, &cmd.Data)
	}
	return nil
}

// Submit 创建或更新人员
func (s *HumanService) Submit(ctx context.Context, cmd *command.HumanSubmitCommand, opts ...idao.CallOptions) error {
	old, err := s.dao.FindById(ctx, cmd.Data.Id, opts...)
	if err != nil {
		return err
	}
	if old != nil {
		return s.Update(ctx, &command.HumanUpdateCommand{CommandId: cmd.CommandId, Data: cmd.Data}, opts...)
	}
	return s.Create(ctx, &command.HumanCreateCommand{CommandId: cmd.CommandId, Data: cmd.Data}, opts...)
}

// DeleteById 按 ID 删除人员
func (s *HumanService) DeleteById(ctx context.Context, id string, opts ...idao.CallOptions) error {
	return s.dao.DeleteById(ctx, id, opts...).GetError()
}

// DeleteByIds 批量删除人员，并级联清理 12 个子实体
func (s *HumanService) DeleteByIds(ctx context.Context, cmd *command.HumanDeleteByIdsCommand, opts ...idao.CallOptions) error {
	if err := s.dao.DeleteByIds(ctx, cmd.Data.Ids, opts...).GetError(); err != nil {
		return err
	}

	rsqlStr := rsql.NewBuilder().In("humanId", cmd.Data.Ids).Build()
	logs.Info(ctx, logs.Fields{"rsql": rsqlStr})

	// 12 张子表级联删除
	if err := s.credentialDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.accountDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.addressDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.capitalDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.companyDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.contractDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.extDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.humanRelDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.linkDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.productDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.reportedAmountDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	if err := s.suspectAmountDao.DeleteByHumanIds(ctx, cmd.Data.Ids, opts...); err != nil {
		return err
	}
	// human_record 不在级联删中

	s.publishFolderDeleteEvent(ctx, cmd.Data.Ids, cmd.Data.CaseId)
	return nil
}

// FindById 按 ID 查询人员
func (s *HumanService) FindById(ctx context.Context, qry *query.HumanFindByIdQuery, opts ...idao.CallOptions) (*model.Human, error) {
	human, err := s.dao.FindById(ctx, qry.Id, opts...)
	return human, err
}

// FindPaging 分页查询人员
func (s *HumanService) FindPaging(ctx context.Context, qry store.FindPagingQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Human] {
	return s.dao.FindPaging(ctx, qry, opts...)
}

// FindPagingByCaseId 按案件 ID 分页查询人员
func (s *HumanService) FindPagingByCaseId(ctx context.Context, qry *query.HumanFindByCaseIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Human] {
	return s.dao.FindPagingByCaseId(ctx, qry, qry.CaseId, opts...)
}

// FindPagingByTagId 按标签 ID 查询人员
func (s *HumanService) FindPagingByTagId(ctx context.Context, qry *query.HumanFindByTagIdQuery, opts ...idao.CallOptions) store.FindPagingResult[*model.Human] {
	paging := store.NewFindPagingQuery()
	paging.SetPageSize(store.PagingMaxPageSize)
	return s.dao.FindPagingByTagId(ctx, paging, qry.CaseId, qry.TagId, opts...)
}

// FindTotal 统计案件下的人员数量
func (s *HumanService) FindTotal(ctx context.Context, qry *query.HumanFindTotalQuery, opts ...idao.CallOptions) (*query.HumanFindTotalResult, error) {
	count, err := s.dao.CountByCaseId(ctx, qry.CaseId, opts...)
	if err != nil {
		return nil, err
	}
	return &query.HumanFindTotalResult{
		Count: count,
		Icon:  qry.Icon,
	}, nil
}

// GetPersonType 获取人员类型列表
func (s *HumanService) GetPersonType() []query.PersonTypeItem {
	return []query.PersonTypeItem{
		{Id: "reporter", Name: "报案人"},
		{Id: "suspect", Name: "嫌疑人"},
		{Id: "victim", Name: "受害人"},
		{Id: "third", Name: "第三人"},
	}
}

// GetPeopleType 获取分析状态列表
func (s *HumanService) GetPeopleType() []query.PeopleTypeItem {
	return []query.PeopleTypeItem{
		{Id: "01", Name: "待分析"},
		{Id: "02", Name: "排除"},
		{Id: "03", Name: "分析中"},
		{Id: "04", Name: "完成"},
	}
}

// AddPeopleType 添加人员分析状态
func (s *HumanService) AddPeopleType(id, name string) []query.PeopleTypeItem {
	return []query.PeopleTypeItem{
		{Id: id, Name: name},
		{Id: "01", Name: "待分析"},
		{Id: "02", Name: "排除"},
		{Id: "03", Name: "分析中"},
		{Id: "04", Name: "完成"},
	}
}

// publishFolderCreateEvent 发布创建目录事件
func (s *HumanService) publishFolderCreateEvent(ctx context.Context, human *model.Human) {
	tenantId := appctx.GetTenantId2(ctx)
	data := &folder_event.FolderCreateEventData{
		Id:                "D" + human.Id,
		CaseId:            human.CaseId,
		BusId:             "case",
		EntityId:          human.CaseId,
		RootId:            fmt.Sprintf("%s_case_%s", tenantId, human.CaseId),
		RootPath:          fmt.Sprintf("/%s/case/%s", tenantId, human.CaseId),
		FolderPath:        fmt.Sprintf("/%s/case/%s/主数据附件/人员/%s", tenantId, human.CaseId, human.Code),
		ParentId:          fmt.Sprintf("%s_case_%s_master_human", tenantId, human.CaseId),
		Name:              human.Code,
		Alias:             human.Name,
		DisabledFrontEdit: true,
		Meta: []*folder_event.FolderCreateEventMeta{
			{
				SourceType: "人员",
				Source:     human.Id,
				Name:       "SourceId",
				Value:      human.Id,
			},
		},
	}
	evt := folder_event.NewFolderCreateEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderCreateEvent error: %s", err.Error())
	}
}

// publishFolderUpdateAliasEvent 发布更新目录别名事件
func (s *HumanService) publishFolderUpdateAliasEvent(ctx context.Context, human *model.Human) {
	data := &folder_event.FolderUpdateAliasEventData{
		Id:    "D" + human.Id,
		Alias: human.Name,
	}
	evt := folder_event.NewFolderUpdateAliasEvent(ctx, folderAppId, data)
	if err := xbase.PublishEvent(ctx, evt); err != nil {
		logs.Errorfmt(ctx, "publish FolderUpdateAliasEvent error: %s", err.Error())
	}
}

// publishFolderDeleteEvent 发布删除目录事件
func (s *HumanService) publishFolderDeleteEvent(ctx context.Context, ids []string, caseId string) {
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
