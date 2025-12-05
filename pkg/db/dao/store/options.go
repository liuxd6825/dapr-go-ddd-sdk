package store

import (
	"context"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
)

type Options interface {
	GetEventType() *string
	SetEventType(v *string) Options

	GetEventVer() *string
	SetEventVer(v *string) Options

	GetCommandId() *string
	SetCommandId(v *string) Options

	//
	// GetTimeout
	// @Description: 超时时间
	// @return *time.Duration
	//
	GetTimeout() *time.Duration
	SetTimeout(v *time.Duration) Options

	//
	// GetSort
	// @Description: 排序字段
	// @return *string
	//
	GetSort() *string
	SetSort(*string) Options

	//
	// GetUpsert
	// @Description: true:如更新记录不存在,则新建记录;
	// @return *bool
	//
	GetUpsert() *bool
	SetUpsert(v bool) Options
	SetUpsertIsNull() Options

	//
	// GetUpdateFields
	// @Description: 更新数据时， 只更新的字段
	// @return *[]string
	//
	GetUpdateFields() []string
	SetUpdateFields([]string) Options

	//
	// GetUpdateCancel
	// @Description: 更新数据时，跳过不更新的字段名
	// @return []string
	//
	GetUpdateCancel() []string
	SetUpdateCancel(v []string) Options
	SetUpdateCancelByDefault() Options

	GetNotUpdateNull() bool
	SetNotUpdateNull(val bool) Options

	SetTenantId(v string) Options
	GetTenantId() string
	GetTenantId2(ctx context.Context) string

	SetUser(user User) Options
	GetUser() User
	GetUser2(ctx context.Context) User

	Merge(opts ...Options) Options
}

type User interface {
	GetId() string
	GetName() string
}

type RepositoryOptions struct {
	eventType *string
	eventVer  *string
	commandId *string

	sort          *string
	timeout       *time.Duration
	updateFields  []string
	updateCancel  []string
	upsert        *bool
	notUpdateNull *bool // 空值是否更新

	tenantId *string
	user     User
}

func NewOptions(o ...Options) Options {
	res := &RepositoryOptions{}
	for _, item := range o {
		if item.GetUpdateCancel() != nil {
			res.updateCancel = item.GetUpdateCancel()
		}
		if item.GetTimeout() != nil {
			res.timeout = item.GetTimeout()
		}
		if item.GetUpsert() != nil {
			res.upsert = item.GetUpsert()
		}
		if item.GetSort() != nil {
			res.sort = item.GetSort()
		}
		if item.GetUpdateFields() != nil {
			res.updateFields = item.GetUpdateFields()
		}
		if item.GetEventType() != nil {
			res.eventType = item.GetEventType()
		}
		if item.GetEventVer() != nil {
			res.eventVer = item.GetEventVer()
		}
		if item.GetCommandId() != nil {
			res.commandId = item.GetCommandId()
		}
		if item.GetTenantId() != "" {
			tenantId := item.GetTenantId()
			res.tenantId = &tenantId
		}
		if item.GetUser() != nil {
			res.user = item.GetUser()
		}
		if !item.GetNotUpdateNull() {
			val := item.GetNotUpdateNull()
			res.notUpdateNull = &val
		}
	}
	return res
}

func GetOptions(opts ...Options) Options {
	if opts == nil {
		return NewOptions()
	}
	if len(opts) == 1 && opts[0] != nil {
		return opts[0]
	}
	return NewOptions(opts...)
}

func (o *RepositoryOptions) GetUser() User {
	return o.user
}

func (o *RepositoryOptions) GetUser2(ctx context.Context) User {
	if o.user == nil {
		if user, ok := appctx.GetAuthUser(ctx); ok {
			return user
		}
	}
	return o.user
}

func (o *RepositoryOptions) SetUser(user User) Options {
	o.user = user
	return o
}

func (o *RepositoryOptions) SetTenantId(v string) Options {
	o.tenantId = &v
	return o
}

func (o *RepositoryOptions) GetTenantId() string {
	if o.tenantId == nil {
		return ""
	}
	return *o.tenantId
}

func (o *RepositoryOptions) GetTenantId2(ctx context.Context) string {
	if o.tenantId == nil {
		return appctx.GetTenantId2(ctx)
	}
	return *o.tenantId
}

func (o *RepositoryOptions) GetEventType() *string {
	return o.eventType
}

func (o *RepositoryOptions) SetEventType(v *string) Options {
	o.eventType = v
	return o
}

func (o *RepositoryOptions) GetEventVer() *string {
	return o.eventVer
}

func (o *RepositoryOptions) SetEventVer(v *string) Options {
	o.eventVer = v
	return o
}

func (o *RepositoryOptions) GetCommandId() *string {
	return o.commandId
}

func (o *RepositoryOptions) SetCommandId(v *string) Options {
	o.commandId = v
	return o
}

func (o *RepositoryOptions) GetNotUpdateNull() bool {
	if o.notUpdateNull == nil {
		return false
	}
	return *o.notUpdateNull
}

func (o *RepositoryOptions) SetNotUpdateNull(val bool) Options {
	b := val
	o.notUpdateNull = &b
	return o
}

func (o *RepositoryOptions) GetTimeout() *time.Duration {
	return o.timeout
}

func (o *RepositoryOptions) SetTimeout(t *time.Duration) Options {
	o.timeout = t
	return o
}

func (o *RepositoryOptions) GetUpdateFields() []string {
	return o.updateFields
}

func (o *RepositoryOptions) SetUpdateFields(updateFields []string) Options {
	o.updateFields = updateFields
	return o
}

func (o *RepositoryOptions) SetUpdateCancelByDefault() Options {
	o.updateCancel = []string{"CreatedTime", "CreatorId", "CreatorName", "Id"}
	return o
}

func (o *RepositoryOptions) SetUpdateCancel(v []string) Options {
	o.updateCancel = v
	return o
}

func (o *RepositoryOptions) GetUpdateCancel() []string {
	return o.updateCancel
}

func (o *RepositoryOptions) SetUpsertIsNull() Options {
	o.upsert = nil
	return o
}

func (o *RepositoryOptions) GetUpsert() *bool {
	return o.upsert
}

func (o *RepositoryOptions) SetUpsert(v bool) Options {
	o.upsert = &v
	return o
}

func (o *RepositoryOptions) GetSort() *string {
	return o.sort
}

func (o *RepositoryOptions) SetSort(s *string) Options {
	o.sort = s
	return o
}

func (o *RepositoryOptions) Merge(opts ...Options) Options {
	res := &RepositoryOptions{}
	var updateCancel []string
	for _, o := range opts {
		if o.GetSort() != nil {
			res.SetSort(o.GetSort())
		}
		if o.GetTimeout() != nil {
			res.SetTimeout(o.GetTimeout())
		}
		if o.GetUpdateFields() != nil {
			res.SetUpdateFields(o.GetUpdateFields())
		}
		if o.GetUpdateCancel() != nil {
			if updateCancel == nil {
				updateCancel = make([]string, 0)
			}
			mask := o.GetUpdateCancel()
			for _, v := range mask {
				updateCancel = append(updateCancel, v)
			}
		}
	}
	res.SetUpdateCancel(updateCancel)
	return res
}

func WithNotUpdateNull(val bool) Options {
	return &RepositoryOptions{
		notUpdateNull: &val,
	}
}
