package appctx

import (
	"context"
	"errors"
)

type authKey struct {
}

var (
	NotFundErr      = errors.New("AuthContext not found")
	ContextIsNilErr = errors.New("context is null")
	authCtxKey      = authKey{}
)

func NewAuthContext(ctx context.Context, jwtText string) (context.Context, error) {
	tk, err := getToken(jwtText)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, authCtxKey, tk), nil
}

func GetAuthUser(ctx context.Context) (AuthUser, error) {
	token, err := GetAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("token is null")
	}
	if token.GetUser() == nil {
		return nil, errors.New("token.AuthUser is null")
	}
	return token.GetUser(), nil
}

func GetAuthToken(ctx context.Context) (AuthToken, error) {
	if ctx == nil {
		return nil, ContextIsNilErr
	}
	val := ctx.Value(authCtxKey)
	if val != nil {
		return val.(AuthToken), nil
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
