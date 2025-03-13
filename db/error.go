package db

type ResultError interface {
	GetError() error
}

func ParseError(result ResultError) {
	if res, ok := result.(ResultError); ok {
		if result != nil && res.GetError() != nil {
			panic(res.GetError())
		}
	}
}
