package subevents

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
)

type SubEvents struct {
	appId    string
	service  *server.Service
	handlers *types.CMap[*Handler]
}

func NewSubEvents(service *server.Service) *SubEvents {
	return &SubEvents{service: service, handlers: types.NewCMap[*Handler]()}
}
