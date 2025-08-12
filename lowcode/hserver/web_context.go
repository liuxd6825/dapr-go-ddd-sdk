package hserver

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	appctx2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils/webjson"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/liuxd6825/k6server/js/modules"
	"io"
	"time"
)

type WebContext struct {
	restapp.RestAssembler
	authToken  appctx2.AuthToken
	ictx       iris.Context
	ctx        context.Context
	vu         modules.VU
	closers    []io.Closer //资源关闭器
	timeFields map[string]any
}

const dateFormat = "2006-01-02"
const dateTimeFormat = "2006-01-02 15:04:05"

func NewWebContext(ctx context.Context, ictx iris.Context) *WebContext {
	authToken, ok := appctx2.GetAuthToken(ctx)
	if !ok {
		panic("auth token not found")
	}

	return &WebContext{
		authToken: authToken,
		ictx:      ictx,
		ctx:       ctx,
		closers:   make([]io.Closer, 0),
	}
}

func (c *WebContext) Ctx() context.Context {
	return c.ctx
}

func (c *WebContext) ICtx() iris.Context {
	return c.ictx
}

// GetTokenUser
//
//	@Description: 取Token中的用户信息
//	@receiver c
//	@return appctx.AuthUser
func (c *WebContext) GetTokenUser() appctx2.AuthUser {
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
func (c *WebContext) GetToken() appctx2.AuthToken {
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

// ReadMap
//
//	@Description: 从body中读取对象
//	@receiver c
//	@param schema 有空：进行验证;  nil:不验证
//	@return map[string]any
func (c *WebContext) ReadMap(sch *jsonschema.Schema) map[string]any {
	var err error

	bytes := c.ReadBytes()
	val, err := jsonutils.UnmarshalTime(bytes, &jsonutils.UnmarshalTimeOptions{
		TimeFields: sch.GetTimeFields(),
		ParseTime:  parseTime,
	})
	if err != nil {
		panic(err)
	}

	schema.ApplyDefaults(sch, val)

	object, ok := val.(map[string]any)
	if !ok {
		panic("ReadObject() invalid object")
	}

	if dataVal, ok := object["data"]; ok {
		if data, ok := dataVal.(map[string]any); ok {
			if timeVal, ok := data["birthday"]; ok {
				if t, ok := timeVal.(*times.Date); ok {
					utc := t.Time()
					fmt.Println(utc.Format(time.DateTime))
				}
			}
		}
	}
	if sch != nil {
		err = schema.Validate(sch, object)
	}
	if err != nil {
		panic(err)
	}
	return object
}

func (c *WebContext) ReadObject(bytes []byte, sch *jsonschema.Schema, data any) any {
	var err error

	val, err := jsonutils.UnmarshalTime(bytes, &jsonutils.UnmarshalTimeOptions{
		TimeFields: sch.GetTimeFields(),
		ParseTime:  parseTime,
	})
	if err != nil {
		panic(err)
	}

	schema.ApplyDefaults(sch, val)

	if sch != nil {
		err = schema.Validate(sch, data)
	}
	if err != nil {
		panic(err)
	}
	return data
}

// Valid
//
//	@Description: 数据认证
//	@receiver c
//	@param runValues
//	@param schema
//	@return error
func (c *WebContext) Valid(data any, sch *jsonschema.Schema) {
	sch.Validate(data)
}

// FormFile
//
//	@Description: 从FormData中读取文件对象
//	@receiver c
//	@param key
//	@return *FormFileResult
func (c *WebContext) FormFile(key string) *common.FormFile {
	file, header, err := c.ictx.FormFile(key)
	if err != nil {
		panic(err)
	}
	formFile := common.NewFormFileResult(file, header)
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
func (c *WebContext) FormObject(name string, required bool, schema *jsonschema.Schema) any {
	if name == "" {
		panic("FormObject(name, schema) name parameter is not empty")
	}

	var err error
	timeFields := map[string]any{}

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
		//getTimeFields2(schema, timeFields, "")
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
		if err = schema.Validate(object); err != nil {
			panic(err)
		}
	}
	return object
}

func (c *WebContext) WriteJson(data any) {
	jsonStr, err := webjson.Marshal(data)
	if err != nil {
		panic(err)
	}
	c.ictx.ContentType("application/json")
	_, err = c.ictx.WriteString(jsonStr)
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

func (c *WebContext) WriteBytes(data []byte) int {
	res, err := c.ictx.Write(data)
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

func (c *WebContext) SetContentType(cType string) {
	c.ictx.ContentType(cType)
}

func (c *WebContext) GetContentType() string {
	return c.ictx.GetContentType()
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
		msgData, err := jsonutils.MarshalBytes(data)
		if err != nil {
			err = errors.New(string(msgData))
		} else {
			err = errors.New(fmt.Sprintf("%s", err.Error()))
		}
	} else {
		err = errors.New("未知的错误类型")
	}

	logs.Error(c.ctx, logs.Fields{"errorId": errId, "error": err})
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

func (c *WebContext) ValueString(key string) string {
	v := c.ictx.Params().GetStringTrim(key)
	var err error
	if v == "" {
		err = errors.New(key + " not found")
		panic(err)
	}
	return v
}

func (c *WebContext) ValueBool(key string) bool {
	v, err := c.ictx.Params().GetBool(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) ValueFloat64(key string) float64 {
	v, err := c.ictx.Params().GetFloat64(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) ValueInt(key string) int {
	v, err := c.ictx.Params().GetInt(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) ValueInt32(key string) int32 {
	v, err := c.ictx.Params().GetInt32(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) ValueInt64(key string) int64 {
	v, err := c.ictx.Params().GetInt64(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (c *WebContext) ValueStrings(key string) []string {
	v := c.ictx.URLParamSlice(key)
	return v
}

func (c *WebContext) GetId() string {
	return c.ictx.Params().GetString("id")
}

func (c *WebContext) GetCaseId() string {
	return c.ictx.Params().GetString("caseId")
}

func (c *WebContext) GetFindPaging() *store.FindPagingQueryRequest {
	v, _ := c.RestAssembler.AsFindPagingRequest(c.ictx)
	return v
}

func (c *WebContext) SetHeader(key string, value string) {
	c.ictx.Header(key, value)
}

var parseTime = func(val string, key any) (timeVal any, err error) {
	if val == "" || val == "null" {
		return nil, nil
	}
	if sch, ok := key.(*jsonschema.Schema); ok {
		if sch.Types.Contains(jsonschema.JsonType_DateTimeType) {
			if tm, err := times.AsTime(val); err != nil {
				return nil, err
			} else if tm != nil {
				timeVal = times.GetTime(tm.PTime())
			}
		} else if sch.Types.Contains(jsonschema.JsonType_DateType) {
			if tm, er := times.AsTime(val); er != nil {
				return nil, er
			} else if tm != nil {
				timeVal = times.GetDate(tm.PTime())
			}
		}
	}
	return timeVal, nil
}

// getTimeFields
//
//	@Description: 从schema中读取date类型定义
//	@param props schema2.Properties
//	@return map[string]any
/*
func getTimeFields(s schema.ISchema) map[string]any {
	if s == nil {
		return nil
	}
	s.Init(s)

	var props schema.Properties
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
		if p.IncludeType(schema.TypeDate) {
			timeFields[k] = schema.TypeDate
		} else if p.IncludeType(schema.TypeDatetime) {
			timeFields[k] = schema.TypeDatetime
		}
		if p.Properties != nil || p.Items != nil {
			if ps := getTimeFields(p); ps != nil {
				timeFields[k] = ps
			}
		}
	}
	return timeFields
}
*/

/*
func getTimeFields2(s *jsonschema.Schema, timeFields map[string]any, propName string) {
	if s == nil {
		return
	}

	if s.Types != nil {
		if s.Types.Contains(jsonschema.JsonType_ArrayType) {
			items := s.Items2020
			if items != nil && items.Types.Contains(jsonschema.JsonType_ObjectType) {
				propsTimeFields := map[string]any{}
				getTimeFieldsProps2(items.Properties, propsTimeFields)
				for k, v := range propsTimeFields {
					timeFields[k] = v
				}
			}
		} else if s.Types.Contains(jsonschema.JsonType_ObjectType) {
			getTimeFieldsProps2(s.Properties, timeFields)
		}
	}
	if s.Ref != nil {
		getTimeFields2(s.Ref, timeFields, propName)
	}

}
*/

/*
func getTimeFieldsProps2(props map[string]*jsonschema.Schema, timeFields map[string]any) {
	if props == nil {
		return
	}

	for k, p := range props {
		if p == nil {
			continue
		}
		if p.Types != nil {
			if p.Types.Contains(jsonschema.JsonType_DateType) {
				timeFields[k] = jsonschema.JsonType_DateType
			} else if p.Types.Contains(jsonschema.JsonType_DateTimeType) {
				timeFields[k] = jsonschema.JsonType_DateTimeType
			}
		}
		if p.Properties != nil || p.Items != nil || p.Ref != nil || p.Items2020 != nil {
			propsFields := map[string]any{}
			getTimeFields2(p, propsFields, k)
			if len(propsFields) > 0 {
				timeFields[k] = propsFields
			}
		}
	}
}
*/
