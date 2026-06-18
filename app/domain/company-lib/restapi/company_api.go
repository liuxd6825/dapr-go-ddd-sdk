package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/ai/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/company-lib/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CompanyAPI struct {
	env      *env.Env
	service  *service.CompanyService
	aiSrv    *service2.AiCompanyService
	rootPath string
}

func NewCompanyAPI(env *env.Env, rootPath string) *CompanyAPI {
	return &CompanyAPI{
		env:      env,
		service:  service.NewCompanyService(),
		aiSrv:    service2.NewAiCompanyService(),
		rootPath: rootPath,
	}
}

func (s *CompanyAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/company-lib", "company.CompanyLibAPI", s)
	ctl.GetData("company/search", "Search")
	ctl.GetData("company/full-text-search", "FullTextSearch")
	ctl.Post("company/full-text-search", "GetCompanyJson")
	return ctl
}

func (s *CompanyAPI) Search(ctx context.Context, query *model.CompanyQuery) (*model.CompanyQueryResult, error) {
	return s.service.Search(ctx, query)
}

func (s *CompanyAPI) FullTextSearch(ctx context.Context, query *model.FullTextSearchQuery) (*model.CompanyQueryResult, error) {
	return s.service.FullTextSearch(ctx, query)
}

//func (s *CompanyAPI) GetCompanyJson(ctx context.Context, ictx iris.Context, qry *model.AiQueryRequest) error {
//	isPrint := logs.GetLevel() >= logs.InfoLevel
//	ictx.Header("Content-Type", "text/event-stream")
//	_, err := s.aiSrv.Query(ctx, qry, func(txt string) {
//		_, _ = ictx.Writef(txt)
//		if isPrint {
//			fmt.Print(txt)
//		}
//		ictx.ResponseWriter().Flush()
//	})
//	if err == nil && isPrint {
//		//fmt.Print("\n")
//	}
//	return err
//}

func (s *CompanyAPI) GetCompanyJson(ctx context.Context, ictx iris.Context, qry *model.AiQueryRequest) string {
	json, err := s.aiSrv.Query(ctx, qry, func(txt string) {})
	if err != nil {
		//fmt.Print("\n")
	}
	return json
}
