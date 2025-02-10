package mongodb

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"

type OperateOptions struct {
	AggId        *string
	EventType    *string
	EventVersion *string
	CommandId    *string
	ddd_repository.RepositoryOptions
}

func NewOperateOptions(opts ...*OperateOptions) *OperateOptions {
	o := new(OperateOptions)
	for _, i := range opts {
		if i.EventVersion != nil {
			o.EventVersion = i.EventVersion
		}
		if i.EventType != nil {
			o.EventType = i.EventType
		}
		if i.AggId != nil {
			o.AggId = i.AggId
		}
	}
	return o
}

func (e *OperateOptions) GetAggId(defVal string) string {
	if e != nil && e.AggId == nil {
		return defVal
	}
	return *e.AggId
}

func (e *OperateOptions) GetEventVersion(defVal string) string {
	if e != nil && e.EventVersion == nil {
		return defVal
	}
	return *e.EventVersion
}

func (e *OperateOptions) GetCommandId(defVal string) string {
	if e != nil && e.CommandId != nil {
		return *e.CommandId
	}
	return defVal
}

func newOptions(opts []*OperateOptions) []ddd_repository.Options {
	var res []ddd_repository.Options
	for _, o := range opts {
		res = append(res, o)
	}
	return res
}
