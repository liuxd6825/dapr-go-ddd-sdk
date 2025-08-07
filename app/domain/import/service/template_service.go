package service

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"sync"
)

type TemplateService struct {
	dao *dao.TemplateDao
	xbase.Service
}

var (
	_templateOnce          sync.Once
	_templateDomainService *TemplateService
)

func NewTemplateService() *TemplateService {
	_templateOnce.Do(func() {
		_templateDomainService = &TemplateService{
			dao: dao.NewTemplateDao(config.DBKey),
		}
	})
	return _templateDomainService
}

func (t *TemplateService) Create(ctx context.Context, cmd *command.TempCreateCommand) {
	temp := &model.Template{}
	temp.Id = cmd.Data.Id
	temp.CaseId = cmd.Data.CaseId
	temp.Name = cmd.Data.Name
	temp.BankName = cmd.Data.BankName
	temp.SheetName = cmd.Data.SheetName

	temp.SchemaId = cmd.Data.SchemaId
	temp.Remark = cmd.Data.Remark
	temp.FileId = cmd.Data.FileId
	temp.FileName = cmd.Data.FileName
	temp.Fields = cmd.Data.Fields
	temp.MapHeads = cmd.Data.MapHeads
	t.dao.Create(ctx, temp)
}

func (t *TemplateService) Delete(ctx context.Context, cmd *command.TempDeleteCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		return t.dao.DeleteById(ctx, cmd.Data.Id).GetError()
	})
}

func (t *TemplateService) CreateMany(ctx context.Context, list []*model.Template) error {
	return t.dao.CreateMany(ctx, list).GetError()
}

func (t *TemplateService) Update(ctx context.Context, cmd *command.TempUpdateCommand) error {
	return xbase.DoCommand(ctx, cmd, func(ctx context.Context) error {
		temp := &model.Template{}
		temp.Id = cmd.Data.Id
		temp.CaseId = cmd.Data.CaseId
		temp.SchemaId = cmd.Data.SchemaId
		temp.Name = cmd.Data.Name
		temp.BankName = cmd.Data.BankName
		temp.Remark = cmd.Data.Remark
		temp.SheetName = cmd.Data.SheetName
		temp.FileId = cmd.Data.FileId
		temp.FileName = cmd.Data.FileName
		temp.MapHeads = cmd.Data.MapHeads
		temp.Fields = cmd.Data.Fields
		return t.dao.Update(ctx, temp).GetError()
	})
}

func (t *TemplateService) FindById(ctx context.Context, qry *query.TemplateFindByIdQuery) (*model.Template, error) {
	if err := qry.Validate(); err != nil {
		panic(err)
	}
	return t.dao.FindById(ctx, qry.Id)
}

func (t *TemplateService) FindPaging(ctx context.Context, caseId string, qry *query.TemplateFindPagingQuery) idao.FindPagingResult[*model.Template] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s'", caseId))
	return t.dao.FindPaging(ctx, qry)
}
