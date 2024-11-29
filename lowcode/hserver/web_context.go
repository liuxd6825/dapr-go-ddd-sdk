package hserver

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
	schema2 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/k6server/js/modules"
	"io"
	"time"
)

type WebContext struct {
	restapp.RestAssembler
	authToken appctx.AuthToken
	ictx      iris.Context
	ctx       context.Context
	params    *Params
	vu        modules.VU
	closers   []io.Closer //资源关闭器

}

func NewWebContext(ctx context.Context, ictx iris.Context) *WebContext {
	authToken, ok := appctx.GetAuthToken(ctx)
	if !ok {
		panic("auth token not found")
	}

	return &WebContext{
		authToken: authToken,
		ictx:      ictx,
		ctx:       ctx,
		params:    NewParams(ictx),
		closers:   make([]io.Closer, 0),
	}
}

func (c *WebContext) Params() *Params {
	return c.params
}

func (c *WebContext) Ictx() iris.Context {
	return c.ictx
}

// GetTokenUser
//
//	@Description: 取Token中的用户信息
//	@receiver c
//	@return appctx.AuthUser
func (c *WebContext) GetTokenUser() appctx.AuthUser {
	return c.authToken.GetUser()
}

// GetTenantId
//
//	@Description: 取租户Id
//	@receiver c
//	@return string
func (c *WebContext) GetTenantId() string {
	return c.authToken.GetUser().GetTenantId()
}

// GetTenantName
//
//	@Description: 取租户名称
//	@receiver c
//	@return string
func (c *WebContext) GetTenantName() string {
	return c.authToken.GetUser().GetTenantName()
}

// GetToken
//
//	@Description: 取token信息
//	@receiver c
//	@return appctx.AuthToken
func (c *WebContext) GetToken() appctx.AuthToken {
	return c.authToken
}

// ReadJson
//
//	@Description: 从body中读取map
//	@receiver c
//	@param runValues
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
	var err error
	var timeFields map[string]any

	bytes := c.ReadBytes()
	if schema != nil {
		timeFields = getTimeFields(schema)
	}

	val, err := jsonutils.UnmarshalTime(bytes, &jsonutils.UnmarshalTimeOptions{
		TimeFields: timeFields,
		ParseTime:  parseTime,
	})
	if err != nil {
		panic(err)
	}

	object := val.(map[string]any)
	if schema != nil {
		schema.Validate(object)
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
//	@param runValues
//	@param schema
//	@return error
func (c *WebContext) Valid(data any, schema *schema.Schema) {
	schema.Validate(data)
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

	var err error
	var timeFields map[string]any

	ictx := c.ictx
	text := ictx.PostValue(name)
	if text == "" && required {
		panic(fmt.Sprintf("%s is required", name))
	}
	var object any
	if text == "" {
		return object
	}

	bytes := []byte(text)
	if schema != nil {
		timeFields = getTimeFields(schema)
	}

	val, err := jsonutils.UnmarshalTime(bytes, &jsonutils.UnmarshalTimeOptions{
		TimeFields: timeFields,
		ParseTime:  parseTime,
	})
	if err != nil {
		panic(err)
	}

	object = val

	if schema != nil {
		schema.Validate(object)
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

func (c *WebContext) SetData(data any) {
	c.WriteJson(data)
}

func (c *WebContext) SetHTML(html string) {
	c.WriteHTML(html)
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
		err = errors.New("runValues")
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

func (c *WebContext) GetCaseId() string {
	return c.ictx.Params().GetString("caseId")
}

func (c *WebContext) GetFindPaging() *ddd_repository.FindPagingQueryRequest {
	v, _ := c.RestAssembler.AsFindPagingRequest(c.ictx)
	return v
}

func (c *WebContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *WebContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *WebContext) Err() error {
	return c.ctx.Err()
}

func (c *WebContext) Value(key any) any {
	return c.ctx.Value(key)
}

const dateFormat = "2006-01-02"
const dateTimeFormat = "2006-01-02 15:04:05"

var parseTime = func(val string, key any) (timeVal any, err error) {
	if val == "" || val == "null" {
		return nil, nil
	}

	typeName, ok := key.(string)
	if !ok || typeName == schema2.TypeDate {
		tm, er := time.Parse(dateFormat, val)
		if er != nil {
			return nil, er
		}
		timeVal = times.NewDate(&tm)
	} else {
		tm, er := time.Parse(dateTimeFormat, val)
		if er != nil {
			return nil, er
		}
		timeVal = times.NewTime(&tm)
	}
	return timeVal, err
}

// getTimeFields
//
//	@Description: 从schema中读取date类型定义
//	@param props schema2.Properties
//	@return map[string]any
func getTimeFields(s schema2.ISchema) map[string]any {
	if s == nil {
		return nil
	}
	s.Init(s)

	var props schema2.Properties
	resFields := map[string]any{}
	timeFields := resFields
	dataType := s.GetType()
	if dataType == "array" {
		timeFields = map[string]any{}
		resFields[s.GetName()] = timeFields
		items := s.GetItems()
		if items != nil && items.Type == "object" {
			props = items.GetProperties()
		}
	} else if dataType == "object" {
		props = s.GetProperties()
	}
	for k, p := range props {
		if p == nil {
			continue
		}
		if p.IncludeType(schema2.TypeDate) {
			timeFields[k] = schema2.TypeDate
		} else if p.IncludeType(schema2.TypeDatetime) {
			timeFields[k] = schema2.TypeDatetime
		}
		if p.Properties != nil || p.Items != nil {
			if ps := getTimeFields(p); ps != nil {
				timeFields[k] = ps
			}
		}
	}
	return timeFields
}
