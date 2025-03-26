package entity

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/xinfra/enums"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/entities"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

type User struct {
	entities.Base
	Account      string           `gorm:"column:account" json:"account" jsonschema:"required"`
	Password     string           `gorm:"column:password" json:"password" jsonschema:"required"`
	Name         string           `gorm:"column:name" json:"name" jsonschema:"required"`
	Gender       enums.Gender     `gorm:"column:gender" json:"gender" json:"gender"`
	Phone        string           `gorm:"column:phone" json:"phone" jsonschema:"phone"`
	Email        string           `gorm:"column:email" json:"email" jsonschema:"email"`
	Work         string           `gorm:"column:work" json:"work" jsonschema:""`
	Address      string           `gorm:"column:address" json:"address" jsonschema:""`
	RegDate      *times.Time      `gorm:"column:reg_date" json:"regDate" jsonschema:""`
	UserStatus   enums.UserStatus `gorm:"column:user_status" json:"userStatus" jsonschema:""`
	HeadPicture  string           `gorm:"column:head_picture" json:"headPicture" jsonschema:""`
	ClientID     string           `gorm:"column:client_id" json:"clientId" jsonschema:""`
	ClientSecret string           `gorm:"column:client_secret" json:"clientSecret" jsonschema:""`
	RedirectUris string           `gorm:"column:redirect_uris" json:"redirectUris" jsonschema:""`
	Admin        bool             `gorm:"column:admin" json:"admin" jsonschema:""`
}
