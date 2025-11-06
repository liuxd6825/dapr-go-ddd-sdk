package appctx

import (
	"context"
	"errors"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
)

type authKey struct {
}
type authUserKey struct{}

var (
	NotFundErr      = errors.New("AuthContext not found")
	ContextIsNilErr = errors.New("context is null")
	authCtxKey      = authKey{}
	authUserCtxKey  = authUserKey{}
)

func NewAuthContext(ctx context.Context, token string) (context.Context, error) {
	tk, err := getAuthToken(token)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, authCtxKey, tk), nil
}

func NewAuthContextUser(ctx context.Context, tk *AuthTokenEntity) (context.Context, error) {
	return context.WithValue(ctx, authCtxKey, tk), nil
}

func NewAuthUserContext(ctx context.Context, authUser AuthUser) (context.Context, error) {
	return context.WithValue(ctx, authUserCtxKey, authUser), nil
}

func SetAuthContext(ctx context.Context, token string) (context.Context, error) {
	newToken, err := getAuthToken(token)
	if err != nil {
		return nil, err
	}
	val := ctx.Value(authCtxKey)
	if val != nil {
		if oldToken, ok := val.(AuthToken); ok {
			oldToken.Copy(newToken)
			return ctx, nil
		}
	}
	return context.WithValue(ctx, authCtxKey, newToken), nil
}
func getAuthUser(ctx context.Context) (AuthUser, bool) {
	val := ctx.Value(authUserCtxKey)
	if val == nil {
		return nil, false
	}
	if authUser, ok := val.(AuthUser); ok {
		return authUser, true
	}
	return nil, false
}

func GetAuthUser(ctx context.Context) (AuthUser, bool) {
	if authUser, ok := getAuthUser(ctx); ok {
		return authUser, true
	}
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

func NewAuthUserEntity(valMap map[string]any) (tk AuthUser) {
	var authUser AuthUserEntity
	maputils.Decode(valMap, &authUser)
	return &authUser
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

func NewAuthTokenEntity(valMap map[string]any) (tk *AuthTokenEntity) {
	maputils.Decode(valMap, &tk)
	return tk
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
