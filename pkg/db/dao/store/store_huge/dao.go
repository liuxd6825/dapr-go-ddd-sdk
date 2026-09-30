package store_huge

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/gremlin"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
)

type Dao[T any] struct {
	Client *hugegraph.CommonClient
	config *Config[T]
}

type Config[T any] struct {
	DBSchema *store2.DBSchema
	Labels   []string
}

func NewClientDao[T any](client *hugegraph.CommonClient, config *Config[T]) *Dao[T] {
	if config == nil {
		panic("store_huge.NewClientDao: config is nil")
	}
	return &Dao[T]{
		Client: client,
		config: config,
	}
}

func (d *Dao[T]) GetDBClient() any {
	return d.Client
}

func (d *Dao[T]) GetSchema() *store2.DBSchema {
	if d.config == nil {
		return nil
	}
	return d.config.DBSchema
}

func (d *Dao[T]) GetDbType() string {
	return "hugedao"
}

func (d *Dao[T]) GetLabels() []string {
	if d.config == nil {
		return nil
	}
	return d.config.Labels
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

// Gremlin 执行 Gremlin 脚本,返回原始 Result.Data (interface{})
func (d *Dao[T]) Gremlin(ctx context.Context, script string, bindings map[string]any) (*gremlin.PostResponseData, error) {
	_ = ctx
	if d.Client == nil {
		return nil, fmt.Errorf("hugegraph client is nil")
	}
	finalScript := formatScript(script, bindings)
	resp, err := d.Client.Gremlin.Post(
		d.Client.Gremlin.Post.WithGremlin(finalScript),
	)
	if err != nil {
		return nil, fmt.Errorf("gremlin post error: %w", err)
	}
	if resp.Body != nil {
		_ = resp.Body.Close()
	}
	if resp.Data == nil {
		return nil, nil
	}
	if resp.Data.Exception != "" || resp.Data.Message != "" {
		return nil, fmt.Errorf("gremlin error: %s, cause: %s", resp.Data.Message, resp.Data.Cause)
	}
	return resp.Data, nil
}

// Write 兼容 store_neo4j.Write 命名,内部直接走 Gremlin 脚本
func (d *Dao[T]) Write(ctx context.Context, script string, bindings map[string]any) (*HugeResult, error) {
	data, err := d.Gremlin(ctx, script, bindings)
	if err != nil {
		return nil, err
	}
	return NewHugeResult(data), nil
}

// Query 读查询,返回 *HugeResult
func (d *Dao[T]) Query(ctx context.Context, script string, bindings map[string]any) (*HugeResult, error) {
	return d.Write(ctx, script, bindings)
}
