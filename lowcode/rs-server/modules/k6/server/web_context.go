package server

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/k6server/js/modules"
)

type WebContext struct {
	ictx iris.Context
	ctx  context.Context
	restapp.RestAssembler
	params *Params
	vu     modules.VU
}

func NewWebContext(ctx context.Context, ictx iris.Context, vu modules.VU) *WebContext {
	return &WebContext{ictx: ictx, ctx: ctx, params: NewParams(ictx), vu: vu}
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
	if err != nil {
		panic(err)
	}
	return common.NewResult[any](v, err)
}

func (c *WebContext) ReadString() string {
	bytes := c.ReadBytes()
	return string(bytes)
}

func (c *WebContext) ReadBytes() []byte {
	bytes, err := c.ictx.GetBody()
	if err != nil {
		panic(err)
	}
	return bytes
}

func (c *WebContext) ReadObject(schema *schema.Schema) map[string]any {
	object := map[string]any{}
	var err error
	if schema != nil {
		if err = c.ictx.ReadJSON(&object); err == nil {
			if err = schema.Validate(object); err == nil {
				object, err = schema.Convertor(object)
			}
		}
	}
	if err != nil {
		panic(errors.NewErr(err, "数据验证失败"))
	}
	return object
}

func (c *WebContext) WriteJson(data any) {
	err := c.ictx.JSON(data)
	if err != nil {
		panic(err)
	}
}

func (c *WebContext) WriteString(body string) int {
	res, err := c.ictx.WriteString(body)
	if err != nil {
		panic(err)
	}
	return res
}

func (c *WebContext) WriteHTML(body string) int {
	res, err := c.ictx.HTML(body)
	if err != nil {
		panic(err)
	}
	return res
}

func (c *WebContext) SetStatus(status int) {
	c.ictx.StatusCode(status)
}

func (c *WebContext) GetStatus() int {
	return c.ictx.GetStatusCode()
}

func (c *WebContext) SetError(errOrMsg any, httpStatus ...int) {
	status := iris.StatusInternalServerError
	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}
	var err error
	errId := idutils.NewId()
	if e, ok := errOrMsg.(error); ok {
		err = e
	} else if msg, ok := errOrMsg.(string); ok {
		err = errors.New(msg)
	} else {
		err = errors.New("未知的错误类型")
	}
	logs.Error(c.ctx, "", logs.Fields{"errId": errId, "error": err})
	if logs.GetLevel() == 0 {
		c.ictx.SetErr(err)
	} else {
		c.ictx.SetErr(errors.ErrorOf("执行时错误，错误编号：%s", errId))
	}
	c.ictx.StatusCode(status)
}

func (c *WebContext) Executor() *Executor {
	return NewExecutor(c)
}

func (c *WebContext) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
