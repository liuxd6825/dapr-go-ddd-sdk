package xbase

type Command[T any] struct {
	CommandId   string   `json:"commandId"  validate:"required" ` // 命令ID
	IsValidOnly bool     `json:"isValidOnly"  validate:"-"`       // 是否仅验证，不执行
	UpdateMask  []string `json:"updateMask"  validate:"-"`
	Data        T        `json:"data"  validate:"required"`
}

type IsValidOnly interface {
	GetIsValidOnly() bool
}

type GetCmdData interface {
	GetCmdData() any
}

type DeleteByIdCommand = Command[IdField]

type IdField struct {
	Id string `json:"id" validate:"required"`
}

type IdsField struct {
	Id []string `json:"ids" validate:"required"`
}

func (c *Command[T]) GetCommandId() string {
	return c.CommandId
}

func (c *Command[T]) GetVersion() string {
	return "v1"
}

func (c *Command[T]) GetIsValidOnly() bool {
	return c.IsValidOnly
}

func (c *Command[T]) SetData(data T) *Command[T] {
	c.Data = data
	return c
}

func (c *Command[T]) GetData() T {
	return c.Data
}

func (c *Command[T]) GeCmdData() any {
	return c.Data
}

func (c *Command[T]) SetUpdateMask(updateMask []string) *Command[T] {
	c.UpdateMask = updateMask
	return c
}

func (c *Command[T]) SetIsValidOnly(isValidOnly bool) *Command[T] {
	c.IsValidOnly = isValidOnly
	return c
}

func (c *Command[T]) SetCommandId(commandId string) *Command[T] {
	c.CommandId = commandId
	return c
}
