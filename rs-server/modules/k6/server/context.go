package server

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/common"
)

type RContext struct {
	ictx iris.Context
	restapp.RestAssembler
}

func NewRContext(ictx iris.Context) *RContext {
	return &RContext{ictx: ictx}
}

func (c *RContext) ParamString(key string) *common.Result[string] {
	v := c.ictx.Params().GetStringTrim(key)
	var err error
	if v == "" {
		err = errors.New(key + " not found")
	}
	return common.NewResult[string](v, err)
}

func (c *RContext) ParamBool(key string) *common.Result[bool] {
	v, e := c.ictx.Params().GetBool(key)
	return common.NewResult[bool](v, e)
}

func (c *RContext) ParamFloat64(key string) *common.Result[float64] {
	v, e := c.ictx.Params().GetFloat64(key)
	return common.NewResult[float64](v, e)
}

func (c *RContext) ParamInt(key string) *common.Result[int] {
	v, e := c.ictx.Params().GetInt(key)
	return common.NewResult[int](v, e)
}

func (c *RContext) ParamInt32(key string) *common.Result[int32] {
	v, e := c.ictx.Params().GetInt32(key)
	return common.NewResult[int32](v, e)
}

func (c *RContext) ParamInt64(key string) *common.Result[int64] {
	v, e := c.ictx.Params().GetInt64(key)
	return common.NewResult[int64](v, e)
}

func (c *RContext) ParamStrings(key string) *common.Result[[]string] {
	val := c.ictx.URLParamSlice(key)
	return common.NewResult[[]string](val, nil)
}

func (c *RContext) GetId() string {
	return c.ictx.Params().GetString("id")
}

func (c *RContext) GetTenantId() string {
	return c.ictx.Params().GetString("tenantId")
}

func (c *RContext) GetCaseId() string {
	return c.ictx.Params().GetString("caseId")
}

func (c *RContext) GetFindPaging() *ddd_repository.FindPagingQueryRequest {
	v, _ := c.RestAssembler.AsFindPagingRequest(c.ictx)
	return v
}

func (c *RContext) Ictx() iris.Context {
	return c.ictx
}

func (c *RContext) ReadJson(data ...any) *common.Result[any] {
	var v any
	for _, d := range data {
		v = d
	}
	isMap := false
	if v == nil {
		v = make(map[string]any)
		isMap = true
	} else if d, ok := v.(map[string]any); ok {
		v = d
		isMap = true
	}
	var err error
	if isMap {
		err = c.ictx.ReadJSON(v)
	} else {
		err = c.ictx.ReadJSON(&v)
	}
	return common.NewResult[any](v, err)
}

func (c *RContext) WriteJson(data any) error {
	c.ictx.StatusCode(iris.StatusOK)
	return c.ictx.JSON(data)
}

func (c *RContext) NewCommand() *common.Result[Command] {
	return common.NewResult[Command](Command{}, nil)
}

/*
func (c *RContext) NewEntity(s *schema.Schema) *common.Result[db.Entity] {
	ent := db.NewEntity()
	res := c.ReadJson(ent)
	if res.Error != nil {
		return common.NewResult[db.Entity](nil, res.Error)
	}
	if err := s.Validate(ent); err != nil {
		return common.NewResult[db.Entity](nil, err)
	}
	if err := s.Copy(res, ent); err != nil {
		return common.NewResult[db.Entity](nil, res.Error)
	}
	return common.NewResult[db.Entity](ent, nil)
}

*/

func (c *RContext) SetError(err error, httpStatus ...int) {
	status := iris.StatusInternalServerError
	for _, s := range httpStatus {
		status = s
	}
	c.ictx.SetErr(err)
	if err != nil {
		c.ictx.StatusCode(status)
		_, _ = c.ictx.WriteString(err.Error())
	}
}

func (c *RContext) Executor() *Executor {
	return NewExecutor(c)
}
