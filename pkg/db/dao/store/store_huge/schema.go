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
