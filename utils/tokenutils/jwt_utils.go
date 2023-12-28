package tokenutils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"strings"
	"time"
)

type LoginUser interface {
	GetId() string
	GetName() string
	GetPhone() string

	GetAccount() string
	GetRegDate() *time.Time
	GetWork() string
	GetStatus() string
	GetUserType() string

	GetTenantId() string
	GetTenantName() string
	GetTenantAccount() string
}

type loginUser struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Phone    string     `json:"phone"`
	Account  string     `json:"account"`
	Regdate  *time.Time `json:"regdate"`
	Work     string     `json:"work"`
	Status   string     `json:"status"`
	UserType string     `json:"userType"`

	TenantId      string `json:"tenantId"`
	TenantName    string `json:"tenantName"`
	TenantAccount string `json:"tenantAccount"`
}

type loginClaim struct {
	User *loginUser
	jwt.RegisteredClaims
}

func decodeJwt(jwtString string) (LoginUser, error) {
	var newJwtString string
	tokenType := jwtString[0:6]
	if strings.ToLower(tokenType) == "bearer" {
		newJwtString = jwtString[7:len(jwtString)]
	} else {
		newJwtString = jwtString
	}

	token, err := jwt.ParseWithClaims(newJwtString, &loginClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("#@!{[duXm-serVice-t0ken]},.(10086)$!"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*loginClaim)
	if !ok || !token.Valid {
		return nil, errors.New("非法的Token内容")
	}
	return claims.User, nil
}

func (u *loginUser) GetId() string {
	return u.Id
}

func (u *loginUser) GetName() string {
	return u.Name
}

func (u *loginUser) GetPhone() string {
	return u.Phone
}

func (u *loginUser) GetAccount() string {
	return u.Account
}

func (u *loginUser) GetRegDate() *time.Time {
	return u.Regdate
}

func (u *loginUser) GetWork() string {
	return u.Work
}

func (u *loginUser) GetStatus() string {
	return u.Status
}

func (u *loginUser) GetUserType() string {
	return u.UserType
}

func (u *loginUser) GetTenantId() string {
	return u.TenantId
}

func (u *loginUser) GetTenantName() string {
	return u.TenantName
}

func (u *loginUser) GetTenantAccount() string {
	return u.TenantAccount
}
