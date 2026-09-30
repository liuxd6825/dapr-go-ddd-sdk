package hugegraph

import (
	"context"
	"testing"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
)

func Test_Client(t *testing.T) {
	clientManager := NewMultiTenantManager(GetPoolConfigDefault(), TenantConfig{
		Host:     "192.168.120.200",
		Port:     18080,
		Graph:    "",
		Username: "admin",
		Password: "admin",
	})

	base := NewClient(clientManager)
	ctx := appctx.NewTenantContext(context.Background(), "dev")
	// g.V().union(identity(), outE())
	// "g.V().both().path().limit(3)"
	data, err := base.gremlin(ctx, "g.V().union(identity(), outE())", nil)
	if err != nil {
		t.Errorf("gremlin err: %s", err.Error())
	}
	t.Log(data.Result.Data)

	result := NewResultWithResponse(data)
	t.Log("paths:", result.Paths)
	t.Log("nodes:", result.Nodes)
	t.Log("edges:", result.Edges)
	t.Log("records:", result.Records)
}

func newClient() *hugegraph.CommonClient {
	client, err := hugegraph.NewCommonClient(hugegraph.Config{
		Host:     "192.168.120.200",
		Port:     18080,
		Graph:    "dev",
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		panic(err)
	}
	return client
}
