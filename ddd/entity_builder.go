package ddd

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"time"
)

type EntityBuilder[T any] interface {
	NewEntity() T
	NewEntityList() []T
	GetTenantId(entity T) string
	SetTenantId(entity T, tenantId string)
	GetId(entity T) string
	SetId(entity T, id string)

	SetCreatedInfo(ctx context.Context, entity any)
	SetUpdatedInfo(ctx context.Context, entity any)
	SetDeletedInfo(ctx context.Context, entity any)
}

type MapEntityBuilder[T map[string]any] struct {
	IsCancelModified   bool
	IsCancelSoftDelete bool
}

func NewMapEntityBuilder() EntityBuilder[map[string]any] {
	return &MapEntityBuilder[map[string]any]{}
}

func (m *MapEntityBuilder[T]) NewEntity() T {
	mv := map[string]any{}
	var a any = mv
	if t, ok := a.(T); ok {
		return t
	}
	panic(fmt.Errorf("cannot convert map to entity"))
}

func (m *MapEntityBuilder[T]) NewEntityList() []T {
	return make([]T, 0)
}

func (m *MapEntityBuilder[T]) GetTenantId(entity T) string {
	return get(entity, fields.TenantId)
}

func (m *MapEntityBuilder[T]) SetTenantId(entity T, id string) {
	m.Set(entity, fields.TenantId, id)
}

func (m *MapEntityBuilder[T]) GetId(entity T) string {
	return get(entity, fields.Id)
}

func (m *MapEntityBuilder[T]) SetId(entity T, id string) {
	m.Set(entity, fields.Id, id)
}

func (m *MapEntityBuilder[T]) Set(entity T, key string, value string) {
	var a any = entity
	if mv, ok := a.(map[string]any); ok {
		mv[key] = value
	}
}

func (m *MapEntityBuilder[T]) as(entity T) (map[string]any, bool) {
	var a any = entity
	if v, ok := a.(map[string]any); ok {
		return v, true
	}
	return nil, false
}

func (m *MapEntityBuilder[T]) SetCreatedInfo(ctx context.Context, entity any) {
	if entity == nil {
		return
	}
	if m.IsCancelModified {
		return
	}
	authUser := m.GetAuthUser(ctx)
	timeNow := time.Now()

	if e, ok := entity.(T); ok {
		e[fields.CreatedTime] = timeNow
		e[fields.CreatorName] = authUser.GetName()
		e[fields.CreatorId] = authUser.GetId()

		e[fields.UpdatedTime] = timeNow
		e[fields.UpdaterName] = authUser.GetName()
		e[fields.UpdaterId] = authUser.GetId()
		//e[fields.IsDeleted] = false
	}

}

func (m *MapEntityBuilder[T]) SetUpdatedInfo(ctx context.Context, entity any) {
	if entity == nil {
		return
	}
	if m.IsCancelModified {
		return
	}
	authUser := m.GetAuthUser(ctx)
	timeNow := time.Now()

	if e, ok := entity.(T); ok {
		//delete(e, fields.CreatedTime)
		//delete(e, fields.CreatorName)
		//delete(e, fields.CreatorId)

		e[fields.UpdatedTime] = times.GetTime(&timeNow)
		e[fields.UpdaterName] = authUser.GetName()
		e[fields.UpdaterId] = authUser.GetId()
	}
}

func (m *MapEntityBuilder[T]) SetDeletedInfo(ctx context.Context, entity any) {
	if entity == nil {
		return
	}
	if m.IsCancelSoftDelete {
		return
	}
	var timeNow *times.Time
	var userName *string
	var userId *string
	var isDeleted = true
	if isDeleted {
		authUser := m.GetAuthUser(ctx)
		name := authUser.GetName()
		id := authUser.GetId()

		timeNow = times.NowTime()
		userName = &name
		userId = &id
	}

	if e, ok := entity.(T); ok {
		e[fields.DeletedTime] = timeNow
		e[fields.DeleterName] = userName
		e[fields.DeleterId] = userId
		e[fields.IsDeleted] = isDeleted
	}
}

func (m *MapEntityBuilder[T]) GetAuthUser(ctx context.Context) appctx.AuthUser {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user
	}
	panic("token is error")
}

// /////////
type StructEntityBuilder[T any] struct {
}

type SetCreator interface {
	SetCreator(ctx context.Context, user appctx.AuthUser, name string)
}

type SetDeleter interface {
	SetDeleter(ctx context.Context, user appctx.AuthUser, isDeleted bool)
}

type SetUpdater interface {
	SetUpdater(ctx context.Context, user appctx.AuthUser)
}

func NewStructEntityBuilder[T any]() EntityBuilder[T] {
	return &StructEntityBuilder[T]{}
}

func (m *StructEntityBuilder[T]) NewEntity() T {
	ent, err := reflectutils.NewStruct[T]()
	if err != nil {
		panic(err)
	}
	return ent
}

func (m *StructEntityBuilder[T]) NewEntityList() []T {
	list, err := reflectutils.NewSlice[[]T]()
	if err != nil {
		panic(err)
	}
	return list
}

func (m *StructEntityBuilder[T]) GetTenantId(entity T) string {
	if e, ok := m.as(entity); ok {
		return e.GetTenantId()
	}
	return ""
}

func (m *StructEntityBuilder[T]) SetTenantId(entity T, tenantId string) {
	if e, ok := m.as(entity); ok {
		e.SetTenantId(tenantId)
	}
}

func (m *StructEntityBuilder[T]) GetId(entity T) string {
	if e, ok := m.as(entity); ok {
		return e.GetId()
	}
	return ""
}

func (m *StructEntityBuilder[T]) SetId(entity T, id string) {
	if e, ok := m.as(entity); ok {
		e.SetId(id)
	}
}

func (m *StructEntityBuilder[T]) SetCreatedInfo(ctx context.Context, entity any) {

}

func (m *StructEntityBuilder[T]) SetUpdatedInfo(ctx context.Context, entity any) {

}

func (m *StructEntityBuilder[T]) SetDeletedInfo(ctx context.Context, entity any) {

}

func (m *StructEntityBuilder[T]) as(entity T) (Entity, bool) {
	var a any = entity
	if e, ok := a.(Entity); ok {
		return e, true
	}
	return nil, false
}

func (m *StructEntityBuilder[T]) GetAuthUser(ctx context.Context) appctx.AuthUser {
	if user, ok := appctx.GetAuthUser(ctx); ok {
		return user
	}
	panic("token is error")
}

func get(entity any, key string) string {
	var a any = entity
	if mv, ok := a.(map[string]any); ok {
		if val, ok := mv[key]; ok {
			if v, ok := val.(string); ok {
				return v
			} else {
				return fmt.Sprintf("%v", val)
			}
		}
	}
	return ""
}
