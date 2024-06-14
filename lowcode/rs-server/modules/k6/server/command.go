package server

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
)

/*
type Command common.Object

func NewCommand() Command {
	return Command{}
}

func (c Command) GetTenantId() string {
	return common.Object(c).GetString(common.TenantId)
}

func (c Command) SetTenantId(val string) {
	_ = common.Object(c).Set(common.TenantId, val)
}

func (c Command) GetCommandId() string {
	return common.Object(c).GetString(common.EventId)
}

func (c Command) SetCommandId(val string) {
	_ = common.Object(c).Set(common.EventId, val)
}

func (c Command) GetTypeName() string {
	return common.Object(c).GetString(common.TypeName)
}

func (c Command) SetTypeName(val string) {
	_ = common.Object(c).Set(common.TypeName, val)
}

func (c Command) GetIsValidOnly() bool {
	b, _ := common.Object(c).GetBool("isValidOnly")
	return b
}

func (c Command) SetIsValidOnly(val bool) {
	_ = common.Object(c).Set("isValidOnly", val)
}

func (c Command) GetData() map[string]any {
	v := common.Object(c).Get("data")
	switch v.(type) {
	case map[string]any:
		return v.(map[string]any)
	case Command:
		return v.(map[string]interface{})
	default:
		return nil
	}
	return nil
}

func (c Command) SetData(val map[string]any) {
	_ = common.Object(c).Set(common.Data, val)
}
*/

type CommandData common.Object

type Command struct {
	CommandId   string
	CommandType string
	TenantId    string
	IsValidOnly bool
	IdKey       string
	Data        map[string]any
}

func NewCommand() *Command {
	return &Command{}
}

func (c *Command) GetAggregateId() ddd.AggregateId {
	id, _ := c.getAggregateId()
	return ddd.NewAggregateId(id)
}

func (c *Command) getAggregateId() (string, error) {
	if c.IdKey == "" {
		return "", errors.New("idKey is required")
	}
	val := c.Data[c.IdKey]
	if val == nil {
		return "", errors.New("idKey is required")
	}
	id := stringutils.AnyToString(val)
	return id, nil
}

func (c *Command) Validate() error {
	_, err := c.getAggregateId()
	return err
}

func (c *Command) GetCommandId() string {
	return c.CommandId
}

func (c *Command) SetCommandId(val string) {
	c.CommandId = val
}

func (c *Command) GetTenantId() string {
	return c.TenantId
}

func (c *Command) SetTenantId(val string) {
	c.TenantId = val
}

func (c *Command) GetCommandType() string {
	return c.CommandType
}

func (c *Command) SetCommandType(val string) {
	c.CommandType = val
}

func (c *Command) GetData() map[string]any {
	return c.Data
}

func (c *Command) SetData(val map[string]any) {
	c.Data = val
}

func (c *Command) GetIsValidOnly() bool {
	return c.IsValidOnly
}

func (c *Command) SetIsValidOnly(val bool) {
	c.IsValidOnly = val
}
