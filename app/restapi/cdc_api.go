package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/neo4jservice"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service/neo4jservice/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
)

type Cdc2Neo4jAPI struct {
	env      *env.Env
	rootPath string
	service  *neo4jservice.Neo4jService
}

func NewCdc2Neo4jAPI(env *env.Env, rootPath string) *Cdc2Neo4jAPI {
	ser := neo4jservice.NewNeo4jService()
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
	gp.Try(func() error {
		var record model.Record

		// 反序列化请求体到结构体
		if err := ctx.ReadJSON(&record); err != nil {
			return err
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
		return nil
	}).Catch(func(e error) {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": e.Error()})
	})

}
