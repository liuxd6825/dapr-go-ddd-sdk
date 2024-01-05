package auth

import (
	"encoding/base64"
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"strings"
	"time"
)

type Token interface {
	GetSub() string
	GetExp() int
	GetUser() User
	GetClientId() string
}

type User interface {
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

type token struct {
	Sub      string `json:"sub"`
	Exp      int    `json:"exp"`
	User     *user  `json:"user"`
	ClientId string `json:"client_id"`
}

type user struct {
	Id            string     `json:"id"`
	Name          string     `json:"name"`
	Phone         string     `json:"phone"`
	Account       string     `json:"account"`
	Regdate       *time.Time `json:"regdate"`
	Work          string     `json:"work"`
	Status        string     `json:"status"`
	UserType      string     `json:"userType"`
	TenantId      string     `json:"tenantId"`
	TenantName    string     `json:"tenantName"`
	TenantAccount string     `json:"tenantAccount"`
}

func getToken(jwtText string) (Token, error) {
	list := strings.Split(jwtText, ".")
	if len(list) != 3 {
		return nil, errors.New("token格式不正确")
	}
	tk := token{}
	bs, err := base64.RawURLEncoding.DecodeString(list[1])
	if err != nil {
		return nil, err
	}
	err = jsonutils.Unmarshal(bs, &tk)
	return &tk, err
}

func newToken() Token {
	return &token{}
}

func (u *token) GetSub() string {
	return u.Sub
}

func (u *token) GetExp() int {
	return u.Exp
}

func (u *token) GetUser() User {
	return u.User
}

func (u *token) GetClientId() string {
	return u.ClientId
}

func (u *user) GetId() string {
	return u.Id
}

func (u *user) GetName() string {
	return u.Name
}

func (u *user) GetPhone() string {
	return u.Phone
}

func (u *user) GetAccount() string {
	return u.Account
}

func (u *user) GetRegDate() *time.Time {
	return u.Regdate
}

func (u *user) GetWork() string {
	return u.Work
}

func (u *user) GetStatus() string {
	return u.Status
}

func (u *user) GetUserType() string {
	return u.UserType
}

func (u *user) GetTenantId() string {
	return u.TenantId
}

func (u *user) GetTenantName() string {
	return u.TenantName
}

func (u *user) GetTenantAccount() string {
	return u.TenantAccount
}
