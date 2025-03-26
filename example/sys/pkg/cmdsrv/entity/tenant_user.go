package entity

import "github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/xinfra/enums"

type TenantUser struct {
	Id            string             `json:"id" gorm:"primary_key"`
	UserId        string             `json:"userId"`
	Account       string             `json:"account"`
	TenantId      string             `json:"tenantId"`
	TenantName    string             `json:"tenantName"`
	TenantAccount string             `json:"tenantAccount"`
	TenantStatus  enums.TenantStatus `json:"tenantStatus"`
	UserStatus    enums.UserStatus   `json:"userStatus"`
	DefaultTenant bool               `json:"defaultTenant"`
	UserType      enums.UserType     `json:"userType"`
}
