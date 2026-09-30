package store_huge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
)

// VertexLabelInfo 描述 HugeGraph 上一个 VertexLabel 的当前状态
type VertexLabelInfo struct {
	Name        string
	Properties  []string
	PrimaryKeys []string
}

// EdgeLabelInfo 描述 HugeGraph 上一个 EdgeLabel 的当前状态
type EdgeLabelInfo struct {
	Name       string
	Properties []string
	Frequency  string
}

// ListPropertyKeys 通过 REST API 列出已存在的 PropertyKey 名称集合
func ListPropertyKeys(ctx context.Context, client *hugegraph.CommonClient) (map[string]struct{}, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	body, err := doGet(ctx, client, "/schema/propertykeys")
	if err != nil {
		return nil, err
	}
	var resp struct {
		PropertyKeys []struct {
			Name string `json:"name"`
		} `json:"propertykeys"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode propertykeys: %w", err)
	}
	out := make(map[string]struct{}, len(resp.PropertyKeys))
	for _, pk := range resp.PropertyKeys {
		out[pk.Name] = struct{}{}
	}
	return out, nil
}

// GetVertexLabel 通过 REST API 拉取单个 VertexLabel 的元数据;
// 不存在时返回 (nil, nil)
func GetVertexLabel(ctx context.Context, client *hugegraph.CommonClient, name string) (*VertexLabelInfo, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	body, status, err := doGetWithStatus(ctx, client, "/schema/vertexlabels/"+name)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	var resp struct {
		Name        string   `json:"name"`
		Properties  []string `json:"properties"`
		PrimaryKeys []string `json:"primary_keys"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode vertexlabel %s: %w", name, err)
	}
	return &VertexLabelInfo{
		Name:        resp.Name,
		Properties:  resp.Properties,
		PrimaryKeys: resp.PrimaryKeys,
	}, nil
}

// GetEdgeLabel 通过 REST API 拉取单个 EdgeLabel 的元数据;
// 不存在时返回 (nil, nil)
func GetEdgeLabel(ctx context.Context, client *hugegraph.CommonClient, name string) (*EdgeLabelInfo, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	body, status, err := doGetWithStatus(ctx, client, "/schema/edgelabels/"+name)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	var resp struct {
		Name       string   `json:"name"`
		Properties []string `json:"properties"`
		Frequency  string   `json:"frequency"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode edgelabel %s: %w", name, err)
	}
	return &EdgeLabelInfo{
		Name:       resp.Name,
		Properties: resp.Properties,
		Frequency:  resp.Frequency,
	}, nil
}

// ListIndexLabels 列出已存在的 IndexLabel 名称集合
func ListIndexLabels(ctx context.Context, client *hugegraph.CommonClient) (map[string]struct{}, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	body, err := doGet(ctx, client, "/schema/indexlabels")
	if err != nil {
		return nil, err
	}
	var resp struct {
		IndexLabels []struct {
			Name string `json:"name"`
		} `json:"indexlabels"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode indexlabels: %w", err)
	}
	out := make(map[string]struct{}, len(resp.IndexLabels))
	for _, il := range resp.IndexLabels {
		out[il.Name] = struct{}{}
	}
	return out, nil
}

func doGet(ctx context.Context, client *hugegraph.CommonClient, path string) ([]byte, error) {
	body, _, err := doGetWithStatus(ctx, client, path)
	return body, err
}

func doGetWithStatus(ctx context.Context, client *hugegraph.CommonClient, path string) ([]byte, int, error) {
	cfg := client.Transport.GetConfig()
	url := fmt.Sprintf("%s://%s/graphs/%s%s", cfg.URL.Scheme, cfg.URL.Host, cfg.Graph, path)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Accept-Encoding", "gzip")
	if cfg.Username != "" {
		req.SetBasicAuth(cfg.Username, cfg.Password)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}()
	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}
