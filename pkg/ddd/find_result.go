package ddd

type OnSuccess func(data interface{}) error
type OnError func(err error) error
type OnIsFond func() error

func onSuccessDefault(data interface{}) error {
	return nil
}

func onErrorDefault(err error) error {
	return nil
}

func onNotFondDefault() error {
	return nil
}

type SetActions struct {
	err       error
	data      interface{}
	OnSuccess OnSuccess
	OnError   OnError
}

type SetOption func(actions *SetActions)

func SetOnSuccess(success OnSuccess) SetOption {
	return func(opt *SetActions) {
		opt.OnSuccess = success
	}
}

func SetOnError(onError OnError) SetOption {
	return func(opt *SetActions) {
		opt.OnError = onError
	}
}

func NewSetActions(onSuccess OnSuccess, onError OnError) *SetActions {
	actions := NewSetActionsNull()
	if onSuccess != nil {
		actions.OnSuccess = onSuccess
	}
	if onError != nil {
		actions.OnError = onError
	}
	return actions
}

func NewSetActionsNull() *SetActions {
	return &SetActions{
		OnSuccess: onSuccessDefault,
		OnError:   onErrorDefault,
	}
}

func (f *SetActions) Error() error {
	return f.err
}

func (f *SetActions) Data() interface{} {
	return f.data
}

func (f *SetActions) SetResult(data interface{}, err error) *SetActions {
	f.data = data
	f.err = err
	return f
}
