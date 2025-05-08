package server

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
)

type Params struct {
	ictx iris.Context
}

func NewParams(ictx iris.Context) *Params {
	return &Params{ictx: ictx}
}

func (c *WebContext) String(key string) string {
	v := c.ictx.Params().GetStringTrim(key)
	var err error
	if v == "" {
		err = errors.New(key + " not found")
		panic(err)
	}
	return v
}

func (c *WebContext) Bool(key string) bool {
	v, err := c.ictx.Params().GetBool(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) Float64(key string) float64 {
	v, err := c.ictx.Params().GetFloat64(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) Int(key string) int {
	v, err := c.ictx.Params().GetInt(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) Int32(key string) int32 {
	v, err := c.ictx.Params().GetInt32(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) Int64(key string) int64 {
	v, err := c.ictx.Params().GetInt64(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) Strings(key string) []string {
	v := c.ictx.URLParamSlice(key)
	return v
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
