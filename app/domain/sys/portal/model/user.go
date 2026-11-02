package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/enum"
)

type User struct {
	Base          `bson:",inline"`
	Account       string      `json:"account" gorm:"account" bson:"account"`
	Name          string      `json:"name" gorm:"name" bson:"name"`
	Password      string      `json:"password" gorm:"password" bson:"password"`
	Phone         string      `json:"phone" gorm:"phone" bson:"phone"`
	Email         string      `json:"email" gorm:"email" bson:"email"`
	Address       string      `json:"address" gorm:"address" bson:"address"`
	Gender        string      `json:"gender" gorm:"gender" bson:"gender"`
	Work          string      `json:"work" gorm:"work" bson:"work"`
	HeadPicture   string      `json:"headPicture" gorm:"head_picture" bson:"head_picture"`
	Status        enum.Status `json:"status" gorm:"status" bson:"status"`
	OryIdentityId string      `json:"oryIdentityId" gorm:"ory_identity_id" bson:"ory_identity_id"`
}

func NewUser() (*User, error) {
	return &User{}, nil
}

type UserView struct {
	Base          `bson:",inline"`
	Account       string      `json:"account" gorm:"account" bson:"account"`
	Name          string      `json:"name" gorm:"name" bson:"name"`
	Password      string      `json:"password" gorm:"password" bson:"password"`
	Phone         string      `json:"phone" gorm:"phone" bson:"phone"`
	Email         string      `json:"email" gorm:"email" bson:"email"`
	Address       string      `json:"address" gorm:"address" bson:"address"`
	Gender        string      `json:"gender" gorm:"gender" bson:"gender"`
	Work          string      `json:"work" gorm:"work" bson:"work"`
	HeadPicture   string      `json:"headPicture" gorm:"head_picture" bson:"head_picture"`
	Status        enum.Status `json:"status" gorm:"status" bson:"status"`
	OryIdentityId string      `json:"oryIdentityId" gorm:"ory_identity_id" bson:"ory_identity_id"`
	TenantUserId  string      `json:"tenantUserId" gorm:"tenant_user_id" bson:"tenant_user_id"`
	TanId         string      `json:"tanId" gorm:"tan_id" bson:"tan_id"`
	IsAdmin       bool        `json:"isAdmin" gorm:"is_admin" bson:"is_admin"`
}

func NewUserView() *UserView {
	return &UserView{}
}
