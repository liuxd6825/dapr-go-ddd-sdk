package main

import (
	"context"
	icontext "github.com/kataras/iris/v12/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
)

func main() {
	envCfg := restapp.NewEnvConfig("webapp")
	envCfg.Mysql["mysql"] = &restapp.MySqlConfig{
		Username: "root",
		Password: "11111111",
		Host:     "127.0.0.1",
		Port:     "3306",
		DbName:   "test",
	}
	runCfg := restapp.NewRunConfig()

	opts := restapp.NewRunOptions().AddOnStartEvent(func(server *restapp.HttpServer) error {
		humanDao := daos.NewDao[*xtest.Human](&daos.NewConfig{})
		humanDao.Table().AutoMigrate(context.Background())

		server.App().Get("/api/v1/human:create", func(ictx *icontext.Context) {
			gp.Try(func() error {
				ctx := xtest.NewContext()
				humans := xtest.NewHumanList(10, "human-test")
				humanDao.CreateMany(ctx, humans)
				return ictx.JSON(humans)
			}).Catch(func(e error) {
				ictx.SetErr(e)
				ictx.StatusCode(504)
			})
		})

		server.App().Get("/api/v1/human:list", func(ictx *icontext.Context) {
			gp.Try(func() error {
				ctx := xtest.NewContext()
				humans := humanDao.FindAll(ctx)
				return ictx.JSON(humans)
			}).Catch(func(e error) {
				ictx.SetErr(e)
				ictx.StatusCode(504)
			})
		})

		return nil
	})
	_, err := restapp.Run(envCfg, runCfg, opts)
	if err != nil {
		return
	}
}
