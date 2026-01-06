package restapi

import (
	"context"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RagApi struct {
	env        *env.Env
	rootPath   string
	ragService *service.RagService
}

func NewRagApi(env *env.Env, rootPath string) *RagApi {
	return &RagApi{
		env:        env,
		rootPath:   rootPath,
		ragService: service.NewRagService(),
	}
}

func (s *RagApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/import", "RagApi", s)
	ctl.Post("rag", "GetFieldJson")
	return ctl
}

func (s *RagApi) GetFieldJson(ctx context.Context, ictx iris.Context, qry *query.RagQueryRequest) error {
	isPrint := logs.GetLevel() >= logs.InfoLevel
	ictx.Header("Content-Type", "text/event-stream")
	_, err := s.ragService.Query(ctx, qry, func(txt string) {
		_, _ = ictx.Writef(txt)
		if isPrint {
			fmt.Print(txt)
		}
		ictx.ResponseWriter().Flush()
	})
	if err == nil && isPrint {
		//fmt.Print("\n")
	}
	return err
}
