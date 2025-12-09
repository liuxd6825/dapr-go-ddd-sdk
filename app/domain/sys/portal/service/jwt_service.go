package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	client "github.com/ory/kratos-client-go"
	"sync"
)

type JwtService struct {
}

var (
	_jwtOnce          sync.Once
	_jwtDomainService *JwtService
)

var secret = "#@!{[duXm-serVice-t0ken]},.(10086)$!"

func NewJwtService() *JwtService {
	_jwtOnce.Do(func() {
		_jwtDomainService = &JwtService{}
	})
	return _jwtDomainService
}

func (s *JwtService) Generate(session *client.Session, user *model.LoginUser) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	identityId := user.OryIdentityId
	if session.Identity != nil && session.Identity.Id != "" {
		identityId = session.Identity.Id
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = identityId
	claims["user"] = user
	claims["exp"] = session.ExpiresAt.Unix()
	claims["iat"] = session.AuthenticatedAt.Unix()
	claims["client_id"] = identityId

	return token.SignedString([]byte(secret))
}
