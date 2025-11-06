package tasks

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	tclient "go.temporal.io/sdk/client"
	tworker "go.temporal.io/sdk/worker"
)

type ClientProvider struct {
	clients *types.CMap[tclient.Client]
	workers *types.CMap[tworker.Worker]
	cfg     Config
	mtx     sync.Mutex
}

func NewClientProvider(cfg Config) *ClientProvider {
	return &ClientProvider{
		cfg:     cfg,
		clients: types.NewCMap[tclient.Client](),
		workers: types.NewCMap[tworker.Worker](),
	}
}

/*
func (c *ClientProvider) GetClient(ctx context.Context) tclient.Client {
	tenantId, ok := appctx.GetTenantId(ctx)
	if !ok {
		panic(errors.New("tenant id not found"))
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.clients.Has(tenantId) {
		client, _ := c.clients.Get(tenantId)
		return client
	}
	client, err := c.newClient(tenantId)
	if err != nil {
		panic(err)
	}
	c.clients.Set(tenantId, client)
	return client
}*/
