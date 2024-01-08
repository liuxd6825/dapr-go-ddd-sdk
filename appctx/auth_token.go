package appctx

import (
	"encoding/base64"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"strings"
)

type AuthToken interface {
	GetSub() string
	GetExp() int
	GetUser() AuthUser
	GetClientId() string
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
	GetTenantAccount() string
}

type authToken struct {
	Sub      string    `json:"sub"`
	Exp      int       `json:"exp"`
	User     *authUser `json:"user"`
	ClientId string    `json:"client_id"`
}

type authUser struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Account       string `json:"account"`
	Regdate       string `json:"regdate"`
	Work          string `json:"work"`
	Status        string `json:"status"`
	UserType      string `json:"userType"`
	TenantId      string `json:"tenantId"`
	TenantName    string `json:"tenantName"`
	TenantAccount string `json:"tenantAccount"`
}

func getToken(jwtText string) (AuthToken, error) {
	list := strings.Split(jwtText, ".")
	if len(list) != 3 {
		return nil, errors.New("token格式不正确")
	}
	tk := authToken{}
	bs, err := base64.RawURLEncoding.DecodeString(list[1])
	if err != nil {
		return nil, err
	}
	err = jsonutils.Unmarshal(bs, &tk)
	return &tk, err
}

func newToken() AuthToken {
	return &authToken{}
}

func (u *authToken) GetSub() string {
	return u.Sub
}

func (u *authToken) GetExp() int {
	return u.Exp
}

func (u *authToken) GetUser() AuthUser {
	return u.User
}

func (u *authToken) GetClientId() string {
	return u.ClientId
}

func (u *authUser) GetId() string {
	return u.Id
}

func (u *authUser) GetName() string {
	return u.Name
}

func (u *authUser) GetPhone() string {
	return u.Phone
}

func (u *authUser) GetAccount() string {
	return u.Account
}

func (u *authUser) GetRegDate() string {
	return u.Regdate
}

func (u *authUser) GetWork() string {
	return u.Work
}

func (u *authUser) GetStatus() string {
	return u.Status
}

func (u *authUser) GetUserType() string {
	return u.UserType
}

func (u *authUser) GetTenantId() string {
	return u.TenantId
}

func (u *authUser) GetTenantName() string {
	return u.TenantName
}

func (u *authUser) GetTenantAccount() string {
	return u.TenantAccount
}
