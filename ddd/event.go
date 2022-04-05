package ddd

import "encoding/json"

type Event interface {
	GetTenantId() string
	GetCommandId() string
	GetEventId() string
}

func DoEvent(eventData map[string]interface{}, event interface{}) *DoEventResult {
	bs, err := json.Marshal(eventData)
	err = json.Unmarshal(bs, event)
	result := &DoEventResult{err: err, event: event}
	return result
}

type DoEventResult struct {
	err   error
	event interface{}
}

func (r *DoEventResult) OnSuccess(fun func(event interface{}) error) *DoEventResult {
	err := fun(r.event)
	if err != nil {
		r.err = err
	}
	return r
}

func (r *DoEventResult) OnError(fun func(err error)) *DoEventResult {
	fun(r.err)
	return r
}

func (r *DoEventResult) Error() error {
	return r.err
}
