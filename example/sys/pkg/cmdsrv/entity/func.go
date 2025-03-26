package entity

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/example/sys/pkg/xinfra/enums"
)

type Func struct {
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Desc   string           `json:"desc"`
	Order  int              `json:"order"`
	Status enums.FuncStatus `json:"status"`
}
