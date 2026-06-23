package restapi

import (
	"context"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// HumanAPI 人员REST API 控制器
type HumanAPI struct {
	service  *service.HumanService
	rootPath string
}

// NewHumanAPI 构造函数
func NewHumanAPI(rootPath string) *HumanAPI {
	return &HumanAPI{
		rootPath: rootPath,
		service:  service.NewHumanService(),
	}
}

// NewAPIController 注册路由
func (s *HumanAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanAPI", s)

	controller.Post("/human:submit", "Submit")
	controller.Post("/human", "Create")
	controller.Put("/human", "Update")
	controller.Delete("/human/{id}", "DeleteById")
	controller.Delete("/human:deleteBatch", "DeleteByIds")
	controller.GetPaging("/human", "FindPagingByCaseId")
	controller.GetData("/human:by-tag", "FindPagingByTagId")
	controller.GetData("/human:total", "FindCardTotal")
	controller.GetOne("/human/{id}", "FindById")
	controller.GetData("/human:personType", "GetPersonType")
	controller.GetData("/people-type", "GetPeopleType")
	// GetBasicHTML 特殊：使用 iris.Context
	controller.Handle("GET", "/basic-html", "GetBasicHTML")

	return controller
}

// Submit 创建或更新人员
func (s *HumanAPI) Submit(ctx context.Context, cmd *command.HumanSubmitCommand) (any, error) {
	return &cmd.Data, s.service.Submit(ctx, cmd)
}

// Create 创建人员
func (s *HumanAPI) Create(ctx context.Context, cmd *command.HumanCreateCommand) (any, error) {
	return &cmd.Data, s.service.Create(ctx, cmd)
}

// Update 更新人员
func (s *HumanAPI) Update(ctx context.Context, cmd *command.HumanUpdateCommand) (any, error) {
	return &cmd.Data, s.service.Update(ctx, cmd)
}

// DeleteById 按 ID 删除人员
func (s *HumanAPI) DeleteById(ctx context.Context, qry *query.HumanFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

// DeleteByIds 批量删除人员
func (s *HumanAPI) DeleteByIds(ctx context.Context, cmd *command.HumanDeleteByIdsCommand) (any, error) {
	res := struct {
		Data    any  `json:"data"`
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}{Data: cmd.Data, Success: true, Error: nil}
	return res, s.service.DeleteByIds(ctx, cmd)
}

// FindById 按 ID 查询人员
func (s *HumanAPI) FindById(ctx context.Context, qry *query.HumanFindByIdQuery) (any, error) {
	return s.service.FindById(ctx, qry)
}

// FindPagingByCaseId 按案件ID分页查询人员
func (s *HumanAPI) FindPagingByCaseId(ctx context.Context, qry *query.HumanFindByCaseIdQuery) (store.FindPagingResult[*model.Human], error) {
	res := s.service.FindPagingByCaseId(ctx, qry)
	return res, res.GetError()
}

// FindPagingByTagId 按标签ID分页查询人员
func (s *HumanAPI) FindPagingByTagId(ctx context.Context, qry *query.HumanFindByTagIdQuery) (store.FindPagingResult[*model.Human], error) {
	res := s.service.FindPagingByTagId(ctx, qry)
	return res, res.GetError()
}

// FindCardTotal 统计案件下的人员数量
func (s *HumanAPI) FindCardTotal(ctx context.Context, ictx iris.Context, qry *query.HumanFindTotalQuery) error {
	res, err := s.service.FindTotal(ctx, qry)
	if err != nil {
		return err
	}
	body := fmt.Sprintf("<div><span style='margin-right: 8px;'><ui5-icon style='height: 32px; width: 32px;' name='%s'></ui5-icon></span><span style='font-size: 42px; color: var(--sapContent_LabelColor)'>%d</span></div>", res.Icon, res.Count)
	_, err = ictx.WriteString(body)
	return err
}

// GetPersonType 获取人员类型列表
func (s *HumanAPI) GetPersonType(ctx context.Context) ([]query.PersonTypeItem, error) {
	return s.service.GetPersonType(), nil
}

// GetPeopleType 获取分析状态列表
func (s *HumanAPI) GetPeopleType(ctx context.Context) ([]query.PeopleTypeItem, error) {
	return s.service.GetPeopleType(), nil
}

// GetBasicHTML 渲染基础HTML模板（特殊端点：使用 iris.Context 直写 HTML）
// 原始 webscript 调用 pkg.template.renderFile("/test/page/parts/basic/basic.part.html")
// Go 端简化为返回基础 HTML 模板内容（实际模板渲染需要完整 hserver 模板引擎支持）
func (s *HumanAPI) GetBasicHTML(ictx iris.Context) {
	const body = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Basic HTML</title>
</head>
<body>
    <h1>Basic HTML Template</h1>
    <p>This is the basic HTML template rendered for human-service.</p>
</body
</html>`
	ictx.HTML(body)
}
