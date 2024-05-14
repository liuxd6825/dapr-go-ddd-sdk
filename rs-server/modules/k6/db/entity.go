package db

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/rs-server/modules/common"
)

type Entity common.Object

func NewEntity() Entity {
	return Entity{}
}

func (e Entity) GetTenantId() string {
	return common.Object(e).GetString(common.TenantId)
}

func (e Entity) SetTenantId(val string) {
	_ = common.Object(e).Set(common.TenantId, val)
}

func (e Entity) GetId() string {
	return common.Object(e).GetString(common.Id)
}

func (e Entity) SetId(val string) {
	_ = common.Object(e).Set(common.Id, val)
}

func (e Entity) GetMapValues() map[string]any {
	return e
}

func (e Entity) SetMapValues(vals map[string]any) {
	for k, v := range vals {
		_ = common.Object(e).Set(k, v)
	}
}
