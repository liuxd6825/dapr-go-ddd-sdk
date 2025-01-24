package subevents

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
)

type Handler struct {
	service *server.Service
}

func NewEventHandler(service *server.Service) *Handler {
	return &Handler{service: service}
}
