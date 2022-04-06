package ddd_repository

type OnSuccess func(data interface{}) error
type OnSuccessPaging func(data *FindPagingData) error
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
