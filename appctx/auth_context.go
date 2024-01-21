package appctx

import (
	"context"
	"errors"
)

type authKey struct {
}

type AuthUser interface {
	GetId() string
	GetName() string
	GetPhone() string

	GetAccount() string
	GetRegDate() string
	GetWork() string
	GetStatus() string
	GetUserType() string

	GetTenantId() string
	GetTenantName() string
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

func GetAuthUser(ctx context.Context) (AuthUser, bool) {
	token, isFound := GetAuthToken(ctx)
	if !isFound {
		return nil, false
	}
	if token == nil {
		return nil, false
	}
	if token.GetUser() == nil {
		return nil, false
	}
	return token.GetUser(), true
}

func GetAuthToken(ctx context.Context) (AuthToken, bool) {
	if ctx == nil {
		return nil, false
	}
	val := ctx.Value(authCtxKey)
	if val != nil {
		return val.(AuthToken), true
	}
	return nil, false
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
