package ddd

import (
	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"
)

const eventActorType = "ddd.EventActorType"

type EventActorService struct {
	actor.ServerImplBase
	daprClient dapr.Client
}
