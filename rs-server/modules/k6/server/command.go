package server

import "github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/common"

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
	return common.Object(c).GetString(common.CommandId)
}

func (c Command) SetCommandId(val string) {
	_ = common.Object(c).Set(common.CommandId, val)
}

func (c Command) GetTypeName() string {
	return common.Object(c).GetString(common.TypeName)
}

func (c Command) SetTypeName(val string) {
	_ = common.Object(c).Set(common.TypeName, val)
}

func (c Command) GetData() map[string]any {
	m, ok := common.Object(c).GetMap(common.Data)
	if !ok {
		return nil
	}
	return m
}

func (c Command) SetData(val map[string]any) {
	_ = common.Object(c).Set(common.Data, val)
}
