package restapi

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// CompanyAPI 公司REST API 控制器
type CompanyAPI struct {
	service  *service.CompanyService
	rootPath string
}

// NewCompanyAPI 构造函数
func NewCompanyAPI(rootPath string) *CompanyAPI {
	return &CompanyAPI{
		rootPath: rootPath,
		service:  service.NewCompanyService(),
	}
}

// NewAPIController 注册路由
func (s *CompanyAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.company.CompanyAPI", s)

	controller.Post("/company:submit", "Submit")
	controller.Post("/company", "Create")
	controller.Put("/company", "Update")
	controller.Delete("/company/{id}", "DeleteById")
	controller.Delete("/company:deleteBatch", "DeleteByIds", restapi.WithParamsInBody(true))
	controller.GetPaging("/company", "FindPagingByCaseId")
	controller.GetData("/company:by-tag", "FindPagingByTagId")
	controller.GetOne("/company/{id}", "FindById")
	controller.GetData("/company:total", "FindTotal")

	return controller
}

// Submit 创建或更新公司
func (s *CompanyAPI) Submit(ctx context.Context, cmd *command.CompanySubmitCommand) (any, error) {
	return &cmd.Data, s.service.Submit(ctx, cmd)
}

// Create 创建公司
func (s *CompanyAPI) Create(ctx context.Context, cmd *command.CompanyCreateCommand) (any, error) {
	return &cmd.Data, s.service.Create(ctx, cmd)
}

// Update 更新公司
func (s *CompanyAPI) Update(ctx context.Context, cmd *command.CompanyUpdateCommand) (any, error) {
	return &cmd.Data, s.service.Update(ctx, cmd)
}

// DeleteById 按 ID 删除公司
func (s *CompanyAPI) DeleteById(ctx context.Context, qry *query.CompanyFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

// DeleteByIds 批量删除公司
func (s *CompanyAPI) DeleteByIds(ctx context.Context, cmd *command.CompanyDeleteByIdsCommand) (any, error) {
	res := struct {
		Data    any  `json:"data"`
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}{Data: cmd.Data, Success: true, Error: nil}
	return res, s.service.DeleteByIds(ctx, cmd)
}

// FindById 按 ID 查询公司
func (s *CompanyAPI) FindById(ctx context.Context, qry *query.CompanyFindByIdQuery) (any, error) {
	return s.service.FindById(ctx, qry)
}

// FindPagingByCaseId 按案件ID分页查询公司
func (s *CompanyAPI) FindPagingByCaseId(ctx context.Context, qry *query.CompanyFindByCaseIdQuery) (store.FindPagingResult[*model.Company], error) {
	res := s.service.FindPagingByCaseId(ctx, qry)
	return res, res.GetError()
}

// FindPagingByTagId 按标签ID分页查询公司
func (s *CompanyAPI) FindPagingByTagId(ctx context.Context, qry *query.CompanyFindByTagIdQuery) (store.FindPagingResult[*model.Company], error) {
	res := s.service.FindPagingByTagId(ctx, qry)
	return res, res.GetError()
}

// FindTotal 统计案件下的公司数量
func (s *CompanyAPI) FindTotal(ictx iris.Context, ctx context.Context, qry *query.CompanyFindTotalQuery) error {
	res, err := s.service.FindTotal(ctx, qry)
	if err != nil {
		return err
	}
	body := fmt.Sprintf("<div><span style='margin-right: 8px;'><ui5-icon style='height: 32px; width: 32px;' name='%s'></ui5-icon></span><span style='font-size: 42px; color: var(--sapContent_LabelColor)'>%d</span></div>", res.Icon, res.Count)
	_, err = ictx.WriteString(body)
	return err
}
