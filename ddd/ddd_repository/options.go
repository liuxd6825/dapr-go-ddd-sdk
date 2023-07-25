package ddd_repository

import "time"

type Options interface {
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
	// GetUpdateMask
	// @Description: 更新数据时，跳过不更新的字段名
	// @return []string
	//
	GetUpdateCancel() []string
	SetUpdateCancel(v []string) Options
	SetUpdateCancelByDefault() Options

	Merge(opts ...Options) Options
}

type options struct {
	sort         *string
	timeout      *time.Duration
	updateFields []string
	updateCancel []string
	upsert       *bool
}

func NewOptions(o ...Options) Options {
	res := &options{}
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
	}
	return res
}

func (o *options) GetTimeout() *time.Duration {
	return o.timeout
}

func (o *options) SetTimeout(t *time.Duration) Options {
	o.timeout = t
	return o
}

func (o *options) GetUpdateFields() []string {
	return o.updateFields
}

func (o *options) SetUpdateFields(updateFields []string) Options {
	o.updateFields = updateFields
	return o
}

func (o *options) SetUpdateCancelByDefault() Options {
	o.updateCancel = []string{"CreatedTime", "CreatorId", "CreatorName", "Id"}
	return o
}

func (o *options) SetUpdateCancel(v []string) Options {
	o.updateCancel = v
	return o
}

func (o *options) GetUpdateCancel() []string {
	return o.updateCancel
}

func (o *options) SetUpsertIsNull() Options {
	o.upsert = nil
	return o
}

func (o *options) GetUpsert() *bool {
	return o.upsert
}

func (o *options) SetUpsert(v bool) Options {
	o.upsert = &v
	return o
}

func (o *options) GetSort() *string {
	return o.sort
}

func (o *options) SetSort(s *string) Options {
	o.sort = s
	return o
}

func (o *options) Merge(opts ...Options) Options {
	res := &options{}
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
