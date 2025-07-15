package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/dao"
	command2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/template/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/template/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/template/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"sync"
)

type TemplateDomainService struct {
	dao *dao.TemplateDao
}

var (
	_templateOnce          sync.Once
	_templateDomainService *TemplateDomainService
)

func NewTemplateDomainService() *TemplateDomainService {
	_templateOnce.Do(func() {
		_templateDomainService = &TemplateDomainService{
			dao: dao.NewTemplateDao(config.DBKey),
		}
	})
	return _templateDomainService
}

func (t *TemplateDomainService) Create(ctx context.Context, cmd *command2.TempCreateCommand) {
	temp := &model.Template{}
	temp.Id = cmd.Data.Id
	temp.CaseId = cmd.Data.CaseId
	temp.TenantId = cmd.Data.TenantId
	temp.Name = cmd.Data.Name
	temp.BankName = cmd.Data.BankName
	temp.SheetName = cmd.Data.SheetName
	temp.MasterType = cmd.Data.MasterType
	temp.Remark = cmd.Data.Remark
	temp.FileId = cmd.Data.FileId
	temp.FileName = cmd.Data.FileName
	temp.Fields = cmd.Data.Fields
	temp.MapHeads = cmd.Data.MapHeads
	t.dao.Create(ctx, temp)
}

func (t *TemplateDomainService) Delete(ctx context.Context, cmd *command2.TempDeleteCommand) {
	if err := cmd.Validate(); err != nil {
		panic(err)
	}
	t.dao.DeleteById(ctx, cmd.Data.Id)
}

func (t *TemplateDomainService) CreateMany(ctx context.Context, list []*model.Template) {
	t.dao.CreateMany(ctx, list)
}

func (t *TemplateDomainService) Update(ctx context.Context, cmd *command2.TempUpdateCommand) {
	temp := &model.Template{}
	temp.Id = cmd.Data.Id
	temp.CaseId = cmd.Data.CaseId
	temp.TenantId = cmd.Data.TenantId
	temp.MasterType = cmd.Data.MasterType
	temp.Name = cmd.Data.Name
	temp.BankName = cmd.Data.BankName
	temp.Remark = cmd.Data.Remark
	temp.SheetName = cmd.Data.SheetName
	temp.FileId = cmd.Data.FileId
	temp.FileName = cmd.Data.FileName
	temp.MapHeads = cmd.Data.MapHeads
	temp.Fields = cmd.Data.Fields
	t.dao.Update(ctx, temp)
}

func (t *TemplateDomainService) FindById(ctx context.Context, qry *query.TemplateFindByIdQuery) (*model.Template, error) {
	if err := qry.Validate(); err != nil {
		panic(err)
	}
	return t.dao.FindById(ctx, qry.Id)
}

func (t *TemplateDomainService) FindPaging(ctx context.Context, caseId string, qry *query.TemplateFindPagingQuery) idao.FindPagingResult[*model.Template] {
	qry.SetMustFilter(fmt.Sprintf("case_id='%s'", caseId))
	return t.dao.FindPaging(ctx, qry)
}
