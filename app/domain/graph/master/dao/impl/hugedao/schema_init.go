package hugedao

import (
	"context"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	store_huge "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/store_huge"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
)

const (
	VertexLabelMaster = "master_node"
	VertexLabelSame   = "master_same"
	EdgeLabelSame     = "master_same"
)

// EnsureSchema 触发 schema 同步(add-only);
// SyncSchema 自身已幂等,无需本地 cache。
func EnsureSchema(ctx context.Context, client *hugegraph.CommonClient, graphMeta *schema.Graph) error {
	if client == nil || graphMeta == nil {
		return nil
	}
	return SyncMasterSchema(ctx, client)
}

// SyncMasterSchema 将默认 master schema spec 同步到 HugeGraph(add-only)。
// 推荐作为 DAO 启动时调用入口,会对比实际 schema 并补齐缺失项。
func SyncMasterSchema(ctx context.Context, client *hugegraph.CommonClient) error {
	return store_huge.SyncSchema(ctx, client, DefaultMasterSchemaSpec())
}

// DefaultMasterSchemaSpec 返回 master graph 默认期望 schema 状态。
// rel_type 改为 EdgeLabel "master_rel" 的属性,通过索引加速按 relType 查询。
// node_type 用于区分主子表节点类型(person/company),即所有主表节点共用
// master VertexLabel,通过 node_type 属性区分业务实体类型。
func DefaultMasterSchemaSpec() store_huge.SchemaSpec {
	//AppendProperties("src_id", "src_type"). // 如果默认是不允许为空，则会报错
	masterVertexProps := []string{
		"id", "node_type",
		"name", "case_id", "tenant_id", "description",
		"type", "table", "keywords", "tags",
		"src_id", "src_type", "src_name", "src_url"}
	masterRelProps := []string{
		"rel_target_id", "rel_source_id", "rel_type", "case_id", "tenant_id",
		"keywords", "tags", "table", "description",
		"src_id", "src_type", "src_name", "src_url",
	}
	return store_huge.SchemaSpec{
		VertexLabels: []store_huge.VertexLabelSpec{
			{
				Name:        VertexLabelMaster,
				Properties:  masterVertexProps,
				PrimaryKeys: []string{"id"},
			},
			{
				Name:        VertexLabelSame,
				Properties:  []string{"name"},
				PrimaryKeys: []string{"name"},
			},
		},
		EdgeLabels: []store_huge.EdgeLabelSpec{
			{
				Name:        EdgeLabelSame,
				SourceLabel: VertexLabelMaster,
				TargetLabel: VertexLabelSame,
				Properties:  []string{"name"},
				Frequency:   store_huge.FrequencySingle,
			},
			{
				Name:        "master_rel",
				SourceLabel: VertexLabelMaster,
				TargetLabel: VertexLabelMaster,
				Properties:  masterRelProps,
				Frequency:   store_huge.FrequencySingle,
			},
		},
	}
}

// EnsureMasterSchema 兼容层:旧 API 内部转调 SyncMasterSchema。
// relTypeName 参数已废弃(rel_type 现在是边的属性,不再独立 EdgeLabel),保留仅为向后兼容。
func EnsureMasterSchema(ctx context.Context, client *hugegraph.CommonClient, relTypeName string) error {
	_ = relTypeName
	return SyncMasterSchema(ctx, client)
}
