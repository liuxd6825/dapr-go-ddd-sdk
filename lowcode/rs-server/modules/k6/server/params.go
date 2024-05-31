package server

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
)

type Params struct {
	ictx iris.Context
}

func NewParams(ictx iris.Context) *Params {
	return &Params{ictx: ictx}
}

func (c *WebContext) String(key string) *common.Result[string] {
	v := c.ictx.Params().GetStringTrim(key)
	var err error
	if v == "" {
		err = errors.New(key + " not found")
	}
	return common.NewResult[string](v, err)
}

func (c *WebContext) Bool(key string) *common.Result[bool] {
	v, e := c.ictx.Params().GetBool(key)
	return common.NewResult[bool](v, e)
}
func (c *WebContext) Float64(key string) *common.Result[float64] {
	v, e := c.ictx.Params().GetFloat64(key)
	return common.NewResult[float64](v, e)
}
func (c *WebContext) Int(key string) *common.Result[int] {
	v, e := c.ictx.Params().GetInt(key)
	return common.NewResult[int](v, e)
}
func (c *WebContext) Int32(key string) *common.Result[int32] {
	v, e := c.ictx.Params().GetInt32(key)
	return common.NewResult[int32](v, e)
}
func (c *WebContext) Int64(key string) *common.Result[int64] {
	v, e := c.ictx.Params().GetInt64(key)
	return common.NewResult[int64](v, e)
}
func (c *WebContext) Strings(key string) *common.Result[[]string] {
	val := c.ictx.URLParamSlice(key)
	return common.NewResult[[]string](val, nil)
}

func (c *WebContext) GetId() string {
	return c.ictx.Params().GetString("id")
}

func (c *WebContext) GetTenantId() string {
	return c.ictx.Params().GetString("tenantId")
}

func (c *WebContext) GetCaseId() string {
	return c.ictx.Params().GetString("caseId")
}

func (c *WebContext) GetFindPaging() *ddd_repository.FindPagingQueryRequest {
	v, _ := c.RestAssembler.AsFindPagingRequest(c.ictx)
	return v
}
