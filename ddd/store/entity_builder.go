package store

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"time"
)

type EntityBuilder[T any] interface {
	NewEntity() T
	NewEntityList() []T

	GetTenantId(entity T) string
	SetTenantId(entity T, tenantId string)

	GetId(entity T) string
	SetId(entity T, id string)

	GetCaseId(entity T) string
	SetCaseId(entity T, id string)

	GetAggId(entity T) string

	SetCreatedInfo(ctx context.Context, entity any)
	SetUpdatedInfo(ctx context.Context, entity any)
	SetDeletedInfo(ctx context.Context, entity any)

	GetAuthUser(ctx context.Context) appctx.AuthUser

	GetConfig() *EntityBuilderConfig
}

type EntityBuilderConfig struct {
	AggIdField       string // 聚合根字段名称
	IsMap            bool   // 是map对象
	IsCancelModified bool   // 取消修改字段的自动化处理
	Fields           Fields // 系统字段的名称
	DbSchema         *DBSchema
}

func NewEntityBuilderConfig() *EntityBuilderConfig {
	fields := newFields()
	return &EntityBuilderConfig{
		AggIdField: "id",
		Fields:     *fields,
	}
}

type AnyEntityBuilder[T any] struct {
	cfg    *EntityBuilderConfig
	schema *DBSchema
}

func NewAnyEntityBuilderWidthConfig[T any](cfg *EntityBuilderConfig) EntityBuilder[T] {
	cfg.IsMap = reflectutils.IsMap[T]()
	return &AnyEntityBuilder[T]{
		cfg: cfg,
	}
}

func NewAnyEntityBuilder[T any](schema *DBSchema) EntityBuilder[T] {
	cfg := NewEntityBuilderConfig()
	cfg.IsMap = reflectutils.IsMap[T]()
	return &AnyEntityBuilder[T]{
		cfg:    cfg,
		schema: schema,
	}
}

func (b *AnyEntityBuilder[T]) GetConfig() *EntityBuilderConfig {
	return b.cfg
}

func (b *AnyEntityBuilder[T]) NewEntity() T {
	v, err := reflectutils.NewObject[T]()
	if err != nil {
		panic(err)
	}
	return v
}

func (b *AnyEntityBuilder[T]) NewEntityList() []T {
	list, err := reflectutils.NewSlice[[]T]()
	if err != nil {
		panic(err)
	}
	return list
}

func (b *AnyEntityBuilder[T]) GetTenantId(entity T) string {
	return b.getFieldString(entity, b.cfg.Fields.TenantId)
}

func (b *AnyEntityBuilder[T]) SetTenantId(entity T, id string) {
	b.setFieldString(entity, b.cfg.Fields.TenantId, id)
}

func (b *AnyEntityBuilder[T]) GetId(entity T) string {
	return reflectutils.GetFieldString(entity, b.cfg.Fields.Id)
}

func (b *AnyEntityBuilder[T]) SetId(entity T, id string) {
	b.setFieldString(entity, b.cfg.Fields.Id, id)
}

func (b *AnyEntityBuilder[T]) GetCaseId(entity T) string {
	return reflectutils.GetFieldString(entity, b.cfg.Fields.CaseId)
}

func (b *AnyEntityBuilder[T]) SetCaseId(entity T, id string) {
	b.setFieldString(entity, fields.CaseId, id)
}

func (b *AnyEntityBuilder[T]) GetAggId(entity T) string {
	return reflectutils.GetFieldString(entity, b.cfg.AggIdField)
}

func (b *AnyEntityBuilder[T]) SetCreatedInfo(ctx context.Context, entity any) {
	e := any(entity)
	if e == nil {
		return
	}
	if b.cfg.IsCancelModified {
		return
	}
	authUser := b.GetAuthUser(ctx)
	timeNow := time.Now().UTC()

	b.setField(entity, fields.CreatedTime, timeNow)
	b.setField(entity, fields.CreatorName, authUser.GetName())
	b.setField(entity, fields.CreatorId, authUser.GetId())

	b.setField(entity, fields.UpdatedTime, timeNow)
	b.setField(entity, fields.UpdaterName, authUser.GetName())
	b.setField(entity, fields.UpdaterId, authUser.GetId())

}

func (b *AnyEntityBuilder[T]) SetUpdatedInfo(ctx context.Context, entity any) {
	e := any(entity)
	if e == nil {
		return
	}
	if b.cfg.IsCancelModified {
		return
	}
	authUser := b.GetAuthUser(ctx)
	timeNow := time.Now().UTC()

	b.setField(entity, fields.UpdatedTime, timeNow)
	b.setField(entity, fields.UpdaterName, authUser.GetName())
	b.setField(entity, fields.UpdaterId, authUser.GetId())
}

func (b *AnyEntityBuilder[T]) SetDeletedInfo(ctx context.Context, entity any) {

}

func (b *AnyEntityBuilder[T]) GetAuthUser(ctx context.Context) appctx.AuthUser {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user
	}
	panic("token is error")
}

func (b *AnyEntityBuilder[T]) setFieldString(entity any, fieldName, val string) {
	fieldName = b.getFieldName(fieldName)
	if field := b.schema.LookedField(fieldName); field != nil {
		reflectutils.SetFieldString(entity, fieldName, val)
	}
}

func (b *AnyEntityBuilder[T]) getFieldString(entity any, fieldName string) string {
	fieldName = b.getFieldName(fieldName)
	return reflectutils.GetFieldString(entity, fieldName)
}

func (b *AnyEntityBuilder[T]) setField(entity any, fieldName string, val any) {
	fieldName = b.getFieldName(fieldName)
	if field := b.schema.LookedField(fieldName); field != nil {
		reflectutils.SetField(entity, fieldName, val)
	}
}

func (b *AnyEntityBuilder[T]) getField(entity any, fieldName string) any {
	fieldName = b.getFieldName(fieldName)
	return reflectutils.GetField(entity, fieldName)
}

func (b *AnyEntityBuilder[T]) getFieldName(fieldName string) string {
	if b.cfg.IsMap {
		fieldName = stringutils.FirstLower(fieldName)
	} else {
		fieldName = stringutils.FirstUpper(fieldName)
	}
	return fieldName
}

func (c *EntityBuilderConfig) SetIsMap(isMap bool) *EntityBuilderConfig {
	c.IsMap = isMap
	return c
}

func (c *EntityBuilderConfig) SetIsCancelModified(isCancelModified bool) *EntityBuilderConfig {
	c.IsCancelModified = isCancelModified
	return c
}

func (c *EntityBuilderConfig) SetFields(fields *Fields) *EntityBuilderConfig {
	c.Fields = *fields
	return c
}
