package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type Cdc2Neo4jAPI struct {
	env      *env.Env
	rootPath string
	service  *service.Neo4jService
}

func NewCdc2Neo4jAPI(env *env.Env, rootPath string) *Cdc2Neo4jAPI {
	ser := service.NewNeo4jService()
	return &Cdc2Neo4jAPI{
		env:      env,
		rootPath: rootPath,
		service:  ser,
	}
}

func (s *Cdc2Neo4jAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/cdc-mysql", "DataChange")
}

func (s *Cdc2Neo4jAPI) DataChange(ctx iris.Context) {
	var record service.RawRecord

	// 反序列化请求体到结构体
	if err := ctx.ReadJSON(&record); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid payload format"})
		return
	}

	// 根据操作类型处理数据
	switch record.OpType {
	case "c":
		s.service.Create(&record)
	case "u":
		s.service.Update(&record)
	case "d":
		s.service.Delete(&record)
	}

	ctx.StatusCode(iris.StatusOK)
}
