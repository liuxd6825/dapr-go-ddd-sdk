package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type UserCreateCommand struct {
	xbase.Command[model.User]
}
type UserUpdateCommand struct {
	xbase.Command[model.User]
}
type UserDeleteCommand struct {
	xbase.Command[model.User]
	//xbase.DeleteByIdCommand
}
type UserDeleteIdentityCommand struct {
	xbase.DeleteByIdCommand
}

type UserDeleteBatchCommand struct {
	xbase.Command[model.User]
}

type UpdatePasswordCommand struct {
	xbase.Command[UpdatePassword]
}

type UpdatePassword struct {
	Account         string `json:"account" gorm:"account" bson:"account"  validate:"required"`
	Password        string `json:"password" gorm:"password" bson:"password"  validate:"required"`
	ConfirmPassword string `json:"confirmPassword" gorm:"confirm_password" bson:"confirm_password"  validate:"required"`
}
