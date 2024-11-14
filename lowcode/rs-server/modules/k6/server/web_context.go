package server

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/k6server/js/modules"
	"io"
)

type WebContext struct {
	ictx iris.Context
	ctx  context.Context
	restapp.RestAssembler
	params  *Params
	vu      modules.VU
	closers []io.Closer //资源关闭器
}

func NewWebContext(ctx context.Context, ictx iris.Context, vu modules.VU) *WebContext {
	return &WebContext{ictx: ictx, ctx: ctx, params: NewParams(ictx), vu: vu, closers: make([]io.Closer, 0)}
}

func (c *WebContext) Params() *Params {
	return c.params
}

func (c *WebContext) Ictx() iris.Context {
	return c.ictx
}

func (c *WebContext) Ctx() context.Context {
	return c.ctx
}

// ReadJson
//
//	@Description: 从body中读取map
//	@receiver c
//	@param data
//	@return any
func (c *WebContext) ReadJson(data ...any) any {
	var v any
	for _, d := range data {
		v = d
	}
	if v == nil {
		v = make(map[string]any)
	}
	var err error
	err = c.ictx.ReadJSON(&v)
	if err != nil {
		panic(err)
	}
	return v
}

// ReadString
//
//	@Description: 从body中读取string
//	@receiver c
//	@return string
func (c *WebContext) ReadString() string {
	bytes := c.ReadBytes()
	return string(bytes)
}

// ReadBytes
//
//	@Description: 从body中读取[]byte
//	@receiver c
//	@return []byte
func (c *WebContext) ReadBytes() []byte {
	bytes, err := c.ictx.GetBody()
	if err != nil {
		panic(err)
	}
	return bytes
}

// ReadObject
//
//	@Description: 从body中读取对象
//	@receiver c
//	@param schema 有空：进行验证;  nil:不验证
//	@return map[string]any
func (c *WebContext) ReadObject(schema *schema.Schema) map[string]any {
	object := map[string]any{}
	var err error
	if err = c.ictx.ReadJSON(&object); err == nil {
		if schema != nil {
			if err = schema.Validate(object); err == nil {
				object, err = schema.Convertor(object)
			}
		}
	}
	if err != nil {
		panic(err)
	}
	return object
}

// Valid
//
//	@Description: 数据认证
//	@receiver c
//	@param data
//	@param schema
//	@return error
func (c *WebContext) Valid(data any, schema *schema.Schema) error {
	if schema != nil {
		return schema.Validate(data)
	}
	return nil
}

// FormFile
//
//	@Description: 从FormData中读取文件对象
//	@receiver c
//	@param key
//	@return *FormFileResult
func (c *WebContext) FormFile(key string) *FormFile {
	file, header, err := c.ictx.FormFile(key)
	if err != nil {
		panic(err)
	}
	formFile := NewFormFileResult(file, header)
	c.closers = append(c.closers, formFile)
	return formFile
}

// FormValue
//
//	@Description: 从FormData中读取string
//	@receiver c
//	@param key
//	@return string
func (c *WebContext) FormValue(name string, required bool) string {
	val := c.ictx.PostValue(name)
	if required && val == "" {
		panic(fmt.Sprintf("%s is required", name))
	}
	return val
}

// FormObject
//
//	@Description: 从FromValue中读取map数据
//	@receiver c
//	@param name formValue中的name名称
//	@param schema 有值:验证数据;
//	@return any
func (c *WebContext) FormObject(name string, required bool, schema *schema.Schema) any {
	if name == "" {
		panic("FormObject(name, schema) name parameter is not empty")
	}
	ictx := c.ictx
	text := ictx.PostValue(name)
	if text == "" && required {
		panic(fmt.Sprintf("%s is required", name))
	}
	object := map[string]any{}
	if text == "" {
		return object
	}

	if err := jsonutils.Unmarshal([]byte(text), &object); err != nil {
		panic(err)
	}
	if schema != nil {
		if err := schema.Validate(object); err == nil {
			object, err = schema.Convertor(object)
			if err != nil {
				panic(err)
			}
		}
	}
	return object
}

func (c *WebContext) WriteJson(data any) {
	jsonData, err := jsonutils.MarshalBytes(data)
	if err != nil {
		panic(err)
	}
	c.ictx.ContentType("application/json")
	_, err = c.ictx.Write(jsonData)
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
	} else if data, ok := errOrMsg.(map[string]any); ok {
		fmt.Println(data)
		err = errors.New("data")
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

func (c *WebContext) Close() {
	for _, closer := range c.closers {
		if err := closer.Close(); err != nil {
			c.Printf(err.Error())
		}
	}
}
