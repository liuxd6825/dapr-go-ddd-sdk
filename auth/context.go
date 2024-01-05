package auth

import (
	"context"
	"errors"
)

type authKey struct {
}

var (
	NotFundErr      = errors.New("AuthContext not found")
	ContextIsNilErr = errors.New("context is null")
	ctxKey          = authKey{}
)

func NewContext(ctx context.Context, jwtText string) (context.Context, error) {
	tk, err := getToken(jwtText)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, ctxKey, tk), nil
}

func GetUser(ctx context.Context) (User, error) {
	token, err := GetToken(ctx)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("token is null")
	}
	if token.GetUser() == nil {
		return nil, errors.New("token.User is null")
	}
	return token.GetUser(), nil
}

func GetToken(ctx context.Context) (Token, error) {
	if ctx == nil {
		return nil, ContextIsNilErr
	}
	val := ctx.Value(ctxKey)
	if val != nil {
		return val.(Token), nil
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
