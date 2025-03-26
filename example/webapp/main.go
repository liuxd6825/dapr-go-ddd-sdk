package main

import (
	"context"
	"github.com/kataras/iris/v12"
	icontext "github.com/kataras/iris/v12/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"net/http"
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
		appInit(server.App())
		return nil
	})
	_, err := restapp.Run(envCfg, runCfg, opts)
	if err != nil {
		return
	}
}

func appInit(app *iris.Application) {
	human_handler(app)
	humanMap_handler(app)
}

func humanMap_handler(app *iris.Application) {
	humanDao := dao.NewDao[map[string]any](&dao.NewConfig{
		DBSchema: dbschema.NewDBSchemaWithJsonSchemaText("humanMap.json", xtest.HumanSchema),
	})
	humanDao.Table().AutoMigrate(context.Background())

	app.Get("/api/v1/human-map:create", func(ictx *icontext.Context) {
		gp.Try(func() error {
			ctx := xtest.NewContext()
			humans := xtest.NewHumanMapList(1, "human-map")
			humanDao.CreateMany(ctx, humans)
			return ictx.JSON(humans)
		}).Catch(func(e error) {
			ictx.SetErr(e)
			ictx.StatusCode(http.StatusInternalServerError)
		})
	})

	app.Get("/api/v1/human-map:list", func(ictx *icontext.Context) {
		gp.Try(func() error {
			ctx := xtest.NewContext()
			humans := humanDao.FindAll(ctx)
			return ictx.JSON(humans)
		}).Catch(func(e error) {
			ictx.SetErr(e)
			ictx.StatusCode(http.StatusInternalServerError)
		})
	})
}

func human_handler(app *iris.Application) {
	humanDao := dao.NewDao[*xtest.Human](&dao.NewConfig{})
	humanDao.Table().AutoMigrate(context.Background())

	app.Get("/api/v1/human:create", func(ictx *icontext.Context) {
		gp.Try(func() error {
			ctx := xtest.NewContext()
			humans := xtest.NewHumanList(1, "human-test")
			humanDao.CreateMany(ctx, humans)
			return ictx.JSON(humans)
		}).Catch(func(e error) {
			ictx.SetErr(e)
			ictx.StatusCode(504)
		})
	})

	app.Get("/api/v1/human:list", func(ictx *icontext.Context) {
		gp.Try(func() error {
			ctx := xtest.NewContext()
			humans := humanDao.FindAll(ctx)
			return ictx.JSON(humans)
		}).Catch(func(e error) {
			ictx.SetErr(e)
			ictx.StatusCode(504)
		})
	})
}
