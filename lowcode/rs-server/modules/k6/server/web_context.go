package server

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type WebContext struct {
	ictx iris.Context
	restapp.RestAssembler
	params *Params
}

func NewWebContext(ictx iris.Context) *WebContext {
	return &WebContext{ictx: ictx, params: NewParams(ictx)}
}

func (c *WebContext) Params() *Params {
	return c.params
}

func (c *WebContext) Ictx() iris.Context {
	return c.ictx
}

func (c *WebContext) ReadJson(data ...any) *common.Result[any] {
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

func (c *WebContext) ReadObject(schema *schema.Schema) (map[string]any, error) {
	object := map[string]any{}
	var err error
	if schema != nil {
		if err = c.ictx.ReadJSON(&object); err == nil {
			if err = schema.Validate(object); err == nil {
				object, err = schema.Convertor(object)
			}
		}
	}
	return object, err
}

func (c *WebContext) WriteJson(data any) error {
	c.ictx.StatusCode(iris.StatusOK)
	return c.ictx.JSON(data)
}

func (c *WebContext) SetError(err error, httpStatus ...int) {
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

func (c *WebContext) Executor() *Executor {
	return NewExecutor(c)
}

func (c *WebContext) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
