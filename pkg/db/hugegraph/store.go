package hugegraph

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/gremlin"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
)

type Store struct {
	pool *ClientPool
}

func NewStore(pool *ClientPool) *Store {
	return &Store{pool: pool}
}

// Write 执行 Gremlin 脚本,返回原始 Result.Data (interface{})
func (d *Store) Write(ctx context.Context, script string, bindings ...map[string]any) (*Result, error) {
	data, err := d.gremlin(ctx, script, getParams(bindings...))
	if err != nil {
		return nil, err
	}
	return NewResultWithResponse(data), nil
}

// Query 执行 Gremlin 脚本,返回原始 Result.Data (interface{})
func (d *Store) Query(ctx context.Context, script string, bindings ...map[string]any) (*Result, error) {
	data, err := d.gremlin(ctx, script, getParams(bindings...))
	if err != nil {
		return nil, err
	}
	return NewResultWithResponse(data), nil
}

// gremlin 执行 Gremlin 脚本,返回原始 Result.Data (interface{})
func (d *Store) gremlin(ctx context.Context, script string, bindings map[string]any) (*gremlin.PostResponseData, error) {
	client, err := d.getClient(ctx)
	if err != nil {
		return nil, err
	}

	tenantId := appctx.GetTenantId2(ctx)
	finalScript := formatScript(script, bindings)
	resp, err1 := client.Gremlin.Post(
		client.Gremlin.Post.WithGremlin(finalScript),
		client.Gremlin.Post.WithGraph(tenantId),
		client.Gremlin.Post.WithGraphSpace("DEFAULT"),
	)
	if err1 != nil {
		return nil, fmt.Errorf("gremlin post error: %w", err)
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.Data == nil {
		return nil, nil
	}
	if resp.Data.Exception != "" || resp.Data.Message != "" {
		return nil, fmt.Errorf("gremlin error: %s, cause: %s", resp.Data.Message, resp.Data.Cause)
	}
	return resp.Data, nil
}

// cypher 执行 cypher 脚本
func (d *Store) cypher(ctx context.Context, script string, bindings map[string]any) (*gremlin.PostResponseData, error) {
	client, err := d.getClient(ctx)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("hugegraph client is nil")
	}
	finalScript := formatScript(script, bindings)
	resp, err := client.Gremlin.Post(
		client.Gremlin.Post.WithGremlin(finalScript),
	)
	if err != nil {
		return nil, fmt.Errorf("gremlin post error: %w", err)
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.Data == nil {
		return nil, nil
	}
	if resp.Data.Exception != "" || resp.Data.Message != "" {
		return nil, fmt.Errorf("gremlin error: %s, cause: %s", resp.Data.Message, resp.Data.Cause)
	}
	return resp.Data, nil
}

func (d *Store) getClient(ctx context.Context) (*hugegraph.CommonClient, error) {
	client, err := d.pool.Get(ctx)
	return client, err
}

func getParams(bindings ...map[string]any) map[string]any {
	var params map[string]any
	if len(bindings) > 0 {
		params = bindings[0]
		for _, binding := range bindings[1:] {
			for k, v := range binding {
				params[k] = v
			}
		}
	}
	return params
}

// formatScript 将 bindings 嵌入 Gremlin 脚本,避免依赖 SDK 中尚未导出的 WithBindings 选项
// bindings 必须是字面量安全的值(string/number/bool);否则请预先 json.Marshal
func formatScript(script string, bindings map[string]any) string {
	if len(bindings) == 0 {
		return script
	}
	for k, v := range bindings {
		placeholder := "$<" + k + ">"
		if !strings.Contains(script, placeholder) {
			continue
		}
		script = strings.ReplaceAll(script, placeholder, literalOf(v))
	}
	return script
}

func literalOf(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		s := strings.ReplaceAll(x, "\\", "\\\\")
		s = strings.ReplaceAll(s, "'", "\\'")
		return "'" + s + "'"
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", x)
	}
}
