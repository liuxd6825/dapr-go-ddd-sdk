package store_huge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/edgelabel"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/propertykey"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/vertexlabel"
)

const (
	VertexLabelMaster = "master"
	VertexLabelSame   = "same"
	EdgeLabelSame     = "same"

	PropDataTypeText      = "TEXT"
	PropCardinalitySingle = "SINGLE"
	PropCardinalityList   = "LIST"
	IDStrategyDefault     = "DEFAULT"
	FrequencySingle       = "SINGLE"
	FrequencyMultiple     = "MULTIPLE"
)

var DefaultPropertyKeys = []propertykey.CreateRequestData{
	{Name: "id", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "name", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "case_id", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "tenant_id", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "source_ids", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "source_type", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "description", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "type", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "table", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "source", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "target", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "keywords", DataType: PropDataTypeText, Cardinality: PropCardinalityList},
	{Name: "rel_type", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
	{Name: "node_type", DataType: PropDataTypeText, Cardinality: PropCardinalitySingle},
}

func EnsurePropertyKeys(_ context.Context, client *hugegraph.CommonClient, list []propertykey.CreateRequestData) error {
	if client == nil {
		return fmt.Errorf("client is nil")
	}
	for _, pk := range list {
		_, err := client.Propertykey.Create(
			client.Propertykey.Create.WithReqData(pk),
		)
		if err != nil {
			if isAlreadyExists(err) {
				continue
			}
			return fmt.Errorf("create propertyKey %s error: %w", pk.Name, err)
		}
	}
	return nil
}

func EnsureVertexLabel(_ context.Context, client *hugegraph.CommonClient, name string, properties []string, primaryKeys []string) error {
	req := vertexlabel.CreateRequestData{
		Name:             name,
		IDStrategy:       IDStrategyDefault,
		Properties:       properties,
		PrimaryKeys:      primaryKeys,
		EnableLabelIndex: true,
	}
	_, err := client.VertexLabel.Create(
		client.VertexLabel.Create.WithReqData(req),
	)
	if err != nil {
		if isAlreadyExists(err) {
			return nil
		}
		return fmt.Errorf("create vertexLabel %s error: %w", name, err)
	}
	return nil
}

func EnsureEdgeLabel(_ context.Context, client *hugegraph.CommonClient, name, sourceLabel, targetLabel string, properties []string, frequency string) error {
	req := edgelabel.CreateRequestData{
		Name:             name,
		SourceLabel:      sourceLabel,
		TargetLabel:      targetLabel,
		Properties:       properties,
		EnableLabelIndex: true,
	}
	v := reflect.ValueOf(&req).Elem().FieldByName("Frequency")
	if v.IsValid() && v.CanSet() {
		v.SetString(frequency)
	}
	_, err := client.EdgeLabel.Create(
		client.EdgeLabel.Create.WithReqData(req),
	)
	if err != nil {
		if isAlreadyExists(err) {
			return nil
		}
		return fmt.Errorf("create edgeLabel %s error: %w", name, err)
	}
	return nil
}

// EnsureVertexIndex 通过 REST API 直接创建索引 label,SDK 未封装该 API
// baseType 取值 "VERTEX_LABEL" 或 "EDGE_LABEL"
func EnsureVertexIndex(ctx context.Context, client *hugegraph.CommonClient, labelName, indexName, property string) error {
	return EnsureIndex(ctx, client, "VERTEX_LABEL", labelName, indexName, property)
}

// EnsureEdgeIndex 创建 EDGE_LABEL 上的索引
func EnsureEdgeIndex(ctx context.Context, client *hugegraph.CommonClient, labelName, indexName, property string) error {
	return EnsureIndex(ctx, client, "EDGE_LABEL", labelName, indexName, property)
}

// EnsureIndex 通用索引创建函数
func EnsureIndex(ctx context.Context, client *hugegraph.CommonClient, baseType, labelName, indexName, property string) error {
	if client == nil {
		return fmt.Errorf("client is nil")
	}
	body := map[string]any{
		"name":       indexName,
		"base_type":  baseType,
		"base_value": labelName,
		"index_type": "SECONDARY",
		"fields":     []string{property},
	}
	data, _ := json.Marshal(body)
	cfg := client.Transport.GetConfig()
	host := cfg.URL.Host
	graph := cfg.Graph
	url := fmt.Sprintf("%s://%s/graphs/%s/schema/indexlabels",
		cfg.URL.Scheme, host, graph)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.Username != "" {
		req.SetBasicAuth(cfg.Username, cfg.Password)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}()
	if resp.StatusCode == http.StatusAccepted ||
		resp.StatusCode == http.StatusConflict ||
		resp.StatusCode == http.StatusCreated ||
		resp.StatusCode == http.StatusOK {
		return nil
	}
	rb, _ := io.ReadAll(resp.Body)
	bodyStr := string(rb)
	if containsCI(bodyStr, "has existed") ||
		containsCI(bodyStr, "already exists") ||
		containsCI(bodyStr, "no need to build index") {
		return nil
	}
	return fmt.Errorf("create indexLabel %s status=%d body=%s", indexName, resp.StatusCode, bodyStr)
}

// DefaultMasterSchemaSpec 返回 master graph 默认期望 schema 状态。
// rel_type 改为 EdgeLabel "master_rel" 的属性,通过索引加速按 relType 查询。
// node_type 用于区分主子表节点类型(person/company),即所有主表节点共用
// master VertexLabel,通过 node_type 属性区分业务实体类型。
func DefaultMasterSchemaSpec() SchemaSpec {
	masterVertexProps := []string{
		"id", "name", "case_id", "tenant_id", "source_ids",
		"source_type", "description", "type", "node_type", "table",
	}
	masterRelProps := []string{
		"id", "rel_type", "source", "target", "keywords",
		"case_id", "tenant_id", "source_ids", "source_type",
		"description", "table",
	}
	return SchemaSpec{
		PropertyKeys: DefaultPropertyKeys,
		VertexLabels: []VertexLabelSpec{
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
		EdgeLabels: []EdgeLabelSpec{
			{
				Name:        EdgeLabelSame,
				SourceLabel: VertexLabelMaster,
				TargetLabel: VertexLabelSame,
				Properties:  []string{"name"},
				Frequency:   FrequencySingle,
			},
			{
				Name:        "master_rel",
				SourceLabel: VertexLabelMaster,
				TargetLabel: VertexLabelMaster,
				Properties:  masterRelProps,
				Frequency:   FrequencySingle,
			},
		},
		IndexLabels: []IndexLabelSpec{
			{
				Name: "masterById", BaseLabel: VertexLabelMaster, BaseType: "VERTEX_LABEL",
				Fields: []string{"id"},
			},
			{
				Name: "masterByName", BaseLabel: VertexLabelMaster, BaseType: "VERTEX_LABEL",
				Fields: []string{"name"},
			},
			{
				Name: "masterByNodeType", BaseLabel: VertexLabelMaster, BaseType: "VERTEX_LABEL",
				Fields: []string{"type"},
			},
			{
				Name: "masterRelByRelType", BaseLabel: "master_rel", BaseType: "EDGE_LABEL",
				Fields: []string{"rel_type"},
			},
			{
				Name: "masterRelById", BaseLabel: "master_rel", BaseType: "EDGE_LABEL",
				Fields: []string{"id"},
			},
			{
				Name: "sameByName", BaseLabel: VertexLabelSame, BaseType: "VERTEX_LABEL",
				Fields: []string{"name"},
			},
		},
	}
}

// SyncMasterSchema 将默认 master schema spec 同步到 HugeGraph(add-only)。
// 推荐作为 DAO 启动时调用入口,会对比实际 schema 并补齐缺失项。
func SyncMasterSchema(ctx context.Context, client *hugegraph.CommonClient) error {
	return SyncSchema(ctx, client, DefaultMasterSchemaSpec())
}

// EnsureMasterSchema 兼容层:旧 API 内部转调 SyncMasterSchema。
// relTypeName 参数已废弃(rel_type 现在是边的属性,不再独立 EdgeLabel),保留仅为向后兼容。
func EnsureMasterSchema(ctx context.Context, client *hugegraph.CommonClient, relTypeName string) error {
	_ = relTypeName
	return SyncMasterSchema(ctx, client)
}

func isAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	for _, key := range []string{"already exists", "ALREADY_EXISTS", "exists", "EXISTS"} {
		if containsCI(s, key) {
			return true
		}
	}
	return false
}

func containsCI(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i+len(substr) <= len(s); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			a, b := s[i+j], substr[j]
			if a >= 'A' && a <= 'Z' {
				a += 32
			}
			if b >= 'A' && b <= 'Z' {
				b += 32
			}
			if a != b {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
