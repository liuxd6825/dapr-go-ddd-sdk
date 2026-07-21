package restapi

import (
	"context"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// ContractAPI 合同REST API 控制器
type ContractAPI struct {
	service  *service.ContractService
	rootPath string
}

// NewContractAPI 构造函数
func NewContractAPI(rootPath string) *ContractAPI {
	return &ContractAPI{
		rootPath: rootPath,
		service:  service.NewContractService(),
	}
}

// NewAPIController 注册路由
func (s *ContractAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.contract.ContractAPI", s)

	controller.Post("/contract:submit", "Submit")
	controller.Post("/contract", "Create")
	controller.Put("/contract", "Update")
	controller.Delete("/contract/{id}", "DeleteById")
	controller.Delete("/contract:deleteBatch", "DeleteByIds", restapi.WithParamsInBody(true))
	controller.GetPaging("/contract", "FindPagingByCaseId")
	controller.GetData("/contract:by-tag", "FindPagingByTagId")
	controller.GetOne("/contract/{id}", "FindById")
	controller.GetData("/contract:total", "FindTotal")

	return controller
}

// Submit 创建或更新合同
func (s *ContractAPI) Submit(ctx context.Context, cmd *command.ContractSubmitCommand) (any, error) {
	return &cmd.Data, s.service.Submit(ctx, cmd)
}

// Create 创建合同
func (s *ContractAPI) Create(ctx context.Context, cmd *command.ContractCreateCommand) (any, error) {
	return &cmd.Data, s.service.Create(ctx, cmd)
}

// Update 更新合同
func (s *ContractAPI) Update(ctx context.Context, cmd *command.ContractUpdateCommand) (any, error) {
	return &cmd.Data, s.service.Update(ctx, cmd)
}

// DeleteById 按 ID 删除合同
func (s *ContractAPI) DeleteById(ctx context.Context, qry *query.ContractFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

// DeleteByIds 批量删除合同
func (s *ContractAPI) DeleteByIds(ctx context.Context, cmd *command.ContractDeleteByIdsCommand) (any, error) {
	res := struct {
		Data    any  `json:"data"`
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}{Data: cmd.Data, Success: true, Error: nil}
	return res, s.service.DeleteByIds(ctx, cmd)
}

// FindById 按 ID 查询合同
func (s *ContractAPI) FindById(ctx context.Context, qry *query.ContractFindByIdQuery) (any, error) {
	return s.service.FindById(ctx, qry)
}

// FindPagingByCaseId 按案件ID分页查询合同
func (s *ContractAPI) FindPagingByCaseId(ctx context.Context, qry *query.ContractFindByCaseIdQuery) (store.FindPagingResult[*model.Contract], error) {
	res := s.service.FindPagingByCaseId(ctx, qry)
	return res, res.GetError()
}

// FindPagingByTagId 按标签ID分页查询合同
func (s *ContractAPI) FindPagingByTagId(ctx context.Context, qry *query.ContractFindByTagIdQuery) (store.FindPagingResult[*model.Contract], error) {
	res := s.service.FindPagingByTagId(ctx, qry)
	return res, res.GetError()
}

// FindTotal 统计案件下的合同数量
func (s *ContractAPI) FindTotal(ictx iris.Context, ctx context.Context, qry *query.ContractFindTotalQuery) error {
	res, err := s.service.FindTotal(ctx, qry)
	if err != nil {
		return err
	}
	body := fmt.Sprintf("<div><span style='margin-right: 8px;'><ui5-icon style='height: 32px; width: 32px;' name='%s'></ui5-icon></span><span style='font-size: 42px; color: var(--sapContent_LabelColor)'>%d</span></div>", res.Icon, res.Count)
	_, err = ictx.WriteString(body)
	return err
}
