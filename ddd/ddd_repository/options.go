package ddd_repository

import "time"

type Options interface {
	GetTimeout() *time.Duration
	SetTimeout(v *time.Duration) Options

	GetUpdateFields() *[]string
	SetUpdateFields(*[]string) Options

	GetSort() *string
	SetSort(*string) Options

	GetUpsert() *bool
	SetUpsert(v bool) Options
	SetUpsertIsNull() Options

	Merge(opts ...Options) Options
}

type options struct {
	sort         *string
	timeout      *time.Duration
	updateFields *[]string
	upsert       *bool
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

func NewOptions() Options {
	return &options{}
}

func (o *options) GetTimeout() *time.Duration {
	return o.timeout
}

func (o *options) SetTimeout(t *time.Duration) Options {
	o.timeout = t
	return o
}

func (o *options) GetUpdateFields() *[]string {
	return o.updateFields
}

func (o *options) SetUpdateFields(updateFields *[]string) Options {
	o.updateFields = updateFields
	return o
}

func (o *options) Merge(opts ...Options) Options {
	res := &options{}
	for _, o := range opts {
		if o.GetTimeout() != nil {
			res.SetTimeout(o.GetTimeout())
		}
		if o.GetUpdateFields() != nil {
			res.SetUpdateFields(o.GetUpdateFields())
		}
	}
	return res
}
