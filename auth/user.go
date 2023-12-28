package auth

import (
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

type user struct {
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
