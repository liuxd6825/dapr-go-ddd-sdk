package model

import (
	client "github.com/ory/kratos-client-go"
)

type Login struct {
	Method     string `json:"method" gorm:"method" bson:"method"`
	Identifier string `json:"identifier" gorm:"identifier" bson:"identifier"`
	Password   string `json:"password" gorm:"password" bson:"password"`
}

type LoginUser struct {
	Id            string `json:"id" bson:"id"`
	TenantId      string `json:"tenantId" gorm:"tenant_id"  bson:"tenant_id"`        // 租户ID
	TenantName    string `json:"tenantName" gorm:"tenant_name"  bson:"tenant_name" ` // 租户ID
	Account       string `json:"account" gorm:"account" bson:"account"`
	Name          string `json:"name" gorm:"name" bson:"name"`
	Phone         string `json:"phone" gorm:"phone" bson:"phone"`
	Email         string `json:"email" gorm:"email" bson:"email"`
	Address       string `json:"address" gorm:"address" bson:"address"`
	Gender        string `json:"gender" gorm:"gender" bson:"gender"`
	Work          string `json:"work" gorm:"work" bson:"work"`
	HeadPicture   string `json:"headPicture" gorm:"head_picture" bson:"head_picture"`
	Status        string `json:"status" gorm:"status" bson:"status"`
	OryIdentityId string `json:"oryIdentityId" gorm:"ory_identity_id" bson:"ory_identity_id"`
}

type LoginResult struct {
	LoginSession *client.SuccessfulNativeLogin `json:"loginSession" bson:"login_session"`
	LoginJwt     string                        `json:"loginJwt" bson:"login_jwt"`
	LoginUser    LoginUser                     `json:"loginUser" bson:"login_user"`
}
