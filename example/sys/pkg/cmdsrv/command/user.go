package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/cmdsrv/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/xinfra/xcmd"
)

// UserCreateCmd 用户创建
type UserCreateCmd struct {
	xcmd.Cmd[entity.User]
}

// UserUpdateCmd 用户更新
type UserUpdateCmd struct {
	xcmd.Cmd[entity.User]
}

// UserSetPasswordCmd 用户更新密码
type UserSetPasswordCmd struct {
	xcmd.Cmd[entity.User]
}

// UserDeleteCmd 用户被删除
type UserDeleteCmd struct {
	xcmd.Cmd[entity.User]
}

func NewUserCreateCmd() *UserCreateCmd {
	c := &UserCreateCmd{}
	c.SetDefaults()
	return c
}

func (c *UserCreateCmd) SetDefaults() {
	c.Cmd.EventType = "user-create"
}

func NewUserUpdateCmd() *UserUpdateCmd {
	c := &UserUpdateCmd{}
	c.SetDefaults()
	return c
}

func (c *UserUpdateCmd) SetDefaults() {
	c.Cmd.EventType = "user-update"
}

func NewUserDeleteCmd() *UserDeleteCmd {
	c := &UserDeleteCmd{}
	c.SetDefaults()
	return c
}

func (c *UserDeleteCmd) SetDefaults() {
	c.Cmd.EventType = "user-delete"
}

func NewUserSetPasswordCmd() *UserSetPasswordCmd {
	res := &UserSetPasswordCmd{}
	res.SetDefaults()
	return res
}

func (c *UserSetPasswordCmd) SetDefaults() {
	c.Cmd.EventType = "user-set-password"
}
