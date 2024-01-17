package restapp

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-sdk/actor"
	"sync"
)

var _actorsFactory []actor.FactoryContext
var _actorsFactoryOnce sync.Once

func RegisterActor(actorServer actor.ServerContext) {
	_actorsFactory = append(getActorsFactory(), func() actor.ServerContext { return actorServer })
}

func GetActors() []actor.FactoryContext {
	return getActorsFactory()
}

func getActorsFactory() []actor.FactoryContext {
	_actorsFactoryOnce.Do(func() {
		_actorsFactory = append(_actorsFactory, newAggregateSnapshotActorFactory)
	})
	return _actorsFactory
}

func newAggregateSnapshotActorFactory() actor.ServerContext {
	client, err := daprclient.GetDaprDDDClient().DaprClient()
	if err != nil {
		panic(err)
	}
	return ddd.NewAggregateSnapshotActorServer(client)
}
