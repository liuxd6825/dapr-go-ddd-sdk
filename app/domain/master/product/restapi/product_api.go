package restapi

import (
	"context"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// ProductAPI 产品REST API 控制器
type ProductAPI struct {
	service  *service.ProductService
	rootPath string
}

// NewProductAPI 构造函数
func NewProductAPI(rootPath string) *ProductAPI {
	return &ProductAPI{
		rootPath: rootPath,
		service:  service.NewProductService(),
	}
}

// NewAPIController 注册路由
func (s *ProductAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.product.ProductAPI", s)

	controller.Post("/product:submit", "Submit")
	controller.Post("/product", "Create")
	controller.Put("/product", "Update")
	controller.Delete("/product/{id}", "DeleteById")
	controller.Delete("/product:deleteBatch", "DeleteByIds")
	controller.GetPaging("/product", "FindPagingByCaseId")
	controller.GetData("/product:by-tag", "FindPagingByTagId")
	controller.GetData("/product:total", "FindTotal")
	controller.GetOne("/product/{id}", "FindById")

	return controller
}

// Submit 创建或更新产品
func (s *ProductAPI) Submit(ctx context.Context, cmd *command.ProductSubmitCommand) (any, error) {
	return &cmd.Data, s.service.Submit(ctx, cmd)
}

// Create 创建产品
func (s *ProductAPI) Create(ctx context.Context, cmd *command.ProductCreateCommand) (any, error) {
	return &cmd.Data, s.service.Create(ctx, cmd)
}

// Update 更新产品
func (s *ProductAPI) Update(ctx context.Context, cmd *command.ProductUpdateCommand) (any, error) {
	return &cmd.Data, s.service.Update(ctx, cmd)
}

// DeleteById 按 ID 删除产品
func (s *ProductAPI) DeleteById(ctx context.Context, qry *query.ProductFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

// DeleteByIds 批量删除产品
func (s *ProductAPI) DeleteByIds(ctx context.Context, cmd *command.ProductDeleteByIdsCommand) (any, error) {
	res := struct {
		Data    any  `json:"data"`
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}{Data: cmd.Data, Success: true, Error: nil}
	return res, s.service.DeleteByIds(ctx, cmd)
}

// FindById 按 ID 查询产品
func (s *ProductAPI) FindById(ctx context.Context, qry *query.ProductFindByIdQuery) (any, error) {
	return s.service.FindById(ctx, qry)
}

// FindPagingByCaseId 按案件ID分页查询产品
func (s *ProductAPI) FindPagingByCaseId(ctx context.Context, qry *query.ProductFindByCaseIdQuery) (store.FindPagingResult[*model.Product], error) {
	res := s.service.FindPagingByCaseId(ctx, qry)
	return res, res.GetError()
}

// FindPagingByTagId 按标签ID分页查询产品
func (s *ProductAPI) FindPagingByTagId(ctx context.Context, qry *query.ProductFindByTagIdQuery) (store.FindPagingResult[*model.Product], error) {
	res := s.service.FindPagingByTagId(ctx, qry)
	return res, res.GetError()
}

// FindTotal 统计案件下的产品数量
func (s *ProductAPI) FindTotal(ictx iris.Context, ctx context.Context, qry *query.ProductFindTotalQuery) error {
	res, err := s.service.FindTotal(ctx, qry)
	if err != nil {
		return err
	}
	body := fmt.Sprintf("<div><span style='margin-right: 8px;'><ui5-icon style='height: 32px; width: 32px;' name='%s'></ui5-icon></span><span style='font-size: 42px; color: var(--sapContent_LabelColor)'>%d</span></div>", res.Icon, res.Count)
	_, err = ictx.WriteString(body)
	return err
}
