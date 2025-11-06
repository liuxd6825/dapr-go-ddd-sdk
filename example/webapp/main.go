package main

import (
	"context"
	"net/http"

	"github.com/kataras/iris/v12"
	icontext "github.com/kataras/iris/v12/context"
	restapp2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func main() {
	envCfg := restapp2.NewEnvConfig("webapp")
	envCfg.Mysql["mysql"] = &restapp2.MySqlConfig{
		User:     "root",
		Password: "11111111",
		Host:     "127.0.0.1",
		Port:     "3306",
		DbName:   "test",
	}
	runCfg := restapp2.NewRunConfig()

	opts := restapp2.NewRunOptions().AddOnStartEvent(func(server *restapp2.HttpServer) error {
		appInit(server.App())
		return nil
	})
	_, err := restapp2.Run(envCfg, runCfg, opts)
	if err != nil {
		return
	}
}

func appInit(app *iris.Application) {
	human_handler(app)
	humanMap_handler(app)
}

func humanMap_handler(app *iris.Application) {
	humanDao := dao.NewDao[map[string]any](&dao.DaoConfig{
		DBSchema: dbschema.NewDBSchemaWithJsonSchemaText("humanMap.json", xtest2.HumanSchema),
	})
	humanDao.Table().AutoMigrate(context.Background())

	app.Get("/api/v1/human-map:create", func(ictx *icontext.Context) {
		gp.Try(func() error {
			ctx := xtest2.NewContext()
			humans := xtest2.NewHumanMapList(1, "human-map")
			humanDao.CreateMany(ctx, humans)
			return ictx.JSON(humans)
		}).Catch(func(e error) {
			ictx.SetErr(e)
			ictx.StatusCode(http.StatusInternalServerError)
		})
	})

	app.Get("/api/v1/human-map:list", func(ictx *icontext.Context) {
		gp.Try(func() error {
			ctx := xtest2.NewContext()
			humans := humanDao.FindAll(ctx)
			return ictx.JSON(humans)
		}).Catch(func(e error) {
			ictx.SetErr(e)
			ictx.StatusCode(http.StatusInternalServerError)
		})
	})
}

func human_handler(app *iris.Application) {
	humanDao := dao.NewDao[*xtest2.Human](&dao.DaoConfig{})
	humanDao.Table().AutoMigrate(context.Background())

	app.Get("/api/v1/human:create", func(ictx *icontext.Context) {
		gp.Try(func() error {
			ctx := xtest2.NewContext()
			humans := xtest2.NewHumanList(1, "human-test")
			humanDao.CreateMany(ctx, humans)
			return ictx.JSON(humans)
		}).Catch(func(e error) {
			ictx.SetErr(e)
			ictx.StatusCode(504)
		})
	})

	app.Get("/api/v1/human:list", func(ictx *icontext.Context) {
		gp.Try(func() error {
			ctx := xtest2.NewContext()
			humans := humanDao.FindAll(ctx)
			return ictx.JSON(humans)
		}).Catch(func(e error) {
			ictx.SetErr(e)
			ictx.StatusCode(504)
		})
	})
}
