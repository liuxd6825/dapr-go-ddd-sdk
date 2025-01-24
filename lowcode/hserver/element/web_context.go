package element

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/jsonschema/v6"
)

type WebContext interface {
	Ictx() iris.Context
	GetTokenUser() appctx.AuthUser
	GetTenantId() string
	GetTenantName() string
	GetToken() appctx.AuthToken
	ReadJson(data ...any) any
	ReadString() string
	ReadBytes() []byte

	ReadObject(schema *jsonschema.Schema) map[string]any
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
	String(key string) string
	Bool(key string) bool
	Float64(key string) float64
	Int(key string) int
	Int32(key string) int32
	Int64(key string) int64
	Strings(key string) []string
	GetId() string
	GetCaseId() string
	GetFindPaging() *ddd_repository.FindPagingQueryRequest
	Err() error
	Value(key any) any
}
