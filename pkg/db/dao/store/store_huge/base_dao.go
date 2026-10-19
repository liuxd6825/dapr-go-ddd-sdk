package store_huge

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/gremlin"
	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
)

type BaseDao[T any] struct {
	client *hugegraph.CommonClient
	config *Config[T]
	label  string
}

type Config[T any] struct {
	DBSchema *store2.DBSchema
	Labels   []string
}

type GremlinOptions struct {
	Graph      string
	GraphSpace string
}

const DbType = "hugegraph"

func NewBaseDao[T any](client *hugegraph.CommonClient, config *Config[T]) *BaseDao[T] {
	if config == nil {
		panic("store_huge.NewClientDao: config is nil")
	}
	return &BaseDao[T]{
		client: client,
		config: config,
	}
}

func (d *BaseDao[T]) GetDBClient() any {
	return d.client
}

func (d *BaseDao[T]) GetSchema() *store2.DBSchema {
	if d.config == nil {
		return nil
	}
	return d.config.DBSchema
}

func (d *BaseDao[T]) GetDbType() string {
	return DbType
}

func (d *BaseDao[T]) GetLabels() []string {
	if d.config == nil {
		return []string{}
	}
	return d.config.Labels
}

func (d *BaseDao[T]) GetLabel() string {
	if d.config == nil || d.config.Labels == nil || len(d.config.Labels) == 0 {
		return ""
	}
	return d.config.Labels[0]
}

func (d *BaseDao[T]) GetHugeClient() *hugegraph.CommonClient {
	return d.client
}

// Gremlin 执行 Gremlin 脚本,返回原始 Result.Data (interface{})

func (d *BaseDao[T]) Gremlin(ctx context.Context, gremlinScript string, bindings map[string]any, options ...GremlinOptions) (*gremlin.PostResponseData, error) {
	_ = ctx
	if d.client == nil {
		return nil, fmt.Errorf("hugegraph client is nil")
	}
	post := d.client.Gremlin.Post
	finalScript := formatScript(gremlinScript, bindings)
	requests := []func(*gremlin.PostRequest){
		post.WithGremlin(finalScript),
	}
	for _, opt := range options {
		if opt.Graph != "" {
			requests = append(requests, post.WithGraph(opt.Graph))
		}
		if opt.GraphSpace != "" {
			requests = append(requests, post.WithGraphSpace(opt.GraphSpace))
		}
	}
	resp, err := post(requests...)

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
func (d *BaseDao[T]) Write(ctx context.Context, script string, bindings map[string]any) (*HugeResult, error) {
	data, err := d.Gremlin(ctx, script, bindings)
	if err != nil {
		return nil, err
	}
	return NewHugeResult(data), nil
}

// Query 读查询,返回 *HugeResult
func (d *BaseDao[T]) Query(ctx context.Context, script string, bindings map[string]any) (*HugeResult, error) {
	return d.Write(ctx, script, bindings)
}

func (d *BaseDao[T]) ToMap(entity T) (map[string]any, error) {
	return d.ToMap2(entity, d.GetSchema())
}

func (d *BaseDao[T]) ToMap2(entity T, schema *store2.DBSchema) (map[string]any, error) {
	data := make(map[string]any)
	for _, field := range schema.Fields {
		val := reflectutils.GetField(entity, field.Name)
		data[field.DBName] = val
	}
	return data, nil
}

func (d *BaseDao[T]) BuildGremlin(data map[string]any) string {
	sb := strings.Builder{}
	count := len(data)
	i := 0
	for k, v := range data {
		i++
		key := k
		if key == "id" {
			key = "T.id"
		} else {
			key = "'" + key + "'"
		}
		val := v.(string)
		sb.WriteString(fmt.Sprintf("property(%s,'%s')", key, escapeGremlin(val)))
		if i < count {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

func escapeGremlin(s string) string {
	out := ""
	for _, r := range s {
		switch r {
		case '\'':
			out += "\\'"
		case '\\':
			out += "\\\\"
		default:
			out += string(r)
		}
	}
	return out
}
