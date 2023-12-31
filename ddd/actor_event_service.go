package ddd

import (
	"github.com/liuxd6825/go-sdk/actor"
	dapr "github.com/liuxd6825/go-sdk/client"
)

const eventActorType = "ddd.EventActorType"

type EventActorService struct {
	actor.ServerImplBase
	daprClient dapr.Client
}
