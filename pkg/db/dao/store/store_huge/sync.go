package store_huge

import (
	"context"
	"fmt"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/propertykey"
)

// SchemaSpec 描述期望的 schema 状态
type SchemaSpec struct {
	PropertyKeys []propertykey.CreateRequestData
	VertexLabels []VertexLabelSpec
	EdgeLabels   []EdgeLabelSpec
	IndexLabels  []IndexLabelSpec
}

// VertexLabelSpec 期望的 VertexLabel 状态
type VertexLabelSpec struct {
	Name        string
	Properties  []string
	PrimaryKeys []string
}

// EdgeLabelSpec 期望的 EdgeLabel 状态
type EdgeLabelSpec struct {
	Name        string
	SourceLabel string
	TargetLabel string
	Properties  []string
	Frequency   string
}

// IndexLabelSpec 期望的 IndexLabel 状态
type IndexLabelSpec struct {
	Name      string
	BaseLabel string
	BaseType  string // "VERTEX_LABEL" 或 "EDGE_LABEL"
	IndexType string // 通常 "SECONDARY"
	Fields    []string
}

// SyncSchema 同步 spec 到 HugeGraph,add-only:
//   1. 缺失的 PropertyKey → Create
//   2. 缺失的 VertexLabel → Create; 已存在但 properties 缺失 → Append
//   3. 缺失的 EdgeLabel → Create; 已存在但 properties 缺失 → Append
//   4. 缺失的 IndexLabel → Create
//
// 不删除任何已有 schema,确保数据安全。
func SyncSchema(ctx context.Context, client *hugegraph.CommonClient, spec SchemaSpec) error {
	if client == nil {
		return fmt.Errorf("client is nil")
	}

	// 1. PropertyKeys
	if err := syncPropertyKeys(ctx, client, spec.PropertyKeys); err != nil {
		return err
	}

	// 2. VertexLabels
	for _, vlSpec := range spec.VertexLabels {
		if err := syncVertexLabel(ctx, client, vlSpec); err != nil {
			return err
		}
	}

	// 3. EdgeLabels
	for _, elSpec := range spec.EdgeLabels {
		if err := syncEdgeLabel(ctx, client, elSpec); err != nil {
			return err
		}
	}

	// 4. IndexLabels
	if err := syncIndexLabels(ctx, client, spec.IndexLabels); err != nil {
		return err
	}

	return nil
}

func syncPropertyKeys(ctx context.Context, client *hugegraph.CommonClient, desired []propertykey.CreateRequestData) error {
	if len(desired) == 0 {
		return nil
	}
	existing, err := ListPropertyKeys(ctx, client)
	if err != nil {
		return fmt.Errorf("list propertykeys: %w", err)
	}
	var missing []propertykey.CreateRequestData
	for _, pk := range desired {
		if _, ok := existing[pk.Name]; !ok {
			missing = append(missing, pk)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return EnsurePropertyKeys(ctx, client, missing)
}

func syncVertexLabel(ctx context.Context, client *hugegraph.CommonClient, vlSpec VertexLabelSpec) error {
	actual, err := GetVertexLabel(ctx, client, vlSpec.Name)
	if err != nil {
		return fmt.Errorf("get vertexlabel %s: %w", vlSpec.Name, err)
	}
	if actual == nil {
		return EnsureVertexLabel(ctx, client, vlSpec.Name, vlSpec.Properties, vlSpec.PrimaryKeys)
	}
	missing := setDiff(vlSpec.Properties, actual.Properties)
	if len(missing) == 0 {
		return nil
	}
	return AppendVertexLabelProperties(ctx, client, vlSpec.Name, missing)
}

func syncEdgeLabel(ctx context.Context, client *hugegraph.CommonClient, elSpec EdgeLabelSpec) error {
	actual, err := GetEdgeLabel(ctx, client, elSpec.Name)
	if err != nil {
		return fmt.Errorf("get edgelabel %s: %w", elSpec.Name, err)
	}
	if actual == nil {
		return EnsureEdgeLabel(ctx, client, elSpec.Name, elSpec.SourceLabel, elSpec.TargetLabel, elSpec.Properties, elSpec.Frequency)
	}
	missing := setDiff(elSpec.Properties, actual.Properties)
	if len(missing) == 0 {
		return nil
	}
	return AppendEdgeLabelProperties(ctx, client, elSpec.Name, missing)
}

func syncIndexLabels(ctx context.Context, client *hugegraph.CommonClient, desired []IndexLabelSpec) error {
	if len(desired) == 0 {
		return nil
	}
	existing, err := ListIndexLabels(ctx, client)
	if err != nil {
		return fmt.Errorf("list indexlabels: %w", err)
	}
	for _, ilSpec := range desired {
		baseType := ilSpec.BaseType
		if baseType == "" {
			baseType = "VERTEX_LABEL"
		}
		indexType := ilSpec.IndexType
		if indexType == "" {
			indexType = "SECONDARY"
		}
		_ = indexType
		// 单字段索引用 ilSpec.Name 自身,与原 EnsureVertexIndex 兼容;
		// 多字段索引加 _<field> 后缀避免冲突
		for _, field := range ilSpec.Fields {
			var indexName string
			if len(ilSpec.Fields) == 1 {
				indexName = ilSpec.Name
			} else {
				indexName = ilSpec.Name + "_" + field
			}
			if _, ok := existing[indexName]; ok {
				continue
			}
			if err := EnsureIndex(ctx, client, baseType, ilSpec.BaseLabel, indexName, field); err != nil {
				return fmt.Errorf("create index %s on %s.%s (%s): %w", indexName, ilSpec.BaseLabel, field, baseType, err)
			}
		}
	}
	return nil
}

// setDiff 返回 desired - actual,即 desired 中存在但 actual 中不存在的元素
func setDiff(desired, actual []string) []string {
	actualSet := make(map[string]struct{}, len(actual))
	for _, a := range actual {
		actualSet[a] = struct{}{}
	}
	var missing []string
	for _, d := range desired {
		if _, ok := actualSet[d]; !ok {
			missing = append(missing, d)
		}
	}
	return missing
}
