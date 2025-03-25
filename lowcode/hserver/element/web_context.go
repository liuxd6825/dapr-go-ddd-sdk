package element

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/jsonschema/v6"
)

type WebContext interface {
	ICtx() iris.Context
	Ctx() context.Context

	GetTokenUser() appctx.AuthUser
	GetTenantId() string
	GetTenantName() string
	GetToken() appctx.AuthToken
	ReadJson(data ...any) any
	ReadString() string
	ReadBytes() []byte
	ReadObject(schema *jsonschema.Schema) map[string]any

	GetId() string
	GetCaseId() string
	GetFindPaging() *store.FindPagingQueryRequest

	Valid(data any, schema *schema.Schema)
	FormFile(key string) *common.FormFile
	FormValue(name string, required bool) string
	FormObject(name string, required bool, schema *jsonschema.Schema) any
	WriteJson(data any)
	WriteString(body string) int
	WriteBytes(data []byte) int
	WriteHTML(body string) int
	SetContentType(cType string)
	GetContentType() string
	SetStatus(status int)
	GetStatus() int
	SetError(errOrMsg any, httpStatus ...int)
	Close()

	ValueString(key string) string
	ValueBool(key string) bool
	ValueFloat64(key string) float64
	ValueInt(key string) int
	ValueInt32(key string) int32
	ValueInt64(key string) int64
	ValueStrings(key string) []string
}
