package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"testing"
)

type api struct {
	app *iris.Application
}

func (a *api) InitController(app *iris.Application) error {
	a.app = app
	h := NewController(app, "/api", a)
	h.GetData("/get-data/{id}", "GetData")
	return nil
}

type GetDataParams struct {
	Id string `json:"id" path:"id" required:"true"`
}

func (a *api) GetData(ctx context.Context, params *GetDataParams) (any, error) {
	return params, nil
}

func Test_Controller_Get(t *testing.T) {
	app := iris.New()
	api := &api{}
	InitController(app, api)
	err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	assert.NoError(t, err)
}
