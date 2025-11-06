package restapi

import (
	"context"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
	"github.com/stretchr/testify/assert"
)

type api struct {
	app *iris.Application
}

func (a *api) InitController(app *iris.Application) error {
	a.app = app
	h := NewController(app, "/api", a)
	//h.GetData("/get-data/{id}", "GetData")
	//h.GetPaging("/paging", "GetPaging")
	//h.GetOne("/get-one/{id}", "GetOne")
	h.Post("/create", "Create")
	return nil
}

type CreateCommand struct {
	CommandId string        `json:"commandId" required:"true" title:"Command ID"`
	Data      GetDataParams `json:"data" required:"true" title:"Data"`
}

type GetDataParams struct {
	Id   string      `json:"id" path:"id" required:"true" title:"ID"`
	Name string      `json:"name" query:"name" required:"true" title:"名称"`
	Type int         `json:"type" query:"type" required:"true" title:"类型"`
	Date *times.Date `json:"date" query:"date" required:"true" title:"日期"`
}

func (a *api) Create(ctx context.Context, ictx iris.Context, params *CreateCommand) (any, error) {
	return params, nil
}

func (a *api) GetOne(ctx context.Context, params *GetDataParams) (any, error) {
	return params, nil
}

func (a *api) GetData(ctx context.Context, params *GetDataParams) (any, error) {
	return params, nil
}

func (a *api) GetPaging(ctx context.Context, query *FindPagingRequest) (any, error) {
	return query, nil
}

func Test_Controller_Get(t *testing.T) {
	app := iris.New()
	api := &api{}
	InitController(app, api)
	err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	assert.NoError(t, err)
}
