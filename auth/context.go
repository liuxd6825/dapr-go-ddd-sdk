package auth

import (
	"context"
	"errors"
)

type authKey struct {
}

var NotFundErr = errors.New("AuthContext not found")
var ContextIsNilErr = errors.New("context is null")

func NewContext(ctx context.Context, jwtText string, jwtKey string) (context.Context, error) {
	user, err := decodeJwt(jwtText, jwtKey)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, authKey{}, user), nil
}

func GetLoginUser(ctx context.Context) (LoginUser, error) {
	if ctx == nil {
		return nil, ContextIsNilErr
	}
	val := ctx.Value(authKey{})
	if val != nil {
		return val.(LoginUser), nil
	}
	return nil, NotFundErr
}

func IsNotFundErr(err error) bool {
	if errors.Is(err, NotFundErr) {
		return true
	}
	return false
}

func IsContextIsNilErr(err error) bool {
	if errors.Is(err, ContextIsNilErr) {
		return true
	}
	return false
}
